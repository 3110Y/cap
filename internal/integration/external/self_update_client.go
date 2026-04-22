package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	domain "github.com/3110Y/cap/internal/domain/repository"
)

const (
	githubReleasesAPI = "https://api.github.com/repos/3110Y/cap/releases/latest"
	githubDownloadURL = "https://github.com/3110Y/cap/releases/download/%s/%s"
	httpTimeout       = 30 * time.Second
)

// GithubSelfUpdateClient fetches and applies CAP releases from GitHub.
type GithubSelfUpdateClient struct {
	httpClient *http.Client
}

// NewGithubSelfUpdateClient creates a GithubSelfUpdateClient.
func NewGithubSelfUpdateClient() domain.SelfUpdateRepository {
	return &GithubSelfUpdateClient{
		httpClient: &http.Client{Timeout: httpTimeout},
	}
}

// LatestVersion queries the GitHub releases API for the latest tag.
func (c *GithubSelfUpdateClient) LatestVersion() (string, error) {
	resp, err := c.httpClient.Get(githubReleasesAPI)
	if err != nil {
		return "", fmt.Errorf("failed to reach GitHub API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned HTTP %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release info: %w", err)
	}
	if release.TagName == "" {
		return "", errors.New("no releases found")
	}
	return release.TagName, nil
}

// Download fetches the binary for version and atomically replaces the running executable.
func (c *GithubSelfUpdateClient) Download(version string) (err error) {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to locate executable: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Asset names match the format produced by scripts/build-all.sh: cap_<os>_<arch>[.exe].
	assetName := fmt.Sprintf("cap_%s_%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		assetName += ".exe"
	}

	url := fmt.Sprintf(githubDownloadURL, version, assetName)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d for %s", resp.StatusCode, url)
	}

	tmp, err := os.CreateTemp(filepath.Dir(execPath), "cap-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmp.Name()
	// Ensure the temp file is always cleaned up; ignore "not exist" after a successful rename.
	defer func() {
		if removeErr := os.Remove(tmpName); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) && err == nil {
			err = fmt.Errorf("failed to remove temp file: %w", removeErr)
		}
	}()

	if _, copyErr := io.Copy(tmp, resp.Body); copyErr != nil {
		_ = tmp.Close()
		return fmt.Errorf("failed to write update: %w", copyErr)
	}
	if closeErr := tmp.Close(); closeErr != nil {
		return fmt.Errorf("failed to close temp file: %w", closeErr)
	}

	if err := os.Chmod(tmpName, 0o755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}
	if err := os.Rename(tmpName, execPath); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}
	return nil
}
