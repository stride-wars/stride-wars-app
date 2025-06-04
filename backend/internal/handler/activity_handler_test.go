package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/dto"
	"stride-wars-app/internal/handler"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/testutil"
	"stride-wars-app/internal/util"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type ActivityAPIResponse struct {
	Success bool                       `json:"success"`
	Data    dto.CreateActivityResponse `json:"data"`
	Error   string                     `json:"error,omitempty"`
}

type UserActivityStatsAPIResponse struct {
	Success bool                             `json:"success"`
	Data    dto.GetUserActivityStatsResponse `json:"data"`
	Error   string                           `json:"error,omitempty"`
}

var validH3Indexes = []string{
	"891e2e6b153ffff",
	"891e2e6b103ffff",
}

func setupTestActivityHandler(t *testing.T) (context.Context, *ent.Client, *handler.ActivityHandler) {
	t.Helper()

	svc := testutil.NewTestServices(t)
	activityHandler := handler.NewActivityHandler(svc.ActivityService, zap.NewExample())
	return svc.Ctx, svc.Client, activityHandler
}

func TestCreateActivity(t *testing.T) {
	t.Parallel()

	// ------------------------
	// Subtest: HappyPath
	// ------------------------
	t.Run("HappyPath", func(t *testing.T) {
		t.Parallel()

		// Arrange: setup handler + DB
		ctx, client, activityHandler := setupTestActivityHandler(t)

		// Seed a valid user
		userRepo := repository.NewUserRepository(client)
		username := "alice"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		// Sanity check: user exists
		foundUser, findErr := userRepo.FindByUsername(ctx, username)
		require.NoError(t, findErr)
		require.Equal(t, createdUser.ID, foundUser.ID)

		// Build request body
		createReq := dto.CreateActivityRequest{
			UserID:    foundUser.ID,
			Duration:  3600,  // 1 hour
			Distance:  10000, // 10 km
			H3Indexes: validH3Indexes,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/activity/create", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		http.HandlerFunc(activityHandler.CreateActivity).ServeHTTP(w, req)

		// Assert status code
		assert.Equal(t, http.StatusCreated, w.Code)

		// Assert response body
		var resp ActivityAPIResponse
		err = util.DecodeJSONBody(bytes.NewBuffer(w.Body.Bytes()), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.NotEqual(t, uuid.Nil, resp.Data.ID)
	})

	// ------------------------
	// Subtest: MissingFields
	// ------------------------
	t.Run("MissingFields", func(t *testing.T) {
		t.Parallel()

		// Arrange: setup handler + DB
		_, _, activityHandler := setupTestActivityHandler(t)

		// Build request with missing Duration
		createReq := dto.CreateActivityRequest{
			UserID:    uuid.New(),
			Distance:  10000,
			H3Indexes: validH3Indexes,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/activity/create", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		http.HandlerFunc(activityHandler.CreateActivity).ServeHTTP(w, req)

		// Assert status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert error message in response
		var resp ActivityAPIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Contains(t, resp.Error, "duration must be positive")
	})

	// ------------------------
	// Subtest: NoUser
	// ------------------------
	t.Run("NoUser", func(t *testing.T) {
		t.Parallel()

		// Arrange: setup handler + DB
		_, _, activityHandler := setupTestActivityHandler(t)

		// Build request with a random (non-existent) user ID
		createReq := dto.CreateActivityRequest{
			UserID:    uuid.New(),
			Duration:  3600,
			Distance:  10000,
			H3Indexes: validH3Indexes,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/activity/create", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		http.HandlerFunc(activityHandler.CreateActivity).ServeHTTP(w, req)

		// Assert status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert error message in response
		var resp ActivityAPIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Contains(t, resp.Error, "user not found")
	})

	// ------------------------
	// Subtest: InvalidH3Indexes
	// ------------------------
	t.Run("InvalidH3Indexes", func(t *testing.T) {
		t.Parallel()

		// Arrange: setup handler + DB
		_, _, activityHandler := setupTestActivityHandler(t)

		// Build request with invalid H3 index
		createReq := dto.CreateActivityRequest{
			UserID:    uuid.New(), // assume no user check here
			Duration:  3600,
			Distance:  10000,
			H3Indexes: []string{"mlody napoleon", "tylko troche wieksze berlo"}, // invalid
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/activity/create", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		http.HandlerFunc(activityHandler.CreateActivity).ServeHTTP(w, req)

		// Assert status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Assert error message in response
		var resp ActivityAPIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Contains(t, resp.Error, "invalid H3 index")
	})
}

func TestGetUserActivityStats(t *testing.T) {
	t.Parallel()

	// ------------------------
	// Subtest: HappyPath (two activities created via the handler)
	// ------------------------
	t.Run("HappyPath", func(t *testing.T) {
		t.Parallel()

		ctx, client, activityHandler := setupTestActivityHandler(t)

		// Seed a valid user
		userRepo := repository.NewUserRepository(client)
		username := "alice"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		// Verify user exists
		foundUser, findErr := userRepo.FindByUsername(ctx, username)
		require.NoError(t, findErr)
		require.Equal(t, createdUser.ID, foundUser.ID)

		// Create two activities via the handler (so CreatedAt == now)
		for i := 0; i < 2; i++ {
			createReq := dto.CreateActivityRequest{
				UserID:    foundUser.ID,
				Duration:  1800, // 30 minutes
				Distance:  5000, // 5 km
				H3Indexes: validH3Indexes,
			}
			reqBody, err := json.Marshal(createReq)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/activity/create", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			activityHandler.CreateActivity(w, req)
			assert.Equal(t, http.StatusCreated, w.Code)

			var createResp ActivityAPIResponse
			err = util.DecodeJSONBody(bytes.NewBuffer(w.Body.Bytes()), &createResp)
			require.NoError(t, err)
			assert.True(t, createResp.Success)
			assert.NotEqual(t, uuid.Nil, createResp.Data.ID)
		}

		// Now call GetUserActivityStats
		statsReq := httptest.NewRequest("GET", "/activity?user_id="+createdUser.ID.String(), nil)
		statsW := httptest.NewRecorder()
		activityHandler.GetUserActivityStats(statsW, statsReq)

		assert.Equal(t, http.StatusOK, statsW.Code)

		var response UserActivityStatsAPIResponse
		err = json.Unmarshal(statsW.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)

		// Expect exactly 2 activities recorded
		assert.Equal(t, int64(2), response.Data.ActivitiesRecorded)

		// Each activity used 2 H3 indexes → HexesVisited = 2 * 2 = 4
		assert.Equal(t, int64(4), response.Data.HexesVisited)

		// Distance: 5000 + 5000 = 10000
		assert.Equal(t, 10000.0, response.Data.DistanceCovered)

		// WeeklyActivities: both are “today”
		week := response.Data.WeeklyActivities
		require.Len(t, week, 7)

		// “Today” index should be 2 (the other six indices must be 0).
		// Our implementation maps “today” → index 6, so:
		for i := 0; i < 6; i++ {
			assert.Equal(t, int64(0), week[i], "expected no activities on day index %d", i)
		}
		assert.Equal(t, int64(2), week[6], "expected 2 activities on 'today' bucket")
	})

	// ------------------------
	// Subtest: MissingUserIDParam
	// ------------------------
	t.Run("MissingUserIDParam", func(t *testing.T) {
		t.Parallel()

		_, _, activityHandler := setupTestActivityHandler(t)

		// Call without any ?user_id=… at all
		req := httptest.NewRequest("GET", "/activity", nil)
		w := httptest.NewRecorder()
		activityHandler.GetUserActivityStats(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp UserActivityStatsAPIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.False(t, resp.Success)
		assert.Contains(t, resp.Error, "User ID is required")
	})

	// ------------------------
	// Subtest: InvalidUserIDParam
	// ------------------------
	t.Run("InvalidUserIDParam", func(t *testing.T) {
		t.Parallel()

		_, _, activityHandler := setupTestActivityHandler(t)

		// user_id is not a valid UUID
		req := httptest.NewRequest("GET", "/activity?user_id=not-a-uuid", nil)
		w := httptest.NewRecorder()
		activityHandler.GetUserActivityStats(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp UserActivityStatsAPIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.False(t, resp.Success)
		assert.Contains(t, resp.Error, "Invalid UUID format for 'user_id'")
	})

	// ------------------------
	// Subtest: NoActivitiesExistingUser
	// ------------------------
	t.Run("NoActivitiesExistingUser", func(t *testing.T) {
		t.Parallel()

		ctx, client, activityHandler := setupTestActivityHandler(t)

		// Create a user, but do NOT create any activities
		userRepo := repository.NewUserRepository(client)
		username := "bob"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		// Sanity check: user exists
		foundUser, findErr := userRepo.FindByUsername(ctx, username)
		require.NoError(t, findErr)
		require.Equal(t, createdUser.ID, foundUser.ID)

		// Call stats
		req := httptest.NewRequest("GET", "/activity?user_id="+createdUser.ID.String(), nil)
		w := httptest.NewRecorder()
		activityHandler.GetUserActivityStats(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp UserActivityStatsAPIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)

		// All counts should be zero
		assert.Equal(t, int64(0), resp.Data.ActivitiesRecorded)
		assert.Equal(t, int64(0), resp.Data.HexesVisited)
		assert.Equal(t, 0.0, resp.Data.DistanceCovered)
		require.Len(t, resp.Data.WeeklyActivities, 7)
		for i := 0; i < 7; i++ {
			assert.Equal(t, int64(0), resp.Data.WeeklyActivities[i], "expected zero at index %d", i)
		}
	})

}
