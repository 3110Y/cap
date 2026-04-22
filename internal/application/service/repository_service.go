package service

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/3110Y/cap/internal/domain/entity"
	"github.com/3110Y/cap/internal/domain/repository"
)

// idLength is the number of hex characters in a repository ID (per cap-technical-notes.md).
const idLength = 8

// RepositoryService manages remote repositories.
type RepositoryService struct {
	repo repository.RepositoryRepository
	// now is swappable for deterministic testing.
	now func() time.Time
}

// NewRepositoryService creates a RepositoryService.
func NewRepositoryService(repo repository.RepositoryRepository) *RepositoryService {
	return &RepositoryService{repo: repo, now: time.Now}
}

// Add registers a new repository by URL and returns the created entity.
func (s *RepositoryService) Add(url string) (*entity.Repository, error) {
	r := &entity.Repository{ID: s.generateID(url), URL: url}
	if err := s.repo.Add(r); err != nil {
		return nil, fmt.Errorf("failed to add repository: %w", err)
	}
	return r, nil
}

// List returns all registered repositories.
func (s *RepositoryService) List() ([]*entity.Repository, error) {
	return s.repo.List()
}

// Delete removes a repository by ID.
func (s *RepositoryService) Delete(id string) error {
	return s.repo.Delete(id)
}

// generateID returns the first 8 hex characters of SHA-1(url + timestamp).
// SHA-1 is used as a fast, non-cryptographic identifier — not for security.
func (s *RepositoryService) generateID(url string) string {
	salt := strconv.FormatInt(s.now().UnixNano(), 10)
	sum := sha1.Sum([]byte(url + salt))
	return hex.EncodeToString(sum[:])[:idLength]
}
