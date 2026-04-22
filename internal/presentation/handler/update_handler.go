package handler

import (
	"github.com/3110Y/cap/internal/application/service"
	"github.com/spf13/cobra"
)

// UpdateHandler handles the cap update command.
type UpdateHandler struct {
	svc *service.UpdateService
}

// NewUpdateHandler creates an UpdateHandler.
func NewUpdateHandler(svc *service.UpdateService) *UpdateHandler {
	return &UpdateHandler{svc: svc}
}

// Update handles "cap update".
func (h *UpdateHandler) Update(cmd *cobra.Command, _ []string) error {
	result, err := h.svc.Update()
	if err != nil {
		return err
	}

	for _, id := range result.Synced {
		cmd.Printf("Synced:  %s\n", id)
	}
	for _, id := range result.Pruned {
		cmd.Printf("Pruned:  %s\n", id)
	}
	for _, id := range result.Skipped {
		cmd.Printf("Skipped: %s (no valid index.json)\n", id)
	}
	cmd.Printf("\nIndexed %d primitives.\n", result.Total)
	return nil
}
