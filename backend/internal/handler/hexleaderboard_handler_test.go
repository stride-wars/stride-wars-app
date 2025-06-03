package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/dto"
	"stride-wars-app/internal/handler"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service"
	"stride-wars-app/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var krakowH3Indexes = []string{
	"891e2e6b153ffff",
	"891e2e6b103ffff",
}

type HexLeaderboardAPIResponse struct {
	Success bool                                        `json:"success"`
	Data    dto.GetAllHexLeaderboardsInsideBBoxResponse `json:"data"`
}

func setupTestHexLeaderboardHandler(t *testing.T) (context.Context, *ent.Client, *handler.HexLeaderboardHandler) {
	t.Helper()
	svc := testutil.NewTestServices(t)

	hexLeaderboardHandler := handler.NewHexLeaderboardHandler(
		svc.HexLeaderboardService,
		zap.NewExample(),
	)

	return svc.Ctx, svc.Client, hexLeaderboardHandler
}

func TestHexLeaderboardHandler(t *testing.T) {
	t.Parallel()

	t.Run("GetLeaderboardsInsideBBox", func(t *testing.T) {
		t.Parallel()
		ctx, client, hexLeaderboardHandler := setupTestHexLeaderboardHandler(t)

		// Seed two users
		userRepo := repository.NewUserRepository(client)
		username1 := "alice"
		externalID1 := uuid.New()
		user1 := &model.User{
			Username:     username1,
			ExternalUser: externalID1,
		}
		createdUser1, err := userRepo.CreateUser(ctx, user1)
		require.NoError(t, err)

		username2 := "bob"
		externalID2 := uuid.New()
		user2 := &model.User{
			Username:     username2,
			ExternalUser: externalID2,
		}
		createdUser2, err := userRepo.CreateUser(ctx, user2)
		require.NoError(t, err)

		// Create two activities with the same H3 indexes (different users)
		activityService := service.NewActivityService(
			repository.NewActivityRepository(client),
			repository.NewHexInfluenceRepository(client),
			repository.NewHexLeaderboardRepository(client),
			repository.NewHexRepository(client),
			service.NewUserService(repository.NewUserRepository(client), zap.NewExample()),
			zap.NewExample(),
		)

		activity1 := dto.CreateActivityRequest{
			UserID:    createdUser1.ID,
			Duration:  3600,
			Distance:  5000,
			H3Indexes: krakowH3Indexes,
		}
		_, err = activityService.CreateActivity(ctx, activity1)
		require.NoError(t, err)

		activity2 := dto.CreateActivityRequest{
			UserID:    createdUser2.ID,
			Duration:  7200,
			Distance:  10000,
			H3Indexes: krakowH3Indexes,
		}
		_, err = activityService.CreateActivity(ctx, activity2)
		require.NoError(t, err)

		// Build bounding box query
		bbox := service.BoundingBox{
			MinLat: 49.9650,
			MinLng: 19.7500,
			MaxLat: 50.1500,
			MaxLng: 20.1000,
		}
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("/hexleaderboards/bbox?min_lat=%f&min_lng=%f&max_lat=%f&max_lng=%f",
				bbox.MinLat, bbox.MinLng, bbox.MaxLat, bbox.MaxLng),
			nil,
		)
		require.NoError(t, err)

		w := httptest.NewRecorder()

		// Act: call the handler
		hexLeaderboardHandler.GetAllLeaderboardsInsideBBox(w, req)

		// Assert status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Assert response body
		var resp HexLeaderboardAPIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Len(t, resp.Data.Leaderboards, 2)
		// assert that leaderboards contain our users
		assert.Equal(t, createdUser1.ID, resp.Data.Leaderboards[0].TopUsers[0].UserID)
		assert.Equal(t, createdUser2.ID, resp.Data.Leaderboards[0].TopUsers[1].UserID)
		assert.Equal(t, createdUser1.ID, resp.Data.Leaderboards[1].TopUsers[0].UserID)
		assert.Equal(t, createdUser2.ID, resp.Data.Leaderboards[1].TopUsers[1].UserID)
	})
}
