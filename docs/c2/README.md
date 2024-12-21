# Использование AWS API для K2

| Upstream Name                                   | Upstream Version                                          |
|-------------------------------------------------|-----------------------------------------------------------|
| [aws-sdk-go](https://github.com/aws/aws-sdk-go) | [1.44.10](https://github.com/aws/aws-sdk-go/tree/v1.44.0) |

- [Использование AWS API для K2](#использование-aws-api-для-c2)
  - [Требования](#требования)
  - [Начало работы](#начало-работы)
    - [Генерация SDK для API](#генерация-sdk-для-api)
    - [Установка и запуск линтеров](#установка-и-запуск-линтеров)
  - [Тесты](#тесты)
  - [Выпуск релиза](#выпуск-релиза)
    - [Версионирование](#версионирование)
    - [Релиз](#релиз)

## Требования

- [Go](https://golang.org/doc/install) 1.21
- [Docker](https://docs.docker.com/get-docker/) (запуск sandbox тестов)

## Начало работы

Клонирование репозитория и проверка сборки модуля:

```
$ git clone git@github.com:C2Devel/aws-sdk-go.git && cd aws-sdk-go
...
$ make build
```

Артефакт в результате проверки не создается.

### Генерация SDK для API

AWS API представлено в виде `.json` файлов в сабдиректориях `models/apis/*`.

Генерация `.go` файлов на основе API: `make generate`.

Результат генерации находится в `service/`.

### Установка и запуск линтеров

Установка линтеров: `make install-deps-verify`

Запуск линтеров: `make verify`

## Тесты

Тесты написаны с помощью [go testing](https://go.dev/doc/code#Testing).

Типы:

- unit: `make unit`
- unit + многопоточность: `make ci-test`
- интеграционные: `make integration`
    - **Важно!** Требуется доработка тестов для запуска на С2.
- unit + разные версии go: `make sandbox-tests`
    - Тесты запускаются в контейнерах.
    - **Важно!** Без ошибок проходит таргет `make sandbox-test-go1.21`

Запуск отдельного теста: `go test -v ./private/protocol/restjson/ -run TestUnmarshalError_SerializationError`

## Выпуск релиза

Модуль **aws-sdk-go** используется в качестве зависимости для разных проектов,
например, [terraform-provider-rockitcloud](https://github.com/C2Devel/terraform-provider-rockitcloud).
Версия модуля указывается в `go.mod` и соответствует тегу релиза.

**Важно!** Go не позволяет изменять версии после их публикации, т.е. не допускается перестановка релизного тега.

### Версионирование

Схема версионирования: <upstream_version>-ROCKITX, где upstream_version - версия форка [aws-sdk-go](https://github.com/aws/aws-sdk-go), X - инкрементируется каждый релиз.

Версия фиксируется только в виде тэга.

### Релиз

**Важно!** Релизы провайдера выпускаются с ветки **develop** (установлена дефолтной).
Ветка **main** используется для получения обновлений с upstream.

1. Создание релизного PR'а в ветку **develop**
   (пример: [Release v1.44.10-ROCKIT9](https://github.com/C2Devel/aws-sdk-go/pull/30))
    - Обновление [CHANGELOG.md](../../CHANGELOG.md)
2. Локальный запуск линтеров и unit тестов на ветке **develop**

    ```
    $ make install-deps-verify
    $ make all
    ```

3. Мердж релизного PR'а
4. Установка релизного тега с версией и его публикация

    ```
    $ git tag v1.2.3
    $ git push <remote> v1.2.3
    ```

5. Создание релиза на github. В описание дублируется запись из CHANGELOG.md
