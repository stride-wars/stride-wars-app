package repository_test

import (
	"context"
	"testing"

	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/testutil"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"time"
)

// TODO to properly test this we probably need to set up local postgres container????

func setupTestActivityHandler(t *testing.T) (context.Context, *ent.Client, repository.ActivityRepository, repository.UserRepository) {
	t.Helper()

	svc := testutil.NewTestServices(t)
	activityRepo := repository.NewActivityRepository(svc.Client)
	userRepo := repository.NewUserRepository(svc.Client)
	return svc.Ctx, svc.Client, activityRepo, userRepo
}

func TestActivityRepository(t *testing.T) {
	t.Parallel()
	t.Run("correctly find activity by id", func(t *testing.T) {
		t.Parallel()
		ctx, _, activityRepo, userRepo := setupTestActivityHandler(t)

		username := "alice"
		externalID := uuid.New()

		new_user := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created_user, err := userRepo.CreateUser(ctx, new_user)
		require.NoError(t, err)

		userID := created_user.ID
		duration := 782.5
		distance := 1500.0
		h3_indexes := []string{
			"8928308280fffff",
			"8928308280bffff",
			"8928308280dffff",
			"89283082807ffff",
			"89283082863ffff",
			"89283082867ffff",
			"8928308286bffff",
			"8928308286fffff",
		}
		created_at := time.Now()

		new_activity := &model.Activity{
			UserID:    userID,
			Duration:  duration,
			Distance:  distance,
			H3Indexes: h3_indexes,
			CreatedAt: created_at,
		}
		created, err := activityRepo.CreateActivity(ctx, new_activity)
		require.NoError(t, err)

		found, err := activityRepo.FindByID(ctx, created.ID)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
		require.Equal(t, userID, found.UserID)
		require.Equal(t, duration, found.DurationSeconds)
		require.Equal(t, distance, found.DistanceMeters)
		require.WithinDuration(t,
			created_at,
			found.CreatedAt, // actual
			time.Second,     // tolerance
		)

	})
}
