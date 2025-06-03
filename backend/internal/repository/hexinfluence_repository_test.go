package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"stride-wars-app/ent/model"
	"stride-wars-app/internal/testutil"

	"github.com/google/uuid"
)

func TestHexInfluenceRepository(t *testing.T) {
	t.Parallel()

	t.Run("correctly update a hex influence", func(t *testing.T) {
		t.Parallel()

		// 1) New in-memory DB:
		tdb := testutil.NewTestServices(t)
		ctx := tdb.Ctx

		// 2) Pull repos from TestServices:
		inflRepo := tdb.HexInfluenceRepo
		userRepo := tdb.UserRepo
		hexRepo := tdb.HexRepo

		// 3) Create a hex row first:
		createdHex, err := hexRepo.CreateHex(ctx, "85283473fffffff")
		require.NoError(t, err)
		require.NotNil(t, createdHex)

		// 4) Create a user:
		username := "alice"
		externalID := uuid.New()
		newUser := &model.User{
			Username:     username,
			ExternalUser: externalID,
		}
		createdUser, err := userRepo.CreateUser(ctx, newUser)
		require.NoError(t, err)

		// 5) Insert a HexInfluence with Score=10.0 and LastUpdated 15 days ago:
		userID := createdUser.ID
		h3Index := "85283473fffffff"
		lastUpdated := time.Now().Add(-15 * 24 * time.Hour)
		score := 10.0

		newInfluence := &model.HexInfluence{
			UserID:      userID,
			H3Index:     h3Index,
			LastUpdated: lastUpdated,
			Score:       score,
		}
		createdInfluence, err := inflRepo.CreateHexInfluence(ctx, newInfluence)
		require.NoError(t, err)

		// 6) Update the influence (should decrement Score to 9.0 and set LastUpdated=now):
		rowsChanged, err := inflRepo.UpdateHexInfluence(ctx, createdInfluence.UserID, createdInfluence.H3Index)
		require.NoError(t, err)
		require.Equal(t, 1, rowsChanged)

		// 7) Fetch the updated row:
		updated, err := inflRepo.FindByUserIDAndHexID(ctx, createdInfluence.UserID, createdInfluence.H3Index)
		require.NoError(t, err)

		require.Equal(t, createdInfluence.ID, updated.ID)
		require.Equal(t, userID, updated.UserID)
		require.Equal(t, h3Index, updated.H3Index)
		require.Equal(t, 9.0, updated.Score)

		// 8) LastUpdated was set to “now”:
		require.WithinDuration(t, time.Now(), updated.LastUpdated, time.Second*2)
	})
}
