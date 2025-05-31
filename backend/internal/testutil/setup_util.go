// testutil/db.go
package testutil

import (
	"context"
	"fmt"
	"testing"

	"stride-wars-app/ent"
	"stride-wars-app/ent/enttest"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TestServices holds a brand-new in-memory ent.Client plus whichever
// repos/services you need for your tests.
type TestServices struct {
	Ctx    context.Context
	Client *ent.Client

	UserRepo           repository.UserRepository
	ActivityRepo       repository.ActivityRepository
	HexRepo            repository.HexRepository
	HexInfluenceRepo   repository.HexInfluenceRepository
	HexLeaderboardRepo repository.HexLeaderboardRepository

	UserService           *service.UserService
	ActivityService       *service.ActivityService
	HexService            *service.HexService
	HexInfluenceService   *service.HexInfluenceService
	HexLeaderboardService *service.HexLeaderboardService
}

// NewTestServices spins up a fresh in-memory SQLite (cache=private) and
// migrates the schema. It then constructs all “standard” Repos/Services
// in one place. Because cache=private + a unique UUID is used, every call
// to NewTestServices(t) truly gets its own private DB.
func NewTestServices(t *testing.T) *TestServices {
	t.Helper()

	dbName := fmt.Sprintf("file:ent_%s?mode=memory&cache=private&_fk=1", uuid.New().String())
	client := enttest.Open(t, "sqlite3", dbName)

	ctx := context.Background()
	if err := client.Schema.Create(ctx, schema.WithForeignKeys(true)); err != nil {
		t.Fatalf("failed to migrate test DB schema: %v", err)
	}

	userRepo := repository.NewUserRepository(client)
	activityRepo := repository.NewActivityRepository(client)
	hexRepo := repository.NewHexRepository(client)
	hexInfluenceRepo := repository.NewHexInfluenceRepository(client)
	hexLeaderboardRepo := repository.NewHexLeaderboardRepository(client)

	logger := zap.NewExample()
	userService := service.NewUserService(userRepo, logger)
	activityService := service.NewActivityService(
		activityRepo,
		hexRepo,
		hexInfluenceRepo,
		hexLeaderboardRepo,
		userRepo,
		*userService,
		logger,
	)
	hexService := service.NewHexService(hexRepo, logger)
	hexInfluenceService := service.NewHexInfluenceService(hexInfluenceRepo, logger)
	hexLeaderboardService := service.NewHexLeaderboardService(
		hexLeaderboardRepo,
		hexInfluenceRepo,
		logger,
	)

	return &TestServices{
		Ctx:                   ctx,
		Client:                client,
		UserRepo:              userRepo,
		ActivityRepo:          activityRepo,
		HexRepo:               hexRepo,
		HexInfluenceRepo:      hexInfluenceRepo,
		HexLeaderboardRepo:    hexLeaderboardRepo,
		UserService:           userService,
		ActivityService:       activityService,
		HexService:            hexService,
		HexInfluenceService:   hexInfluenceService,
		HexLeaderboardService: hexLeaderboardService,
	}
}
