# Использование AWS API для C2

| Upstream Name                                   | Upstream Version                                          |
|-------------------------------------------------|-----------------------------------------------------------|
| [aws-sdk-go](https://github.com/aws/aws-sdk-go) | [1.44.10](https://github.com/aws/aws-sdk-go/tree/v1.44.0) |

- [Использование AWS API для C2](#использование-aws-api-для-c2)
  - [Требования](#требования)
  - [Начало работы](#начало-работы)
    - [Генерация SDK для API](#генерация-sdk-для-api)
    - [Работа с endpoints.json](#работа-с-endpointsjson)
    - [Установка и запуск линтеров](#установка-и-запуск-линтеров)
  - [Уровень control plane в EKS v3](#уровень-control-plane-в-eks-v3)
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

### Работа с endpoints.json

Облачная инсталляция может быть описана в виде партиции с регионами и доступными сервисами в `models/endpoints/endpoints.json`.
На основе описания партиции SDK формирует адреса для API эндпойнтов сервисов.

После внесения изменений в `endpoints.json` требуется перегенерация кода: `make gen-endpoints`.

Результат генерации находится в `aws/endpoints/defaults.go`.

### Установка и запуск линтеров

Установка линтеров: `make install-deps-verify`

Запуск линтеров: `make verify`

## Уровень control plane в EKS v3

В C2 EKS v3 поле `controlPlaneScalingConfig.tier` задаёт лимиты ресурсов
контейнеров control plane. Это расширение API C2. Доступны значения
`tier-large`, `tier-2large`, `tier-4large` и `tier-8large`.
Если поле не передано при создании, сервис выбирает уровень по умолчанию.
У существующих кластеров может возвращаться `standard`; его нельзя выбрать
при создании нового кластера.

При создании добавьте `ControlPlaneScalingConfig` в `eks.CreateClusterInput`:

```go
input.ControlPlaneScalingConfig = &eks.ControlPlaneScalingConfig{
    Tier: aws.String("tier-large"),
}
```

Уровень существующего кластера меняется без пересоздания через
`UpdateClusterConfig`. Дождитесь успешного завершения операции по ID из
`output.Update.Id`, затем прочитайте фактический уровень через
`DescribeCluster`:

```go
output, err := client.UpdateClusterConfig(&eks.UpdateClusterConfigInput{
    Name: aws.String(clusterName),
    ControlPlaneScalingConfig: &eks.ControlPlaneScalingConfig{
        Tier: aws.String("tier-2large"),
    },
})
if err != nil {
    return err
}
updateID := aws.StringValue(output.Update.Id)
// Опросите DescribeUpdate с updateID до статуса Successful.
```

Обновление уровня отправляйте отдельно от изменения security groups: C2 не
принимает их в одном запросе `UpdateClusterConfig`. Для EKS v2 это поле не
поддерживается.

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
