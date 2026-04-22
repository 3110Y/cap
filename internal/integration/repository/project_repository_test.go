package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindDotClaudeFrom_FindsInParent(t *testing.T) {
	tmp := t.TempDir()

	// Create .claude in tmp root.
	dotClaude := filepath.Join(tmp, ".claude")
	require.NoError(t, os.MkdirAll(dotClaude, 0755))

	// Search from a deeply nested subdirectory.
	subDir := filepath.Join(tmp, "sub", "sub2")
	require.NoError(t, os.MkdirAll(subDir, 0755))

	got, err := findDotClaudeFrom(subDir)

	require.NoError(t, err)
	assert.Equal(t, dotClaude, got)
}

func TestFindDotClaudeFrom_FindsInSameDir(t *testing.T) {
	tmp := t.TempDir()

	dotClaude := filepath.Join(tmp, ".claude")
	require.NoError(t, os.MkdirAll(dotClaude, 0755))

	got, err := findDotClaudeFrom(tmp)

	require.NoError(t, err)
	assert.Equal(t, dotClaude, got)
}

func TestFindDotClaudeFrom_NotFound_ReturnsError(t *testing.T) {
	// Use a fresh temp dir that has no .claude anywhere in the hierarchy.
	// We search from an isolated subtree so we never hit a real .claude in the tree.
	tmp := t.TempDir()

	_, err := findDotClaudeFrom(tmp)

	require.Error(t, err)
	assert.Contains(t, err.Error(), ".claude directory not found")
}

func TestFileProjectRepository_FindDotClaude_ViaChdir(t *testing.T) {
	tmp := t.TempDir()

	dotClaude := filepath.Join(tmp, ".claude")
	require.NoError(t, os.MkdirAll(dotClaude, 0755))

	sub := filepath.Join(tmp, "a", "b")
	require.NoError(t, os.MkdirAll(sub, 0755))

	original, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(original) })

	require.NoError(t, os.Chdir(sub))

	repo := NewFileProjectRepository()
	got, err := repo.FindDotClaude()

	require.NoError(t, err)
	assert.Equal(t, dotClaude, got)
}

func TestFileProjectRepository_Install_CopiesFiles(t *testing.T) {
	tmp := t.TempDir()

	// Set up a .claude directory in tmp and chdir there.
	dotClaude := filepath.Join(tmp, ".claude")
	require.NoError(t, os.MkdirAll(dotClaude, 0755))

	original, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(original) })
	require.NoError(t, os.Chdir(tmp))

	// Create a source directory with a sample file.
	sourceDir := filepath.Join(tmp, "source")
	require.NoError(t, os.MkdirAll(sourceDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "skill.md"), []byte("# Skill"), 0644))

	repo := NewFileProjectRepository()
	err = repo.Install(sourceDir, "skills", "anthropic", "code-review", "1.0.0")
	require.NoError(t, err)

	installed := filepath.Join(dotClaude, "skills", "anthropic", "code-review", "1.0.0", "skill.md")
	_, statErr := os.Stat(installed)
	assert.NoError(t, statErr, "installed file should exist")
}
