---
name: generate-main
description: Создаёт точку входа cmd/{project_name}/main.go.
---

# Skill: generate-main

## Описание

Создаёт файл точки входа, вызывающий `wire.InitializeCLI` и запускающий исполнитель команд.

## Параметры

- `project_name` — имя проекта.
- `module_path` — путь модуля Go.

## Создаваемый файл

- `cmd/{project_name}/main.go`

## Шаблон

Вызывает `InitializeCLI` из пакета Wire и запускает возвращённый executor.

## Результат

- Файл `main.go` создан и готов к сборке.
