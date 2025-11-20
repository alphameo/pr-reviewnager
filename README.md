# _pr-reviewnager_

Backend для cервиса менеджмента Pull Request'ов ([тестовое задание](https://github.com/avito-tech/tech-internship/tree/main/Tech%20Internships/Backend/Backend-trainee-assignment-autumn-2025))

## Стэк

Написан на языке [Go](https://go.dev/)

В качестве базы данных используется [PostgreSQL](https://www.postgresql.org/)

Для генерации кода взаимодействия с базой данных использован инструмент [sqlc](https://sqlc.dev/)

Для генерации кода по спецификации [openapi.yaml](https://github.com/alphameo/pr-reviewnager/blob/main/openapi.yaml)
использовался инструмент [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).
Сама генерация выполняется через команду `make generate-api` ([Makefile](https://github.com/alphameo/pr-reviewnager/blob/main/Makefile))

Для выполнения миграций используется [migrate](https://github.com/golang-migrate/migrate)

Основные зависимости в [`go.mod`](https://github.com/alphameo/pr-reviewnager/blob/main/go.mod):

- [UUID](https://github.com/google/uuid) в качестве ID для сущностей и БД
- [Драйвер для `PostgreSQL`](https://github.com/jackc/pgx)
- [Web framework](https://github.com/labstack/echo)
- [oapi-codegen](https://github.com/oapi-codegen/)

## Запуск сервиса

Запуск выполняется вместе с [миграциями](https://github.com/alphameo/pr-reviewnager/tree/main/migrations)

```bash
docker-compose up --build
```
