package command

import (
	"github.com/3110Y/cap/internal/presentation/handler"
	"github.com/spf13/cobra"
)

// UpgradeCommand оборачивает команду cap upgrade.
type UpgradeCommand struct {
	*cobra.Command
}

// NewUpgradeCommand создаёт команду cap upgrade.
func NewUpgradeCommand(h *handler.InstalledHandler) *UpgradeCommand {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Обновить все установленные примитивы до последних версий",
		RunE:  h.UpgradeAll,
	}
	return &UpgradeCommand{Command: cmd}
}
