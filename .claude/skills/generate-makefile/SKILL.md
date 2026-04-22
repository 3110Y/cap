---
name: generate-makefile
description: Создаёт Makefile с целями init, wire, build, build-all, test, clean.
---

# Skill: generate-makefile

## Описание

Создаёт Makefile со стандартными правилами и целями для сборки, тестирования и кросс-компиляции под Windows, macOS и Linux.

## Параметры

- `project_name` — имя проекта.

## Создаваемый файл

- `Makefile`

## Цели

- `init` — установка зависимостей.
- `wire` — генерация `wire_gen.go`.
- `build` — сборка под текущую платформу.
- `build-all` — кросс-компиляция под Windows, macOS и Linux.
- `test` — запуск тестов.
- `clean` — очистка артефактов сборки.

## Результат

- Makefile готов к использованию разработчиком и CI.
