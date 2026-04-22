package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/3110Y/cap/internal/domain/entity"
)

func newTestCacheRepo(t *testing.T) (*FileCacheRepository, string) {
	t.Helper()
	dir := t.TempDir()
	repo, err := newFileCacheRepositoryAt(dir)
	require.NoError(t, err)
	return repo.(*FileCacheRepository), dir
}

func TestFileCacheRepository_ListCachedIDs_Empty(t *testing.T) {
	repo, _ := newTestCacheRepo(t)

	ids, err := repo.ListCachedIDs()

	require.NoError(t, err)
	assert.Nil(t, ids)
}

func TestFileCacheRepository_ListCachedIDs_ReturnsSubdirNames(t *testing.T) {
	repo, dir := newTestCacheRepo(t)

	require.NoError(t, os.MkdirAll(filepath.Join(dir, "repos", "abc12345"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "repos", "def67890"), 0755))

	ids, err := repo.ListCachedIDs()

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"abc12345", "def67890"}, ids)
}

func TestFileCacheRepository_ListCachedIDs_IgnoresFiles(t *testing.T) {
	repo, dir := newTestCacheRepo(t)

	require.NoError(t, os.MkdirAll(filepath.Join(dir, "repos", "myrepo"), 0755))
	// Create a file (not a directory) — it must be ignored.
	require.NoError(t, os.WriteFile(filepath.Join(dir, "repos", "notadir.txt"), []byte("x"), 0644))

	ids, err := repo.ListCachedIDs()

	require.NoError(t, err)
	assert.Equal(t, []string{"myrepo"}, ids)
}

func TestFileCacheRepository_RemoveRepo(t *testing.T) {
	repo, dir := newTestCacheRepo(t)

	repoPath := filepath.Join(dir, "repos", "abc12345")
	require.NoError(t, os.MkdirAll(repoPath, 0755))

	err := repo.RemoveRepo("abc12345")

	require.NoError(t, err)
	_, statErr := os.Stat(repoPath)
	assert.True(t, os.IsNotExist(statErr))
}

func TestFileCacheRepository_RepoDir(t *testing.T) {
	repo, dir := newTestCacheRepo(t)

	got := repo.RepoDir("abc12345")

	assert.Equal(t, filepath.Join(dir, "repos", "abc12345"), got)
}

func TestFileCacheRepository_ReadIndex_FileAbsent(t *testing.T) {
	repo, dir := newTestCacheRepo(t)

	// Create repo dir without index.json.
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "repos", "r1"), 0755))

	primitives, err := repo.ReadIndex("r1")

	require.NoError(t, err)
	assert.Nil(t, primitives)
}

func TestFileCacheRepository_ReadIndex_ValidJSON(t *testing.T) {
	repo, dir := newTestCacheRepo(t)

	repoDir := filepath.Join(dir, "repos", "r1")
	require.NoError(t, os.MkdirAll(repoDir, 0755))

	data, err := json.Marshal([]*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.0.0"},
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(repoDir, "index.json"), data, 0644))

	primitives, err := repo.ReadIndex("r1")

	require.NoError(t, err)
	require.Len(t, primitives, 1)
	assert.Equal(t, "code-review", primitives[0].Name)
}

func TestFileCacheRepository_ReadIndex_SetsSourceRepoID(t *testing.T) {
	repo, dir := newTestCacheRepo(t)

	repoDir := filepath.Join(dir, "repos", "myid")
	require.NoError(t, os.MkdirAll(repoDir, 0755))

	data, err := json.Marshal([]*entity.Primitive{
		{Type: "agents", Vendor: "v", Name: "n", Version: "1.0.0"},
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(repoDir, "index.json"), data, 0644))

	primitives, err := repo.ReadIndex("myid")

	require.NoError(t, err)
	require.Len(t, primitives, 1)
	assert.Equal(t, "myid", primitives[0].SourceRepoID)
}

func TestFileCacheRepository_WriteAndReadIndexes(t *testing.T) {
	repo, _ := newTestCacheRepo(t)

	input := []*entity.Primitive{
		{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.2.0", Description: "desc"},
		{Type: "agents", Vendor: "community", Name: "pdf-parser", Version: "0.5.1"},
	}

	err := repo.WriteIndexes(input)
	require.NoError(t, err)

	got, err := repo.ReadIndexes()
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "code-review", got[0].Name)
	assert.Equal(t, "pdf-parser", got[1].Name)
}

func TestFileCacheRepository_ReadIndexes_FileAbsent(t *testing.T) {
	repo, _ := newTestCacheRepo(t)

	got, err := repo.ReadIndexes()

	require.NoError(t, err)
	assert.Nil(t, got)
}
