package handler_test

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/3110Y/cap/internal/application/service"
	"github.com/3110Y/cap/internal/domain/entity"
	"github.com/3110Y/cap/internal/integration/mocks"
	"github.com/3110Y/cap/internal/presentation/handler"
)

func newInstalledHandler(
	t *testing.T,
	mockInstalled *mocks.MockInstalledRepository,
	mockMarket *mocks.MockMarketplaceLister,
	mockInstaller *mocks.MockPrimitiveInstaller,
) *handler.InstalledHandler {
	t.Helper()
	svc := service.NewInstalledServiceWith(mockInstalled, mockMarket, mockInstaller)
	return handler.NewInstalledHandler(svc)
}

func newCmdWithFlags() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().Int("page", 1, "")
	cmd.Flags().Int("size", 50, "")
	cmd.Flags().String("version", "", "")
	return cmd
}

// TestInstalledHandler_List_PrintsTable verifies the NAME/VERSION table.
func TestInstalledHandler_List_PrintsTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	primitives := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.2.0"},
	}
	mockInstalled.EXPECT().List("skills").Return(primitives, nil)

	h := newInstalledHandler(t, mockInstalled, mockMarket, mockInstaller)
	cmd := newCmdWithFlags()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.List("skills")(cmd, nil)

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "NAME")
	assert.Contains(t, out, "VERSION")
	assert.Contains(t, out, "code-review")
	assert.Contains(t, out, "1.2.0")
}

// TestInstalledHandler_List_Empty prints the "no X installed" message.
func TestInstalledHandler_List_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	mockInstalled.EXPECT().List("skills").Return([]*entity.Primitive{}, nil)

	h := newInstalledHandler(t, mockInstalled, mockMarket, mockInstaller)
	cmd := newCmdWithFlags()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.List("skills")(cmd, nil)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "No skills installed")
}

// TestInstalledHandler_Info_PrintsDetails verifies info output.
func TestInstalledHandler_Info_PrintsDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	installedPrims := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.2.0"},
	}
	availablePrims := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.2.0", Description: "Best skill"},
	}
	mockInstalled.EXPECT().List("skills").Return(installedPrims, nil)
	mockMarket.EXPECT().Info("skills", "anthropic/code-review", "").Return(availablePrims, nil)

	h := newInstalledHandler(t, mockInstalled, mockMarket, mockInstaller)
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.Info("skills")(cmd, []string{"code-review"})

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Name:")
	assert.Contains(t, out, "code-review")
	assert.Contains(t, out, "Vendor:")
	assert.Contains(t, out, "anthropic")
	assert.Contains(t, out, "Namespace:")
	assert.Contains(t, out, "cap:anthropic:code-review")
	assert.Contains(t, out, "Installed versions:")
}

// TestInstalledHandler_Delete_AllVersions verifies the "all versions" message.
func TestInstalledHandler_Delete_AllVersions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	installed := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.2.0"},
	}
	mockInstalled.EXPECT().List("skills").Return(installed, nil)
	mockInstalled.EXPECT().Delete("skills", "anthropic", "code-review", "").Return(nil)

	h := newInstalledHandler(t, mockInstalled, mockMarket, mockInstaller)
	cmd := newCmdWithFlags()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.Delete("skills")(cmd, []string{"code-review"})

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "all versions")
}

// TestInstalledHandler_Upgrade_PrintsUpgraded verifies the upgrade output.
func TestInstalledHandler_Upgrade_PrintsUpgraded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstalled := mocks.NewMockInstalledRepository(ctrl)
	mockMarket := mocks.NewMockMarketplaceLister(ctrl)
	mockInstaller := mocks.NewMockPrimitiveInstaller(ctrl)

	upgraded := &entity.Primitive{
		Type:    "skills",
		Vendor:  "anthropic",
		Name:    "code-review",
		Version: "2.0.0",
	}
	mockInstaller.EXPECT().Add("skills", "code-review", "").Return(upgraded, nil)

	h := newInstalledHandler(t, mockInstalled, mockMarket, mockInstaller)
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.Upgrade("skills")(cmd, []string{"code-review"})

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Upgraded:")
	assert.Contains(t, out, "cap:anthropic:code-review")
	assert.Contains(t, out, "2.0.0")
}
