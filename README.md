# CRUD-приложение, предоставляющее Web API к данным.

## Work in progress...

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)

## Features

- Получение по id (GET \product)
- Получение списка всех сущностей (GET \products)
- Создание (POST/PUT)
- Обновление сущности по id (POST/PUT)
- Удаление сущности по id (DELETE)

## Using Table

```sql
create table products
(
    id      serial       not null unique,
    title   varchar(255) not null,
    count   integer      not null,
    price   real         not null,
    created timestamp    not null default now(),
    updated timestamp    not null default now()
);
```

```sql
create table logs
(
    id     serial unique,
    entity varchar(255) not null,
    action varchar(255) not null,
    time   timestamp    not null default now()
);
```

## .env

```text
DRIVER=...
HOST=...
PORT=...
USER=...
DBNAME=...
SSLMODE=...
PASSWORD=...
```

## Others

- События insert/update/delete пишутся в таблицу logs
- Реализована отмена транзакций при ошибках во время insert/delete
- При операциях get/update, если id не найден, в лог-файл (logs/log.log) пишется соответствующее сообщение (INFO)
- Фатальные ошибки пишутся в лог-файл (ERROR)
