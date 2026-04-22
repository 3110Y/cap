package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/3110Y/cap/internal/application/service"
	"github.com/3110Y/cap/internal/domain/entity"
	"github.com/3110Y/cap/internal/integration/mocks"
)

// newInstalledSvc is a helper that wires up InstalledService with mocks.
func newInstalledSvc(
	t *testing.T,
	installed *mocks.MockInstalledRepository,
	marketplace *mocks.MockMarketplaceLister,
	installer *mocks.MockPrimitiveInstaller,
) *service.InstalledService {
	t.Helper()
	return service.NewInstalledServiceWith(installed, marketplace, installer)
}

// --- List ---

func TestInstalledService_List_DelegatesToRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	want := []*entity.Primitive{
		{Type: "skills", Vendor: "v", Name: "n", Version: "1.0.0"},
	}
	mockInstalled.EXPECT().List("skills").Return(want, nil)

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	got, err := svc.List("skills")

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestInstalledService_List_PropagatesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	mockInstalled.EXPECT().List("hooks").Return(nil, errors.New("fs error"))

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	_, err := svc.List("hooks")

	require.Error(t, err)
}

// --- Info ---

func TestInstalledService_Info_CombinesInstalledAndMarketplace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	installedPrims := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.1.0"},
	}
	availablePrims := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.2.0", Description: "Latest"},
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.1.0", Description: "Older"},
	}
	mockInstalled.EXPECT().List("skills").Return(installedPrims, nil)
	mockMarket.EXPECT().Info("skills", "anthropic/code-review", "").Return(availablePrims, nil)

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	info, err := svc.Info("skills", "code-review")

	require.NoError(t, err)
	assert.Equal(t, "code-review", info.Name)
	assert.Equal(t, "anthropic", info.Vendor)
	assert.Equal(t, []string{"1.1.0"}, info.InstalledVersions)
	assert.ElementsMatch(t, []string{"1.2.0", "1.1.0"}, info.AvailableVersions)
	assert.Equal(t, "Latest", info.Description)
}

func TestInstalledService_Info_ByVendorName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	installedPrims := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.0.0"},
	}
	mockInstalled.EXPECT().List("skills").Return(installedPrims, nil)
	mockMarket.EXPECT().Info("skills", "anthropic/code-review", "").Return(installedPrims, nil)

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	info, err := svc.Info("skills", "anthropic/code-review")

	require.NoError(t, err)
	assert.Equal(t, "code-review", info.Name)
	assert.Equal(t, "anthropic", info.Vendor)
}

func TestInstalledService_Info_NotInstalled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	mockInstalled.EXPECT().List("skills").Return([]*entity.Primitive{}, nil)

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	_, err := svc.Info("skills", "missing-skill")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not installed")
}

// --- Delete ---

func TestInstalledService_Delete_WithoutVendor_FindsVendor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	installed := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.0.0"},
	}
	mockInstalled.EXPECT().List("skills").Return(installed, nil)
	mockInstalled.EXPECT().Delete("skills", "anthropic", "code-review", "").Return(nil)

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	err := svc.Delete("skills", "code-review", "")

	require.NoError(t, err)
}

func TestInstalledService_Delete_WithVendor_DeletesDirectly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	// When vendor is given, List is not called.
	mockInstalled.EXPECT().Delete("skills", "anthropic", "code-review", "1.0.0").Return(nil)

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	err := svc.Delete("skills", "anthropic/code-review", "1.0.0")

	require.NoError(t, err)
}

func TestInstalledService_Delete_Ambiguous_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	installed := []*entity.Primitive{
		{Type: "skills", Vendor: "vendor-a", Name: "skill", Version: "1.0.0"},
		{Type: "skills", Vendor: "vendor-b", Name: "skill", Version: "1.0.0"},
	}
	mockInstalled.EXPECT().List("skills").Return(installed, nil)

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	err := svc.Delete("skills", "skill", "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ambiguous")
}

func TestInstalledService_Delete_NotInstalled_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	mockInstalled.EXPECT().List("skills").Return([]*entity.Primitive{}, nil)

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	err := svc.Delete("skills", "nonexistent", "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not installed")
}

// --- Upgrade ---

func TestInstalledService_Upgrade_CallsInstallAddWithoutVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	want := &entity.Primitive{Type: "skills", Vendor: "v", Name: "code-review", Version: "2.0.0"}
	mockInstaller.EXPECT().Add("skills", "code-review", "").Return(want, nil)

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	got, err := svc.Upgrade("skills", "code-review")

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

// --- UpgradeAll ---

func TestInstalledService_UpgradeAll_UpgradesUniqueNames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	// skills has two versions of the same primitive (same vendor/name), agents has one primitive.
	skillsPrims := []*entity.Primitive{
		{Type: "skills", Vendor: "v", Name: "skill-a", Version: "1.0.0"},
		{Type: "skills", Vendor: "v", Name: "skill-a", Version: "1.1.0"}, // duplicate key
	}
	agentsPrims := []*entity.Primitive{
		{Type: "agents", Vendor: "v", Name: "agent-b", Version: "1.0.0"},
	}

	// All 5 types are iterated; 3 return empty, skills and agents return data.
	mockInstalled.EXPECT().List("skills").Return(skillsPrims, nil)
	mockInstalled.EXPECT().List("agents").Return(agentsPrims, nil)
	mockInstalled.EXPECT().List("commands").Return([]*entity.Primitive{}, nil)
	mockInstalled.EXPECT().List("hooks").Return([]*entity.Primitive{}, nil)
	mockInstalled.EXPECT().List("rules").Return([]*entity.Primitive{}, nil)

	// Each unique vendor/name upgraded once.
	mockInstaller.EXPECT().Add("skills", "v/skill-a", "").Return(
		&entity.Primitive{Type: "skills", Vendor: "v", Name: "skill-a", Version: "1.2.0"}, nil)
	mockInstaller.EXPECT().Add("agents", "v/agent-b", "").Return(
		&entity.Primitive{Type: "agents", Vendor: "v", Name: "agent-b", Version: "2.0.0"}, nil)

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	upgraded, errs := svc.UpgradeAll()

	assert.Empty(t, errs)
	assert.Len(t, upgraded, 2)
}

func TestInstalledService_UpgradeAll_CollectsErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	prim := []*entity.Primitive{{Type: "skills", Vendor: "v", Name: "bad-skill", Version: "1.0.0"}}
	mockInstalled.EXPECT().List("skills").Return(prim, nil)
	mockInstalled.EXPECT().List("agents").Return([]*entity.Primitive{}, nil)
	mockInstalled.EXPECT().List("commands").Return([]*entity.Primitive{}, nil)
	mockInstalled.EXPECT().List("hooks").Return([]*entity.Primitive{}, nil)
	mockInstalled.EXPECT().List("rules").Return([]*entity.Primitive{}, nil)

	mockInstaller.EXPECT().Add("skills", "v/bad-skill", "").Return(nil, errors.New("not found in marketplace"))

	svc := newInstalledSvc(t, mockInstalled, mockMarket, mockInstaller)
	upgraded, errs := svc.UpgradeAll()

	assert.Empty(t, upgraded)
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "not found in marketplace")
}
