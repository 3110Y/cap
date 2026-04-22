package repository

import "github.com/3110Y/cap/internal/domain/entity"

// CacheRepository manages the local CAP cache (~/.cache/cap/).
type CacheRepository interface {
	// ListCachedIDs returns IDs of all repos currently in the cache.
	ListCachedIDs() ([]string, error)
	// RemoveRepo deletes the cached clone for the given repo ID.
	RemoveRepo(id string) error
	// RepoDir returns the local path for a repo's cached clone.
	RepoDir(id string) string
	// ReadIndex reads index.json from a cached repo; returns nil if absent.
	ReadIndex(id string) ([]*entity.Primitive, error)
	// WriteIndexes persists the merged indexes.json.
	WriteIndexes(primitives []*entity.Primitive) error
	// ReadIndexes reads the merged indexes.json.
	ReadIndexes() ([]*entity.Primitive, error)
}
