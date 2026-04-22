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

func TestUpdateService_PrunesStaleCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepos := mocks.NewMockRepositoryRepository(ctrl)
	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockGit := mocks.NewMockGitClient(ctrl)

	// Only "aaa" is registered; "bbb" is stale in cache.
	mockRepos.EXPECT().List().Return([]*entity.Repository{
		{ID: "aaa", URL: "https://example.com/a"},
	}, nil)
	mockCache.EXPECT().ListCachedIDs().Return([]string{"aaa", "bbb"}, nil)
	mockCache.EXPECT().RemoveRepo("bbb").Return(nil)

	mockCache.EXPECT().RepoDir("aaa").Return("/cache/aaa")
	mockGit.EXPECT().CloneOrFetch("https://example.com/a", "/cache/aaa").Return(nil)
	mockCache.EXPECT().ReadIndex("aaa").Return([]*entity.Primitive{
		{Type: "skills", Vendor: "v", Name: "n", Version: "1.0.0"},
	}, nil)
	mockCache.EXPECT().WriteIndexes(gomock.Any()).Return(nil)

	svc := service.NewUpdateService(mockRepos, mockCache, mockGit)
	result, err := svc.Update()

	require.NoError(t, err)
	assert.Equal(t, []string{"bbb"}, result.Pruned)
	assert.Equal(t, []string{"aaa"}, result.Synced)
}

func TestUpdateService_ClonesOrFetchesEachRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepos := mocks.NewMockRepositoryRepository(ctrl)
	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockGit := mocks.NewMockGitClient(ctrl)

	mockRepos.EXPECT().List().Return([]*entity.Repository{
		{ID: "r1", URL: "https://example.com/1"},
		{ID: "r2", URL: "https://example.com/2"},
	}, nil)
	mockCache.EXPECT().ListCachedIDs().Return([]string{}, nil)

	mockCache.EXPECT().RepoDir("r1").Return("/cache/r1")
	mockCache.EXPECT().RepoDir("r2").Return("/cache/r2")
	mockGit.EXPECT().CloneOrFetch("https://example.com/1", "/cache/r1").Return(nil)
	mockGit.EXPECT().CloneOrFetch("https://example.com/2", "/cache/r2").Return(nil)

	p1 := []*entity.Primitive{{Type: "skills", Vendor: "v", Name: "a", Version: "1.0.0"}}
	p2 := []*entity.Primitive{{Type: "agents", Vendor: "v", Name: "b", Version: "2.0.0"}}
	mockCache.EXPECT().ReadIndex("r1").Return(p1, nil)
	mockCache.EXPECT().ReadIndex("r2").Return(p2, nil)
	mockCache.EXPECT().WriteIndexes(gomock.Any()).Return(nil)

	svc := service.NewUpdateService(mockRepos, mockCache, mockGit)
	result, err := svc.Update()

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"r1", "r2"}, result.Synced)
	assert.Equal(t, 2, result.Total)
}

func TestUpdateService_SkipsRepoWithNilIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepos := mocks.NewMockRepositoryRepository(ctrl)
	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockGit := mocks.NewMockGitClient(ctrl)

	mockRepos.EXPECT().List().Return([]*entity.Repository{
		{ID: "r1", URL: "https://example.com/1"},
	}, nil)
	mockCache.EXPECT().ListCachedIDs().Return(nil, nil)
	mockCache.EXPECT().RepoDir("r1").Return("/cache/r1")
	mockGit.EXPECT().CloneOrFetch("https://example.com/1", "/cache/r1").Return(nil)
	// ReadIndex returns nil (missing index.json)
	mockCache.EXPECT().ReadIndex("r1").Return(nil, nil)
	mockCache.EXPECT().WriteIndexes(gomock.Any()).Return(nil)

	svc := service.NewUpdateService(mockRepos, mockCache, mockGit)
	result, err := svc.Update()

	require.NoError(t, err)
	assert.Contains(t, result.Skipped, "r1")
	assert.Equal(t, 0, result.Total)
}

func TestUpdateService_SkipsRepoWithIndexError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepos := mocks.NewMockRepositoryRepository(ctrl)
	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockGit := mocks.NewMockGitClient(ctrl)

	mockRepos.EXPECT().List().Return([]*entity.Repository{
		{ID: "r1", URL: "https://example.com/1"},
	}, nil)
	mockCache.EXPECT().ListCachedIDs().Return(nil, nil)
	mockCache.EXPECT().RepoDir("r1").Return("/cache/r1")
	mockGit.EXPECT().CloneOrFetch("https://example.com/1", "/cache/r1").Return(nil)
	mockCache.EXPECT().ReadIndex("r1").Return(nil, errors.New("parse error"))
	mockCache.EXPECT().WriteIndexes(gomock.Any()).Return(nil)

	svc := service.NewUpdateService(mockRepos, mockCache, mockGit)
	result, err := svc.Update()

	require.NoError(t, err)
	assert.Contains(t, result.Skipped, "r1")
}

func TestUpdateService_WriteIndexesCalledWithMergedPrimitives(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepos := mocks.NewMockRepositoryRepository(ctrl)
	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockGit := mocks.NewMockGitClient(ctrl)

	p1 := []*entity.Primitive{{Type: "skills", Vendor: "v", Name: "a", Version: "1.0.0"}}
	p2 := []*entity.Primitive{{Type: "agents", Vendor: "v", Name: "b", Version: "2.0.0"}}

	mockRepos.EXPECT().List().Return([]*entity.Repository{
		{ID: "r1", URL: "https://a.com"},
		{ID: "r2", URL: "https://b.com"},
	}, nil)
	mockCache.EXPECT().ListCachedIDs().Return(nil, nil)
	mockCache.EXPECT().RepoDir("r1").Return("/cache/r1")
	mockCache.EXPECT().RepoDir("r2").Return("/cache/r2")
	mockGit.EXPECT().CloneOrFetch("https://a.com", "/cache/r1").Return(nil)
	mockGit.EXPECT().CloneOrFetch("https://b.com", "/cache/r2").Return(nil)
	mockCache.EXPECT().ReadIndex("r1").Return(p1, nil)
	mockCache.EXPECT().ReadIndex("r2").Return(p2, nil)

	mockCache.EXPECT().WriteIndexes(gomock.InAnyOrder([]*entity.Primitive{p1[0], p2[0]})).Return(nil)

	svc := service.NewUpdateService(mockRepos, mockCache, mockGit)
	result, err := svc.Update()

	require.NoError(t, err)
	assert.Equal(t, 2, result.Total)
}

func TestUpdateService_ReturnsErrorWhenListReposFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepos := mocks.NewMockRepositoryRepository(ctrl)
	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockGit := mocks.NewMockGitClient(ctrl)

	mockRepos.EXPECT().List().Return(nil, errors.New("yaml broken"))

	svc := service.NewUpdateService(mockRepos, mockCache, mockGit)
	_, err := svc.Update()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list repositories")
}

func TestUpdateService_GitFailureMarksRepoSkipped(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepos := mocks.NewMockRepositoryRepository(ctrl)
	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockGit := mocks.NewMockGitClient(ctrl)

	mockRepos.EXPECT().List().Return([]*entity.Repository{
		{ID: "r1", URL: "https://example.com/1"},
	}, nil)
	mockCache.EXPECT().ListCachedIDs().Return(nil, nil)
	mockCache.EXPECT().RepoDir("r1").Return("/cache/r1")
	mockGit.EXPECT().CloneOrFetch("https://example.com/1", "/cache/r1").Return(errors.New("network timeout"))
	mockCache.EXPECT().WriteIndexes(gomock.Any()).Return(nil)

	svc := service.NewUpdateService(mockRepos, mockCache, mockGit)
	result, err := svc.Update()

	require.NoError(t, err)
	assert.Contains(t, result.Skipped, "r1")
	assert.Empty(t, result.Synced)
}
