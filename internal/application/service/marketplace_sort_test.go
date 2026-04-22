package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/3110Y/cap/internal/domain/entity"
)

// TestMarketplaceService_Info_SortsByLatestFirst verifies that Info returns
// versions in descending SemVer order so callers can take result[0] as "latest".
func TestMarketplaceService_Info_SortsByLatestFirst(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, []*entity.Primitive{
		{Type: "skills", Vendor: "v", Name: "n", Version: "1.0.0"},
		{Type: "skills", Vendor: "v", Name: "n", Version: "1.10.0"},
		{Type: "skills", Vendor: "v", Name: "n", Version: "1.2.0"},
		{Type: "skills", Vendor: "v", Name: "n", Version: "2.0.0"},
	})
	defer ctrl.Finish()

	got, err := svc.Info("skills", "n", "")

	require.NoError(t, err)
	require.Len(t, got, 4)
	assert.Equal(t, "2.0.0", got[0].Version, "latest version must come first")
	assert.Equal(t, "1.10.0", got[1].Version, "1.10.0 must sort above 1.2.0")
	assert.Equal(t, "1.2.0", got[2].Version)
	assert.Equal(t, "1.0.0", got[3].Version)
}

// TestMarketplaceService_Info_HandlesVPrefix confirms that "v1.2.3" and "1.2.3" compare equal.
func TestMarketplaceService_Info_HandlesVPrefix(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, []*entity.Primitive{
		{Type: "skills", Vendor: "v", Name: "n", Version: "v1.0.0"},
		{Type: "skills", Vendor: "v", Name: "n", Version: "2.0.0"},
	})
	defer ctrl.Finish()

	got, err := svc.Info("skills", "n", "")

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "2.0.0", got[0].Version)
}
