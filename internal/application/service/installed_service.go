package service

import (
	"fmt"
	"sort"

	"github.com/3110Y/cap/internal/domain/entity"
	"github.com/3110Y/cap/internal/domain/repository"
)

var allPrimitiveTypes = []string{"skills", "agents", "commands", "hooks", "rules"}

// marketplaceLister defines the marketplace read operations required by InstalledService.
type marketplaceLister interface {
	List(primitiveType string) ([]*entity.Primitive, error)
	Info(primitiveType, name, version string) ([]*entity.Primitive, error)
}

// primitiveInstaller defines the install operation required by InstalledService.
type primitiveInstaller interface {
	Add(primitiveType, name, version string) (*entity.Primitive, error)
}

// InstalledInfo combines installed and marketplace data for a single primitive.
type InstalledInfo struct {
	Name              string
	Vendor            string
	Description       string
	InstalledVersions []string
	AvailableVersions []string
}

// InstalledService manages primitives installed in the project's .claude directory.
type InstalledService struct {
	installed   repository.InstalledRepository
	marketplace marketplaceLister
	install     primitiveInstaller
}

// NewInstalledService creates an InstalledService.
func NewInstalledService(
	installed repository.InstalledRepository,
	marketplace *MarketplaceService,
	install *InstallService,
) *InstalledService {
	return &InstalledService{installed: installed, marketplace: marketplace, install: install}
}

// NewInstalledServiceWith creates an InstalledService with explicit interface dependencies.
// Prefer this constructor in tests to inject mocks.
func NewInstalledServiceWith(
	installed repository.InstalledRepository,
	marketplace marketplaceLister,
	install primitiveInstaller,
) *InstalledService {
	return &InstalledService{installed: installed, marketplace: marketplace, install: install}
}

// List returns all primitives of the given type installed in the project.
func (s *InstalledService) List(primitiveType string) ([]*entity.Primitive, error) {
	return s.installed.List(primitiveType)
}

// Info returns combined installed + marketplace info for a primitive.
func (s *InstalledService) Info(primitiveType, name string) (*InstalledInfo, error) {
	vendor, primName := ParseName(name)

	installed, err := s.installed.List(primitiveType)
	if err != nil {
		return nil, err
	}

	info := &InstalledInfo{Name: primName}
	for _, p := range installed {
		if p.Name != primName {
			continue
		}
		if vendor != "" && p.Vendor != vendor {
			continue
		}
		info.Vendor = p.Vendor
		info.InstalledVersions = append(info.InstalledVersions, p.Version)
	}
	if len(info.InstalledVersions) == 0 {
		return nil, fmt.Errorf("primitive %q is not installed", name)
	}

	sort.SliceStable(info.InstalledVersions, func(i, j int) bool {
		return compareSemVer(info.InstalledVersions[i], info.InstalledVersions[j]) > 0
	})

	available, _ := s.marketplace.Info(primitiveType, info.Vendor+"/"+info.Name, "")
	seen := map[string]bool{}
	for _, p := range available {
		if !seen[p.Version] {
			seen[p.Version] = true
			info.AvailableVersions = append(info.AvailableVersions, p.Version)
			if info.Description == "" {
				info.Description = p.Description
			}
		}
	}
	return info, nil
}

// Delete removes an installed primitive; version="" removes all versions.
func (s *InstalledService) Delete(primitiveType, name, version string) error {
	vendor, primName := ParseName(name)

	if vendor == "" {
		installed, err := s.installed.List(primitiveType)
		if err != nil {
			return err
		}
		vendors := uniqueVendors(installed, primName)
		switch len(vendors) {
		case 0:
			return fmt.Errorf("primitive %q is not installed", name)
		case 1:
			vendor = vendors[0]
		default:
			return fmt.Errorf("ambiguous: multiple vendors for %q %v — use vendor/name", primName, vendors)
		}
	}
	return s.installed.Delete(primitiveType, vendor, primName, version)
}

// Upgrade installs the latest available version of a primitive.
func (s *InstalledService) Upgrade(primitiveType, name string) (*entity.Primitive, error) {
	return s.install.Add(primitiveType, name, "")
}

// UpgradeAll upgrades every installed primitive across all types.
// Returns the successfully upgraded primitives and a slice of per-primitive errors.
func (s *InstalledService) UpgradeAll() ([]*entity.Primitive, []error) {
	var upgraded []*entity.Primitive
	var errs []error
	for _, t := range allPrimitiveTypes {
		installed, err := s.installed.List(t)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", t, err))
			continue
		}
		seen := map[string]bool{}
		for _, p := range installed {
			key := p.Vendor + "/" + p.Name
			if seen[key] {
				continue
			}
			seen[key] = true
			result, err := s.install.Add(t, key, "")
			if err != nil {
				errs = append(errs, fmt.Errorf("%s/%s: %w", t, key, err))
				continue
			}
			upgraded = append(upgraded, result)
		}
	}
	return upgraded, errs
}

func uniqueVendors(primitives []*entity.Primitive, name string) []string {
	seen := map[string]bool{}
	var vendors []string
	for _, p := range primitives {
		if p.Name == name && !seen[p.Vendor] {
			seen[p.Vendor] = true
			vendors = append(vendors, p.Vendor)
		}
	}
	return vendors
}
