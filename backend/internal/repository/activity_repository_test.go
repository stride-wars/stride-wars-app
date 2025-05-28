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

	"time"
)

// TODO to properly test this we probably need to set up local postgres container????

func TestActivityRepository_FindByID(t *testing.T) {
	t.Parallel()
	t.Run("correctly find activity by id", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer func() {
			if err := client.Close(); err != nil {
				t.Fatalf("closing client: %v", err)
			}
		}()

		ctx := context.Background()
		activity_repo := repository.NewActivityRepository(client)
		user_repo := repository.NewUserRepository(client)

		username := "alice"
		externalID := uuid.New()

		new_user := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		created_user, err := user_repo.CreateUser(ctx, new_user)
		require.NoError(t, err)

		userID := created_user.ID
		duration := 782.5
		distance := 1500.0
		h3_indexes := []int64{
			617700169958293503,
			617700169957507071,
			617700169957769215,
			617700169958031359,
			617700169958162431,
			617700169958096895,
			617700169957375999,
			617700169957441535,
		}
		created_at := time.Now()

		new_activity := &model.Activity{
			UserID:    userID,
			Duration:  duration,
			Distance:  distance,
			H3Indexes: h3_indexes,
			CreatedAt: created_at,
		}
		created, err := activity_repo.CreateActivity(ctx, new_activity)
		require.NoError(t, err)

		found, err := activity_repo.FindByID(ctx, created.ID)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
		require.Equal(t, userID, found.UserID)
		require.Equal(t, duration, found.DurationSeconds)
		require.Equal(t, distance, found.DistanceMeters)
		require.WithinDuration(t,
			created_at,
			found.CreatedAt, // actual
			time.Second,     // tolerance
		)

	})
}
