---
name: cap-cli-commands
description: Справочник CLI-команд CAP — репозитории, кеш, установленные примитивы, маркетплейс, служебные команды.
---

# Rule: cap-cli-commands

В командах `<type>` — одно из: `skills`, `agents`, `commands`, `hooks`, `rules`.

## Управление репозиториями

| Команда                       | Описание                                                           |
|-------------------------------|--------------------------------------------------------------------|
| `cap repository add <url>`    | Добавить удалённый репозиторий. Генерирует случайный `id`.         |
| `cap repository list`         | Список сохранённых репозиториев (пагинация: `-p`, `-t`).           |
| `cap repository del <id>`     | Удалить репозиторий по `id`.                                       |

## Кеш и синхронизация

| Команда        | Описание                                                                                                                                                                                                                                                                |
|----------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `cap update`   | 1. Удаляет из `~/.cache/cap/repos/` клоны репозиториев, отсутствующих в `repositories.yml`. 2. Для каждого репозитория выполняет `git clone` (если нет) или `git fetch --prune` (если есть). 3. Читает `index.json` из каждого клона и объединяет их в `indexes.json`.  |
| `cap upgrade`  | Обновить все установленные примитивы в проекте до последних версий из кеша (по `indexes.json`).                                                                                                                                                                         |

## Управление установленными примитивами

| Команда                         | Описание                                                                                  |
|---------------------------------|-------------------------------------------------------------------------------------------|
| `cap <type> list`               | Список установленных примитивов указанного типа в проекте (пагинация).                    |
| `cap <type> upgrade <name>`     | Обновить конкретный примитив до последней версии.                                         |
| `cap <type> info <name>`        | Описание и список версий установленного примитива.                                        |
| `cap <type> del <name>`         | Удалить примитив из проекта (все версии или конкретную, если указана).                    |

`<name>` может быть простым именем (если оно уникально) или полным идентификатором `vendor/name`.

## Маркетплейс (поиск и установка)

| Команда                                          | Описание                                                                                                                     |
|--------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------|
| `cap <type> marketplace list`                    | Список доступных примитивов из всех репозиториев (по `~/.cache/cap/indexes.json`).                                           |
| `cap <type> marketplace info <name> [<version>]` | Информация о примитиве из маркетплейса (все версии или конкретная).                                                          |
| `cap <type> marketplace search <regex>`          | Поиск по регулярному выражению в поле `description`.                                                                         |
| `cap <type> marketplace add <name> [<version>]`  | Установить примитив в текущий проект (по умолчанию последняя версия). Регистрируется в Claude под `cap:<vendor>:<name>`.     |

## Служебные команды

| Команда            | Описание                                   |
|--------------------|--------------------------------------------|
| `cap help`         | Справка по командам.                       |
| `cap version`      | Версия утилиты CAP.                        |
| `cap ping`         | Проверка связи (отвечает `pong`).          |
| `cap self-update`  | Обновить саму утилиту CAP.                 |

## Форматы вывода

### `cap repository list`

```
ID        URL
a1b2c3d4  https://github.com/user/claude-skills
e5f6g7h8  https://gitlab.com/team/agents
...
Страница 1 из 3
```

### `cap <type> list` / `cap <type> marketplace list`

```
NAME                VERSION
code-review         1.2.0
pdf-parser          0.5.1
...
Страница 1 из 2
```

### `cap <type> info <name>`

```
Name: code-review
Vendor: anthropic
Namespace: cap:anthropic:code-review
Description: Автоматический анализ кода и поиск проблем.

Installed versions:
  - 1.2.0
  - 1.1.0

Available versions (marketplace):
  - 1.2.0
  - 1.1.0
  - 1.0.0
```

### `cap <type> marketplace info <name> <version>`

```
Name: code-review
Version: 1.2.0
Vendor: anthropic
Namespace: cap:anthropic:code-review
Description: Автоматический анализ кода и поиск проблем. В этой версии добавлена поддержка Python 3.12.
```

## Примеры использования

```bash
# Добавить репозиторий
cap repository add https://github.com/example/claude-cap-repo

# Обновить кеш метаданных
cap update

# Посмотреть доступные skills
cap skills marketplace list -t 50

# Найти skills по описанию
cap skills marketplace search "(?i)python|script"

# Установить skill последней версии
cap skills marketplace add code-review

# Установить skill конкретной версии
cap skills marketplace add anthropic/code-review 1.1.0

# Показать информацию об установленном skill
cap skills info code-review

# Обновить все примитивы в проекте
cap upgrade

# Удалить hook
cap hooks del pre-commit-check
```
