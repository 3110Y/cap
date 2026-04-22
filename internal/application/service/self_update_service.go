package service

import (
	"fmt"
	"strings"

	"github.com/3110Y/cap/internal/domain/repository"
)

// CurrentVersion is the running version of CAP, injected at wire time.
type CurrentVersion string

// SelfUpdateResult describes the outcome of a self-update attempt.
type SelfUpdateResult struct {
	AlreadyLatest bool
	OldVersion    string
	NewVersion    string
}

// SelfUpdateService checks for and applies updates to the CAP binary.
type SelfUpdateService struct {
	repo    repository.SelfUpdateRepository
	current CurrentVersion
}

// NewSelfUpdateService creates a SelfUpdateService.
func NewSelfUpdateService(repo repository.SelfUpdateRepository, current CurrentVersion) *SelfUpdateService {
	return &SelfUpdateService{repo: repo, current: current}
}

// Update fetches the latest version and replaces the binary if it's newer than the running one.
func (s *SelfUpdateService) Update() (*SelfUpdateResult, error) {
	latest, err := s.repo.LatestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	current := strings.TrimPrefix(string(s.current), "v")
	latestClean := strings.TrimPrefix(latest, "v")

	if compareSemVer(latestClean, current) <= 0 {
		return &SelfUpdateResult{AlreadyLatest: true, OldVersion: current, NewVersion: latestClean}, nil
	}

	if err := s.repo.Download(latest); err != nil {
		return nil, fmt.Errorf("failed to download update: %w", err)
	}
	return &SelfUpdateResult{OldVersion: current, NewVersion: latestClean}, nil
}
