package command

import (
	"github.com/3110Y/cap/internal/presentation/handler"
	"github.com/spf13/cobra"
)

// UpdateCommand оборачивает команду cap update.
type UpdateCommand struct {
	*cobra.Command
}

// NewUpdateCommand создаёт команду cap update.
func NewUpdateCommand(h *handler.UpdateHandler) *UpdateCommand {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Синхронизировать кеш с подключёнными репозиториями",
		RunE:  h.Update,
	}
	return &UpdateCommand{Command: cmd}
}
