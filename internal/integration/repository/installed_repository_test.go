package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubProjectRepository is a minimal inline implementation of domain.ProjectRepository
// that returns a fixed .claude path for testing purposes.
type stubProjectRepository struct {
	dotClaudePath string
	err           error
}

func (s *stubProjectRepository) FindDotClaude() (string, error) {
	return s.dotClaudePath, s.err
}

func (s *stubProjectRepository) Install(sourceDir, primitiveType, vendor, name, version string) error {
	return nil
}

// mkInstalledDir creates .claude/<type>/<vendor>/<name>/<version>/ in tmpDir.
func mkInstalledDir(t *testing.T, base, primType, vendor, name, version string) {
	t.Helper()
	path := filepath.Join(base, primType, vendor, name, version)
	require.NoError(t, os.MkdirAll(path, 0755))
}

func TestFileInstalledRepository_List_NoDotClaude(t *testing.T) {
	stub := &stubProjectRepository{err: assert.AnError}
	repo := NewFileInstalledRepository(stub)

	primitives, err := repo.List("skills")

	require.NoError(t, err)
	assert.Nil(t, primitives)
}

func TestFileInstalledRepository_List_CorrectStructure(t *testing.T) {
	tmp := t.TempDir()
	dotClaude := filepath.Join(tmp, ".claude")
	require.NoError(t, os.MkdirAll(dotClaude, 0755))

	mkInstalledDir(t, dotClaude, "skills", "anthropic", "code-review", "1.2.0")
	mkInstalledDir(t, dotClaude, "skills", "anthropic", "code-review", "1.1.0")
	mkInstalledDir(t, dotClaude, "skills", "community", "pdf-skill", "0.5.0")

	stub := &stubProjectRepository{dotClaudePath: dotClaude}
	repo := NewFileInstalledRepository(stub)

	primitives, err := repo.List("skills")

	require.NoError(t, err)
	require.Len(t, primitives, 3)
	for _, p := range primitives {
		assert.Equal(t, "skills", p.Type)
	}
}

func TestFileInstalledRepository_List_EmptyTypeDir(t *testing.T) {
	tmp := t.TempDir()
	dotClaude := filepath.Join(tmp, ".claude")
	require.NoError(t, os.MkdirAll(dotClaude, 0755))
	// No skills dir at all.

	stub := &stubProjectRepository{dotClaudePath: dotClaude}
	repo := NewFileInstalledRepository(stub)

	primitives, err := repo.List("skills")

	require.NoError(t, err)
	assert.Nil(t, primitives)
}

func TestFileInstalledRepository_Delete_AllVersions(t *testing.T) {
	tmp := t.TempDir()
	dotClaude := filepath.Join(tmp, ".claude")

	mkInstalledDir(t, dotClaude, "skills", "anthropic", "code-review", "1.2.0")
	mkInstalledDir(t, dotClaude, "skills", "anthropic", "code-review", "1.1.0")

	stub := &stubProjectRepository{dotClaudePath: dotClaude}
	repo := NewFileInstalledRepository(stub)

	err := repo.Delete("skills", "anthropic", "code-review", "")

	require.NoError(t, err)
	nameDir := filepath.Join(dotClaude, "skills", "anthropic", "code-review")
	_, statErr := os.Stat(nameDir)
	assert.True(t, os.IsNotExist(statErr))
}

func TestFileInstalledRepository_Delete_SpecificVersion(t *testing.T) {
	tmp := t.TempDir()
	dotClaude := filepath.Join(tmp, ".claude")

	mkInstalledDir(t, dotClaude, "skills", "anthropic", "code-review", "1.2.0")
	mkInstalledDir(t, dotClaude, "skills", "anthropic", "code-review", "1.1.0")

	stub := &stubProjectRepository{dotClaudePath: dotClaude}
	repo := NewFileInstalledRepository(stub)

	err := repo.Delete("skills", "anthropic", "code-review", "1.1.0")

	require.NoError(t, err)

	// 1.1.0 must be gone.
	_, statOld := os.Stat(filepath.Join(dotClaude, "skills", "anthropic", "code-review", "1.1.0"))
	assert.True(t, os.IsNotExist(statOld))

	// 1.2.0 must still exist.
	_, statNew := os.Stat(filepath.Join(dotClaude, "skills", "anthropic", "code-review", "1.2.0"))
	assert.NoError(t, statNew)
}
