package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/api/middleware"
	"stride-wars-app/internal/handler"
	"stride-wars-app/internal/repository"
	"stride-wars-app/internal/service"
	"stride-wars-app/internal/testutil"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type UserAPIResponse struct {
	Success bool     `json:"success"`
	Data    UserData `json:"data"`
}

type UserData struct {
	ID           uuid.UUID     `json:"id"`
	Username     string        `json:"username"`
	ExternalUser uuid.UUID     `json:"external_user"`
	Edges        ent.UserEdges `json:"edges"`
}

func setupTestUserHandler(t *testing.T) (context.Context, *ent.Client, *handler.UserHandler) {
	t.Helper()

	svc := testutil.NewTestServices(t)
	userHandler := handler.NewUserHandler(svc.UserService, zap.NewExample())

	return svc.Ctx, svc.Client, userHandler
}

func TestUserHandler(t *testing.T) {
	t.Parallel()

	// ------------------------
	// Subtest: GetUserByUsername/HappyPath
	// ------------------------
	t.Run("GetUserByUsername/HappyPath", func(t *testing.T) {
		t.Parallel()

		ctx, client, userHandler := setupTestUserHandler(t)

		// Seed a user
		repo := repository.NewUserRepository(client)
		username := "alice"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := repo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		// Ensure user is in DB
		foundUser, findErr := repo.FindByUsername(ctx, username)
		require.NoError(t, findErr)
		require.Equal(t, createdUser.ID, foundUser.ID)

		req := httptest.NewRequest("GET", "/user/username?username=alice", nil)
		w := httptest.NewRecorder()

		userHandler.GetUserByUsername(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response UserAPIResponse
		unmarshalErr := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, unmarshalErr)

		assert.Equal(t, "alice", response.Data.Username)
		assert.NotEqual(t, uuid.Nil, response.Data.ID)
	})

	// ------------------------
	// Subtest: GetUserByUsername/NotFound
	// ------------------------
	t.Run("GetUserByUsername/NotFound", func(t *testing.T) {
		t.Parallel()

		_, _, userHandler := setupTestUserHandler(t)

		req := httptest.NewRequest("GET", "/user/username?username=alice", nil)
		w := httptest.NewRecorder()

		userHandler.GetUserByUsername(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response UserAPIResponse
		unmarshalErr := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, unmarshalErr)

		assert.Equal(t, "", response.Data.Username)
		assert.Equal(t, uuid.Nil, response.Data.ID)
	})

	// ------------------------
	// Subtest: GetUserByUsername/BadRequest
	// ------------------------
	t.Run("GetUserByUsername/BadRequest", func(t *testing.T) {
		t.Parallel()

		_, _, userHandler := setupTestUserHandler(t)

		req := httptest.NewRequest("GET", "/user/username?usernname=alice", nil)
		// notice the typo in the query parameter
		w := httptest.NewRecorder()

		userHandler.GetUserByUsername(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// ------------------------
	// Subtest: GetUserByID/HappyPath
	// ------------------------
	t.Run("GetUserByID/HappyPath", func(t *testing.T) {
		t.Parallel()

		ctx, client, userHandler := setupTestUserHandler(t)

		repo := repository.NewUserRepository(client)
		username := "alice"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := repo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		foundUser, findErr := repo.FindByID(ctx, createdUser.ID)
		require.NoError(t, findErr)
		require.Equal(t, createdUser.ID, foundUser.ID)

		req := httptest.NewRequest("GET", fmt.Sprintf("/user/id?id=%s", createdUser.ID), nil)
		w := httptest.NewRecorder()

		userHandler.GetUserByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response UserAPIResponse
		unmarshalErr := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, unmarshalErr)

		assert.Equal(t, createdUser.ID, response.Data.ID)
		assert.Equal(t, "alice", response.Data.Username)
	})

	// ------------------------
	// Subtest: GetUserByID/NotFound
	// ------------------------
	t.Run("GetUserByID/NotFound", func(t *testing.T) {
		t.Parallel()

		_, _, userHandler := setupTestUserHandler(t)

		nonExistentID := uuid.New()
		req := httptest.NewRequest("GET", fmt.Sprintf("/user/id?id=%s", nonExistentID), nil)
		w := httptest.NewRecorder()

		userHandler.GetUserByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response UserAPIResponse
		unmarshalErr := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, unmarshalErr)

		assert.Equal(t, "", response.Data.Username)
		assert.Equal(t, uuid.Nil, response.Data.ID)
	})

	// ------------------------
	// Subtest: GetUserByID/BadRequest
	// ------------------------
	t.Run("GetUserByID/BadRequest", func(t *testing.T) {
		t.Parallel()

		_, _, userHandler := setupTestUserHandler(t)

		req := httptest.NewRequest("GET", "/user/id?id=invalid-uuid", nil)
		w := httptest.NewRecorder()

		userHandler.GetUserByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// ------------------------
	// Subtest: UpdateUsername/HappyPath
	// ------------------------
	t.Run("UpdateUsername/HappyPath", func(t *testing.T) {
		t.Parallel()

		ctx, client, userHandler := setupTestUserHandler(t)

		// Seed a user
		repo := repository.NewUserRepository(client)
		username := "alice"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := repo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		updateReq := service.UpdateUsernameRequest{
			OldUsername: "alice",
			NewUsername: "bob",
		}
		reqBody, err := json.Marshal(updateReq)
		require.NoError(t, err)

		req := httptest.NewRequest("PUT", "/user/update", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlerWithMiddleware := middleware.ParseJSON(http.HandlerFunc(userHandler.UpdateUsername))
		handlerWithMiddleware.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response UserAPIResponse
		unmarshalErr := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, unmarshalErr)

		assert.Equal(t, createdUser.ID, response.Data.ID)
		assert.Equal(t, "bob", response.Data.Username)
	})

	// ------------------------
	// Subtest: UpdateUsername/UserNotFound
	// ------------------------
	t.Run("UpdateUsername/UserNotFound", func(t *testing.T) {
		t.Parallel()

		_, _, userHandler := setupTestUserHandler(t)

		updateReq := service.UpdateUsernameRequest{
			OldUsername: "nonexistent",
			NewUsername: "bob",
		}
		reqBody, err := json.Marshal(updateReq)
		require.NoError(t, err)

		req := httptest.NewRequest("PUT", "/user/update", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlerWithMiddleware := middleware.ParseJSON(http.HandlerFunc(userHandler.UpdateUsername))
		handlerWithMiddleware.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response UserAPIResponse
		unmarshalErr := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, unmarshalErr)

		assert.Equal(t, "", response.Data.Username)
		assert.Equal(t, uuid.Nil, response.Data.ID)
	})

	// ------------------------
	// Subtest: UpdateUsername/BadRequest
	// ------------------------
	t.Run("UpdateUsername/BadRequest", func(t *testing.T) {
		t.Parallel()

		_, _, userHandler := setupTestUserHandler(t)

		req := httptest.NewRequest("PUT", "/user/update", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlerWithMiddleware := middleware.ParseJSON(http.HandlerFunc(userHandler.UpdateUsername))
		handlerWithMiddleware.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
