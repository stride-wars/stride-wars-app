package repository_test

import (
	"context"
	"testing"
	"time"

	"stride-wars-app/ent/enttest"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/repository"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// TODO to properly test this we probably need to set up local postgres container????

func TestHexInfluenceRepository_UpdateHexInfluence(t *testing.T) {
	t.Parallel()
	t.Run("correctly update a hex influence", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer func() {
			if err := client.Close(); err != nil {
				t.Fatalf("closing client: %v", err)
			}
		}()

		ctx := context.Background()
		influence_repo := repository.NewHexInfluenceRepository(client)
		user_repo := repository.NewUserRepository(client)
		hex_repo := repository.NewHexRepository(client)

		created_hex, err := hex_repo.CreateHex(ctx, 617700169958293503)
		require.NoError(t, err)
		require.NotNil(t, created_hex)

		username := "alice"
		externalID := uuid.New()

		new_user := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created_user, err := user_repo.CreateUser(ctx, new_user)
		require.NoError(t, err)

		userID := created_user.ID
		h3_index := int64(617700169958293503)
		last_updated := time.Now().Add(-14 * 24 * time.Hour)
		score := 10.0

		new_hex_influence := &model.HexInfluence{
			UserID:      userID,
			H3Index:     h3_index,
			LastUpdated: last_updated,
			Score:       score,
		}
		created, err := influence_repo.CreateHexInfluence(ctx, new_hex_influence)
		require.NoError(t, err)

		rows_changed, err := influence_repo.UpdateHexInfluence(ctx, created.UserID, created.H3Index)
		require.NoError(t, err)
		updated_influence, err := influence_repo.FindByUserIDAndHexID(ctx, created.UserID, created.H3Index)
		require.NoError(t, err)

		require.Equal(t, 1, rows_changed)
		require.Equal(t, created.ID, updated_influence.ID)
		require.Equal(t, 8.0, updated_influence.Score)
		require.Equal(t, created.UserID, updated_influence.UserID)
		require.Equal(t, created.H3Index, h3_index)
		require.WithinDuration(t,
			time.Now(),
			updated_influence.LastUpdated, // actual
			time.Second*2,                 // tolerance
		)

	})
}
