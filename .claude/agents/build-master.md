---
name: build-master
description: Агент для создания Makefile и скриптов кросс-платформенной сборки.
---

# Agent: build-master

## Назначение

Создаёт Makefile и скрипты сборки для кросс-компиляции под Windows, macOS и Linux.

## Используемые скиллы

- `generate-makefile`

## Процесс

1. Создаёт Makefile с целями: `init`, `wire`, `build`, `build-all`, `test`, `clean`.
2. Создаёт скрипт `scripts/build-all.sh` для кросс-компиляции.
3. Связывает цели корректными зависимостями.

## Пример вызова

Используй агента `build-master` для генерации сборочной инфраструктуры проекта.
