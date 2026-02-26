# Subscription Service

REST-сервис для агрегации данных об онлайн подписках пользователей.

## Требования

- Go 1.21+
- Docker и Docker Compose
- PostgreSQL 15+

## Установка и запуск

### С помощью Docker Compose (рекомендуется)

1. Клонируйте репозиторий
2. Запустите сервис:

```bash
docker-compose up -d
```

Сервис будет доступен по адресу `http://localhost:8080`

### Локальный запуск

1. Установите зависимости:

```bash
go mod download
```

2. Настройте базу данных PostgreSQL

3. Создайте файл `.env` на основе `.env.example`:

```bash
cp .env.example .env
```

4. Запустите миграции и сервер:

```bash
go run main.go
```

## API Документация

После запуска сервиса Swagger документация доступна по адресу:
- http://localhost:8080/swagger/index.html

## Эндпоинты

### Подписки (CRUDL)

- `POST /api/v1/subscriptions` - Создать подписку
- `GET /api/v1/subscriptions` - Список подписок (с пагинацией)
- `GET /api/v1/subscriptions/:id` - Получить подписку по ID
- `PUT /api/v1/subscriptions/:id` - Обновить подписку
- `DELETE /api/v1/subscriptions/:id` - Удалить подписку

### Расчет стоимости

- `GET /api/v1/subscriptions/total-cost` - Рассчитать суммарную стоимость подписок

Параметры запроса:
- `start_date` (опционально) - Начало периода в формате MM-YYYY
- `end_date` (опционально) - Конец периода в формате MM-YYYY
- `user_id` (опционально) - UUID пользователя
- `service_name` (опционально) - Название сервиса

## Примеры запросов

### Создание подписки

```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

### Расчет стоимости

```bash
curl "http://localhost:8080/api/v1/subscriptions/total-cost?start_date=01-2025&end_date=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

## Конфигурация

Конфигурация может быть задана через:
1. Файл `config.yaml`
2. Переменные окружения (`.env` файл)
3. Переменные окружения системы

Приоритет: переменные окружения > config.yaml > значения по умолчанию

## Структура проекта

```
.
├── main.go                 # Точка входа
├── config.yaml             # Конфигурационный файл
├── docker-compose.yml      # Docker Compose конфигурация
├── Dockerfile              # Docker образ
├── go.mod                  # Go модули
├── internal/
│   ├── config/            # Конфигурация
│   ├── database/          # Подключение к БД
│   ├── handlers/          # HTTP обработчики
│   ├── migrations/       # Миграции БД
│   ├── models/           # Модели данных
│   └── router/           # Роутинг
└── README.md
```

## Логирование

Все операции логируются в стандартный вывод (stdout). Логи включают:
- Создание, обновление, удаление подписок
- Ошибки валидации и обработки запросов
- Ошибки подключения к БД

## Технологии

- **Go 1.21** - Язык программирования
- **Gin** - HTTP веб-фреймворк
- **GORM** - ORM для работы с БД
- **PostgreSQL** - База данных
- **Swagger** - Документация API
- **Docker** - Контейнеризация
# test-service
