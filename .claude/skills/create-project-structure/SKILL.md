---
name: create-project-structure
description: Создаёт дерево каталогов слоистой архитектуры для нового Go-проекта.
---

# Skill: create-project-structure

## Описание

Создаёт корневую директорию проекта и дерево каталогов согласно слоистой архитектуре (Presentation, Application, Domain, Integration).

## Параметры

- `project_name` — имя корневой директории проекта.

## Алгоритм

1. Создать корневую директорию `./<project_name>`.
2. Внутри создать дерево каталогов (пустые файлы не нужны — только директории):

```
<project_name>/
├── cmd/
│   └── <project_name>/
├── internal/
│   ├── presentation/
│   │   ├── command/
│   │   └── handler/
│   ├── application/
│   │   ├── service/
│   │   ├── validation/
│   │   └── error/
│   ├── domain/
│   │   ├── entity/
│   │   └── repository/
│   └── integration/
│       ├── di/
│       ├── repository/
│       └── external/
├── scripts/
└── pkg/          (опционально)
```

## Результат

- Структура папок создана и готова к наполнению кодом.
