//go:build wireinject
// +build wireinject

package di

import (
	"github.com/3110Y/cap/internal/application/service"
	"github.com/3110Y/cap/internal/integration/external"
	integrationrepo "github.com/3110Y/cap/internal/integration/repository"
	"github.com/3110Y/cap/internal/presentation/command"
	"github.com/3110Y/cap/internal/presentation/handler"
	"github.com/google/wire"
)

// InitializeCLI собирает граф зависимостей и возвращает корневую команду CAP.
func InitializeCLI() (*command.RootCommand, error) {
	wire.Build(
		command.NewRootCommand,
		command.NewVersionCommand,
		command.NewPingCommand,
		// Repository
		integrationrepo.NewFileRepositoryRepository,
		service.NewRepositoryService,
		handler.NewRepositoryHandler,
		command.NewRepositoryCommand,
		// Update
		integrationrepo.NewFileCacheRepository,
		external.NewExternalGitClient,
		service.NewUpdateService,
		handler.NewUpdateHandler,
		command.NewUpdateCommand,
		// Marketplace & Install
		integrationrepo.NewFileProjectRepository,
		service.NewMarketplaceService,
		service.NewInstallService,
		handler.NewMarketplaceHandler,
		// Installed & Upgrade
		integrationrepo.NewFileInstalledRepository,
		service.NewInstalledService,
		handler.NewInstalledHandler,
		command.NewUpgradeCommand,
		// Self-update
		external.NewGithubSelfUpdateClient,
		wire.Value(service.CurrentVersion(command.Version)),
		service.NewSelfUpdateService,
		handler.NewSelfUpdateHandler,
		command.NewSelfUpdateCommand,
		// All primitive type commands
		command.NewPrimitivesCommand,
	)
	return nil, nil
}
