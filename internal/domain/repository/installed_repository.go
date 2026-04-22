package repository

import "github.com/3110Y/cap/internal/domain/entity"

// InstalledRepository scans and mutates installed primitives in the project's .claude directory.
type InstalledRepository interface {
	List(primitiveType string) ([]*entity.Primitive, error)
	Delete(primitiveType, vendor, name, version string) error
}
