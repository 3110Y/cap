package service

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/3110Y/cap/internal/domain/entity"
	"github.com/3110Y/cap/internal/domain/repository"
)

// MarketplaceService provides read-only access to the cached primitive index.
type MarketplaceService struct {
	cache repository.CacheRepository
}

// NewMarketplaceService creates a MarketplaceService.
func NewMarketplaceService(cache repository.CacheRepository) *MarketplaceService {
	return &MarketplaceService{cache: cache}
}

// List returns all primitives of the given type from the merged index.
func (s *MarketplaceService) List(primitiveType string) ([]*entity.Primitive, error) {
	all, err := s.cache.ReadIndexes()
	if err != nil {
		return nil, fmt.Errorf("failed to read indexes: %w", err)
	}
	return filterByType(all, primitiveType), nil
}

// Search returns primitives matching pattern (case-insensitive regex on description).
func (s *MarketplaceService) Search(primitiveType, pattern string) ([]*entity.Primitive, error) {
	list, err := s.List(primitiveType)
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex %q: %w", pattern, err)
	}
	var result []*entity.Primitive
	for _, p := range list {
		if re.MatchString(p.Description) {
			result = append(result, p)
		}
	}
	return result, nil
}

// Info returns matching primitives for the given name and optional version.
// Results are sorted by SemVer descending, so result[0] is always the latest.
// name may be "vendor/name" or just "name".
func (s *MarketplaceService) Info(primitiveType, name, version string) ([]*entity.Primitive, error) {
	vendor, primName := ParseName(name)
	list, err := s.List(primitiveType)
	if err != nil {
		return nil, err
	}
	var result []*entity.Primitive
	for _, p := range list {
		if p.Name != primName {
			continue
		}
		if vendor != "" && p.Vendor != vendor {
			continue
		}
		if version != "" && p.Version != version {
			continue
		}
		result = append(result, p)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("primitive %q not found in marketplace", name)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return compareSemVer(result[i].Version, result[j].Version) > 0
	})
	return result, nil
}

// ParseName splits "vendor/name" into its parts; returns ("", name) for bare names.
func ParseName(name string) (vendor, primitive string) {
	vendor, primitive, found := strings.Cut(name, "/")
	if !found {
		return "", name
	}
	return vendor, primitive
}

func filterByType(primitives []*entity.Primitive, t string) []*entity.Primitive {
	var result []*entity.Primitive
	for _, p := range primitives {
		if p.Type == t {
			result = append(result, p)
		}
	}
	return result
}

// compareSemVer returns >0 if a > b, <0 if a < b, 0 if equal.
// Handles an optional "v" prefix and falls back to lexical comparison for
// non-numeric segments (e.g. pre-release tags like "1.0.0-rc1").
func compareSemVer(a, b string) int {
	aParts := strings.Split(strings.TrimPrefix(a, "v"), ".")
	bParts := strings.Split(strings.TrimPrefix(b, "v"), ".")
	n := len(aParts)
	if len(bParts) > n {
		n = len(bParts)
	}
	for i := 0; i < n; i++ {
		var ap, bp string
		if i < len(aParts) {
			ap = aParts[i]
		}
		if i < len(bParts) {
			bp = bParts[i]
		}
		ai, aErr := strconv.Atoi(ap)
		bi, bErr := strconv.Atoi(bp)
		if aErr == nil && bErr == nil {
			if ai != bi {
				return ai - bi
			}
			continue
		}
		if ap != bp {
			return strings.Compare(ap, bp)
		}
	}
	return 0
}
