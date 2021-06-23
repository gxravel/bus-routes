# Сервис для отслеживания маршрутов общественного транспорта

## Описание

На данном этапе проект представляет собой простой CRUD для таких сущностей, как 

```
type Bus struct {
	ID     int64  `json:"id,omitempty"`
	City   string `json:"city,omitempty"`
	CityID string `json:"city_id,omitempty"`
	Num    string `json:"num"`
}
```
```
type City struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}
```
```
type Stop struct {
	ID      int64  `json:"id,omitempty"`
	City    string `json:"city,omitempty"`
	CityID  string `json:"city_id,omitempty"`
	Address string `json:"address"`
}
```
```
type Route struct {
	BusID  int64 `json:"bus_id"`
	StopID int64 `json:"stop_id"`
	Step   int8  `json:"step"`
}
```

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
