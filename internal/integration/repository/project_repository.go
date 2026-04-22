package repository

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	domain "github.com/3110Y/cap/internal/domain/repository"
)

// FileProjectRepository locates and writes to the nearest .claude directory.
type FileProjectRepository struct{}

// NewFileProjectRepository creates a FileProjectRepository.
func NewFileProjectRepository() domain.ProjectRepository {
	return &FileProjectRepository{}
}

// FindDotClaude walks upward from cwd until it finds a .claude directory.
func (r *FileProjectRepository) FindDotClaude() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	return findDotClaudeFrom(dir)
}

// findDotClaudeFrom walks upward from startDir until it finds a .claude directory.
// Exposed as a package-level helper so tests can call it directly without touching os.Getwd.
func findDotClaudeFrom(startDir string) (string, error) {
	dir := startDir
	for {
		candidate := filepath.Join(dir, ".claude")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("failed to stat %s: %w", candidate, err)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf(".claude directory not found from %s upward", startDir)
		}
		dir = parent
	}
}

// Install copies sourceDir into .claude/<type>/<vendor>/<name>/<version>/.
func (r *FileProjectRepository) Install(sourceDir, primitiveType, vendor, name, version string) error {
	dotClaude, err := r.FindDotClaude()
	if err != nil {
		return err
	}
	target := filepath.Join(dotClaude, primitiveType, vendor, name, version)
	if err := os.MkdirAll(target, 0o755); err != nil {
		return fmt.Errorf("failed to create install dir: %w", err)
	}
	return copyDir(sourceDir, target)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to compute relative path: %w", err)
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() {
		if closeErr := in.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close source file: %w", closeErr)
		}
	}()

	info, err := in.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close destination file: %w", closeErr)
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}
	return nil
}
