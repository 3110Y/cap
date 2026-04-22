# CAP — Claude Agentic Primitives

A command-line utility for managing Claude agentic primitives: **skills**, **agents**, **commands**, **hooks**, **rules**. It lets you connect Git repositories of primitives, sync their metadata into a local cache, install/update/remove primitives in a project (the `.claude` directory), and search the marketplace.

Once installed, a primitive is automatically available to Claude in-dialog under the name `cap:<vendor>:<name>`.

## Features

- Connect an arbitrary number of remote Git repositories with primitives.
- Unified metadata cache aggregating all connected repositories (`~/.cache/cap/indexes.json`).
- Marketplace: listing, regex search (over descriptions), detailed primitive info.
- Install a specific version or the latest one (SemVer-sorted).
- Per-item and bulk upgrade of installed primitives.
- Self-updating binary via GitHub Releases.
- Cross-compilation for macOS, Linux, Windows (amd64/arm64).

## Requirements

- Go 1.21+ — to build from source.
- `git` in `PATH` — used to clone and update repositories.

## Installation

### One-liner (Linux / macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/3110Y/cap/main/scripts/install.sh | sh
```

The script detects your OS and architecture, downloads the matching release asset from GitHub and places the `cap` binary into the first writable directory among `/usr/local/bin`, `~/.local/bin`, `~/bin` (falling back to `/usr/local/bin` via `sudo`).

Options and environment variables:

| Flag / env                     | Purpose                                                  |
|--------------------------------|----------------------------------------------------------|
| `--version <tag>` / `CAP_VERSION`     | Install a specific release (default: latest).     |
| `--dir <path>` / `CAP_INSTALL_DIR`    | Target directory (default: auto-detected).        |
| `CAP_REPO=<owner/name>`               | Override the GitHub repository.                   |

Examples:

```bash
# Install a specific version
curl -fsSL https://raw.githubusercontent.com/3110Y/cap/main/scripts/install.sh | sh -s -- --version v0.2.0

# Install into a custom directory
curl -fsSL https://raw.githubusercontent.com/3110Y/cap/main/scripts/install.sh | sh -s -- --dir "$HOME/.local/bin"
```

Windows is not supported by this script — download the binary manually from the [Releases page](https://github.com/3110Y/cap/releases).

### From source

```bash
git clone https://github.com/3110Y/cap.git
cd cap
make init     # downloads dependencies and installs google/wire
make build    # builds ./bin/cap
```

You'll get a `./bin/cap` binary for the current platform. Put it on your `PATH` or call it by path.

### Cross-compilation

```bash
make build-all
```

Artifacts for darwin/linux/windows × amd64/arm64 will be placed in `./bin/`.

### Self-update

```bash
cap self-update
```

Downloads the latest release from GitHub and atomically replaces the current executable.

## Uninstallation

```bash
curl -fsSL https://raw.githubusercontent.com/3110Y/cap/main/scripts/uninstall.sh | sh
```

The script looks for `cap` in the usual install locations (`/usr/local/bin`, `~/.local/bin`, `~/bin`, `/usr/bin`, `/opt/homebrew/bin`) as well as whatever `command -v cap` resolves to, and removes what it finds (escalating with `sudo` when needed).

Options:

| Flag / env                            | Purpose                                                          |
|---------------------------------------|------------------------------------------------------------------|
| `--dir <path>` / `CAP_INSTALL_DIR`    | Remove only from the given directory.                            |
| `--purge`                             | Additionally delete `~/.config/cap` (config) and `~/.cache/cap` (clones + merged index). |

By default the script keeps configuration and cache intact so you can reinstall without losing connected repositories. Use `--purge` for a full wipe:

```bash
curl -fsSL https://raw.githubusercontent.com/3110Y/cap/main/scripts/uninstall.sh | sh -s -- --purge
```

Installed primitives inside a project's `.claude/` directory are **not** touched — remove them per-project with `cap <type> del <name>` or by deleting the directory manually.

## Quick start

```bash
# 1. Connect a primitives repository
cap repository add https://github.com/example/claude-cap-repo

# 2. Sync the cache
cap update

# 3. See what's available
cap skills marketplace list

# 4. Search skills by description (regex, case-insensitive)
cap skills marketplace search "(?i)python|script"

# 5. Install the latest version
cap skills marketplace add code-review

# 6. Or a specific version
cap skills marketplace add anthropic/code-review 1.1.0

# 7. Inspect what is installed
cap skills info code-review

# 8. Upgrade everything
cap upgrade
```

## Command reference

In the commands below, `<type>` is one of: `skills`, `agents`, `commands`, `hooks`, `rules`.

### Repositories

| Command                      | What it does                                        |
|------------------------------|-----------------------------------------------------|
| `cap repository add <url>`   | Add a repository, generate an ID (SHA-1 of URL+timestamp, first 8 characters) |
| `cap repository list`        | List connected repositories (`-p`, `-t` for pagination) |
| `cap repository del <id>`    | Remove a repository by ID                           |

### Cache and synchronization

| Command        | What it does                                                                                     |
|----------------|--------------------------------------------------------------------------------------------------|
| `cap update`   | `git clone`/`git fetch --prune` every repo, drop orphan clones, merge all `index.json` files into `indexes.json` |
| `cap upgrade`  | Upgrade every installed primitive to the latest SemVer available in the cache                    |

### Installed primitives

| Command                        | What it does                                                |
|--------------------------------|-------------------------------------------------------------|
| `cap <type> list`              | List installed primitives of the given type (paginated)     |
| `cap <type> info <name>`       | Description + installed versions + versions in the marketplace |
| `cap <type> upgrade <name>`    | Upgrade a specific primitive to the latest version          |
| `cap <type> del <name>`        | Remove a primitive (all versions, or a specific one via `--version`) |

### Marketplace

| Command                                          | What it does                                        |
|--------------------------------------------------|-----------------------------------------------------|
| `cap <type> marketplace list`                    | All available primitives of the given type, from the cache |
| `cap <type> marketplace search <regex>`          | Search by `description` (Go RE2, case-insensitive)  |
| `cap <type> marketplace info <name> [<version>]` | Detailed information + list of versions             |
| `cap <type> marketplace add <name> [<version>]`  | Install into the project (latest SemVer by default) |

### Utility

| Command            | What it does                             |
|--------------------|------------------------------------------|
| `cap help`         | Help                                     |
| `cap version`      | CAP version                              |
| `cap ping`         | `pong` — checks that the binary is alive |
| `cap self-update`  | Update CAP from GitHub Releases          |

`<name>` is accepted either as `vendor/name` or just `name` (when the name is unique).

## Storage

- Repository configuration: `~/.config/cap/repositories.yml`
- Clone cache: `~/.cache/cap/repos/<repo-id>/`
- Merged index: `~/.cache/cap/indexes.json`
- Installed primitives: `<project>/.claude/<type>/<vendor>/<name>/<version>/`

`.claude` is resolved by walking upward from the current working directory.

## Primitives repository layout

A CAP-compatible Git repository follows a strict directory hierarchy and **must** contain an `index.json` file at its root. The owner maintains `index.json`; CAP only reads it. A repository without `index.json` is considered invalid and is skipped by `cap update` with a warning.

### Directory tree

```
<repository-root>/
├── index.json                         # MANDATORY aggregated index
├── skills/
│   └── <vendor>/
│       └── <name>/
│           └── <version>/             # SemVer, e.g. 1.2.0
│               ├── <name>.md          # metafile with YAML frontmatter
│               ├── README.md          # optional
│               └── ...                # any additional assets
├── agents/
│   └── <vendor>/<name>/<version>/<name>.md
├── commands/
│   └── <vendor>/<name>/<version>/<name>.md
├── hooks/
│   └── <vendor>/<name>/<version>/<name>.md
└── rules/
    └── <vendor>/<name>/<version>/<name>.md
```

Example full path to a primitive version:

```
skills/anthropic/code-review/1.2.0/code-review.md
```

Rules:

- `<type>` is one of: `skills`, `agents`, `commands`, `hooks`, `rules`.
- `<vendor>` identifies the owner (e.g. `anthropic`, `community`, `my-org`).
- `<version>` must be [SemVer](https://semver.org)-compliant.
- The version directory must contain a metafile named exactly `<name>.md`.

### `index.json` format

An array of entries — one per primitive version — aggregated across all types of the repository:

```json
[
  {
    "type": "skills",
    "vendor": "anthropic",
    "name": "code-review",
    "version": "1.2.0",
    "description": "Automatic code analysis and issue detection."
  },
  {
    "type": "agents",
    "vendor": "community",
    "name": "pdf-parser",
    "version": "0.5.1",
    "description": "Extracts text and metadata from PDF files."
  }
]
```

All five fields (`type`, `vendor`, `name`, `version`, `description`) are required.

### Metafile frontmatter

The `<name>.md` file at the root of each version carries YAML frontmatter with extended metadata and free-form Markdown documentation:

```markdown
---
name: code-review
vendor: anthropic
version: 1.2.0
description: Automatic code analysis and issue detection.
dependencies:
  - python>=3.10
author: Anthropic Team
license: MIT
---
# Code Review Skill

This skill performs static analysis...
```

Required frontmatter fields: `name`, `vendor`, `version`, `description`. Optional: `dependencies`, `author`, `license`, `homepage`, `tags`, `requires`, etc.

Full format specifications live in `.claude/rules/cap-*.md`.

## Architecture

The project follows a layered architecture with dependency inversion:

```
Presentation (Cobra)  →  Application (services)  →  Domain (entities, interfaces)
                                                         ↑
                                   Integration (Git, filesystem, HTTP)
```

- `internal/domain/` — entities and repository interfaces.
- `internal/application/service/` — use cases.
- `internal/integration/` — adapters: Git, FS, HTTP, repository implementations, Wire graph (`di/`).
- `internal/presentation/` — Cobra commands and their handlers.
- `cmd/cap/` — entry point.

DI is assembled via [Google Wire](https://github.com/google/wire): `make wire` regenerates `internal/integration/di/wire_gen.go` from the declaration in `wire.go`.

## Development

```bash
make init        # dependencies + install wire
make fmt         # go fmt ./...
make vet         # go vet ./...
make test        # go test -race ./...
make wire        # regenerate wire_gen.go
make build       # local build
make build-all   # cross-compilation
make clean       # remove bin/
```

Commit message format, Go style, and rules for adding new entities — see `.claude/rules/`.

## License

See the [LICENSE](LICENSE) file.
