# Architecture

## Overview

QW Trading Platform — виртуальная биржа криптовалют, построенная по принципам микросервисной архитектуры.

## Стек технологий

| Компонент | Технология | Обоснование |
|-----------|-----------|-------------|
| User Service | Go | CRUD операции, аутентификация |
| Account Service | Go | Финансовые операции с ACID |
| Order Service | Go | Управление жизненным циклом ордеров |
| Matching Engine | C++ | Максимальная производительность |
| Portfolio Service | Go | Управление позициями |
| API Gateway | Go | Маршрутизация, аутентификация |
| База данных | PostgreSQL | ACID, финансовыe данные |
| Кэш | Redis | Быстрый доступ, сессии |

## Структура проекта

```
QW_Trading_Platform/
├── cmd/                        # Точка входа каждого сервиса
│   ├── api-gateway/            # API Gateway (Go)
│   ├── user-service/           # User Service (Go)
│   ├── account-service/        # Account Service (Go)
│   ├── order-service/          # Order Service (Go)
│   ├── portfolio-service/      # Portfolio Service (Go)
│   └── matching-engine/        # Matching Engine (C++)
├── internal/                   # Бизнес-логика сервисов
│   ├── models/                 # Доменные модели
│   ├── db/                     # Подключение к БД
│   ├── user/
│   │   ├── handler/            # HTTP обработчики
│   │   └── repository/         # Работа с БД
│   ├── account/
│   │   ├── handler/
│   │   └── repository/
│   ├── order/
│   │   ├── handler/
│   │   └── repository/
│   ├── portfolio/
│   │   ├── handler/
│   │   └── repository/
│   └── matching/               # Matching engine Go wrapper
├── pkg/                        # Общие пакеты
│   ├── config/                 # Конфигурация
│   ├── errors/                 # Обработка ошибок
│   ├── middleware/             # HTTP middleware
│   └── response/              # Формат ответов
├── api/                        # API контракты
│   └── proto/                  # gRPC proto файлы
├── deploy/                     # Деплой
│   ├── migrations/             # SQL миграции
│   └── scripts/               # Скрипты
├── docs/                       # Документация
├── docker-compose.yml
├── Dockerfile.go
├── Dockerfile.cpp
└── go.mod
```

## Сервисы

### API Gateway (Go)
- Единая точка входа для клиентов
- Маршрутизация запросов к сервисам
- Аутентификация через JWT
- CORS, Rate Limiting
- Порт: 8080

### User Service (Go)
- Регистрация и аутентификация пользователей
- Управление профилем
- Генерация JWT токенов
- Порт: 8081

### Account Service (Go)
- Управление виртуальными счетами
- Баланс, заморозка средств
- История транзакций
- Порт: 8082

### Order Service (Go)
- Создание и отмена ордеров
- Валидация параметров
- История ордеров
- Порт: 8083

### Portfolio Service (Go)
- Учет открытых позиций
- Расчет средней цены
- Стоимость портфеля
- Порт: 8084

### Matching Engine (C++)
- Книга заявок в памяти
- Price-Time Priority
- Limit и Market ордера
- Частичное и полное исполнение
- Расчет комиссий
- Порт: 50051 (gRPC)

## Потоки данных

### Размещение ордера
```
Client → API Gateway → Order Service → Matching Engine
                                          ↓
                                    Trade Executed
                                          ↓
                              Account Service (обновление баланса)
```

### Рыночные данные
```
Matching Engine → Market Data Service → WebSocket → Client
```

## Консистентность

| Сервис | Консистентность | Причина |
|--------|----------------|---------|
| Matching Engine | Strong | Точность торгов критична |
| Account Service | Strong | Целостность балансов |
| Order Service | Strong | Управление состоянием ордеров |
| Portfolio Service | Eventual | Может иметь задержку |

## Требования к производительности

| Метрика | Цель |
|---------|------|
| Латентность размещения ордера | < 100ms |
| Латентность рыночных данных | < 50ms |
| Ордеров в секунду | 1000+ |
| Одновременных пользователей | 10,000+ |
| Uptime | 99.9% |
