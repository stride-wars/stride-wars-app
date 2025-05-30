package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"stride-wars-app/ent/enttest"

	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/mattn/go-sqlite3"

	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service"

	"fmt"

	"github.com/google/uuid"
)

var validH3Indexes = []int64{
	617420388352917503,
	618094571271487487,
}

// setupTest initializes an in-memory SQLite ent client, applies schema,
// and returns an ActivityService wired with concrete repositories.
func setupTest(t *testing.T) (context.Context, *ent.Client, *service.ActivityService) {
	t.Helper()
	// Use enttest to create a transient in-memory SQLite DB
	dbName := fmt.Sprintf("file:ent_%s?mode=memory&cache=shared&_fk=1", uuid.New().String())
	client := enttest.Open(t, "sqlite3", dbName)

	// Run auto migration to create schema
	err := client.Schema.Create(context.Background(), schema.WithForeignKeys(true))
	require.NoError(t, err)

	// Wire repositories
	activityRepo := repository.NewActivityRepository(client)
	hexRepo := repository.NewHexRepository(client)
	hexInfluenceRepo := repository.NewHexInfluenceRepository(client)
	hexLeaderboardRepo := repository.NewHexLeaderboardRepository(client)
	userRepo := repository.NewUserRepository(client)
	logger := zap.NewExample()
	// Create service
	svc := service.NewActivityService(
		activityRepo,
		hexRepo,
		hexInfluenceRepo,
		hexLeaderboardRepo,
		userRepo,
		logger,
	)
	return context.Background(), client, svc
}

func TestFindByID(t *testing.T) {
	ctx, client, svc := setupTest(t)

	username := "karolnawrocki"
	externalID := uuid.New()
	user_repo := repository.NewUserRepository(client)

	new_user := &model.User{
		Username:     username,
		ExternalUser: externalID,
	}
	created_user, err := user_repo.CreateUser(ctx, new_user)
	require.NoError(t, err)

	act, err := client.Activity.Create().
		SetUserID(created_user.ID).
		SetDurationSeconds(120.5).
		SetDistanceMeters(3000.0).
		SetH3Indexes(validH3Indexes).
		Save(ctx)
	require.NoError(t, err)

	// Now use service to fetch
	fetched, err := svc.FindByID(ctx, act.ID)
	require.NoError(t, err)
	require.Equal(t, act.ID, fetched.ID)
	require.Equal(t, created_user.ID, fetched.UserID)
}

func TestFindByUserID(t *testing.T) {
	ctx, client, svc := setupTest(t)
	username := "karolnawrocki"
	externalID := uuid.New()
	user_repo := repository.NewUserRepository(client)

	new_user := &model.User{
		Username:     username,
		ExternalUser: externalID,
	}
	created_user, user_err := user_repo.CreateUser(ctx, new_user)
	require.NoError(t, user_err)
	// Create two activities for same user
	for i := 0; i < 2; i++ {
		_, err := client.Activity.Create().
			SetUserID(created_user.ID).
			SetDurationSeconds(60).
			SetDistanceMeters(1000).
			SetH3Indexes(validH3Indexes).
			Save(ctx)
		require.NoError(t, err)
	}

	results, err := svc.FindByUserID(ctx, created_user.ID)
	require.NoError(t, err)
	require.Len(t, results, 2)
}

func TestCreateOneActivity_WithH3Indexes(t *testing.T) {
	ctx, client, svc := setupTest(t)

	username := "grzegorzbraun"
	externalID := uuid.New()
	user_repo := repository.NewUserRepository(client)

	new_user := &model.User{
		Username:     username,
		ExternalUser: externalID,
	}
	created_user, user_err := user_repo.CreateUser(ctx, new_user)
	require.NoError(t, user_err)

	input := service.CreateActivityRequest{
		UserID:    created_user.ID,
		Duration:  150.0,
		Distance:  2500.0,
		H3Indexes: validH3Indexes,
	}
	created, err := svc.CreateActivity(ctx, input)
	require.NoError(t, err)
	// Verify activity saved
	stored, err := client.Activity.Get(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, validH3Indexes, stored.H3Indexes)

	hex_repo := repository.NewHexRepository((client))
	// Verify hexes created
	for _, idx := range validH3Indexes {
		hex, err := hex_repo.FindByID(ctx, idx)
		require.NoError(t, err)
		require.Equal(t, idx, hex.ID)
	}

	hexinfluence_repo := repository.NewHexInfluenceRepository(client)
	all, err := client.HexInfluence.Query().All(ctx)
	require.NoError(t, err)
	for _, hi := range all {
		t.Logf("HexInfluence: ID=%d, UserID=%d, Score=%f", hi.ID, hi.UserID, hi.Score)
	}
	// Verify influences created
	for _, idx := range validH3Indexes {
		hx_influence, err := hexinfluence_repo.FindByUserIDAndHexID(ctx, created.UserID, idx)
		require.NoError(t, err)
		require.Equal(t, idx, hx_influence.H3Index)
		require.Equal(t, 1.0, hx_influence.Score)
	}

}

func TestCreateTwoActivities_WithTheSameH3Indexes(t *testing.T) {
	ctx, client, svc := setupTest(t)

	username := "grzegorzbraun"
	externalID := uuid.New()
	user_repo := repository.NewUserRepository(client)

	new_user := &model.User{
		Username:     username,
		ExternalUser: externalID,
	}
	created_user, user_err := user_repo.CreateUser(ctx, new_user)
	require.NoError(t, user_err)

	input := service.CreateActivityRequest{
		UserID:    created_user.ID,
		Duration:  150.0,
		Distance:  2500.0,
		H3Indexes: validH3Indexes,
	}
	first_activity, err := svc.CreateActivity(ctx, input)
	require.NoError(t, err)

	_, a_err := svc.CreateActivity(ctx, input)
	require.NoError(t, a_err)

	hex_repo := repository.NewHexRepository((client))
	// Verify hexes created
	for _, idx := range validH3Indexes {
		hex, err := hex_repo.FindByID(ctx, idx)
		require.NoError(t, err)
		require.Equal(t, idx, hex.ID)
	}

	hexinfluence_repo := repository.NewHexInfluenceRepository(client)

	// Verify influences created
	for _, idx := range validH3Indexes {
		hx_influence, err := hexinfluence_repo.FindByUserIDAndHexID(ctx, first_activity.UserID, idx)
		require.NoError(t, err)
		require.Equal(t, idx, hx_influence.H3Index)
		require.Equal(t, 2.0, hx_influence.Score)
	}

}

func TestIfActivitiesAffectLeaderboardsCorrectly(t *testing.T) {
	ctx, client, svc := setupTest(t)

	userRepo := repository.NewUserRepository(client)
	hexRepo := repository.NewHexRepository(client)
	hexInfluenceRepo := repository.NewHexInfluenceRepository(client)
	hexLeadearboardRepo := repository.NewHexLeaderboardRepository(client)

	usernames := []string{
		"grzegorzbraun",
		"januszkorwinmikke",
		"krzysztofbosak",
		"jaroslawkaczynski",
		"robertbiedron",
		"andrzejleper",
	}

	var createdUsers []*ent.User
	for _, name := range usernames {
		user := &model.User{
			Username:     name,
			ExternalUser: uuid.New(),
		}
		created, err := userRepo.CreateUser(ctx, user)
		require.NoError(t, err)
		createdUsers = append(createdUsers, created)

		activity := service.CreateActivityRequest{
			UserID:    created.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: validH3Indexes,
		}
		_, err = svc.CreateActivity(ctx, activity)
		require.NoError(t, err)
	}

	for _, user := range createdUsers {

		activity := service.CreateActivityRequest{
			UserID:    user.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: validH3Indexes,
		}
		_, err := svc.CreateActivity(ctx, activity)
		require.NoError(t, err)
	}
	// Verify hexes were created
	for _, idx := range validH3Indexes {
		hex, err := hexRepo.FindByID(ctx, idx)
		require.NoError(t, err)
		require.Equal(t, idx, hex.ID)
	}

	// Verify each user's influence
	for _, user := range createdUsers {
		for _, idx := range validH3Indexes {
			infl, err := hexInfluenceRepo.FindByUserIDAndHexID(ctx, user.ID, idx)
			require.NoError(t, err)
			require.Equal(t, idx, infl.H3Index)
			require.Equal(t, 2.0, infl.Score)

			leaderboard, l_err := hexLeadearboardRepo.FindByH3Index(ctx, idx)
			require.NoError(t, l_err)
			require.Equal(t, 5, len(leaderboard.TopUsers))
		}
	}
}

func TestIfActivitiesAffectLeaderboardsIncremetally(t *testing.T) {
	ctx, client, svc := setupTest(t)

	userRepo := repository.NewUserRepository(client)
	hexLeadearboardRepo := repository.NewHexLeaderboardRepository(client)

	usernames := []string{
		"grzegorzbraun",
		"januszkorwinmikke",
		"krzysztofbosak",
		"jaroslawkaczynski",
		"robertbiedron",
	}

	for i, name := range usernames {
		user := &model.User{
			Username:     name,
			ExternalUser: uuid.New(),
		}
		created, err := userRepo.CreateUser(ctx, user)
		require.NoError(t, err)

		activity := service.CreateActivityRequest{
			UserID:    created.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: validH3Indexes,
		}
		_, err = svc.CreateActivity(ctx, activity)
		require.NoError(t, err)

		leaderboard, l_err := hexLeadearboardRepo.FindByH3Index(ctx, validH3Indexes[0])
		require.NoError(t, l_err)
		temp := len(leaderboard.TopUsers)
		require.Equal(t, i+1, temp)
	}

}

func TestIfLeaderboardChangesCorrectly(t *testing.T) {
	ctx, client, svc := setupTest(t)

	userRepo := repository.NewUserRepository(client)
	hexLeaderboardRepo := repository.NewHexLeaderboardRepository(client)

	usernames := []string{
		"grzegorzbraun",
		"januszkorwinmikke",
		"krzysztofbosak",
		"jaroslawkaczynski",
		"robertbiedron",
		"andrzejleper",
	}

	var createdUsers []*ent.User
	for _, name := range usernames {
		user := &model.User{
			Username:     name,
			ExternalUser: uuid.New(),
		}
		created, err := userRepo.CreateUser(ctx, user)
		require.NoError(t, err)
		createdUsers = append(createdUsers, created)

		activity := service.CreateActivityRequest{
			UserID:    created.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: validH3Indexes,
		}
		if name != "andrzejleper" {
			_, err = svc.CreateActivity(ctx, activity)
			require.NoError(t, err)
		}
	}
	leper := createdUsers[5]

	for i := range 5 {
		activity := service.CreateActivityRequest{
			UserID:    leper.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: validH3Indexes,
		}
		_, err := svc.CreateActivity(ctx, activity)
		require.NoError(t, err)

		positionPtr, err := hexLeaderboardRepo.GetUserPositionInLeaderboard(ctx, validH3Indexes[0], leper.ID)
		require.NoError(t, err)

		if i == 0 {
			// we expected “not found” → nil pointer
			require.Nil(t, positionPtr, "user should not be on the leaderboard for i == 0")
		} else {
			// we expected position 1 → non-nil, value == 1
			require.NotNil(t, positionPtr, "user should be on the leaderboard for i != 0")
			require.Equal(t, 1, *positionPtr, "expected position 1")
		}

	}

}
