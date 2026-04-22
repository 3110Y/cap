package command

import (
	"github.com/3110Y/cap/internal/presentation/handler"
	"github.com/spf13/cobra"
)

// RepositoryCommand оборачивает группу команд cap repository.
type RepositoryCommand struct {
	*cobra.Command
}

// NewRepositoryCommand создаёт команду cap repository с подкомандами add, list, del.
func NewRepositoryCommand(h *handler.RepositoryHandler) *RepositoryCommand {
	cmd := &cobra.Command{
		Use:   "repository",
		Short: "Управление подключёнными репозиториями",
	}

	addCmd := &cobra.Command{
		Use:   "add <url>",
		Short: "Добавить удалённый репозиторий",
		Args:  cobra.ExactArgs(1),
		RunE:  h.Add,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Список подключённых репозиториев",
		RunE:  h.List,
	}
	listCmd.Flags().IntP("page", "p", 1, "Номер страницы")
	listCmd.Flags().IntP("size", "t", 20, "Количество строк на странице")

	delCmd := &cobra.Command{
		Use:   "del <id>",
		Short: "Удалить репозиторий по ID",
		Args:  cobra.ExactArgs(1),
		RunE:  h.Delete,
	}

	cmd.AddCommand(addCmd, listCmd, delCmd)
	return &RepositoryCommand{Command: cmd}
}
