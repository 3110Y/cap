package command

import (
	"github.com/3110Y/cap/internal/presentation/handler"
	"github.com/spf13/cobra"
)

// SelfUpdateCommand оборачивает команду cap self-update.
type SelfUpdateCommand struct {
	*cobra.Command
}

// NewSelfUpdateCommand создаёт команду cap self-update.
func NewSelfUpdateCommand(h *handler.SelfUpdateHandler) *SelfUpdateCommand {
	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Обновить утилиту CAP до последней версии",
		RunE:  h.SelfUpdate,
	}
	return &SelfUpdateCommand{Command: cmd}
}
