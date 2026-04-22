---
name: wire-di
description: Использование Google Wire для внедрения зависимостей в Go-проекте.
---

# Skill: wire-di

## Описание

Google Wire — инструмент генерации кода для внедрения зависимостей. Скилл описывает структуру провайдеров и порядок генерации `wire_gen.go`.

## Файл wire.go

Размещается в `internal/integration/di/wire.go`:

```go
//go:build wireinject

package di

import (
    "github.com/google/wire"
    "{module_path}/internal/application/service"
    "{module_path}/internal/integration/repository"
    "{module_path}/internal/presentation/command"
    "{module_path}/internal/presentation/handler"
)

func InitializeCLI() (*command.RootCommand, error) {
    wire.Build(
        // Integration
        repository.NewUserRepository,

        // Application
        service.NewUserService,

        // Presentation
        handler.NewUserHandler,
        command.NewRootCommand,
    )
    return nil, nil
}
```

## Генерация кода

Выполните в папке `internal/integration/di`:

```bash
wire
```

Wire создаст `wire_gen.go` с реализацией конструктора.

## Внедрение в main.go

```go
func main() {
    rootCmd, err := di.InitializeCLI()
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to initialize CLI: %v\n", err)
        os.Exit(1)
    }
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

## Провайдеры

Каждый конструктор из слоя Integration должен возвращать интерфейс из Domain:

```go
type GitClient interface { ... }

func NewGitClient() (domain.GitClient, error) { ... }
```

Это позволяет Wire автоматически связывать зависимости.

## Рекомендации

- Держите wire-файлы в `internal/integration/di`.
- Используйте отдельные wire-наборы для тестов.
- Обновляйте `wire_gen.go` при изменении зависимостей.

## Ссылки

- [Wire Tutorial](https://github.com/google/wire/blob/main/_tutorial/README.md)
