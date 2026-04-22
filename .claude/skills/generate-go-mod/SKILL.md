---
name: generate-go-mod
description: Инициализирует go.mod и добавляет необходимые зависимости.
---

# Skill: generate-go-mod

## Описание

Инициализирует модуль Go и добавляет базовые зависимости, необходимые для CLI-проекта.

## Параметры

- `module_path` — путь модуля Go.

## Действия

```bash
go mod init {module_path}
go get github.com/spf13/cobra@latest github.com/google/wire@latest
```

## Результат

- Файл `go.mod` создан, зависимости добавлены.
