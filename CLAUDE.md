# CAP — Claude Agentic Primitives

Консольная утилита для управления агентными примитивами экосистемы Claude: **skills**, **agents**, **commands**, **hooks**, **rules**. Позволяет подключать Git-репозитории, синхронизировать метаданные, устанавливать, обновлять и искать примитивы в проектах (директория `.claude`).

## Глоссарий

| Термин               | Значение                                                                                                                                     |
|----------------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| Agentic Primitives   | Расширения для Claude: `skills`, `agents`, `commands`, `hooks`, `rules`.                                                                     |
| Repository           | Удалённый Git-репозиторий с фиксированной структурой каталогов и обязательным `index.json` в корне.                                          |
| Vendor               | Идентификатор владельца или организации (например, `anthropic`, `community`, `my-org`).                                                      |
| Пространство имён    | После установки примитив регистрируется в Claude под именем `cap:<vendor>:<name>` — используется для вызова из диалога.                      |
| Кеш                  | Локальное хранилище: клоны репозиториев в `~/.cache/cap/repos/<repo-id>/` и объединённый индекс `~/.cache/cap/indexes.json`.                 |
| `index.json`         | Файл в корне каждого репозитория с метаданными всех его примитивов.                                                                          |
| `indexes.json`       | Объединённый файл в корне кеша CAP со всеми записями из всех активных `index.json`.                                                          |
| `.claude`            | Директория проекта, куда устанавливаются примитивы (поиск рекурсивный от текущей рабочей директории).                                        |
| Метафайл             | Файл `<primitive>.md` в корне версии примитива с YAML frontmatter.                                                                           |

## Технологический стек

- **Язык:** Go 1.21+
- **Целевые платформы сборки:** macOS, Windows, Ubuntu (кросс-компиляция через `make build-all`).
- **CLI-фреймворк:** Cobra.
- **Внедрение зависимостей:** Google Wire.
- **Git:** `go-git` или внешний `git`.
- **JSON:** стандартный `encoding/json`.
- **YAML:** `gopkg.in/yaml.v3` (для метафайлов).
- **Markdown frontmatter:** `adrg/frontmatter` (опционально, для расширенных метаданных).

## Слоистая архитектура

```
Presentation (CLI commands)
       ↓
Application (Use Cases)
       ↓
Integration (Infrastructure)
       ↓
Domain (Entities & Interfaces)
```

- **Presentation** — команды Cobra, RPC, обработчики (`handler`, `rpc`, `command`): форматирование вывода, разбор аргументов.
- **Application** — сценарии использования (`service`, `validation`, `error`): управление репозиториями, синхронизация кеша, установка/удаление примитивов, поиск.
- **Integration** — адаптеры: Git-клиент, файловая система (кеш, проект, конфиг), реализации репозиториев, объединение `index.json` (`di`, `repository`, ...).
- **Domain** — сущности (`Repository`, `Primitive`, `Version`) и интерфейсы репозиториев для инверсии зависимостей.

Все зависимости внедряются через конструкторы; прикладной слой зависит только от интерфейсов из `domain`, реализации живут в `integration`.

## Структура каталогов проекта

```
cap/
├── cmd/cap/
│   └── main.go
├── internal/
│   ├── presentation/        # command (Cobra), handler, rpc
│   ├── application/         # service, validation, error
│   ├── domain/              # entity, repository (interfaces)
│   └── integration/
│       ├── di/              # Wire (wire.go, wire_gen.go)
│       ├── repository/      # реализации репозиториев
│       └── external/        # Git-клиент, FS-адаптеры
├── pkg/                     # публичные утилиты (опционально)
├── scripts/
├── go.mod
└── go.sum
```

## Рабочие соглашения

Стек и архитектура зафиксированы — см. `.claude/rules/project-conventions.md`:

- Cobra и Google Wire — не заменять альтернативами без явного запроса.
- Новые сущности и CLI-команды добавляются полным набором слоёв через команду `/add-entity`.
- `index.json` в подключаемых репозиториях обязателен; CAP его только читает.

## Документация по темам

- `.claude/rules/cap-storage.md` — пути, `repositories.yml`, структура удалённого репозитория, структура кеша, `.claude` и пространства имён.
- `.claude/rules/cap-index-format.md` — формат `index.json`, `indexes.json`, YAML frontmatter метафайла.
- `.claude/rules/cap-cli-commands.md` — справочник CLI-команд, форматы вывода, примеры использования.
- `.claude/rules/cap-technical-notes.md` — технические решения и ограничения реализации.
- `.claude/rules/project-conventions.md` — рабочие соглашения: фиксированный стек, правила добавления сущностей, работа с `index.json`.
- `.claude/rules/go-style-guide.md` — стиль кода Go.
- `.claude/rules/commit-message.md` — формат сообщений коммитов.
