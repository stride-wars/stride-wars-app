package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"stride-wars-app/ent/model"
	"stride-wars-app/internal/testutil"
)

func TestUserRepository(t *testing.T) {
	t.Parallel()

	t.Run("correctly find user by id", func(t *testing.T) {
		t.Parallel()

		tdb := testutil.NewTestServices(t)
		ctx := tdb.Ctx

		userRepo := tdb.UserRepo

		username := "alice"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		found, err := userRepo.FindByID(ctx, created.ID)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
		require.Equal(t, username, found.Username)
		require.Equal(t, externalID, found.ExternalUser)
	})

	t.Run("correctly find user by username", func(t *testing.T) {
		t.Parallel()

		tdb := testutil.NewTestServices(t)
		ctx := tdb.Ctx

		userRepo := tdb.UserRepo

		username := "alice"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		found, err := userRepo.FindByUsername(ctx, username)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
		require.Equal(t, username, found.Username)
		require.Equal(t, externalID, found.ExternalUser)
	})

	t.Run("correctly create user", func(t *testing.T) {
		t.Parallel()

		tdb := testutil.NewTestServices(t)
		ctx := tdb.Ctx

		userRepo := tdb.UserRepo

		username := "alice"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		found, err := userRepo.FindByUsername(ctx, username)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
		require.Equal(t, username, found.Username)
		require.Equal(t, externalID, found.ExternalUser)
	})

	t.Run("correctly update user", func(t *testing.T) {
		t.Parallel()

		tdb := testutil.NewTestServices(t)
		ctx := tdb.Ctx

		userRepo := tdb.UserRepo

		username := "alice"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		updatedUsername := "bob"
		updatedUser := &model.User{
			ID:           created.ID,
			Username:     updatedUsername,
			ExternalUser: externalID,
		}
		rowsChanged, err := userRepo.UpdateUsername(ctx, updatedUser)
		require.NoError(t, err)
		require.Equal(t, 1, rowsChanged)

		found, err := userRepo.FindByUsername(ctx, updatedUsername)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
		require.Equal(t, updatedUsername, found.Username)
		require.Equal(t, externalID, found.ExternalUser)
	})
}
