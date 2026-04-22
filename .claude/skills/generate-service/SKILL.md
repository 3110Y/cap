---
name: generate-service
description: Генерирует сервис приложения в internal/application/service/.
---

# Skill: generate-service

## Описание

Генерирует сервис прикладного слоя с зависимостью от интерфейса репозитория и методами бизнес-логики (`Get`, `Create` и т.д.).

## Параметры

- `name` — имя сущности в PascalCase.
- `module_path` — путь модуля Go.

## Создаваемый файл

- `internal/application/service/{name_lower}_service.go`

## Шаблон

Структура сервиса принимает интерфейс репозитория из `internal/domain/repository/` и реализует методы бизнес-логики.

## Результат

- Сервис создан и готов к регистрации в Wire.
