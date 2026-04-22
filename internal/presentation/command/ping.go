package command

import "github.com/spf13/cobra"

// PingCommand оборачивает команду cap ping.
type PingCommand struct {
	*cobra.Command
}

// NewPingCommand создаёт подкоманду cap ping, отвечающую pong.
func NewPingCommand() *PingCommand {
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Проверка связи (отвечает pong)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("pong")
		},
	}
	return &PingCommand{Command: cmd}
}
