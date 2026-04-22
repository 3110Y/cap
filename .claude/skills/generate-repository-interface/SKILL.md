---
name: generate-repository-interface
description: Создаёт интерфейс репозитория в internal/domain/repository/.
---

# Skill: generate-repository-interface

## Описание

Создаёт интерфейс репозитория с базовыми CRUD-методами в пакете `internal/domain/repository/`.

## Параметры

- `entity_name` — имя сущности в PascalCase.
- `module_path` — путь модуля Go.

## Создаваемый файл

- `internal/domain/repository/{entity_name_lower}_repository.go`

## Шаблон

```go
package repository

import "{module_path}/internal/domain/entity"

type {EntityName}Repository interface {
    FindByID(id string) (*entity.{EntityName}, error)
    Save({entity_var} *entity.{EntityName}) error
    Delete(id string) error
    List() ([]*entity.{EntityName}, error)
}
```

## Результат

- Интерфейс репозитория создан в доменном слое.
