package repository_test

import (
	"context"
	"testing"

	"stride-wars-app/ent/enttest"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// TODO to properly test this we probably need to set up local postgres container????

func TestUserRepository_FindByID(t *testing.T) {
	t.Parallel()
	t.Run("correctly find user by id", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer func() {
			if err := client.Close(); err != nil {
				t.Fatalf("closing client: %v", err)
			}
		}()

		ctx := context.Background()
		repo := repository.NewUserRepository(client)

		username := "alice"
		externalID := uuid.New()

		new_user := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created, err := repo.CreateUser(ctx, new_user)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, created.ID)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
		require.Equal(t, username, found.Username)
		require.Equal(t, externalID, found.ExternalUser)
	})
}

func TestUserRepository_FindByUsername(t *testing.T) {
	t.Parallel()

	t.Run("correctly find user by username", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer func() {
			if err := client.Close(); err != nil {
				t.Fatalf("closing client: %v", err)
			}
		}()

		ctx := context.Background()
		repo := repository.NewUserRepository(client)

		username := "alice"
		externalID := uuid.New()

		new_user := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created, err := repo.CreateUser(ctx, new_user)
		require.NoError(t, err)

		found, err := repo.FindByUsername(ctx, created.Username)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
		require.Equal(t, username, found.Username)
		require.Equal(t, externalID, found.ExternalUser)
	})
}

func TestUserRepository_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("correctly create user", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer func() {
			if err := client.Close(); err != nil {
				t.Fatalf("closing client: %v", err)
			}
		}()

		ctx := context.Background()
		repo := repository.NewUserRepository(client)

		username := "alice"
		externalID := uuid.New()

		new_user := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created, err := repo.CreateUser(ctx, new_user)
		require.NoError(t, err)

		found, err := repo.FindByUsername(ctx, created.Username)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
		require.Equal(t, username, found.Username)
		require.Equal(t, externalID, found.ExternalUser)
	})
}

func TestUserRepository_UpdateUsername(t *testing.T) {
	t.Parallel()

	t.Run("correctly update user", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer func() {
			if err := client.Close(); err != nil {
				t.Fatalf("closing client: %v", err)
			}
		}()

		ctx := context.Background()
		repo := repository.NewUserRepository(client)

		username := "alice"
		externalID := uuid.New()

		new_user := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created, err := repo.CreateUser(ctx, new_user)
		require.NoError(t, err)

		updated_username := "bob" // sex change
		updated_user := &model.User{
			ID:           created.ID,
			Username:     updated_username,
			ExternalUser: externalID,
		}
		no_rows, err := repo.UpdateUsername(ctx, updated_user)
		require.NoError(t, err)

		found, err := repo.FindByUsername(ctx, updated_username)
		require.NoError(t, err)

		require.Equal(t, 1, no_rows)
		require.Equal(t, created.ID, found.ID)
		require.Equal(t, updated_username, found.Username)
		require.Equal(t, externalID, found.ExternalUser)
	})
}
