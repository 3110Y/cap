package command

import "github.com/spf13/cobra"

// Version задаёт текущую версию CAP. Замещается линковщиком при релизной сборке.
var Version = "0.1.0-dev"

// VersionCommand оборачивает команду cap version.
type VersionCommand struct {
	*cobra.Command
}

// NewVersionCommand создаёт подкоманду cap version.
func NewVersionCommand() *VersionCommand {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Показать версию CAP",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("cap %s\n", Version)
		},
	}
	return &VersionCommand{Command: cmd}
}
