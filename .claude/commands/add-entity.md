---
name: add-entity
description: Добавляет полный набор слоёв (entity, repository, service, handler, CLI) для новой сущности и обновляет Wire.
---

# Command: add-entity

## Описание

Добавляет полный набор слоёв (entity, repository interface, repository impl, service, handler, CLI command) для новой сущности и обновляет конфигурацию Wire.

## Синтаксис

```text
/add-entity entity_name=<ИмяСущности> fields=<поля> [module_path=<путь>]
```

## Параметры

| Параметр      | Обязательный | Описание                                                                              |
|---------------|--------------|---------------------------------------------------------------------------------------|
| `entity_name` | Да           | Имя сущности в PascalCase (например, `Task`)                                          |
| `fields`      | Да           | Список полей в формате `"Поле:Тип,Поле:Тип"` (например, `"ID:string,Title:string"`)   |
| `module_path` | Нет          | Если не указан, извлекается из существующего `go.mod`                                 |

## Пример вызова

```text
/add-entity entity_name=Task fields="ID:string,Title:string,Completed:bool"
```

## Последовательность выполнения

| Шаг | Агент                    | Действие                                                                               |
|-----|--------------------------|----------------------------------------------------------------------------------------|
| 1   | `domain-expert`          | Генерация `internal/domain/entity/{entity_name_lower}.go`                              |
| 2   | `domain-expert`          | Генерация `internal/domain/repository/{entity_name_lower}_repository.go`               |
| 3   | `infrastructure-engineer`| Генерация `internal/integration/repository/{entity_name_lower}_repository.go`          |
| 4   | `application-developer`  | Генерация `internal/application/service/{entity_name_lower}_service.go`                |
| 5   | `cli-integrator`         | Генерация `internal/presentation/handler/{entity_name_lower}_handler.go`               |
| 6   | `cli-integrator`         | Генерация `internal/presentation/command/{entity_name_lower}.go` и обновление `root.go`|
| 7   | `dependency-injector`    | Обновление `wire.go` (добавление провайдеров новой сущности) и запуск `wire`           |

## Результат

- Новая сущность полностью интегрирована в проект.
- Пользователь может сразу использовать команды вида `{project_name} {entity_name_lower} get <id>`.
