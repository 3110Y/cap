package command

import (
	"fmt"

	"github.com/3110Y/cap/internal/presentation/handler"
	"github.com/spf13/cobra"
)

var primitiveTypes = []string{"skills", "agents", "commands", "hooks", "rules"}

// PrimitivesCommand groups command subtrees for all five primitive types.
type PrimitivesCommand struct {
	cmds []*cobra.Command
}

// Commands returns the individual type command trees to be added to root.
func (p *PrimitivesCommand) Commands() []*cobra.Command {
	return p.cmds
}

// NewPrimitivesCommand builds cap <type> list/info/del/upgrade/marketplace for each type.
func NewPrimitivesCommand(mp *handler.MarketplaceHandler, installed *handler.InstalledHandler) *PrimitivesCommand {
	var cmds []*cobra.Command
	for _, t := range primitiveTypes {
		cmds = append(cmds, buildTypeCommand(t, mp, installed))
	}
	return &PrimitivesCommand{cmds: cmds}
}

func buildTypeCommand(typeName string, mp *handler.MarketplaceHandler, installed *handler.InstalledHandler) *cobra.Command {
	typeCmd := &cobra.Command{
		Use:   typeName,
		Short: fmt.Sprintf("Управление примитивами типа %s", typeName),
	}

	// Installed primitive commands.
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Список установленных примитивов",
		RunE:  installed.List(typeName),
	}
	listCmd.Flags().IntP("page", "p", 1, "Номер страницы")
	listCmd.Flags().IntP("size", "t", 20, "Количество строк на странице")

	infoCmd := &cobra.Command{
		Use:   "info <name>",
		Short: "Информация об установленном примитиве",
		Args:  cobra.ExactArgs(1),
		RunE:  installed.Info(typeName),
	}

	delCmd := &cobra.Command{
		Use:   "del <name>",
		Short: "Удалить установленный примитив (все версии или конкретную)",
		Args:  cobra.ExactArgs(1),
		RunE:  installed.Delete(typeName),
	}
	delCmd.Flags().String("version", "", "Удалить конкретную версию (по умолчанию все)")

	upgradeCmd := &cobra.Command{
		Use:   "upgrade <name>",
		Short: "Обновить примитив до последней версии",
		Args:  cobra.ExactArgs(1),
		RunE:  installed.Upgrade(typeName),
	}

	// Marketplace subcommand group.
	marketplaceCmd := &cobra.Command{
		Use:   "marketplace",
		Short: "Операции с маркетплейсом",
	}

	mpListCmd := &cobra.Command{
		Use:   "list",
		Short: "Список доступных примитивов в маркетплейсе",
		RunE:  mp.List(typeName),
	}
	mpListCmd.Flags().IntP("page", "p", 1, "Номер страницы")
	mpListCmd.Flags().IntP("size", "t", 20, "Количество строк на странице")

	mpSearchCmd := &cobra.Command{
		Use:   "search <regex>",
		Short: "Поиск по описанию (regex, регистронезависимый)",
		Args:  cobra.ExactArgs(1),
		RunE:  mp.Search(typeName),
	}

	mpInfoCmd := &cobra.Command{
		Use:   "info <name> [<version>]",
		Short: "Информация о примитиве из маркетплейса",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  mp.Info(typeName),
	}

	mpAddCmd := &cobra.Command{
		Use:   "add <name> [<version>]",
		Short: "Установить примитив из маркетплейса в проект",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  mp.Add(typeName),
	}

	marketplaceCmd.AddCommand(mpListCmd, mpSearchCmd, mpInfoCmd, mpAddCmd)
	typeCmd.AddCommand(listCmd, infoCmd, delCmd, upgradeCmd, marketplaceCmd)
	return typeCmd
}
