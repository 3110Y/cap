---
name: generate-wire
description: Создаёт internal/integration/di/wire.go и генерирует wire_gen.go.
---

# Skill: generate-wire

## Описание

Создаёт файл конфигурации Google Wire с провайдерами всех сущностей и запускает генерацию `wire_gen.go`.

## Параметры

- `entities` — список имён сущностей (например, `["User"]`).
- `module_path` — путь модуля Go.

## Действия

1. Создать `internal/integration/di/wire.go` с вызовами `wire.Build`.
2. Выполнить `cd internal/integration/di && wire`.

## Результат

- Файлы `wire.go` и `wire_gen.go` созданы и готовы к использованию в `main.go`.
