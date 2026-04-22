package handler

import (
	"fmt"

	"github.com/3110Y/cap/internal/application/service"
	"github.com/spf13/cobra"
)

// RepositoryHandler handles repository CLI commands.
type RepositoryHandler struct {
	svc *service.RepositoryService
}

// NewRepositoryHandler creates a RepositoryHandler.
func NewRepositoryHandler(svc *service.RepositoryService) *RepositoryHandler {
	return &RepositoryHandler{svc: svc}
}

// Add handles "cap repository add <url>".
func (h *RepositoryHandler) Add(cmd *cobra.Command, args []string) error {
	repo, err := h.svc.Add(args[0])
	if err != nil {
		return err
	}
	cmd.Printf("Added: %s  %s\n", repo.ID, repo.URL)
	return nil
}

// List handles "cap repository list".
func (h *RepositoryHandler) List(cmd *cobra.Command, _ []string) error {
	page, _ := cmd.Flags().GetInt("page")
	size, _ := cmd.Flags().GetInt("size")

	repos, err := h.svc.List()
	if err != nil {
		return err
	}
	if len(repos) == 0 {
		cmd.Println("No repositories registered.")
		return nil
	}

	total := len(repos)
	start := (page - 1) * size
	if start >= total {
		return fmt.Errorf("page %d out of range", page)
	}
	end := start + size
	if end > total {
		end = total
	}

	cmd.Printf("%-10s  %s\n", "ID", "URL")
	for _, r := range repos[start:end] {
		cmd.Printf("%-10s  %s\n", r.ID, r.URL)
	}
	pages := (total + size - 1) / size
	cmd.Printf("\nСтраница %d из %d\n", page, pages)
	return nil
}

// Delete handles "cap repository del <id>".
func (h *RepositoryHandler) Delete(cmd *cobra.Command, args []string) error {
	if err := h.svc.Delete(args[0]); err != nil {
		return err
	}
	cmd.Printf("Removed: %s\n", args[0])
	return nil
}
