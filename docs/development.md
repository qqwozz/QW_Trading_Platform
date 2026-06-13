# Developer Guide

## Быстрый старт

### Предварительные требования

- Go 1.22+
- C++17 (GCC/Clang)
- CMake 3.14+
- PostgreSQL 16+
- Docker & Docker Compose

### Запуск через Docker Compose

```bash
# Клонировать репозиторий
git clone https://github.com/qw-trading/platform.git
cd platform

# Запустить все сервисы
docker-compose up --build

# Или в фоновом режиме
docker-compose up -d --build
```

### Локальная разработка

```bash
# 1. Запустить PostgreSQL
docker run -d --name postgres -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=qw_trading \
  postgres:16-alpine

# 2. Применить миграции
psql -h localhost -U postgres -d qw_trading -f deploy/migrations/001_init.sql

# 3. Запустить сервисы (в отдельных терминалах)
export DB_HOST=localhost DB_USER=postgres DB_PASSWORD=postgres DB_NAME=qw_trading JWT_SECRET=dev-secret

go run cmd/user-service/main.go      # :8081
go run cmd/account-service/main.go   # :8082
go run cmd/order-service/main.go     # :8083
go run cmd/portfolio-service/main.go # :8084
go run cmd/api-gateway/main.go       # :8080

# 4. Matching Engine (C++)
cd cmd/matching-engine
mkdir build && cd build
cmake .. && make
./matching-engine
```

## API Endpoints

### Аутентификация

```bash
# Регистрация
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","username":"trader","password":"secret123"}'

# Вход
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'
# Response: {"access_token":"...","refresh_token":"..."}
```

### Счета

```bash
# Список счетов
curl http://localhost:8080/accounts \
  -H "Authorization: Bearer <token>"

# Пополнение
curl -X POST http://localhost:8080/accounts/deposit \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"currency":"USDT","amount":10000}'
```

### Ордера

```bash
# Создание лимитного ордера
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTC/USDT","side":"BUY","type":"LIMIT","price":50000,"quantity":0.1}'

# Создание маркет ордера
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTC/USDT","side":"BUY","type":"MARKET","quantity":0.1}'

# Список ордеров
curl "http://localhost:8080/orders?symbol=BTC/USDT&limit=20" \
  -H "Authorization: Bearer <token>"

# Отмена ордера
curl -X DELETE http://localhost:8080/orders/<order-id> \
  -H "Authorization: Bearer <token>"
```

### Портфель

```bash
# Сводка портфеля
curl http://localhost:8080/portfolio \
  -H "Authorization: Bearer <token>"

# Позиции
curl http://localhost:8080/positions \
  -H "Authorization: Bearer <token>"

# Балансы
curl http://localhost:8080/balances \
  -H "Authorization: Bearer <token>"
```

## Структура кода

### Добавление нового эндпоинта

1. Создать handler в `internal/<service>/handler/handler.go`:
```go
func (h *Handler) MyEndpoint(w http.ResponseWriter, r *http.Request) {
    userID, ok := middleware.GetUserID(r)
    if !ok {
        response.Unauthorized(w, "unauthorized")
        return
    }
    // Business logic
    response.Success(w, result)
}
```

2. Зарегистрировать маршрут в `cmd/<service>/main.go`:
```go
mux.HandleFunc("GET /my-endpoint", h.MyEndpoint)
```

### Добавление нового репозитория

1. Создать файл `internal/<service>/repository/repository.go`:
```go
type MyRepository struct {
    db *db.Database
}

func New(db *db.Database) *MyRepository {
    return &MyRepository{db: db}
}

func (r *MyRepository) FindByID(id uuid.UUID) (*MyModel, error) {
    // SQL query
}
```

### Обработка ошибок

Используйте пакет `pkg/errors`:
```go
import apperr "github.com/qw-trading/platform/pkg/errors"

// Возвращаемые ошибки автоматически конвертируются в HTTP статусы
return apperr.NotFound("resource not found")
return apperr.BadRequest("invalid input")
return apperr.Internal("something went wrong")
```

## Конфигурация

Все сервисы настраиваются через переменные окружения:

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DB_HOST` | Хост PostgreSQL | localhost |
| `DB_PORT` | Порт PostgreSQL | 5432 |
| `DB_USER` | Пользователь БД | postgres |
| `DB_PASSWORD` | Пароль БД | postgres |
| `DB_NAME` | Имя БД | qw_trading |
| `JWT_SECRET` | Секрет для JWT | (обязательно) |
| `PORT` | Порт сервиса | 8080 |
| `APP_ENV` | Окружение | development |

## Тестирование

```bash
# Unit тесты
go test ./...

# С интеграцией
go test -tags=integration ./...

# С покрытием
go test -cover ./...
```

## Деплой

### Docker

```bash
# Сборка образов
docker-compose build

# Запуск
docker-compose up -d

# Проверка статуса
docker-compose ps

# Логи
docker-compose logs -f <service>
```

### Kubernetes

Деплой через Helm chart или kubectl:
```bash
kubectl apply -f deploy/k8s/
```
