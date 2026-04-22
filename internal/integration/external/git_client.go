package external

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	domain "github.com/3110Y/cap/internal/domain/repository"
)

// ExternalGitClient runs git as an external process.
type ExternalGitClient struct{}

// NewExternalGitClient creates an ExternalGitClient.
func NewExternalGitClient() domain.GitClient {
	return &ExternalGitClient{}
}

// CloneOrFetch clones the repo into dir, or fetches if the clone already exists.
func (c *ExternalGitClient) CloneOrFetch(url, dir string) error {
	info, err := os.Stat(dir)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("cache path %s exists and is not a directory", dir)
		}
		return c.run("git", "-C", dir, "fetch", "--prune")
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to stat %s: %w", dir, err)
	}
	return c.run("git", "clone", "--depth=1", url, dir)
}

func (c *ExternalGitClient) run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git command %q failed: %w", name+" "+strings.Join(args, " "), err)
	}
	return nil
}
