package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/3110Y/cap/internal/domain/entity"
	domain "github.com/3110Y/cap/internal/domain/repository"
	"gopkg.in/yaml.v3"
)

type yamlRepository struct {
	ID  string `yaml:"id"`
	URL string `yaml:"url"`
}

// FileRepositoryRepository stores repositories in ~/.config/cap/repositories.yml.
type FileRepositoryRepository struct {
	filePath string
}

// newFileRepositoryRepositoryAt creates a FileRepositoryRepository backed by the given YAML file path.
// Intended for use in tests to avoid touching the real home directory.
func newFileRepositoryRepositoryAt(filePath string) domain.RepositoryRepository {
	return &FileRepositoryRepository{filePath: filePath}
}

// NewFileRepositoryRepository creates a repository backed by the YAML config file.
func NewFileRepositoryRepository() (domain.RepositoryRepository, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve home dir: %w", err)
	}
	dir := filepath.Join(home, ".config", "cap")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create config dir: %w", err)
	}
	return &FileRepositoryRepository{filePath: filepath.Join(dir, "repositories.yml")}, nil
}

func (r *FileRepositoryRepository) read() ([]*yamlRepository, error) {
	data, err := os.ReadFile(r.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read repositories file: %w", err)
	}
	var repos []*yamlRepository
	if err := yaml.Unmarshal(data, &repos); err != nil {
		return nil, fmt.Errorf("failed to parse repositories file: %w", err)
	}
	return repos, nil
}

func (r *FileRepositoryRepository) write(repos []*yamlRepository) error {
	data, err := yaml.Marshal(repos)
	if err != nil {
		return fmt.Errorf("failed to marshal repositories: %w", err)
	}
	if err := os.WriteFile(r.filePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write repositories file: %w", err)
	}
	return nil
}

// Add appends a new repository entry.
func (r *FileRepositoryRepository) Add(repo *entity.Repository) error {
	repos, err := r.read()
	if err != nil {
		return err
	}
	repos = append(repos, &yamlRepository{ID: repo.ID, URL: repo.URL})
	return r.write(repos)
}

// List returns all registered repositories.
func (r *FileRepositoryRepository) List() ([]*entity.Repository, error) {
	yamlRepos, err := r.read()
	if err != nil {
		return nil, err
	}
	result := make([]*entity.Repository, len(yamlRepos))
	for i, yr := range yamlRepos {
		result[i] = &entity.Repository{ID: yr.ID, URL: yr.URL}
	}
	return result, nil
}

// Delete removes a repository by ID.
func (r *FileRepositoryRepository) Delete(id string) error {
	repos, err := r.read()
	if err != nil {
		return err
	}
	filtered := repos[:0]
	found := false
	for _, yr := range repos {
		if yr.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, yr)
	}
	if !found {
		return fmt.Errorf("repository %q not found", id)
	}
	return r.write(filtered)
}
