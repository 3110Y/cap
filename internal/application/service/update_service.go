package service

import (
	"fmt"

	"github.com/3110Y/cap/internal/domain/entity"
	"github.com/3110Y/cap/internal/domain/repository"
)

// UpdateResult summarises the outcome of a cache sync.
type UpdateResult struct {
	Synced  []string
	Pruned  []string
	Skipped []string
	Total   int
}

// UpdateService synchronises the local cache with registered repositories.
type UpdateService struct {
	repos repository.RepositoryRepository
	cache repository.CacheRepository
	git   repository.GitClient
}

// NewUpdateService creates an UpdateService.
func NewUpdateService(
	repos repository.RepositoryRepository,
	cache repository.CacheRepository,
	git repository.GitClient,
) *UpdateService {
	return &UpdateService{repos: repos, cache: cache, git: git}
}

// Update performs a full cache sync and returns the result.
func (s *UpdateService) Update() (*UpdateResult, error) {
	registered, err := s.repos.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	result := &UpdateResult{}

	// Step 1: prune stale cached clones.
	activeIDs := make(map[string]struct{}, len(registered))
	for _, r := range registered {
		activeIDs[r.ID] = struct{}{}
	}
	cachedIDs, err := s.cache.ListCachedIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to list cached repos: %w", err)
	}
	for _, id := range cachedIDs {
		if _, ok := activeIDs[id]; !ok {
			if err := s.cache.RemoveRepo(id); err != nil {
				return nil, err
			}
			result.Pruned = append(result.Pruned, id)
		}
	}

	// Step 2: clone or fetch each registered repository.
	var all []*entity.Primitive
	for _, r := range registered {
		dir := s.cache.RepoDir(r.ID)
		if err := s.git.CloneOrFetch(r.URL, dir); err != nil {
			result.Skipped = append(result.Skipped, r.ID)
			continue
		}

		// Step 3: read index.json; skip repos without one.
		primitives, err := s.cache.ReadIndex(r.ID)
		if err != nil || primitives == nil {
			result.Skipped = append(result.Skipped, r.ID)
			continue
		}
		result.Synced = append(result.Synced, r.ID)
		all = append(all, primitives...)
	}

	// Step 4: write merged indexes.json.
	if err := s.cache.WriteIndexes(all); err != nil {
		return nil, fmt.Errorf("failed to write indexes.json: %w", err)
	}
	result.Total = len(all)
	return result, nil
}
