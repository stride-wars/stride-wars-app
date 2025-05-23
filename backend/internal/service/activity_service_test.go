package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"entgo.io/ent/dialect/sql/schema"
	"stride-wars-app/ent/enttest"
	_ "github.com/mattn/go-sqlite3"

	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service" 

	"github.com/google/uuid"
	"fmt"
)

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
	logger := zap.NewExample()
	// Create service
	svc := service.NewActivityService(
		activityRepo,
		hexRepo,
		hexInfluenceRepo,
		hexLeaderboardRepo,
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
		SetH3Indexes([]int64{}).
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
			SetH3Indexes([]int64{}).
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

	indexes := []int64{12345, 67890}
	input := &model.Activity{
		UserID:    created_user.ID,
		Duration:  150.0,
		Distance:  2500.0,
		H3Indexes: indexes,
	}
	created, err := svc.CreateActivity(ctx, input)
	require.NoError(t, err)
	// Verify activity saved
	stored, err := client.Activity.Get(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, indexes, stored.H3Indexes)

	hex_repo := repository.NewHexRepository((client))
	// Verify hexes created
	for _, idx := range indexes {
		hex, err := hex_repo.FindByID(ctx,idx)
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
	for _, idx := range indexes {
		hx_influence, err := hexinfluence_repo.FindByUserIDAndHexID(ctx,created.UserID,idx)
		require.NoError(t, err)
		require.Equal(t, idx, hx_influence.H3Index)
		require.Equal(t,1.0,hx_influence.Score)
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

	indexes := []int64{12345, 67890}
	input := &model.Activity{
		UserID:    created_user.ID,
		Duration:  150.0,
		Distance:  2500.0,
		H3Indexes: indexes,
	}
	first_activity, err := svc.CreateActivity(ctx, input)
	require.NoError(t, err)

	svc.CreateActivity(ctx, input)
	require.NoError(t, err)



	hex_repo := repository.NewHexRepository((client))
	// Verify hexes created
	for _, idx := range indexes {
		hex, err := hex_repo.FindByID(ctx,idx)
		require.NoError(t, err)
		require.Equal(t, idx, hex.ID)
	}

	hexinfluence_repo := repository.NewHexInfluenceRepository(client)


	// Verify influences created
	for _, idx := range indexes {
		hx_influence, err := hexinfluence_repo.FindByUserIDAndHexID(ctx,first_activity.UserID,idx)
		require.NoError(t, err)
		require.Equal(t, idx, hx_influence.H3Index)
		require.Equal(t,2.0,hx_influence.Score)
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

	indexes := []int64{12345, 67890}

	var createdUsers []*ent.User
	for _, name := range usernames {
		user := &model.User{
			Username:     name,
			ExternalUser: uuid.New(),
		}
		created, err := userRepo.CreateUser(ctx, user)
		require.NoError(t, err)
		createdUsers = append(createdUsers, created)

		activity := &model.Activity{
			UserID:    created.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: indexes,
		}
		_, err = svc.CreateActivity(ctx, activity)
		require.NoError(t, err)
	}

	for _, user := range createdUsers {

			activity := &model.Activity{
			UserID:    user.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: indexes,
		}
		_, err := svc.CreateActivity(ctx, activity)
		require.NoError(t, err)
	}
	// Verify hexes were created
	for _, idx := range indexes {
		hex, err := hexRepo.FindByID(ctx, idx)
		require.NoError(t, err)
		require.Equal(t, idx, hex.ID)
	}

	// Verify each user's influence
	for _, user := range createdUsers {
		for _, idx := range indexes {
			infl, err := hexInfluenceRepo.FindByUserIDAndHexID(ctx, user.ID, idx)
			require.NoError(t, err)
			require.Equal(t, idx, infl.H3Index)
			require.Equal(t, 2.0, infl.Score)  
			
			
			leaderboard, l_err := hexLeadearboardRepo.FindByH3Index(ctx,idx)
			require.NoError(t,l_err)
			require.Equal(t,5,len(leaderboard.TopUsers))
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

	indexes := []int64{12345, 67890}

	for i, name := range usernames {
		user := &model.User{
			Username:     name,
			ExternalUser: uuid.New(),
		}
		created, err := userRepo.CreateUser(ctx, user)
		require.NoError(t, err)

		activity := &model.Activity{
			UserID:    created.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: indexes,
		}
		_, err = svc.CreateActivity(ctx, activity)
		require.NoError(t, err)

		leaderboard, l_err := hexLeadearboardRepo.FindByH3Index(ctx,indexes[0])
		require.NoError(t,l_err)
		temp := len(leaderboard.TopUsers)
		require.Equal(t,i+1,temp)
	}



}

func TestIfLeaderboardChangesCorrectly(t *testing.T) {
	ctx, client, svc := setupTest(t)

	userRepo := repository.NewUserRepository(client)
	hexLeadearboardRepo := repository.NewHexLeaderboardRepository(client)

	usernames := []string{
		"grzegorzbraun",
		"januszkorwinmikke",
		"krzysztofbosak",
		"jaroslawkaczynski",
		"robertbiedron",
		"andrzejleper",
	}

	indexes := []int64{12345, 67890}

	var createdUsers []*ent.User
	for _, name := range usernames {
		user := &model.User{
			Username:     name,
			ExternalUser: uuid.New(),
		}
		created, err := userRepo.CreateUser(ctx, user)
		require.NoError(t, err)
		createdUsers = append(createdUsers, created)

		activity := &model.Activity{
			UserID:    created.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: indexes,
		}
		if name != "andrzejleper" {
		_, err = svc.CreateActivity(ctx, activity)
		require.NoError(t, err)
		}
	}
	leper := createdUsers[5]

	for i := range 5 {
		activity := &model.Activity{
			UserID:    leper.ID,
			Duration:  150.0,
			Distance:  2500.0,
			H3Indexes: indexes,
		}
		_, err := svc.CreateActivity(ctx, activity)
		require.NoError(t, err)

		position, l_err := hexLeadearboardRepo.GetUserPositionInLeaderboard(ctx,indexes[0],leper.ID)
		require.NoError(t,l_err)
		expected_postion := 0
		if i == 0 {
			expected_postion = -1
		} else {
			expected_postion = 1
		}
		require.Equal(t,expected_postion,position)
	}

}