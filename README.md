# MedodsTechTask

Микросервис аутентификации с JWT-токенами, реализованный на Go с PostgreSQL.

## Технологии
- Go
- PostgreSQL
- Docker
- JWT (SHA512)
- Swagger

## Требования к запуску
- Docker
- Docker Compose

## Быстрый старт
```bash
docker-compose -f docker-compose.yml up -d
```
или
```bash
sh run-local.sh
```

Сервис будет доступен на http://localhost:8080

Swagger документация: http://localhost:8080/swagger/index.html

API Endpoints

* POST /api/v1/user/auth/signup/email - Регистрация аккаунта
* POST /api/v1/user/auth/confirm/email - Подтверждение регистрации
* POST /api/v1/user/login/email - Вход в аккаунт
* POST /api/v1/user/refresh/token - Обновление jwt токена

## Особенности реализации

* Access Token:
    - Формат JWT
    - Алгоритм SHA512
    - Не хранится в БД
* Refresh Token:
    - Произвольный формат
    - Передается только в base64
    - Хранится как bcrypt хеш
    - Защищен от повторного использования и изменений
* Безопастность:
    - Проверка User-Agent при refresh
    - Автоматическая деавторизация при изменении параметров пользователя

## Конфигурация

Настройки по умолчанию заданы в .env файле:
- Для работы достаточно у `.env.example` убрать `.example`

## Структура проекта

```bash
MedodsTechTask
├── docker-compose.yml
├── Dockerfile
├── migrations
│   ├── 0000-init.sql
│   └── 0001-auth-db.sql
├── README.md
├── run-local.sh
└── src
    ├── app
    │   ├── core
    │   │   ├── endpoints.go
    │   │   ├── exceptions.go
    │   │   └── pg.go
    │   ├── main.go
    │   └── user
    │       └── auth
    │           ├── auth_api.go
    │           ├── auth_uc.go
    │           ├── configs
    │           │   └── config.go
    │           ├── repo
    │           │   ├── auth_repo.go
    │           │   └── auth_xdao.go
    │           ├── security.go
    │           └── share
    │               └── auth_dto.go
    ├── docs
    │   ├── docs.go
    │   ├── swagger.json
    │   └── swagger.yaml
    ├── gen_swagger_if_needed.sh
    ├── go.mod
    └── go.sum
```