---
name: generate-handler
description: Генерирует обработчик CLI-команд в internal/presentation/handler/.
---

# Skill: generate-handler

## Описание

Генерирует обработчик CLI-команд, принимающий сервис прикладного слоя и вызывающий его методы.

## Параметры

- `name` — имя сущности в PascalCase.
- `module_path` — путь модуля Go.

## Создаваемый файл

- `internal/presentation/handler/{name_lower}_handler.go`

## Шаблон

Структура с зависимостью от сервиса и методами вида `HandleGet`, `HandleCreate`, выводящими результаты в консоль.

## Результат

- Обработчик создан и готов к использованию в Cobra-командах.
