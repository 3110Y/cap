package handler

import (
	"github.com/3110Y/cap/internal/application/service"
	"github.com/spf13/cobra"
)

// SelfUpdateHandler handles the cap self-update command.
type SelfUpdateHandler struct {
	svc *service.SelfUpdateService
}

// NewSelfUpdateHandler creates a SelfUpdateHandler.
func NewSelfUpdateHandler(svc *service.SelfUpdateService) *SelfUpdateHandler {
	return &SelfUpdateHandler{svc: svc}
}

// SelfUpdate handles "cap self-update".
func (h *SelfUpdateHandler) SelfUpdate(cmd *cobra.Command, _ []string) error {
	cmd.Println("Checking for updates...")

	result, err := h.svc.Update()
	if err != nil {
		return err
	}

	if result.AlreadyLatest {
		cmd.Printf("Already up to date (version %s).\n", result.OldVersion)
		return nil
	}

	cmd.Printf("Updated: %s → %s\n", result.OldVersion, result.NewVersion)
	cmd.Println("Restart cap to use the new version.")
	return nil
}
