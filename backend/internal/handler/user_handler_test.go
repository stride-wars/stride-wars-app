package handler_test

import (
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

	"bytes"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/api/middleware"
	"stride-wars-app/internal/handler"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service"
)

type UserAPIResponse struct {
	Success bool     `json:"success"`
	Data    UserData `json:"data"`
}

type UserData struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	ExternalUser uuid.UUID `json:"external_user"`
	Edges        struct{}  `json:"edges"`
}

func setupTestUserHandler(t *testing.T) (context.Context, *ent.Client, *handler.UserHandler) {
	t.Helper()
	// Use enttest to create a transient in-memory SQLite DB
	dbName := fmt.Sprintf("file:ent_%s?mode=memory&cache=shared&_fk=1", uuid.New().String())
	client := enttest.Open(t, "sqlite3", dbName)

	// Run auto migration to create schema
	err := client.Schema.Create(context.Background(), schema.WithForeignKeys(true))
	require.NoError(t, err)

	// Wire repositories
	userRepo := repository.NewUserRepository(client)
	logger := zap.NewExample()
	// Create service
	userService := service.NewUserService(
		userRepo,
		logger,
	)
	// Create handler
	userHandler := handler.NewUserHandler(
		userService,
		logger,
	)
	return context.Background(), client, userHandler
}

func TestGetUserByUsernameHappyPath(t *testing.T) {
	// Setup
	ctx, client, userHandler := setupTestUserHandler(t)

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

	// Test request
	req := httptest.NewRequest("GET", "/user/by-username?username=alice", nil)
	w := httptest.NewRecorder()

	userHandler.GetUserByUsername(w, req)

	// DEBUG: Print response body
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response status: %d", w.Code)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response UserAPIResponse
	unmarshal_err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, unmarshal_err)

	// DEBUG: Print parsed response
	t.Logf("Parsed response: ID=%s, Username=%s", response.Data.ID, response.Data.Username)

	assert.Equal(t, "alice", response.Data.Username)
	assert.NotEqual(t, uuid.Nil, response.Data.ID)
}

func TestGetUserByUsernameNotFound(t *testing.T) {
	// Setup
	_, _, userHandler := setupTestUserHandler(t)

	// Test request
	req := httptest.NewRequest("GET", "/user/by-username?username=alice", nil)
	w := httptest.NewRecorder()

	userHandler.GetUserByUsername(w, req)

	// DEBUG: Print response body
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response status: %d", w.Code)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response UserAPIResponse
	unmarshal_err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, unmarshal_err)

	// DEBUG: Print parsed response
	t.Logf("Parsed response: ID=%s, Username=%s", response.Data.ID, response.Data.Username)

	assert.Equal(t, "", response.Data.Username)
	assert.Equal(t, uuid.Nil, response.Data.ID)
}

func TestGetUserByUsernameBadRequest(t *testing.T) {
	// Setup
	_, _, userHandler := setupTestUserHandler(t)

	// Test request
	req := httptest.NewRequest("GET", "/user/by-username?usernnname=alice", nil)
	w := httptest.NewRecorder()

	userHandler.GetUserByUsername(w, req)

	// DEBUG: Print response body
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response status: %d", w.Code)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserByIDHappyPath(t *testing.T) {
	// Setup
	ctx, client, userHandler := setupTestUserHandler(t)

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
	found_user, find_err := repo.FindByID(ctx, created_user.ID)
	require.NoError(t, find_err)
	t.Logf("Found user in repo: ID=%s, Username=%s", found_user.ID, found_user.Username)

	// Test request
	req := httptest.NewRequest("GET", fmt.Sprintf("/user/by-id?id=%s", created_user.ID), nil)
	w := httptest.NewRecorder()

	userHandler.GetUserByID(w, req)

	// DEBUG: Print response body
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response status: %d", w.Code)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response UserAPIResponse
	unmarshal_err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, unmarshal_err)

	// DEBUG: Print parsed response
	t.Logf("Parsed response: ID=%s, Username=%s", response.Data.ID, response.Data.Username)

	assert.Equal(t, created_user.ID, response.Data.ID)
	assert.Equal(t, "alice", response.Data.Username)
}

func TestGetUserByIDNotFound(t *testing.T) {
	// Setup
	_, _, userHandler := setupTestUserHandler(t)

	// Test request with non-existent ID
	nonExistentID := uuid.New()
	req := httptest.NewRequest("GET", fmt.Sprintf("/user/by-id?id=%s", nonExistentID), nil)
	w := httptest.NewRecorder()

	userHandler.GetUserByID(w, req)

	// DEBUG: Print response body
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response status: %d", w.Code)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response UserAPIResponse
	unmarshal_err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, unmarshal_err)

	// DEBUG: Print parsed response
	t.Logf("Parsed response: ID=%s, Username=%s", response.Data.ID, response.Data.Username)

	assert.Equal(t, "", response.Data.Username)
	assert.Equal(t, uuid.Nil, response.Data.ID)
}

func TestGetUserByIDBadRequest(t *testing.T) {
	// Setup
	_, _, userHandler := setupTestUserHandler(t)

	// Test request with invalid ID format
	req := httptest.NewRequest("GET", "/user/by-id?id=invalid-uuid", nil)
	w := httptest.NewRecorder()

	userHandler.GetUserByID(w, req)

	// DEBUG: Print response body
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response status: %d", w.Code)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUsernameHappyPath(t *testing.T) {
	// Setup
	ctx, client, userHandler := setupTestUserHandler(t)

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

	// Test request
	updateReq := service.UpdateUsernameRequest{
		OldUsername: "alice",
		NewUsername: "bob",
	}
	reqBody, _ := json.Marshal(updateReq)

	// DEBUG: Print the request body
	t.Logf("Request body: %s", string(reqBody))

	// Create request
	req := httptest.NewRequest("PUT", "/user/update-username", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// DEBUG: Print headers to see what's actually being sent
	t.Logf("Request headers: %v", req.Header)
	t.Logf("Content-Type header: '%s'", req.Header.Get("Content-Type"))
	t.Logf("Request method: %s", req.Method)
	t.Logf("Request body length: %d", len(reqBody))

	w := httptest.NewRecorder()

	// Test the middleware directly first to isolate the issue
	middleware_obj := middleware.ParseJSON(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("Middleware passed successfully!")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))

	middleware_obj.ServeHTTP(w, req)

	t.Logf("Middleware test - Response status: %d", w.Code)
	t.Logf("Middleware test - Response body: %s", w.Body.String())

	// If middleware test passes, test the actual handler
	if w.Code == http.StatusOK {
		reqBody2, _ := json.Marshal(updateReq)
		req2 := httptest.NewRequest("PUT", "/user/update-username", bytes.NewBuffer(reqBody2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()

		// Apply middleware to handler
		handler := middleware.ParseJSON(http.HandlerFunc(userHandler.UpdateUsername))
		handler.ServeHTTP(w2, req2)

		t.Logf("Handler test - Response status: %d", w2.Code)
		t.Logf("Handler test - Response body: %s", w2.Body.String())

		assert.Equal(t, http.StatusOK, w2.Code)

		var response UserAPIResponse
		unmarshal_err := json.Unmarshal(w2.Body.Bytes(), &response)
		assert.NoError(t, unmarshal_err)

		assert.Equal(t, created_user.ID, response.Data.ID)
		assert.Equal(t, "bob", response.Data.Username)
	}
}

func TestUpdateUsernameUserNotFound(t *testing.T) {
	// Setup
	_, _, userHandler := setupTestUserHandler(t)

	// Test request with non-existent user
	updateReq := service.UpdateUsernameRequest{
		OldUsername: "nonexistent",
		NewUsername: "bob",
	}
	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/user/update-username", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := middleware.ParseJSON(http.HandlerFunc(userHandler.UpdateUsername))
	handler.ServeHTTP(w, req)

	// DEBUG: Print response body
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response status: %d", w.Code)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response UserAPIResponse
	unmarshal_err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, unmarshal_err)

	assert.Equal(t, "", response.Data.Username)
	assert.Equal(t, uuid.Nil, response.Data.ID)
}

func TestUpdateUsernameBadRequest(t *testing.T) {
	// Setup
	_, _, userHandler := setupTestUserHandler(t)

	// Test request with invalid JSON body
	req := httptest.NewRequest("PUT", "/user/update-username", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := middleware.ParseJSON(http.HandlerFunc(userHandler.UpdateUsername))
	handler.ServeHTTP(w, req)

	// DEBUG: Print response body
	t.Logf("Response body: %s", w.Body.String())
	t.Logf("Response status: %d", w.Code)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
