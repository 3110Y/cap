---
name: generate-entity
description: Генерирует файл сущности в internal/domain/entity/.
---

# Skill: generate-entity

## Описание

Генерирует файл сущности в пакете `internal/domain/entity/` согласно переданным полям.

## Параметры

- `entity_name` — имя сущности в PascalCase (например, `User`).
- `fields` — словарь `{"ИмяПоля": "Тип"}` (например, `{"ID": "string", "Name": "string"}`).
- `module_path` — путь модуля Go.

## Создаваемый файл

- `internal/domain/entity/{entity_name_lower}.go`

## Шаблон

```go
package entity

type {EntityName} struct {
    {range $name, $type := .Fields}
    {$name} {$type}
    {end}
}

func New{EntityName}({...}) *{EntityName} {
    return &{EntityName}{...}
}
```

## Результат

- Файл сущности создан в доменном слое.
