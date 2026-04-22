---
name: generate-repository
description: Генерирует реализацию репозитория в internal/integration/repository/.
---

# Skill: generate-repository

## Описание

Генерирует реализацию репозитория для сущности в пакете `internal/integration/repository/`. Конструктор должен возвращать тип, удовлетворяющий интерфейсу из `internal/domain/repository/`.

## Параметры

- `entity_name` — имя сущности в PascalCase.
- `module_path` — путь модуля Go.
- `table_name` — имя таблицы (по умолчанию `{entity_name_lower}s`).

## Создаваемый файл

- `internal/integration/repository/{entity_name_lower}_repository.go`

## Результат

- Реализация репозитория создана и готова к регистрации в Wire.
