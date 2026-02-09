# Subscription Service

REST API сервис для агрегации данных об онлайн-подписках пользователей.

## Возможности

- CRUDL-операции над записями о подписках
- Подсчёт суммарной стоимости подписок за выбранный период с фильтрацией
- PostgreSQL с миграциями для инициализации БД
- Swagger-документация
- Запуск через Docker Compose
- Логирование с помощью logrus (JSON-формат)
- Конфигурация через `.env` файл

## Требования

- Docker и Docker Compose

## Запуск

```bash
docker-compose up --build
```

После запуска:
- Сервис доступен по адресу: `http://localhost:8080`
- Swagger-документация: `http://localhost:8080/swagger/index.html`

## API Endpoints

| Метод    | URL                          | Описание                                  |
|----------|------------------------------|-------------------------------------------|
| `POST`   | `/subscriptions`             | Создать подписку                          |
| `GET`    | `/subscriptions`             | Получить список всех подписок             |
| `GET`    | `/subscriptions/{id}`        | Получить подписку по ID                   |
| `PUT`    | `/subscriptions/{id}`        | Обновить подписку                         |
| `DELETE` | `/subscriptions/{id}`        | Удалить подписку                          |
| `GET`    | `/subscriptions/aggregate`   | Подсчитать суммарную стоимость подписок   |

## Пример запроса на создание подписки

```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

## Пример запроса на агрегацию

```bash
curl "http://localhost:8080/subscriptions/aggregate?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&start_date=01-2025&end_date=12-2025"
```

## Конфигурация

Переменные окружения (файл `.env`):

| Переменная     | Описание                    | По умолчанию    |
|----------------|-----------------------------|-----------------|
| `SERVER_PORT`  | Порт сервера                | `8080`          |
| `DB_HOST`      | Хост базы данных            | `localhost`     |
| `DB_PORT`      | Порт базы данных            | `5432`          |
| `DB_USER`      | Пользователь БД             | `postgres`      |
| `DB_PASSWORD`  | Пароль БД                   | `postgres`      |
| `DB_NAME`      | Имя базы данных             | `subscription_db` |
| `DB_SSL_MODE`  | SSL режим                   | `disable`       |
| `LOG_LEVEL`    | Уровень логирования         | `info`          |

## Миграции

Миграции расположены в директории `migrations/` и автоматически применяются при запуске через Docker Compose (через `docker-entrypoint-initdb.d`).

## Структура проекта

```
.
├── cmd/server/               # Точка входа приложения
├── internal/
│   ├── config/               # Загрузка конфигурации
│   ├── database/             # Подключение к БД
│   ├── handlers/             # HTTP-обработчики
│   ├── models/               # Модели данных
│   ├── repositories/         # Слой доступа к данным
│   └── services/             # Бизнес-логика
├── migrations/               # SQL-миграции
├── docs/                     # Swagger-документация
├── .env                      # Переменные окружения
├── Dockerfile                # Сборка Docker-образа
└── docker-compose.yml        # Конфигурация Docker Compose