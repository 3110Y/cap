package repository

import "github.com/3110Y/cap/internal/domain/entity"

// RepositoryRepository defines persistence operations for Repository entities.
type RepositoryRepository interface {
	Add(repo *entity.Repository) error
	List() ([]*entity.Repository, error)
	Delete(id string) error
}
