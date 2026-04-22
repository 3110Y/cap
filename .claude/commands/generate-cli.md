---
name: generate-cli
description: Создаёт готовый к сборке Go CLI-проект с многослойной архитектурой, примером сущности User, Wire и кросс-компиляцией.
---

# Command: generate-cli

## Описание

Создаёт полностью готовый к сборке проект с многослойной архитектурой, примером сущности `User`, настроенным Google Wire и кросс-компиляцией.

## Синтаксис

```text
/generate-cli project_name=<имя_проекта> [module_path=<путь_модуля>]
```

## Параметры

| Параметр       | Обязательный | По умолчанию                        | Описание                                    |
|----------------|--------------|-------------------------------------|---------------------------------------------|
| `project_name` | Да           | —                                   | Имя корневой папки и исполняемого файла     |
| `module_path`  | Нет          | `github.com/example/{project_name}` | Полный путь модуля Go                       |

## Пример вызова

```text
/generate-cli project_name=taskmanager module_path=github.com/company/taskmanager
```

## Последовательность выполнения

| Шаг | Агент                    | Действие                                                         | Скилл                          |
|-----|--------------------------|------------------------------------------------------------------|--------------------------------|
| 1   | `architect`              | Создание структуры папок                                         | `create-project-structure`     |
| 2   | `architect`              | Инициализация `go.mod` и зависимостей                            | `generate-go-mod`              |
| 3   | `domain-expert`          | Генерация сущности `User` (`ID`, `Name`)                         | `generate-entity`              |
| 4   | `domain-expert`          | Генерация интерфейса `UserRepository`                            | `generate-repository-interface`|
| 5   | `infrastructure-engineer`| Генерация реализации `UserRepository`                            | `generate-repository`          |
| 6   | `application-developer`  | Генерация `UserService` и базовых ошибок                         | `generate-service`             |
| 7   | `cli-integrator`         | Генерация `UserHandler`                                          | `generate-handler`             |
| 8   | `cli-integrator`         | Генерация команд Cobra (`user get`, `user create`) и `root.go`   | `generate-cli-command`         |
| 9   | `dependency-injector`    | Генерация `wire.go` и запуск `wire`                              | `generate-wire`                |
| 10  | `build-master`           | Генерация `Makefile` и `scripts/build-all.sh`                    | `generate-makefile`            |
| 11  | `architect`              | Генерация `main.go` с инициализацией Wire                        | `generate-main`                |
| 12  | —                        | Выполнение `make init` в папке проекта                           | —                              |

## Результат

- Полноценный проект в папке `{project_name}`.
- После завершения пользователь может выполнить `make build-all` и получить бинарники для Windows, macOS и Linux.
