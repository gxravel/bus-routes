# Сервис для отслеживания маршрутов общественного транспорта

## Зависимости

* Go 1.16
* MySQL 8+

Подготовка окружения - установка необходимых программ, генерирование кода и конфигурации:

```shell script
$ make prepare.tools
$ make gen.config.local
```

## Сборка

```shell script
$ make build
```

## Запуск

### Локальный

Запуск бинарного файла с предварительной сборкой:

```shell script
$ make run
```

## Проверки (запуск линтеров)

Проверка спецификации swagger:
```shell script
make check.swagger
```

Проверка кода линтерами:
```shell script
make lint
```

Проверка всего:
```shell script
make check
```

## Генерация миграций

```shell script
make migration name=add_some_column
```
