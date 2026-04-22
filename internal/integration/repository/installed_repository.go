package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/3110Y/cap/internal/domain/entity"
	domain "github.com/3110Y/cap/internal/domain/repository"
)

// FileInstalledRepository scans .claude/<type>/<vendor>/<name>/<version>/ for installed primitives.
type FileInstalledRepository struct {
	project domain.ProjectRepository
}

// NewFileInstalledRepository creates a FileInstalledRepository.
func NewFileInstalledRepository(project domain.ProjectRepository) domain.InstalledRepository {
	return &FileInstalledRepository{project: project}
}

// List returns all installed primitives of the given type.
// If no .claude directory is found, returns (nil, nil) — there's simply nothing installed.
func (r *FileInstalledRepository) List(primitiveType string) ([]*entity.Primitive, error) {
	dotClaude, err := r.project.FindDotClaude()
	if err != nil {
		return nil, nil
	}
	return scanTypeDir(filepath.Join(dotClaude, primitiveType), primitiveType)
}

// Delete removes the installed primitive directory.
// If version is empty, removes all versions (the name directory).
func (r *FileInstalledRepository) Delete(primitiveType, vendor, name, version string) error {
	dotClaude, err := r.project.FindDotClaude()
	if err != nil {
		return err
	}
	var target string
	if version == "" {
		target = filepath.Join(dotClaude, primitiveType, vendor, name)
	} else {
		target = filepath.Join(dotClaude, primitiveType, vendor, name, version)
	}
	return os.RemoveAll(target)
}

// scanTypeDir walks .claude/<type>/<vendor>/<name>/<version>/ and returns entities.
func scanTypeDir(typeDir, primitiveType string) ([]*entity.Primitive, error) {
	vendors, err := os.ReadDir(typeDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", typeDir, err)
	}
	var result []*entity.Primitive
	for _, vendorEntry := range vendors {
		if !vendorEntry.IsDir() {
			continue
		}
		vendorDir := filepath.Join(typeDir, vendorEntry.Name())
		names, err := os.ReadDir(vendorDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", vendorDir, err)
		}
		for _, nameEntry := range names {
			if !nameEntry.IsDir() {
				continue
			}
			nameDir := filepath.Join(vendorDir, nameEntry.Name())
			versions, err := os.ReadDir(nameDir)
			if err != nil {
				return nil, fmt.Errorf("failed to read %s: %w", nameDir, err)
			}
			for _, vEntry := range versions {
				if !vEntry.IsDir() {
					continue
				}
				result = append(result, &entity.Primitive{
					Type:    primitiveType,
					Vendor:  vendorEntry.Name(),
					Name:    nameEntry.Name(),
					Version: vEntry.Name(),
				})
			}
		}
	}
	return result, nil
}
