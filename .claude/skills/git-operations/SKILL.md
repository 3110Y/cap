---
name: git-operations
description: "Работа с Git через go-git и внешние команды: клонирование, fetch, обновление."
---

# Skill: git-operations

## Описание

Операции с Git для синхронизации удалённых репозиториев. Скилл содержит примеры реализации на базе `go-git` и вызова внешнего `git`.

## Клонирование репозитория

```go
import "github.com/go-git/go-git/v5"

func Clone(url, path string) error {
    _, err := git.PlainClone(path, false, &git.CloneOptions{
        URL:      url,
        Progress: os.Stdout,
    })
    return err
}
```

## Fetch и обновление
Для существующего репозитория:

```go
func Fetch(path string) error {
    r, err := git.PlainOpen(path)
    if err != nil {
        return err
    }
    err = r.Fetch(&git.FetchOptions{
        RemoteName: "origin",
        Force:      true,
    })
    if err != nil && err != git.NoErrAlreadyUpToDate {
        return err
    }
    // После fetch можно выполнить reset до origin/main (если нужно)
    w, err := r.Worktree()
    if err != nil {
        return err
    }
    return w.Reset(&git.ResetOptions{
        Mode: git.HardReset,
    })
}
```

## Использование внешнего Git

Альтернативно можно вызывать системный `git` через `os/exec`. Это предпочтительнее, если `go-git` не поддерживает какие-то опции (например, `--prune`).

```go
func FetchWithExternal(path string) error {
    cmd := exec.Command("git", "-C", path, "fetch", "--prune")
    return cmd.Run()
}
```

## Обработка ошибок
Всегда проверяйте наличие сети и доступность удалённого репозитория. При ошибках выводите понятные сообщения пользователю.

## Ссылки
- [go-git documentation](https://github.com/go-git/go-git)
