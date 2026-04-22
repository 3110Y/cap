package repository

// ProjectRepository locates and writes to the project's .claude directory.
type ProjectRepository interface {
	// FindDotClaude searches upward from cwd for the nearest .claude directory.
	FindDotClaude() (string, error)
	// Install copies the primitive version directory from sourceDir into the project.
	Install(sourceDir, primitiveType, vendor, name, version string) error
}
