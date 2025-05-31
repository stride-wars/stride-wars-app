package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"stride-wars-app/ent/enttest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"

	"fmt"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/dto"
	"stride-wars-app/internal/handler"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service"
)

var krakowH3Indexes = []int64{
	617524104371896319,
	617524104366653439,
}

type HexLeaderboardAPIResponse struct {
	Success bool                                        `json:"success"`
	Data    dto.GetAllHexLeaderboardsInsideBBoxResponse `json:"data"`
}

func setupTestHexLeaderboardHandler(t *testing.T) (context.Context, *ent.Client, *handler.HexLeaderboardHandler) {
	t.Helper()
	// Use enttest to create a transient in-memory SQLite DB
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")

	// Run auto migration to create schema
	err := client.Schema.Create(context.Background())
	require.NoError(t, err)

	// Wire repositories
	hexLeaderboardRepo := repository.NewHexLeaderboardRepository(client)
	hexInfluenceRepo := repository.NewHexInfluenceRepository(client)
	logger := zap.NewExample()

	// Create service
	hexLeaderboardService := service.NewHexLeaderboardService(hexLeaderboardRepo, hexInfluenceRepo, logger)

	// Create handler
	handler := handler.NewHexLeaderboardHandler(hexLeaderboardService, logger)

	return context.Background(), client, handler
}

// this test creates two activites, both for the same two hexes with different users
func TestGetHexLeaderboardsInBBox(t *testing.T) {
	ctx, client, handler := setupTestHexLeaderboardHandler(t)
	// create two users
	repo := repository.NewUserRepository(client)
	username := "alice"
	externalID := uuid.New()
	newUser := &model.User{
		Username:     username,
		ExternalUser: externalID,
	}
	createdUser, userErr := repo.CreateUser(ctx, newUser)
	require.NoError(t, userErr)

	// Create a second user
	username2 := "bob"
	externalID2 := uuid.New()
	newUser2 := &model.User{
		Username:     username2,
		ExternalUser: externalID2,
	}
	createdUser2, userErr := repo.CreateUser(ctx, newUser2)
	require.NoError(t, userErr)

	// Create activities for both users
	activity1 := service.CreateActivityRequest{
		UserID:    createdUser.ID,
		Duration:  3600, // 1 hour
		Distance:  5000, // 5 km
		H3Indexes: krakowH3Indexes,
	}
	activity2 := service.CreateActivityRequest{
		UserID:    createdUser2.ID,
		Duration:  7200,  // 2 hours
		Distance:  10000, // 10 km
		H3Indexes: krakowH3Indexes,
	}

	logger := zap.NewExample()

	// Create the first activity
	activityService := service.NewActivityService(
		repository.NewActivityRepository(client),
		repository.NewHexRepository(client),
		repository.NewHexInfluenceRepository(client),
		repository.NewHexLeaderboardRepository(client),
		repository.NewUserRepository(client),
		*service.NewUserService(repository.NewUserRepository(client), logger),
		logger,
	)
	_, err := activityService.CreateActivity(ctx, activity1)
	require.NoError(t, err)
	// Create the second activity
	_, err = activityService.CreateActivity(ctx, activity2)
	require.NoError(t, err)

	bbox := service.BoundingBox{
		MinLat: 49.9650, // southern boundary
		MinLng: 19.7500, // western boundary
		MaxLat: 50.1500, // northern boundary
		MaxLng: 20.1000,
	}
	// Prepare the request with the arguments from bounding box
	req, err := http.NewRequest("GET",
		fmt.Sprintf("/hexleaderboards?min_lat=%f&min_lng=%f&max_lat=%f&max_lng=%f",
			bbox.MinLat, bbox.MinLng, bbox.MaxLat, bbox.MaxLng),
		nil)
	require.NoError(t, err)
	// Create a ResponseRecorder to capture the response
	recorder := httptest.NewRecorder()
	// Call the handler
	handler.GetAllLeaderboardsInsideBBox(recorder, req)
	t.Logf("Response body: %s", recorder.Body.String())
	// Check the response status code
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
	// Check the response body
	var response HexLeaderboardAPIResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err, "Expected valid JSON response")
	assert.NotEmpty(t, response.Data.Leaderboards, "Expected non-empty leaderboard response")
	// Check that both users are in the response
	assert.Len(t, response.Data.Leaderboards, 2, "Expected two users in the leaderboard response")

}
