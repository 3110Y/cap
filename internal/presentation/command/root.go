package command

import "github.com/spf13/cobra"

// RootCommand оборачивает корневую команду Cobra.
type RootCommand struct {
	*cobra.Command
}

// NewRootCommand собирает корневую команду cap и регистрирует подкоманды.
func NewRootCommand(version *VersionCommand, ping *PingCommand, repository *RepositoryCommand, update *UpdateCommand, upgrade *UpgradeCommand, selfUpdate *SelfUpdateCommand, primitives *PrimitivesCommand) *RootCommand {
	cmd := &cobra.Command{
		Use:           "cap",
		Short:         "CAP — Claude Agentic Primitives CLI",
		Long:          "cap — утилита для управления агентными примитивами Claude: skills, agents, commands, hooks, rules.",
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	cmd.AddCommand(version.Command, ping.Command, repository.Command, update.Command, upgrade.Command, selfUpdate.Command)
	for _, c := range primitives.Commands() {
		cmd.AddCommand(c)
	}
	return &RootCommand{Command: cmd}
}
