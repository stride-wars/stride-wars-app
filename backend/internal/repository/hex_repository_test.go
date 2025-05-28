package repository_test

import (
	"context"
	"testing"

	"stride-wars-app/ent/enttest"
	"stride-wars-app/internal/repository"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/uber/h3-go/v4"
)

// TODO to properly test this we probably need to set up local postgres container????

func TestHexRepository_FindByID(t *testing.T) {
	t.Parallel()
	t.Run("correctly find hex by id", func(t *testing.T) {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		defer func() {
			if err := client.Close(); err != nil {
				t.Fatalf("closing client: %v", err)
			}
		}()

		ctx := context.Background()
		repo := repository.NewHexRepository(client)

		latLng := h3.NewLatLng(37.775938728915946, -122.41795063018799)
		resolution := 9

		h3_index, err := h3.LatLngToCell(latLng, resolution)
		require.NoError(t, err)

		created, err := repo.CreateHex(ctx, int64(h3_index))
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, created.ID)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
	})
}
