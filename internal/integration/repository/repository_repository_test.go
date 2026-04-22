package repository

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/3110Y/cap/internal/domain/entity"
)

func newTestRepoRepository(t *testing.T) *FileRepositoryRepository {
	t.Helper()
	path := filepath.Join(t.TempDir(), "repositories.yml")
	return newFileRepositoryRepositoryAt(path).(*FileRepositoryRepository)
}

func TestFileRepositoryRepository_AddAndList(t *testing.T) {
	repo := newTestRepoRepository(t)

	r1 := &entity.Repository{ID: "aaa", URL: "https://example.com/r1"}
	r2 := &entity.Repository{ID: "bbb", URL: "https://example.com/r2"}

	require.NoError(t, repo.Add(r1))
	require.NoError(t, repo.Add(r2))

	list, err := repo.List()

	require.NoError(t, err)
	require.Len(t, list, 2)

	ids := []string{list[0].ID, list[1].ID}
	assert.ElementsMatch(t, []string{"aaa", "bbb"}, ids)
}

func TestFileRepositoryRepository_Delete_RemovesEntry(t *testing.T) {
	repo := newTestRepoRepository(t)

	require.NoError(t, repo.Add(&entity.Repository{ID: "aaa", URL: "https://example.com/a"}))
	require.NoError(t, repo.Add(&entity.Repository{ID: "bbb", URL: "https://example.com/b"}))

	err := repo.Delete("aaa")
	require.NoError(t, err)

	list, err := repo.List()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "bbb", list[0].ID)
}

func TestFileRepositoryRepository_Delete_NotFound_ReturnsError(t *testing.T) {
	repo := newTestRepoRepository(t)

	require.NoError(t, repo.Add(&entity.Repository{ID: "aaa", URL: "https://example.com/a"}))

	err := repo.Delete("nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestFileRepositoryRepository_List_FileAbsent_ReturnsEmpty(t *testing.T) {
	repo := newTestRepoRepository(t)

	list, err := repo.List()

	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestFileRepositoryRepository_List_PreservesURLs(t *testing.T) {
	repo := newTestRepoRepository(t)

	require.NoError(t, repo.Add(&entity.Repository{ID: "x1", URL: "https://github.com/user/repo"}))

	list, err := repo.List()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "https://github.com/user/repo", list[0].URL)
}
