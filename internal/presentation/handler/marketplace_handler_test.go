package handler_test

import (
	"bytes"
	"path/filepath"
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

func newMarketplaceHandler(
	t *testing.T,
	mockCache *mocks.MockCacheRepository,
	mockProject *mocks.MockProjectRepository,
) *handler.MarketplaceHandler {
	t.Helper()
	mp := service.NewMarketplaceService(mockCache)
	install := service.NewInstallService(mockCache, mockProject, mp)
	return handler.NewMarketplaceHandler(mp, install)
}

func newPaginatedCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().Int("page", 1, "")
	cmd.Flags().Int("size", 50, "")
	return cmd
}

// TestMarketplaceHandler_List_PrintsTable verifies the NAME/VERSION table output.
func TestMarketplaceHandler_List_PrintsTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	primitives := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.2.0"},
		{Type: "skills", Vendor: "community", Name: "pdf-skill", Version: "0.5.1"},
	}
	mockCache.EXPECT().ReadIndexes().Return(primitives, nil)

	h := newMarketplaceHandler(t, mockCache, mockProject)
	cmd := newPaginatedCmd()
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

// TestMarketplaceHandler_List_Empty prints the "no primitives" message.
func TestMarketplaceHandler_List_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	mockCache.EXPECT().ReadIndexes().Return([]*entity.Primitive{}, nil)

	h := newMarketplaceHandler(t, mockCache, mockProject)
	cmd := newPaginatedCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.List("skills")(cmd, nil)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "No primitives available")
}

// TestMarketplaceHandler_Search_ReturnsMatches prints NAME/VERSION/DESCRIPTION.
func TestMarketplaceHandler_Search_ReturnsMatches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	primitives := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.2.0", Description: "Analyse Go code"},
	}
	mockCache.EXPECT().ReadIndexes().Return(primitives, nil)

	h := newMarketplaceHandler(t, mockCache, mockProject)
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.Search("skills")(cmd, []string{"Go"})

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "code-review")
	assert.Contains(t, out, "1.2.0")
	assert.Contains(t, out, "Analyse Go code")
}

// TestMarketplaceHandler_Search_NoMatches prints the "no matches" message.
func TestMarketplaceHandler_Search_NoMatches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	mockCache.EXPECT().ReadIndexes().Return([]*entity.Primitive{
		{Type: "skills", Vendor: "v", Name: "n", Version: "1.0.0", Description: "something"},
	}, nil)

	h := newMarketplaceHandler(t, mockCache, mockProject)
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.Search("skills")(cmd, []string{"ничего_нет"})

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "No matches found")
}

// TestMarketplaceHandler_Info_PrintsDetails verifies info output fields.
func TestMarketplaceHandler_Info_PrintsDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	primitives := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.2.0", Description: "Great skill"},
	}
	mockCache.EXPECT().ReadIndexes().Return(primitives, nil)

	h := newMarketplaceHandler(t, mockCache, mockProject)
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.Info("skills")(cmd, []string{"code-review"})

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Name:")
	assert.Contains(t, out, "code-review")
	assert.Contains(t, out, "anthropic")
	assert.Contains(t, out, "cap:anthropic:code-review")
	assert.Contains(t, out, "Great skill")
}

// TestMarketplaceHandler_Add_PrintsInstalled verifies the installed output.
func TestMarketplaceHandler_Add_PrintsInstalled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	repoID := "abc12345"
	primitives := []*entity.Primitive{
		{
			Type:         "skills",
			Vendor:       "anthropic",
			Name:         "code-review",
			Version:      "1.2.0",
			SourceRepoID: repoID,
		},
	}
	repoDir := "/fake/cache/repos/" + repoID
	expectedSourceDir := filepath.Join(repoDir, "skills", "anthropic", "code-review", "1.2.0")

	mockCache.EXPECT().ReadIndexes().Return(primitives, nil)
	mockCache.EXPECT().RepoDir(repoID).Return(repoDir)
	mockProject.EXPECT().Install(expectedSourceDir, "skills", "anthropic", "code-review", "1.2.0").Return(nil)

	h := newMarketplaceHandler(t, mockCache, mockProject)
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.Add("skills")(cmd, []string{"code-review"})

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Installed:")
	assert.Contains(t, out, "cap:anthropic:code-review")
}

// TestMarketplaceHandler_Add_Error propagates install errors.
func TestMarketplaceHandler_Add_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	mockCache.EXPECT().ReadIndexes().Return([]*entity.Primitive{}, nil)

	h := newMarketplaceHandler(t, mockCache, mockProject)
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.Add("skills")(cmd, []string{"nonexistent"})

	require.Error(t, err)
}
