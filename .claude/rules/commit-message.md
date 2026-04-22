---
name: commit-message
description: Формат сообщений коммитов для проекта.
---

# Rule: commit-message

## Формат
<type>(<scope>): <subject>

[body]

[footer]


## Типы

- `feat` — новая функциональность.
- `fix` — исправление ошибки.
- `docs` — изменения в документации.
- `style` — форматирование, отступы (без изменения логики).
- `refactor` — рефакторинг кода.
- `test` — добавление или изменение тестов.
- `chore` — рутинные задачи (обновление зависимостей, настройка сборки).

## Scope

Область изменений: `cli`, `cache`, `repository`, `marketplace`, `wire` и т.п.

## Примеры

- `feat(cli): add --prune flag to update`
- `fix(cache): handle missing index.json gracefully`
- `docs(readme): update installation instructions`