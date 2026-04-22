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

func newUpdateHandler(
	t *testing.T,
	mockRepos *mocks.MockRepositoryRepository,
	mockCache *mocks.MockCacheRepository,
	mockGit *mocks.MockGitClient,
) *handler.UpdateHandler {
	t.Helper()
	svc := service.NewUpdateService(mockRepos, mockCache, mockGit)
	return handler.NewUpdateHandler(svc)
}

func TestUpdateHandler_Update_PrintsSyncedAndPruned(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepos := mocks.NewMockRepositoryRepository(ctrl)
	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockGit := mocks.NewMockGitClient(ctrl)

	mockRepos.EXPECT().List().Return([]*entity.Repository{
		{ID: "r1", URL: "https://example.com/r1"},
	}, nil)
	mockCache.EXPECT().ListCachedIDs().Return([]string{"r1", "stale"}, nil)
	mockCache.EXPECT().RemoveRepo("stale").Return(nil)
	mockCache.EXPECT().RepoDir("r1").Return("/cache/r1")
	mockGit.EXPECT().CloneOrFetch("https://example.com/r1", "/cache/r1").Return(nil)
	mockCache.EXPECT().ReadIndex("r1").Return([]*entity.Primitive{
		{Type: "skills", Vendor: "v", Name: "n", Version: "1.0.0"},
	}, nil)
	mockCache.EXPECT().WriteIndexes(gomock.Any()).Return(nil)

	h := newUpdateHandler(t, mockRepos, mockCache, mockGit)
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.Update(cmd, nil)

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Synced:")
	assert.Contains(t, out, "r1")
	assert.Contains(t, out, "Pruned:")
	assert.Contains(t, out, "stale")
}

func TestUpdateHandler_Update_PrintsSkipped(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepos := mocks.NewMockRepositoryRepository(ctrl)
	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockGit := mocks.NewMockGitClient(ctrl)

	mockRepos.EXPECT().List().Return([]*entity.Repository{
		{ID: "r1", URL: "https://example.com/r1"},
	}, nil)
	mockCache.EXPECT().ListCachedIDs().Return(nil, nil)
	mockCache.EXPECT().RepoDir("r1").Return("/cache/r1")
	mockGit.EXPECT().CloneOrFetch("https://example.com/r1", "/cache/r1").Return(nil)
	// nil index.json causes r1 to be skipped.
	mockCache.EXPECT().ReadIndex("r1").Return(nil, nil)
	mockCache.EXPECT().WriteIndexes(gomock.Any()).Return(nil)

	h := newUpdateHandler(t, mockRepos, mockCache, mockGit)
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.Update(cmd, nil)

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Skipped:")
	assert.Contains(t, out, "r1")
}
