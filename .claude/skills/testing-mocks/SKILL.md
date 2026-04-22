---
name: testing-mocks
description: Написание модульных тестов с использованием интерфейсов и моков.
---

# Skill: testing-mocks

## Описание

Архитектура построена вокруг интерфейсов, что упрощает модульное тестирование. Скилл описывает подходы к созданию тестов с моками.

## Генерация моков

Используйте `mockgen` из [gomock](https://github.com/uber-go/mock) или встроенную генерацию Go 1.18+.

Установка:

```bash
go install go.uber.org/mock/mockgen@latest
```

Генерация мока для интерфейса:

```bash
mockgen -source=internal/domain/repository/user_repository.go \
        -destination=internal/integration/mocks/mock_user_repository.go \
        -package=mocks
```

## Пример теста

```go
func TestUserService_Create(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    mockRepo.EXPECT().Save(gomock.Any()).Return(nil)

    service := service.NewUserService(mockRepo)

    err := service.Create("a1b2c3d4", "Alice")
    assert.NoError(t, err)
}
```

## Тестирование файловой системы

Используйте временную директорию `t.TempDir()` для изоляции тестов:

```go
func TestCacheFS_WriteIndex(t *testing.T) {
    tmpDir := t.TempDir()
    fs := integration.NewCacheFS(tmpDir)
    // ...
}
```

## Рекомендации

- Тестируйте слой Application отдельно от Integration.
- Для Git-операций создайте мок-интерфейс и подменяйте реализацию в тестах.
- Используйте `testify/assert` для удобных проверок.

## Ссылки

- [Testify](https://github.com/stretchr/testify)
- [gomock](https://github.com/uber-go/mock)
