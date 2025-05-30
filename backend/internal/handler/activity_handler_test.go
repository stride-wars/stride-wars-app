package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"stride-wars-app/ent"
	"stride-wars-app/ent/enttest"
	"testing"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"stride-wars-app/ent/model"
	"stride-wars-app/internal/handler"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service"
)

type ActivityAPIResponse struct {
	Success bool                           `json:"success"`
	Data    service.CreateActivityResponse `json:"data"`
	Error   string                         `json:"error,omitempty"`
}

var validH3Indexes = []int64{
	618094571073044479,
	618094571271487487,
}

func setupTestActivityHandler(t *testing.T) (context.Context, *ent.Client, *handler.ActivityHandler) {
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
	activityService := service.NewActivityService(
		activityRepo,
		hexRepo,
		hexInfluenceRepo,
		hexLeaderboardRepo,
		userRepo,
		logger,
	)
	// Create handler
	activityHandler := handler.NewActivityHandler(
		activityService,
		logger,
	)
	return context.Background(), client, activityHandler
}

func TestCreateActivityHappyPath(t *testing.T) {
	// Setup
	ctx, client, activityHandler := setupTestActivityHandler(t)

	// Create test data
	repo := repository.NewUserRepository(client)
	username := "alice"
	externalID := uuid.New()
	new_user := &model.User{
		Username:     username,
		ExternalUser: externalID,
	}

	created_user, err := repo.CreateUser(ctx, new_user)
	require.NoError(t, err)

	// DEBUG: Print created user details
	t.Logf("Created user: ID=%s, Username=%s", created_user.ID, created_user.Username)

	// DEBUG: Verify user exists in database
	found_user, find_err := repo.FindByUsername(ctx, username)
	require.NoError(t, find_err)
	t.Logf("Found user in repo: ID=%s, Username=%s", found_user.ID, found_user.Username)

	// Prepare request body
	create_req := service.CreateActivityRequest{
		UserID:    created_user.ID,
		Duration:  3600,  // 1 hour in seconds
		Distance:  10000, // 10 km in meters
		H3Indexes: validH3Indexes,
	}
	req_body, err := json.Marshal(create_req)
	require.NoError(t, err)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest("POST", "/activity/create", bytes.NewBuffer(req_body))
	req.Header.Set("Content-Type", "application/json")
	// DEBUG: Print request body
	t.Logf("Request body: %s", req_body)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(activityHandler.CreateActivity)
	handler.ServeHTTP(w, req)

	// DEBUG: Print response body
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response status: %d", w.Code)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response ActivityAPIResponse
	unmarshal_err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, unmarshal_err)

	// DEBUG: Print parsed response
	t.Logf("Parsed response: ID=%s", response.Data.ID)

	assert.NotEqual(t, uuid.Nil, response.Data.ID)
}

/*
func TestCreateActivityHappyPath(t *testing.T) {
	// Setup
	ctx, client, activityHandler := setupTestActivityHandler(t)

	// Create test data
	repo := repository.NewUserRepository(client)
	username := "alice"
	externalID := uuid.New()
	new_user := &model.User{
		Username:     username,
		ExternalUser: externalID,
	}

	created_user, err := repo.CreateUser(ctx, new_user)
	require.NoError(t, err)

	// DEBUG: Print created user details
	t.Logf("Created user: ID=%s, Username=%s", created_user.ID, created_user.Username)

	// DEBUG: Verify user exists in database
	found_user, find_err := repo.FindByUsername(ctx, username)
	require.NoError(t, find_err)
	t.Logf("Found user in repo: ID=%s, Username=%s", found_user.ID, found_user.Username)

	// Prepare request body
	create_req := service.CreateActivityRequest{
		UserID:    created_user.ID,
		Duration:  3600, // 1 hour in seconds
		Distance:  10000, // 10 km in meters
		H3Indexes: validH3Indexes,
	}
	req_body, err := json.Marshal(create_req)
	require.NoError(t, err)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest("POST", "/activity/create", bytes.NewBuffer(req_body))
	req.Header.Set("Content-Type", "application/json")
	// DEBUG: Print request body
	t.Logf("Request body: %s", req_body)
	w := httptest.NewRecorder()

	handler := middleware.ParseJSON(http.HandlerFunc(activityHandler.CreateActivity))
    handler.ServeHTTP(w, req)

	// DEBUG: Print response body
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response status: %d", w.Code)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response ActivityAPIResponse
	unmarshal_err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, unmarshal_err)

	// DEBUG: Print parsed response
	t.Logf("Parsed response: ID=%s", response.Data.ID)

	assert.NotEqual(t, uuid.Nil, response.Data.ID)
}
*/

func TestCreateActivityMissingFields(t *testing.T) {
	// Setup
	_, _, activityHandler := setupTestActivityHandler(t)

	// Prepare request with missing fields
	create_req := service.CreateActivityRequest{
		UserID: uuid.New(),
		// Missing  Duration
		Distance:  10000, // 10 km in meters
		H3Indexes: validH3Indexes,
	}
	req_body, err := json.Marshal(create_req)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/activity/create", bytes.NewBuffer(req_body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(activityHandler.CreateActivity)
	handler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response ActivityAPIResponse
	unmarshal_err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, unmarshal_err)
	assert.Contains(t, response.Error, "duration must be positive")
}

func TestCreateActivityNoUser(t *testing.T) {
	// Setup
	_, _, activityHandler := setupTestActivityHandler(t)

	// Prepare request with non-existent user ID
	create_req := service.CreateActivityRequest{
		UserID:    uuid.New(), // Non-existent user ID
		Duration:  3600,       // 1 hour in seconds
		Distance:  10000,      // 10 km in meters
		H3Indexes: validH3Indexes,
	}
	req_body, err := json.Marshal(create_req)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/activity/create", bytes.NewBuffer(req_body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(activityHandler.CreateActivity)
	handler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response ActivityAPIResponse
	unmarshal_err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, unmarshal_err)
	assert.Contains(t, response.Error, "user not found")
}

func TestCreateActivityInvalidH3Indexes(t *testing.T) {
	// Setup
	_, _, activityHandler := setupTestActivityHandler(t)

	// Prepare request with invalid H3 indexes
	create_req := service.CreateActivityRequest{
		UserID:    uuid.New(),   // Assuming this user exists
		Duration:  3600,         // 1 hour in seconds
		Distance:  10000,        // 10 km in meters
		H3Indexes: []int64{122}, // Empty H3 indexes
	}
	req_body, err := json.Marshal(create_req)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/activity/create", bytes.NewBuffer(req_body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(activityHandler.CreateActivity)
	handler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response ActivityAPIResponse
	unmarshal_err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, unmarshal_err)
	assert.Contains(t, response.Error, "invalid H3 index")
}
