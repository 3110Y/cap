---
name: dependency-injector
description: Агент для настройки Google Wire и генерации wire_gen.go.
---

# Agent: dependency-injector

## Назначение

Собирает провайдеры всех слоёв в конфигурации Wire и генерирует `wire_gen.go`. Функция `InitializeCLI` должна возвращать исполнитель команд, запускаемый из `main.go`.

## Используемые скиллы

- `generate-wire`

## Процесс

1. Собирает все провайдеры в файле `internal/integration/di/wire.go`: репозитории, сервисы, обработчики, команды.
2. Выполняет команду `wire` для генерации `wire_gen.go`.
3. Проверяет, что `InitializeCLI` возвращает исполнитель команд.

## Пример вызова

Используй агента `dependency-injector` после добавления новой сущности для обновления Wire.
