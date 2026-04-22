package service

import (
	"fmt"
	"path/filepath"

	"github.com/3110Y/cap/internal/domain/entity"
	"github.com/3110Y/cap/internal/domain/repository"
)

// InstallService installs primitives from the cache into the project.
type InstallService struct {
	cache       repository.CacheRepository
	project     repository.ProjectRepository
	marketplace *MarketplaceService
}

// NewInstallService creates an InstallService.
func NewInstallService(
	cache repository.CacheRepository,
	project repository.ProjectRepository,
	marketplace *MarketplaceService,
) *InstallService {
	return &InstallService{cache: cache, project: project, marketplace: marketplace}
}

// Add installs a primitive into the project's .claude directory.
// name may be "vendor/name" or just "name"; when version is empty, the latest
// SemVer from the marketplace is used (MarketplaceService.Info sorts descending).
func (s *InstallService) Add(primitiveType, name, version string) (*entity.Primitive, error) {
	matches, err := s.marketplace.Info(primitiveType, name, version)
	if err != nil {
		return nil, err
	}
	target := matches[0]

	if target.SourceRepoID == "" {
		return nil, fmt.Errorf("primitive %q has no source repo; run 'cap update' first", name)
	}

	sourceDir := filepath.Join(
		s.cache.RepoDir(target.SourceRepoID),
		target.Type, target.Vendor, target.Name, target.Version,
	)
	if err := s.project.Install(sourceDir, target.Type, target.Vendor, target.Name, target.Version); err != nil {
		return nil, fmt.Errorf("failed to install %s: %w", name, err)
	}
	return target, nil
}
