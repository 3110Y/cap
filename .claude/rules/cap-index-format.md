---
name: cap-index-format
description: Форматы index.json, indexes.json и YAML frontmatter метафайлов примитивов.
---

# Rule: cap-index-format

## Формат `index.json`

Файл в корне репозитория. Содержит массив объектов — по одному на версию примитива.

**Поля:**

- `type` (строка) — `skills`, `agents`, `commands`, `hooks`, `rules`.
- `vendor` (строка) — идентификатор владельца.
- `name` (строка) — имя примитива.
- `version` (строка, SemVer) — версия.
- `description` (строка) — краткое описание из метафайла.

**Пример:**

```json
[
  {
    "type": "skills",
    "vendor": "anthropic",
    "name": "code-review",
    "version": "1.2.0",
    "description": "Автоматический анализ кода и поиск проблем."
  },
  {
    "type": "skills",
    "vendor": "anthropic",
    "name": "code-review",
    "version": "1.1.0",
    "description": "Анализ кода (предыдущая версия)."
  },
  {
    "type": "agents",
    "vendor": "community",
    "name": "pdf-parser",
    "version": "0.5.1",
    "description": "Извлечение текста и метаданных из PDF."
  }
]
```

**Правила:**

- `index.json` обновляется владельцем репозитория при добавлении/изменении примитивов.
- CAP **не** модифицирует этот файл — только читает его при синхронизации.
- Если `index.json` отсутствует, репозиторий считается невалидным и пропускается при синхронизации с предупреждением.

## Формат `indexes.json`

Объединённый файл в корне кеша (`~/.cache/cap/indexes.json`). Структура совпадает с `index.json`: массив тех же объектов, но агрегированный по всем активным репозиториям. Формируется командой `cap update` и используется всеми операциями маркетплейса (`marketplace list/search/info`) как единый источник правды.

## Формат метафайла (`<primitive-name>.md`)

Markdown-файл в корне версии примитива. Содержит YAML frontmatter (отделён тройными дефисами) и произвольный Markdown-контент (документацию). Используется для получения расширенных метаданных, не входящих в `index.json`.

**Пример:**

```markdown
---
name: code-review
vendor: anthropic
version: 1.2.0
description: Автоматический анализ кода и поиск проблем.
dependencies:
  - python>=3.10
author: Anthropic Team
license: MIT
---
# Code Review Skill

This skill performs static analysis...
```

**Обязательные поля:**

- `name` — имя примитива.
- `vendor` — идентификатор владельца.
- `version` — версия (SemVer).
- `description` — краткое описание.

**Опциональные поля:** `dependencies`, `author`, `license`, `homepage`, `tags`, `requires` и другие.
