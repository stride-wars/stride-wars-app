package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/dto"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service"
	"stride-wars-app/internal/testutil"

	_ "github.com/mattn/go-sqlite3"
)

var validH3Indexes = []int64{
	617420388352917503,
	618094571271487487,
}

// setupTest initializes an in-memory SQLite ent client, applies schema,
// and returns an ActivityService wired with concrete repositories.
func setupTest(t *testing.T) (context.Context, *ent.Client, *service.ActivityService) {
	t.Helper()
	svc := testutil.NewTestServices(t)
	client := svc.Client
	ctx := svc.Ctx
	activityService := service.NewActivityService(
		svc.ActivityRepo,
		svc.HexInfluenceRepo,
		svc.HexLeaderboardRepo,
		svc.HexRepo,
		svc.UserService,
		zap.NewExample(),
	)
	return ctx, client, activityService
}

func TestActivityService(t *testing.T) {
	t.Parallel()

	// ------------------------
	// Subtest: FindByID
	// ------------------------
	t.Run("FindByID", func(t *testing.T) {
		t.Parallel()
		ctx, client, svc := setupTest(t)

		username := "karolnawrocki"
		externalID := uuid.New()
		userRepo := repository.NewUserRepository(client)

		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		act, err := client.Activity.Create().
			SetUserID(createdUser.ID).
			SetDurationSeconds(120.5).
			SetDistanceMeters(3000.0).
			SetH3Indexes(validH3Indexes).
			Save(ctx)
		require.NoError(t, err)

		fetched, err := svc.FindByID(ctx, act.ID)
		require.NoError(t, err)
		require.Equal(t, act.ID, fetched.ID)
		require.Equal(t, createdUser.ID, fetched.UserID)
	})

	// ------------------------
	// Subtest: FindByUserID
	// ------------------------
	t.Run("FindByUserID", func(t *testing.T) {
		t.Parallel()
		ctx, client, svc := setupTest(t)

		username := "karolnawrocki"
		externalID := uuid.New()
		userRepo := repository.NewUserRepository(client)

		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		for i := 0; i < 2; i++ {
			_, err := client.Activity.Create().
				SetUserID(createdUser.ID).
				SetDurationSeconds(60).
				SetDistanceMeters(1000).
				SetH3Indexes(validH3Indexes).
				Save(ctx)
			require.NoError(t, err)
		}

		results, err := svc.FindByUserID(ctx, createdUser.ID)
		require.NoError(t, err)
		require.Len(t, results, 2)
	})

	// ------------------------
	// Subtest: CreateOneActivity_WithH3Indexes
	// ------------------------
	t.Run("CreateOneActivity_WithH3Indexes", func(t *testing.T) {
		t.Parallel()
		ctx, client, svc := setupTest(t)

		username := "grzegorzbraun"
		externalID := uuid.New()
		userRepo := repository.NewUserRepository(client)

		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		input := dto.CreateActivityRequest{
			UserID:    createdUser.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: validH3Indexes,
		}
		created, err := svc.CreateActivity(ctx, input)
		require.NoError(t, err)

		stored, err := client.Activity.Get(ctx, created.ID)
		require.NoError(t, err)
		require.Equal(t, validH3Indexes, stored.H3Indexes)

		hexRepo := repository.NewHexRepository(client)
		for _, idx := range validH3Indexes {
			hex, err := hexRepo.FindByID(ctx, idx)
			require.NoError(t, err)
			require.Equal(t, idx, hex.ID)
		}

		hexInfluenceRepo := repository.NewHexInfluenceRepository(client)
		all, err := client.HexInfluence.Query().All(ctx)
		require.NoError(t, err)
		for _, hi := range all {
			t.Logf("HexInfluence: ID=%d, UserID=%d, Score=%f", hi.ID, hi.UserID, hi.Score)
		}

		for _, idx := range validH3Indexes {
			hxInfluence, err := hexInfluenceRepo.FindByUserIDAndHexID(ctx, created.UserID, idx)
			require.NoError(t, err)
			require.Equal(t, idx, hxInfluence.H3Index)
			require.Equal(t, 1.0, hxInfluence.Score)
		}
	})

	// ------------------------
	// Subtest: CreateTwoActivities_WithTheSameH3Indexes
	// ------------------------
	t.Run("CreateTwoActivities_WithTheSameH3Indexes", func(t *testing.T) {
		t.Parallel()
		ctx, client, svc := setupTest(t)

		username := "grzegorzbraun"
		externalID := uuid.New()
		userRepo := repository.NewUserRepository(client)

		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		input := dto.CreateActivityRequest{
			UserID:    createdUser.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: validH3Indexes,
		}
		firstActivity, err := svc.CreateActivity(ctx, input)
		require.NoError(t, err)

		_, err = svc.CreateActivity(ctx, input)
		require.NoError(t, err)

		hexRepo := repository.NewHexRepository(client)
		for _, idx := range validH3Indexes {
			hex, err := hexRepo.FindByID(ctx, idx)
			require.NoError(t, err)
			require.Equal(t, idx, hex.ID)
		}

		hexInfluenceRepo := repository.NewHexInfluenceRepository(client)
		for _, idx := range validH3Indexes {
			hxInfluence, err := hexInfluenceRepo.FindByUserIDAndHexID(ctx, firstActivity.UserID, idx)
			require.NoError(t, err)
			require.Equal(t, idx, hxInfluence.H3Index)
			require.Equal(t, 2.0, hxInfluence.Score)
		}
	})

	// ------------------------
	// Subtest: IfActivitiesAffectLeaderboardsCorrectly
	// ------------------------
	t.Run("IfActivitiesAffectLeaderboardsCorrectly", func(t *testing.T) {
		t.Parallel()
		ctx, client, svc := setupTest(t)

		userRepo := repository.NewUserRepository(client)
		hexRepo := repository.NewHexRepository(client)
		hexInfluenceRepo := repository.NewHexInfluenceRepository(client)
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

			activityReq := dto.CreateActivityRequest{
				UserID:    created.ID,
				Duration:  150.0,
				Distance:  2500.0,
				H3Indexes: validH3Indexes,
			}
			_, err = svc.CreateActivity(ctx, activityReq)
			require.NoError(t, err)
		}

		for _, user := range createdUsers {
			activityReq := dto.CreateActivityRequest{
				UserID:    user.ID,
				Duration:  150.0,
				Distance:  2500.0,
				H3Indexes: validH3Indexes,
			}
			_, err := svc.CreateActivity(ctx, activityReq)
			require.NoError(t, err)
		}

		for _, idx := range validH3Indexes {
			hex, err := hexRepo.FindByID(ctx, idx)
			require.NoError(t, err)
			require.Equal(t, idx, hex.ID)
		}

		for _, user := range createdUsers {
			for _, idx := range validH3Indexes {
				infl, err := hexInfluenceRepo.FindByUserIDAndHexID(ctx, user.ID, idx)
				require.NoError(t, err)
				require.Equal(t, idx, infl.H3Index)
				require.Equal(t, 2.0, infl.Score)

				leaderboard, err := hexLeaderboardRepo.FindByH3Index(ctx, idx)
				require.NoError(t, err)
				require.Equal(t, 5, len(leaderboard.TopUsers))
			}
		}
	})

	// ------------------------
	// Subtest: IfActivitiesAffectLeaderboardsIncrementally
	// ------------------------
	t.Run("IfActivitiesAffectLeaderboardsIncrementally", func(t *testing.T) {
		t.Parallel()
		ctx, client, svc := setupTest(t)

		userRepo := repository.NewUserRepository(client)
		hexLeaderboardRepo := repository.NewHexLeaderboardRepository(client)

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

			activityReq := dto.CreateActivityRequest{
				UserID:    created.ID,
				Duration:  150.0,
				Distance:  2500.0,
				H3Indexes: validH3Indexes,
			}
			_, err = svc.CreateActivity(ctx, activityReq)
			require.NoError(t, err)

			leaderboard, err := hexLeaderboardRepo.FindByH3Index(ctx, validH3Indexes[0])
			require.NoError(t, err)
			require.Equal(t, i+1, len(leaderboard.TopUsers))
		}
	})

	// ------------------------
	// Subtest: IfLeaderboardChangesCorrectly
	// ------------------------
	t.Run("IfLeaderboardChangesCorrectly", func(t *testing.T) {
		t.Parallel()
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

			activityReq := dto.CreateActivityRequest{
				UserID:    created.ID,
				Duration:  150.0,
				Distance:  2500.0,
				H3Indexes: validH3Indexes,
			}
			if name != "andrzejleper" {
				_, err = svc.CreateActivity(ctx, activityReq)
				require.NoError(t, err)
			}
		}

		leper := createdUsers[5]
		for i := range createdUsers[:5] {
			activityReq := dto.CreateActivityRequest{
				UserID:    leper.ID,
				Duration:  150.0,
				Distance:  2500.0,
				H3Indexes: validH3Indexes,
			}
			_, err := svc.CreateActivity(ctx, activityReq)
			require.NoError(t, err)

			positionPtr, err := hexLeaderboardRepo.GetUserPositionInLeaderboard(ctx, validH3Indexes[0], leper.ID)
			require.NoError(t, err)

			if i == 0 {
				require.Nil(t, positionPtr, "user should not be on the leaderboard for i == 0")
			} else {
				require.NotNil(t, positionPtr, "user should be on the leaderboard for i != 0")
				require.Equal(t, 1, *positionPtr, "expected position 1")
			}
		}
	})
}
