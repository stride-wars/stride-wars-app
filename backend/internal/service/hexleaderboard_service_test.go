package service_test

import (
	"math/rand"
	"strconv"
	"stride-wars-app/ent"
	"stride-wars-app/ent/model"
	"stride-wars-app/internal/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHexLeaderboardService_GlobalLeaderboard_ManyUsers(t *testing.T) {
	t.Parallel()

	tdb := testutil.NewTestServices(t)
	ctx := tdb.Ctx
	hexLeaderboardService := tdb.HexLeaderboardService
	userRepo := tdb.UserRepo
	hexSvc := tdb.HexService

	var users []*ent.User
	for i := 0; i < 25; i++ {
		user := &model.User{
			Username:     "user" + string(rune('A'+i)),
			ExternalUser: uuid.New(),
		}
		created, err := userRepo.CreateUser(ctx, user)
		require.NoError(t, err)
		users = append(users, created)
	}

	for h := 0; h < 30; h++ {
		topUserID := users[rand.Intn(len(users))].ID
		topUserUsername := users[rand.Intn(len(users))].Username

		// Now using string H3 index instead of int64
		idxStr := strconv.Itoa(1000 + h)
		_, err := hexSvc.CreateHex(ctx, idxStr)
		require.NoError(t, err)

		_, err = hexLeaderboardService.CreateHexLeaderboard(ctx, &model.HexLeaderboard{
			H3Index:  idxStr,
			TopUsers: []model.TopUser{{UserID: topUserID, UserName: topUserUsername, Score: float64(rand.Intn(100))}},
		})
		require.NoError(t, err)
	}

	entries, err := hexLeaderboardService.GetGlobalLeaderboard(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.LessOrEqual(t, len(entries), 10)

	require.GreaterOrEqual(t, entries[0].TopCount, entries[len(entries)-1].TopCount)
}
