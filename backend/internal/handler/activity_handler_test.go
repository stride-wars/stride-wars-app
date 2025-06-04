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

var validH3Indexes = []int64{
	618094571073044479,
	618094571271487487,
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
			H3Indexes: []int64{122}, // invalid
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
