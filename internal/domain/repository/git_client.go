package repository

// GitClient abstracts git clone and fetch operations.
type GitClient interface {
	// CloneOrFetch clones the repo into dir if it doesn't exist, otherwise fetches.
	CloneOrFetch(url, dir string) error
}
