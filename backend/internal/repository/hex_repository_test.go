package repository_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uber/h3-go/v4"

	"stride-wars-app/internal/testutil"
)

func TestHexRepository(t *testing.T) {
	t.Parallel()

	t.Run("correctly find hex by id", func(t *testing.T) {
		t.Parallel()

		// 1) Spin up a new in-memory DB + schema migration:
		tdb := testutil.NewTestServices(t)
		ctx := tdb.Ctx

		// 2) Use the injected HexRepo from TestServices:
		hexRepo := tdb.HexRepo

		// 3) Create a valid H3 index (using h3 library):
		latLng := h3.NewLatLng(37.775938728915946, -122.41795063018799)
		resolution := 9
		h3Index, err := h3.LatLngToCell(latLng, resolution)
		require.NoError(t, err)
		H3String := h3Index.String()

		// 4) Insert a new Hex row:
		created, err := hexRepo.CreateHex(ctx, H3String)
		require.NoError(t, err)
		require.NotNil(t, created)

		// 5) Find by ID:
		found, err := hexRepo.FindByID(ctx, created.ID)
		require.NoError(t, err)

		require.Equal(t, created.ID, found.ID)
	})
}
