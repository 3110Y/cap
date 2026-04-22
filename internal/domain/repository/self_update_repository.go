package repository

// SelfUpdateRepository checks and applies updates to the CAP binary.
type SelfUpdateRepository interface {
	// LatestVersion returns the latest released version tag (e.g. "v0.2.0").
	LatestVersion() (string, error)
	// Download fetches the release binary for version and replaces the running executable.
	Download(version string) error
}
