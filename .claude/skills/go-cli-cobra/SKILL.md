---
name: go-cli-cobra
description: Руководство по созданию CLI-приложений на Go с использованием Cobra.
---

# Skill: go-cli-cobra

## Описание

Шаблоны и лучшие практики для построения многоуровневого интерфейса командной строки с использованием Cobra.

## Шаблон команды Cobra

```go
package cli

import (
    "github.com/spf13/cobra"
)

func NewSkillsCommand(service *application.SkillService) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "skills",
        Short: "Управление скиллами",
        Long:  `Установка, обновление, удаление и поиск скиллов.`,
    }

    cmd.AddCommand(
        newSkillsListCommand(service),
        newSkillsInfoCommand(service),
        newSkillsUpgradeCommand(service),
        newSkillsDelCommand(service),
        newSkillsMarketplaceCommand(service),
    )

    return cmd
}
```

## Обработка флагов и аргументов

Используйте `cmd.Flags()` для добавления флагов, например `-p` для пагинации или `-t` для количества элементов на странице.

## Внедрение зависимостей

Сервисы приложения (Use Cases) внедряются через конструктор командного уровня с помощью Wire.

## Вывод результатов

Используйте стандартный `fmt` или библиотеку `tablewriter` для табличного вывода.

## Ссылки

- [Cobra User Guide](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md)
- [layered-architecture](../layered-architecture/SKILL.md)
