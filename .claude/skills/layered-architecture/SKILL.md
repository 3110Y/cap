---
name: layered-architecture
description: Шаблоны и принципы слоистой архитектуры (Presentation, Application, Integration, Domain).
---

# Skill: layered-architecture

## Описание

Описывает реализацию слоистой архитектуры для Go CLI-приложений.

## Слои

1. **Presentation** — команды Cobra, форматирование вывода, обработка аргументов.
2. **Application** — сценарии использования (Use Cases, Service): бизнес-логика, валидация, ошибки приложения.
3. **Integration** — адаптеры: Git-клиент, работа с файловой системой (кеш, проект, конфиг), реализации репозиториев, DI (Wire).
4. **Domain** — сущности и интерфейсы репозиториев (для инверсии зависимостей).

## Структура каталогов

```
internal/
├── presentation/   # CLI-команды, handler
├── application/    # service, validation, error
├── integration/    # Git, FS, repository, di
└── domain/         # entity, repository interfaces
```

## Интерфейсы в Domain

```go
type RepositoryStorage interface {
    List() ([]Repository, error)
    Add(url string) (string, error)
    Remove(id string) error
}

type GitClient interface {
    Clone(url, dest string) error
    Fetch(dest string) error
}
```

## Инверсия зависимостей

Все зависимости внедряются через конструкторы. Прикладной слой (`application`) зависит только от интерфейсов из `domain`. Реализации живут в `integration`.

## Преимущества

- Легко тестировать (моки интерфейсов).
- Заменяемые компоненты (например, переход с `go-git` на вызов внешнего `git`).
- Чёткое разделение ответственности.

## Связанные скиллы

- [wire-di](../wire-di/SKILL.md)
- [testing-mocks](../testing-mocks/SKILL.md)
