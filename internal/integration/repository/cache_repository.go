package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/3110Y/cap/internal/domain/entity"
	domain "github.com/3110Y/cap/internal/domain/repository"
)

// FileCacheRepository manages the CAP cache at ~/.cache/cap/.
type FileCacheRepository struct {
	cacheDir string
}

// newFileCacheRepositoryAt creates a FileCacheRepository rooted at the given cacheDir.
// Intended for use in tests to avoid touching the real home directory.
func newFileCacheRepositoryAt(cacheDir string) (domain.CacheRepository, error) {
	if err := os.MkdirAll(filepath.Join(cacheDir, "repos"), 0o755); err != nil {
		return nil, fmt.Errorf("failed to create cache dir: %w", err)
	}
	return &FileCacheRepository{cacheDir: cacheDir}, nil
}

// NewFileCacheRepository creates a FileCacheRepository rooted at ~/.cache/cap/.
func NewFileCacheRepository() (domain.CacheRepository, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve home dir: %w", err)
	}
	return newFileCacheRepositoryAt(filepath.Join(home, ".cache", "cap"))
}

// ListCachedIDs returns the names of all subdirectories in ~/.cache/cap/repos/.
func (r *FileCacheRepository) ListCachedIDs() ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(r.cacheDir, "repos"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read repos dir: %w", err)
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			ids = append(ids, e.Name())
		}
	}
	return ids, nil
}

// RemoveRepo deletes the cached clone directory for the given ID.
func (r *FileCacheRepository) RemoveRepo(id string) error {
	if err := os.RemoveAll(r.RepoDir(id)); err != nil {
		return fmt.Errorf("failed to remove cached repo %s: %w", id, err)
	}
	return nil
}

// RepoDir returns the local cache path for the given repo ID.
func (r *FileCacheRepository) RepoDir(id string) string {
	return filepath.Join(r.cacheDir, "repos", id)
}

// ReadIndex parses index.json from the cached clone for the given repo ID.
// Returns nil without error if index.json is absent.
func (r *FileCacheRepository) ReadIndex(id string) ([]*entity.Primitive, error) {
	path := filepath.Join(r.RepoDir(id), "index.json")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read index.json for repo %s: %w", id, err)
	}
	var primitives []*entity.Primitive
	if err := json.Unmarshal(data, &primitives); err != nil {
		return nil, fmt.Errorf("failed to parse index.json for repo %s: %w", id, err)
	}
	for _, p := range primitives {
		p.SourceRepoID = id
	}
	return primitives, nil
}

// WriteIndexes persists the merged list to ~/.cache/cap/indexes.json.
// A nil slice is serialised as an empty JSON array rather than "null".
func (r *FileCacheRepository) WriteIndexes(primitives []*entity.Primitive) error {
	if primitives == nil {
		primitives = []*entity.Primitive{}
	}
	data, err := json.MarshalIndent(primitives, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal indexes: %w", err)
	}
	path := filepath.Join(r.cacheDir, "indexes.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write indexes.json: %w", err)
	}
	return nil
}

// ReadIndexes reads the merged ~/.cache/cap/indexes.json.
func (r *FileCacheRepository) ReadIndexes() ([]*entity.Primitive, error) {
	path := filepath.Join(r.cacheDir, "indexes.json")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read indexes.json: %w", err)
	}
	var primitives []*entity.Primitive
	if err := json.Unmarshal(data, &primitives); err != nil {
		return nil, fmt.Errorf("failed to parse indexes.json: %w", err)
	}
	return primitives, nil
}
