# QW Trading Platform

Виртуальная биржа криптовалют с высокопроизводительным торговым ядром, построенная на микросервисной архитектуре.

## Возможности

- **Торговое ядро** — C++ matching engine с поддержкой Limit/Market ордеров и Price-Time Priority
- **Микросервисы** — независимые Go сервисы для каждой бизнес-области
- **Высокая производительность** — обработка 1000+ ордеров/сек
- **Безопасность** — JWT аутентификация, изоляция данных пользователей
- **API** — REST API с полной документацией

## Архитектура

```
┌─────────────┐     ┌─────────────────┐
│   Client    │────▶│   API Gateway   │
└─────────────┘     └────────┬────────┘
                             │
            ┌────────────────┼────────────────┐
            │                │                │
    ┌───────▼──────┐ ┌──────▼───────┐ ┌──────▼───────┐
    │ User Service │ │Order Service │ │  Portfolio   │
    │   (Go)       │ │   (Go)       │ │  Service     │
    └───────┬──────┘ └──────┬───────┘ └──────┬───────┘
            │                │                │
            │         ┌──────▼───────┐        │
            │         │  Matching    │        │
            │         │  Engine      │        │
            │         │   (C++)      │        │
            │         └──────┬───────┘        │
            │                │                │
            └────────────────┼────────────────┘
                             │
                    ┌────────▼────────┐
                    │   PostgreSQL    │
                    └─────────────────┘
```

## Сервисы

| Сервис | Язык | Порт | Описание |
|--------|------|------|----------|
| API Gateway | Go | 8080 | Маршрутизация, аутентификация |
| User Service | Go | 8081 | Регистрация, профиль |
| Account Service | Go | 8082 | Счета, баланс |
| Order Service | Go | 8083 | Ордера |
| Portfolio Service | Go | 8084 | Позиции, портфель |
| Matching Engine | C++ | 50051 | Торговое ядро |

## Быстрый старт

### Docker (рекомендуется)

```bash
git clone https://github.com/qw-trading/platform.git
cd platform
docker-compose up --build
```

### Локальная разработка

Требования: Go 1.22+, PostgreSQL 16+

```bash
# Запустить PostgreSQL
docker run -d --name postgres -p 5432:5432 \
  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=qw_trading postgres:16-alpine

# Применить миграции
psql -h localhost -U postgres -d qw_trading -f deploy/migrations/001_init.sql

# Запустить сервисы
export DB_HOST=localhost DB_USER=postgres DB_PASSWORD=postgres DB_NAME=qw_trading JWT_SECRET=dev

go run cmd/api-gateway/main.go &
go run cmd/user-service/main.go &
go run cmd/account-service/main.go &
go run cmd/order-service/main.go &
go run cmd/portfolio-service/main.go &
```

## API

### Регистрация

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"trader@example.com","username":"trader","password":"secret123"}'
```

### Вход

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"trader@example.com","password":"secret123"}'
# Response: {"access_token":"eyJ...","refresh_token":"eyJ..."}
```

### Создание ордера

```bash
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTC/USDT",
    "side": "BUY",
    "type": "LIMIT",
    "price": 50000,
    "quantity": 0.1
  }'
```

### Просмотр портфеля

```bash
curl http://localhost:8080/portfolio \
  -H "Authorization: Bearer <token>"
```

Полная документация API: [docs/api-contracts.md](docs/api-contracts.md)

## Структура проекта

```
├── cmd/                    # Точки входа сервисов
│   ├── api-gateway/
│   ├── user-service/
│   ├── account-service/
│   ├── order-service/
│   ├── portfolio-service/
│   └── matching-engine/
├── internal/               # Бизнес-логика
│   ├── models/
│   ├── db/
│   ├── user/
│   ├── account/
│   ├── order/
│   └── portfolio/
├── pkg/                    # Общие пакеты
│   ├── config/
│   ├── errors/
│   ├── middleware/
│   └── response/
├── deploy/                 # Деплой и миграции
├── docs/                   # Документация
├── docker-compose.yml
└── go.mod
```

## Документация

- [Архитектура](docs/architecture.md)
- [Руководство разработчика](docs/development.md)
- [Доменная модель](docs/01-domain-model.md)
- [ER диаграмма](docs/02-er-diagram.md)
- [API контракты](docs/04-api-contracts.md)
- [Событийная модель](docs/05-event-flow.md)

## Требования к производительности

| Метрика | Цель |
|---------|------|
| Латентность ордера | < 100ms |
| Латентность рыночных данных | < 50ms |
| Throughput | 1000+ ордеров/сек |
| Одновременные пользователи | 10,000+ |
| Uptime | 99.9% |

## Тестирование

```bash
# Все тесты
go test ./...

# С покрытием
go test -cover ./...

# Интеграционные
go test -tags=integration ./...
```

## Технологии

- **Go 1.22** — микросервисы
- **C++17** — matching engine
- **PostgreSQL 16** — хранилище данных
- **Redis** — кэширование
- **Docker** — контейнеризация
- **gRPC** — межсервисное взаимодействие

## Лицензия

MIT License
