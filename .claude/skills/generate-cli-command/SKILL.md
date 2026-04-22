---
name: generate-cli-command
description: Генерирует команду Cobra для сущности и регистрирует её в root.go.
---

# Skill: generate-cli-command

## Описание

Генерирует команду Cobra с подкомандами `get` и `create` для сущности и регистрирует её в корневой команде.

## Параметры

- `name` — имя сущности в PascalCase.
- `module_path` — путь модуля Go.
- `project_name` — имя проекта.

## Создаваемые файлы

- `internal/presentation/command/{name_lower}.go` (подкоманды `get` / `create`).
- `internal/presentation/command/root.go` (создаётся, если ещё не существует).

## Результат

- Команда зарегистрирована в `root.go` и доступна через CLI.
