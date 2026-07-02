# Как я создал QW Trading Platform: Полная история разработки

## Введение

Эта статья — подробный рассказ о том, как я создал QW Trading Platform — полноценную биржевую платформу с микросервисной архитектурой. Проект включает в себя высокопроизводительный matching engine на C++, микросервисы на Go, React-интерфейс для торговли, и всю необходимую инфраструктуру для работы в production.

**Цель проекта:** Создать не просто учебное приложение, а полноценный симулятор криптобиржи, который можно было бы показать работодателю как пример серьезной инженерной работы.

**Стек технологий:**
- Backend: Go 1.25, C++17
- Frontend: React 19, TypeScript 6, Tailwind CSS 4
- Database: PostgreSQL 16, Redis 7
- Infrastructure: Docker Compose, GitHub Actions CI/CD
- Тестирование: go test с race detection, моки

**Итоговые цифры:**
- ~10,000 строк Go + ~2,000 строк C++ + ~3,000 строк TypeScript
- 94 файла изменены за все время разработки
- 20 коммитов за 3 недели (13 июня — 2 июля 2026)
- 9 контейнеров в Docker Compose
- 7 независимых микросервисов

---

## Глава 1: Начало пути — от идеи к архитектуре

### Идея

Все началось с идеи создать не просто веб-приложение, а полноценную биржевую платформу. Я хотел, чтобы это был production-grade проект с настоящей микросервисной архитектурой, который можно было бы масштабировать и развивать.

Первым делом я создал техническое задание в файле `TECH_EX.md`. Это был не просто список требований — это была полная документация с архитектурными решениями, компромиссами и roadmap из 14 фаз.

### Архитектурные решения

Основные принципы, которые я заложил в архитектуру:

1. **Разделение ответственности:** Каждый сервис отвечает только за свою предметную область
2. **Независимость:** Сервисы могут развертываться и масштабироваться независимо друг от друга
3. **Производительность:** Matching Engine выделен в отдельный сервис на C++ для максимальной скорости
4. **Надежность:** Все критические данные хранятся в PostgreSQL с ACID-гарантиями

**Ключевое решение:** Matching Engine на C++. Почему C++? Потому что торговое ядро — это критически важный компонент, где каждая миллисекунда на счету. Хранение книги заявок в оперативной памяти обеспечивает минимальную задержку поиска встречных ордеров.

### Структура проекта

Я создал четкую структуру проекта:

```
qw-trading-platform/
├── cmd/                  # Точки входа сервисов
├── internal/             # Бизнес-логика
├── pkg/                  # Общие пакеты
├── deploy/               # Docker и миграции
├── docs/                 # Документация
└── matching-engine/      # C++ движок
```

Такая структура позволяет легко ориентироваться в проекте и понимать, где что находится.

---

## Глава 2: Закладка фундамента — модели данных и база

### Доменные модели

Сначала я создал `internal/models/models.go` — все доменные сущности системы:

```go
type User struct {
    ID           uuid.UUID
    Email        string
    Username     string
    PasswordHash string
    Status       UserStatus
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type Account struct {
    ID            uuid.UUID
    UserID        uuid.UUID
    Type          AccountType
    Balance       float64
    FrozenBalance float64
    Currency      string
    Status        AccountStatus
}

type Order struct {
    ID             uuid.UUID
    UserID         uuid.UUID
    AccountID      uuid.UUID
    Symbol         string
    Side           OrderSide
    Type           OrderType
    Price          *float64
    Quantity       float64
    FilledQuantity float64
    Status         OrderStatus
    TimeInForce    TimeInForce
}
```

**Важно:** Все финансовые поля используют `NUMERIC(20,8)` в PostgreSQL для точности. Это критично для финансовых систем — нельзя терять копейки при округлении.

### Схема базы данных

Создал три файла миграций:

**001_init.sql — основные таблицы:**

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'ACTIVE'
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL DEFAULT 'CASH',
    balance NUMERIC(20,8) DEFAULT 0,
    frozen_balance NUMERIC(20,8) DEFAULT 0,
    currency VARCHAR(10) NOT NULL,
    -- ...
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) NOT NULL,
    type VARCHAR(10) NOT NULL,
    price NUMERIC(20,8),
    quantity NUMERIC(20,8) NOT NULL,
    filled_quantity NUMERIC(20,8) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'PENDING',
    -- ...
);

CREATE TABLE trades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) NOT NULL,
    buyer_order_id UUID NOT NULL REFERENCES orders(id),
    seller_order_id UUID NOT NULL REFERENCES orders(id),
    price NUMERIC(20,8) NOT NULL,
    quantity NUMERIC(20,8) NOT NULL,
    buyer_fee NUMERIC(20,8) DEFAULT 0,
    seller_fee NUMERIC(20,8) DEFAULT 0,
    executed_at TIMESTAMPTZ DEFAULT NOW()
);
```

**002_audit.sql — таблицы аудита:**

```sql
CREATE TABLE balance_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    account_id UUID NOT NULL REFERENCES accounts(id),
    currency VARCHAR(10) NOT NULL,
    amount NUMERIC(20,8) NOT NULL,
    balance_before NUMERIC(20,8) NOT NULL,
    balance_after NUMERIC(20,8) NOT NULL,
    reason VARCHAR(50) NOT NULL,
    reference_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE position_history (
    -- Аналогично для позиций
);
```

**003_market_data.sql — таблицы рыночных данных:**

```sql
CREATE TABLE market_tickers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(20) UNIQUE NOT NULL,
    last_price NUMERIC(20,8) NOT NULL,
    best_bid NUMERIC(20,8),
    best_ask NUMERIC(20,8),
    volume_24h NUMERIC(20,8),
    high_24h NUMERIC(20,8),
    low_24h NUMERIC(20,8),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Интересный момент:** Я использовал `pgcrypto` для генерации UUID в PostgreSQL — это обеспечивает уникальность идентификаторов без дополнительных зависимостей на стороне приложения.

---

## Глава 3: Микросервисы — 7 независимых сервисов

### Обзор архитектуры

Я создал 7 независимых Go-сервисов, каждый в своем порту:

| Сервис | Порт | Ответственность |
|--------|------|-----------------|
| API Gateway | :8080 | Единая точка входа, CORS, rate limiting, request ID |
| User Service | :8081 | Регистрация, логин, JWT |
| Account Service | :8082 | Виртуальные счета, депозиты, выводы |
| Order Service | :8083 | CRUD ордеров, валидация |
| Portfolio Service | :8084 | Позиции, портфель, расчет equity |
| Market Data Service | :8085 | Тикеры, книга ордеров, WebSocket |
| History Service | :8086 | Аудит-трейл операций |

### API Gateway

API Gateway — это единая точка входа для всех клиентов. Он реализован как HTTP reverse proxy с маршрутизацией по префиксам:

```go
type Gateway struct {
    routes []route
    client *http.Client
}

func NewGateway(_ *config.Config) *Gateway {
    routes := []route{
        {"/v1/auth", userSvc},
        {"/v1/users", userSvc},
        {"/v1/accounts", accountSvc},
        {"/v1/orders", orderSvc},
        {"/v1/positions", portfolioSvc},
        {"/v1/portfolio", portfolioSvc},
        {"/v1/market", marketSvc},
        {"/v1/history", historySvc},
    }
    // Сортировка по длине префикса для корректного матчинга
    sort.Slice(routes, func(i, j int) bool {
        return len(routes[i].prefix) > len(routes[j].prefix)
    })
    return &Gateway{routes: routes, client: &http.Client{Timeout: 30 * time.Second}}
}
```

**Важная деталь:** Маршруты сортируются по длине префикса. Это нужно, чтобы `/v1/auth/login` матчился раньше, чем `/v1/auth`.

**Middleware стек:**

```go
wrapped := middleware.RequestID(
    middleware.Logger(logger)(
        rl.Middleware(
            middleware.CORS(cfg.AllowedOrigins)(gateway),
        ),
    ),
)
```

Цепочка middleware:
1. **Request ID** — генерирует уникальный ID для каждого запроса
2. **Logger** — логирует информацию о запросе
3. **Rate Limiter** — ограничивает количество запросов
4. **CORS** — обрабатывает跨域 запросы

### User Service

User Service отвечает за регистрацию и аутентификацию пользователей. Вот ключевые моменты:

**Регистрация с валидацией:**

```go
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    // Ограничение размера тела запроса
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.BadRequest(w, "invalid request body")
        return
    }

    if err := validateRegister(&req); err != nil {
        response.BadRequest(w, err.Error())
        return
    }

    // Проверка уникальности email
    exists, err := h.repo.EmailExists(r.Context(), req.Email)
    if exists {
        response.Conflict(w, "email already exists")
        return
    }

    // Хеширование пароля с bcrypt
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        response.InternalError(w, "failed to hash password")
        return
    }

    user := &models.User{
        ID:           uuid.New(),
        Email:        req.Email,
        Username:     req.Username,
        PasswordHash: string(hash),
        Status:       models.UserStatusActive,
    }

    if err := h.repo.Create(r.Context(), user); err != nil {
        response.InternalError(w, "failed to create user")
        return
    }

    response.Created(w, UserResponse{
        ID:        user.ID.String(),
        Email:     user.Email,
        Username:  user.Username,
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
    })
}
```

**Валидация:**

```go
func validateRegister(req *RegisterRequest) error {
    if req.Email == "" || req.Username == "" || req.Password == "" {
        return errors.BadRequest("email, username, and password are required")
    }
    if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
        return errors.BadRequest("invalid email format")
    }
    if utf8.RuneCountInString(req.Username) < 3 || utf8.RuneCountInString(req.Username) > 32 {
        return errors.BadRequest("username must be 3-32 characters")
    }
    if utf8.RuneCountInString(req.Password) < 8 {
        return errors.BadRequest("password must be at least 8 characters")
    }
    return nil
}
```

**Логин с защитой от brute-force:**

```go
var (
    failedLogins = make(map[string]int)
    failedMu     sync.Mutex
)

const (
    maxFailedLogins = 5
    lockoutDuration = 15 * time.Minute
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    // Проверка лимита неудачных попыток
    failedMu.Lock()
    if failedLogins[req.Email] >= maxFailedLogins {
        failedMu.Unlock()
        response.TooManyRequests(w, "account temporarily locked, try again later")
        return
    }
    failedMu.Unlock()

    user, err := h.repo.GetByEmail(r.Context(), req.Email)
    if err != nil {
        response.Unauthorized(w, "invalid credentials")
        return
    }

    // Проверка пароля
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        failedMu.Lock()
        failedLogins[req.Email]++
        failedMu.Unlock()
        response.Unauthorized(w, "invalid credentials")
        return
    }

    // Успешный вход — сброс счетчика
    failedMu.Lock()
    delete(failedLogins, req.Email)
    failedMu.Unlock()

    // Генерация JWT токенов
    accessToken, _ := h.generateToken(user.ID, time.Duration(h.jwtExpiry)*time.Hour)
    refreshToken, _ := h.generateToken(user.ID, 7*24*time.Hour)

    response.Success(w, AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    h.jwtExpiry * 3600,
    })
}
```

**Guest Login для быстрого входа:**

```go
func (h *Handler) GuestLogin(w http.ResponseWriter, r *http.Request) {
    // Если пользователь уже существует — логиним
    exists, _ := h.repo.EmailExists(r.Context(), req.Email)
    if exists {
        user, _ := h.repo.GetByEmail(r.Context(), req.Email)
        accessToken, _ := h.generateToken(user.ID, time.Duration(h.jwtExpiry)*time.Hour)
        refreshToken, _ := h.generateToken(user.ID, 7*24*time.Hour)
        response.Success(w, AuthResponse{
            AccessToken:  accessToken,
            RefreshToken: refreshToken,
            ExpiresIn:    h.jwtExpiry * 3600,
        })
        return
    }

    // Иначе — регистрируем и создаем счет с 1000 USDT
    hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    user := &models.User{
        ID:           uuid.New(),
        Email:        req.Email,
        Username:     req.Username,
        PasswordHash: string(hash),
        Status:       models.UserStatusActive,
    }
    h.repo.Create(r.Context(), user)

    // Создаем счет с 1000 USDT
    accountID := uuid.New()
    h.db.ExecContext(r.Context(),
        `INSERT INTO accounts (id, user_id, type, balance, frozen_balance, currency, status) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
        accountID, user.ID, models.AccountTypeCash, 1000.0, 0.0, "USDT", models.AccountStatusActive,
    )

    // Генерируем токены
    accessToken, _ := h.generateToken(user.ID, time.Duration(h.jwtExpiry)*time.Hour)
    refreshToken, _ := h.generateToken(user.ID, 7*24*time.Hour)
    response.Success(w, AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    h.jwtExpiry * 3600,
    })
}
```

### Order Service

Order Service отвечает за создание, чтение, обновление и удаление ордеров.

**Создание ордера с валидацией:**

```go
type CreateOrderRequest struct {
    Symbol      string   `json:"symbol"`
    Side        string   `json:"side"`
    Type        string   `json:"type"`
    Price       *float64 `json:"price,omitempty"`
    Quantity    float64  `json:"quantity"`
    TimeInForce string   `json:"time_in_force"`
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    userID, ok := middleware.GetUserID(r)
    if !ok {
        response.Unauthorized(w, "unauthorized")
        return
    }

    var req CreateOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.BadRequest(w, "invalid request body")
        return
    }

    if err := validateOrder(&req); err != nil {
        response.BadRequest(w, err.Error())
        return
    }

    order := &models.Order{
        ID:          uuid.New(),
        UserID:      userID,
        Symbol:      req.Symbol,
        Side:        models.OrderSide(req.Side),
        Type:        models.OrderType(req.Type),
        Price:       req.Price,
        Quantity:    req.Quantity,
        Status:      models.OrderStatusOpen,
        TimeInForce: models.TimeInForce(req.TimeInForce),
    }

    if order.TimeInForce == "" {
        order.TimeInForce = models.TimeInForceGTC
    }

    if err := h.repo.Create(r.Context(), order); err != nil {
        response.InternalError(w, "failed to create order")
        return
    }

    response.Created(w, toOrderResponse(order))
}

func validateOrder(req *CreateOrderRequest) error {
    if req.Symbol == "" {
        return errors.BadRequest("symbol is required")
    }
    if req.Side != "BUY" && req.Side != "SELL" {
        return errors.BadRequest("side must be BUY or SELL")
    }
    if req.Type != "LIMIT" && req.Type != "MARKET" {
        return errors.BadRequest("type must be LIMIT or MARKET")
    }
    if req.Quantity <= 0 {
        return errors.BadRequest("quantity must be positive")
    }
    if req.Type == "LIMIT" && (req.Price == nil || *req.Price <= 0) {
        return errors.BadRequest("price is required for limit orders")
    }
    tif := req.TimeInForce
    if tif == "" {
        tif = "GTC"
    }
    if tif != "GTC" && tif != "IOC" && tif != "FOK" {
        return errors.BadRequest("time_in_force must be GTC, IOC, or FOK")
    }
    return nil
}
```

**Получение ордера с проверкой владельца:**

```go
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
    userID, ok := middleware.GetUserID(r)
    if !ok {
        response.Unauthorized(w, "unauthorized")
        return
    }

    orderID, err := uuid.Parse(r.PathValue("id"))
    if err != nil {
        response.BadRequest(w, "invalid order id")
        return
    }

    order, err := h.repo.GetByID(r.Context(), orderID)
    if err != nil {
        if appErr, ok := err.(*errors.AppError); ok && appErr.Code == http.StatusNotFound {
            response.NotFound(w, "order not found")
            return
        }
        response.InternalError(w, "failed to get order")
        return
    }

    // Проверка, что ордер принадлежит пользователю
    if order.UserID != userID {
        response.Forbidden(w, "access denied")
        return
    }

    response.Success(w, toOrderResponse(order))
}
```

**Отмена ордера:**

```go
func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
    userID, ok := middleware.GetUserID(r)
    if !ok {
        response.Unauthorized(w, "unauthorized")
        return
    }

    orderID, err := uuid.Parse(r.PathValue("id"))
    if err != nil {
        response.BadRequest(w, "invalid order id")
        return
    }

    order, err := h.repo.GetByID(r.Context(), orderID)
    if err != nil {
        response.NotFound(w, "order not found")
        return
    }

    if order.UserID != userID {
        response.Forbidden(w, "access denied")
        return
    }

    // Можно отменить только OPEN или PARTIAL ордера
    if order.Status != models.OrderStatusOpen && order.Status != models.OrderStatusPartiallyFilled {
        response.Conflict(w, "order cannot be cancelled")
        return
    }

    if err := h.repo.UpdateStatus(r.Context(), orderID, models.OrderStatusCancelled, order.FilledQuantity); err != nil {
        response.InternalError(w, "failed to cancel order")
        return
    }

    order.Status = models.OrderStatusCancelled
    response.Success(w, toOrderResponse(order))
}
```

---

## Глава 4: Matching Engine на C++ — сердце системы

### Архитектура движка

Это была самая сложная и интересная часть проекта. Я написал полноценный торговый движок на C++17.

**Файлы:**
- `types.h` — определения типов (Order, Trade, OrderBook, PriceLevel)
- `orderbook.h/.cpp` — логика книги заявок
- `engine.h/.cpp` — обертка с TCP-сервером

### Структуры данных

```cpp
enum class OrderType { LIMIT, MARKET };
enum class OrderSide { BUY, SELL };
enum class OrderStatus { OPEN, PARTIAL, FILLED, CANCELLED };
enum class TimeInForce { GTC, IOC, FOK };

struct Order {
    std::string id;
    std::string user_id;
    std::string symbol;
    OrderSide side;
    OrderType type;
    TimeInForce tif;
    double price;
    double quantity;
    double filled = 0;
    OrderStatus status;
    std::chrono::steady_clock::time_point created_at;

    double remaining() const { return quantity - filled; }
    bool is_buy() const { return side == OrderSide::BUY; }
};

struct Trade {
    std::string id;
    std::string symbol;
    std::string buy_order_id;
    std::string sell_order_id;
    std::string buyer_id;
    std::string seller_id;
    double price;
    double quantity;
    double buyer_fee;
    double seller_fee;
    std::chrono::steady_clock::time_point executed_at;
};

struct PriceLevel {
    double price;
    double total_qty = 0;
    std::list<std::string> order_ids;  // FIFO очередь
};

class OrderBook {
    std::map<double, PriceLevel, std::greater<double>> bids_;  //DESC
    std::map<double, PriceLevel, std::less<double>> asks_;    //ASC
    std::unordered_map<std::string, Order> orders_;
    std::unordered_map<std::string, double> remaining_;
    double last_price_ = 0;
    std::string symbol_;
    mutable std::mutex mu_;
    std::function<void(const Trade&)> on_trade_;
    std::function<void(const Order&)> on_order_;
};
```

### Алгоритм исполнения ордеров

Ключевая функция — `match_against`:

```cpp
std::vector<Trade> OrderBook::match_against(Order& incoming, bool check_price) {
    std::vector<Trade> trades;

    auto do_match = [&](auto& opposite, auto end) {
        for (auto it = opposite.begin(); it != end && incoming.remaining() > 0;) {
            auto& level = it->second;

            // Проверка цены для LIMIT ордеров
            if (check_price) {
                if (incoming.is_buy() && incoming.price < level.price) return;
                if (!incoming.is_buy() && incoming.price > level.price) return;
            }

            // Проходим по очереди ордеров на этом уровне цен
            for (auto oid = level.order_ids.begin(); oid != level.order_ids.end() && incoming.remaining() > 0;) {
                auto oit = orders_.find(*oid);
                if (oit == orders_.end()) {
                    oid = level.order_ids.erase(oid);
                    continue;
                }

                Order& counter = oit->second;
                double avail = remaining_[*oid];
                double qty = std::min(incoming.remaining(), avail);

                // Создаем сделку
                Trade trade;
                trade.id = std::to_string(++trade_counter);
                trade.symbol = symbol_;
                trade.price = level.price;
                trade.quantity = qty;

                if (incoming.is_buy()) {
                    trade.buy_order_id = incoming.id;
                    trade.sell_order_id = counter.id;
                    trade.buyer_id = incoming.user_id;
                    trade.seller_id = counter.user_id;
                } else {
                    trade.buy_order_id = counter.id;
                    trade.sell_order_id = incoming.id;
                    trade.buyer_id = counter.user_id;
                    trade.seller_id = incoming.user_id;
                }

                // Расчет комиссий
                double amount = trade.price * trade.quantity;
                trade.buyer_fee = calc_fee(amount, true);
                trade.seller_fee = calc_fee(amount, false);
                trade.executed_at = std::chrono::steady_clock::now();

                // Обновляем состояния
                incoming.filled += qty;
                counter.filled += qty;
                remaining_[counter.id] = counter.remaining();
                level.total_qty -= qty;
                last_price_ = trade.price;

                if (counter.remaining() <= 0) {
                    counter.status = OrderStatus::FILLED;
                    oid = level.order_ids.erase(oid);
                    if (on_order_) on_order_(counter);
                } else {
                    counter.status = OrderStatus::PARTIAL;
                    ++oid;
                }

                if (on_trade_) on_trade_(trade);
                trades.push_back(trade);
            }

            if (level.order_ids.empty() || level.total_qty <= 0) {
                it = opposite.erase(it);
            } else {
                ++it;
            }
        }
    };

    if (incoming.is_buy()) {
        do_match(asks_, asks_.end());
    } else {
        do_match(bids_, bids_.end());
    }

    return trades;
}
```

### Time-in-Force политики

```cpp
if (order.remaining() > 0 && order.status != OrderStatus::CANCELLED) {
    if (order.type == OrderType::MARKET && order.tif == TimeInForce::FOK) {
        order.status = OrderStatus::CANCELLED;
    } else if (order.tif == TimeInForce::IOC) {
        order.status = OrderStatus::CANCELLED;
    } else {
        insert_to_book(order);
        orders_[order.id] = order;
        remaining_[order.id] = order.remaining();
    }
}
```

- **GTC** (Good Till Cancelled) — ордер остается в книге до исполнения или отмены
- **IOC** (Immediate or Cancel) — исполняется немедленно, остаток отменяется
- **FOK** (Fill or Kill) — исполняется полностью или отменяется целиком

### Потокобезопасность

Каждая книга заявок имеет свой мьютекс:

```cpp
class OrderBook {
    mutable std::mutex mu_;
    // ...
};

double OrderBook::best_bid() const {
    std::lock_guard<std::mutex> lock(mu_);
    return bids_.empty() ? 0 : bids_.begin()->first;
}

double OrderBook::best_ask() const {
    std::lock_guard<std::mutex> lock(mu_);
    return asks_.empty() ? 0 : asks_.begin()->first;
}
```

### Комиссии

```cpp
double OrderBook::calc_fee(double amount, bool maker) {
    return amount * (maker ? 0.001 : 0.002);
}
```

- **Maker** (предложил ликвидность): 0.10% (10 bps)
- **Taker** (взял ликвидность): 0.20% (20 bps)

---

## Глава 5: WebSocket Hub — реальное время

### Архитектура WebSocket

WebSocket Hub реализован по паттерну "Hub and Spoke":

```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

type Client struct {
    hub    *Hub
    conn   *websocket.Conn
    send   chan []byte
    symbol string
}

type WSMessage struct {
    Type   string          `json:"type"`
    Symbol string          `json:"symbol,omitempty"`
    Data   json.RawMessage `json:"data,omitempty"`
}
```

### Основной цикл Hub

```go
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            total := len(h.clients)
            h.mu.Unlock()
            log.Printf("WebSocket client connected, symbol: %s, total: %d", client.symbol, total)

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            total := len(h.clients)
            h.mu.Unlock()
            log.Printf("WebSocket client disconnected, symbol: %s, total: %d", client.symbol, total)

        case message := <-h.broadcast:
            h.broadcastMessage(message, "")
        }
    }
}
```

### Broadcast с фильтрацией по символу

```go
func (h *Hub) broadcastMessage(payload []byte, symbol string) {
    h.mu.RLock()
    var stale []*Client
    for client := range h.clients {
        // Отправляем только клиентам, подписанным на нужный символ
        if symbol != "" && client.symbol != "" && client.symbol != symbol {
            continue
        }
        select {
        case client.send <- payload:
        default:
            // Канал полон — клиент считается устаревшим
            stale = append(stale, client)
        }
    }
    h.mu.RUnlock()

    // Удаляем устаревших клиентов
    if len(stale) > 0 {
        h.mu.Lock()
        for _, client := range stale {
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
        }
        h.mu.Unlock()
    }
}
```

### Аутентификация WebSocket

WebSocket не поддерживает заголовки Authorization, поэтому JWT передается в query parameter:

```go
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    tokenStr := r.URL.Query().Get("token")
    if tokenStr == "" {
        http.Error(w, `{"error":"token required"}`, http.StatusUnauthorized)
        return
    }

    secret := r.URL.Query().Get("secret")
    if secret == "" {
        http.Error(w, `{"error":"secret required"}`, http.StatusUnauthorized)
        return
    }

    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(secret), nil
    })
    if err != nil || !token.Valid {
        http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }

    client := &Client{
        hub:  h,
        conn: conn,
        send: make(chan []byte, 256),
    }

    h.register <- client
    go client.writePump()
    go client.readPump()
}
```

### Подписка на символы

```go
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        var msg WSMessage
        if err := json.Unmarshal(message, &msg); err != nil {
            continue
        }
        switch msg.Type {
        case "subscribe":
            c.symbol = msg.Symbol
        case "unsubscribe":
            c.symbol = ""
        }
    }
}
```

---

## Глава 6: Инфраструктура и DevOps

### Docker Compose

Создал `docker-compose.yml` с 9 контейнерами:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    container_name: qw_postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: qw_trading
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./deploy/migrations/001_init.sql:/docker-entrypoint-initdb.d/01_init.sql
      - ./deploy/migrations/003_market_data.sql:/docker-entrypoint-initdb.d/03_market_data.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 3s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: qw_redis
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

  user-service:
    build:
      context: .
      dockerfile: deploy/Dockerfile-go
      args:
        SERVICE: cmd/user-service
    container_name: qw_user_service
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: qw_trading
      JWT_SECRET: ${JWT_SECRET:-change-me-in-production}
      PORT: 8081
    ports:
      - "8081:8081"
    depends_on:
      postgres: { condition: service_healthy }

  # ... остальные сервисы аналогично

  matching-engine:
    build:
      context: .
      dockerfile: deploy/Dockerfile-cpp
    container_name: qw_matching_engine
    ports:
      - "50051:50051"

  api-gateway:
    build:
      context: .
      dockerfile: deploy/Dockerfile-go
      args:
        SERVICE: cmd/api-gateway
    container_name: qw_api_gateway
    environment:
      PORT: 8080
      USER_SERVICE_URL: http://user-service:8081
      ACCOUNT_SERVICE_URL: http://account-service:8082
      ORDER_SERVICE_URL: http://order-service:8083
      PORTFOLIO_SERVICE_URL: http://portfolio-service:8084
      MARKET_SERVICE_URL: http://market-data-service:8085
      HISTORY_SERVICE_URL: http://history-service:8086
    ports:
      - "8080:8080"
    depends_on:
      - user-service
      - account-service
      - order-service
      - portfolio-service
      - market-data-service
      - history-service

  frontend:
    build:
      context: .
      dockerfile: deploy/Dockerfile-frontend
    container_name: qw_frontend
    ports:
      - "3000:80"
    depends_on:
      - api-gateway

volumes:
  pgdata:
```

### Dockerfile для Go

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /service ./${SERVICE}/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /service .
EXPOSE 8080
CMD ["./service"]
```

### Dockerfile для C++

```dockerfile
FROM gcc:13 AS builder
WORKDIR /app
RUN apt-get update && apt-get install -y cmake
COPY . .
RUN mkdir build && cd build && cmake .. && make -j$(nproc)

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends libstdc++6 && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/build/matching-engine .
EXPOSE 50051
CMD ["./matching-engine"]
```

### Dockerfile для Frontend

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### CI/CD Pipeline

GitHub Actions workflow (`.github/workflows/ci.yml`):

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: go vet ./...
      - run: go install honnef.co/go/tools/cmd/staticcheck@latest
      - run: staticcheck ./...

  test:
    runs-on: ubuntu-latest
    needs: lint
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: qw_trading_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: go test -v -race -coverprofile=coverage.out ./...
      - run: go tool cover -html=coverage.out -o coverage.html
      - uses: actions/upload-artifact@v4
        with:
          name: coverage
          path: coverage.html

  build:
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: |
          CGO_ENABLED=0 go build -o bin/api-gateway ./cmd/api-gateway
          CGO_ENABLED=0 go build -o bin/user-service ./cmd/user-service
          CGO_ENABLED=0 go build -o bin/account-service ./cmd/account-service
          CGO_ENABLED=0 go build -o bin/order-service ./cmd/order-service
          CGO_ENABLED=0 go build -o bin/portfolio-service ./cmd/portfolio-service
          CGO_ENABLED=0 go build -o bin/market-data-service ./cmd/market-data-service
          CGO_ENABLED=0 go build -o bin/history-service ./cmd/history-service

  docker:
    runs-on: ubuntu-latest
    needs: build
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          push: true
          tags: ghcr.io/${{ github.repository }}:latest
          file: deploy/Dockerfile-go
          build-args: SERVICE=cmd/api-gateway
```

---

## Глава 7: React Frontend — торговый терминал

### Архитектура фронтенда

Создал React SPA с TypeScript:

```
frontend/
├── src/
│   ├── components/
│   │   ├── Layout.tsx
│   │   └── CandlestickChart.tsx
│   ├── pages/
│   │   ├── LoginPage.tsx
│   │   ├── RegisterPage.tsx
│   │   ├── DashboardPage.tsx
│   │   ├── TradePage.tsx
│   │   ├── PortfolioPage.tsx
│   │   └── HistoryPage.tsx
│   ├── lib/
│   │   └── api.ts
│   ├── App.tsx
│   └── main.tsx
├── package.json
└── vite.config.ts
```

### Стек технологий

| Компонент | Технология |
|-----------|------------|
| Framework | React 19 + TypeScript 6 |
| Build Tool | Vite 8 |
| Styling | Tailwind CSS 4 |
| Charts | lightweight-charts (TradingView) |
| Icons | Lucide React |
| Routing | React Router 7 |
| Linting | ESLint 10 + typescript-eslint |

### API клиент

Создал единый API клиент с автоматическим обновлением токенов:

```typescript
import axios from 'axios';

const BASE_URL = 'http://localhost:8080';

let authToken: string | null = null;

export function setToken(token: string) {
  authToken = token;
  localStorage.setItem('token', token);
}

export function clearToken() {
  authToken = null;
  localStorage.removeItem('token');
}

export function getToken(): string | null {
  if (!authToken) {
    authToken = localStorage.getItem('token');
  }
  return authToken;
}

const client = axios.create({
  baseURL: BASE_URL,
});

client.interceptors.request.use((config) => {
  const token = getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

client.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      clearToken();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const api = {
  auth: {
    login: (email: string, password: string) =>
      client.post('/v1/auth/login', { email, password }).then(r => r.data),
    register: (email: string, username: string, password: string) =>
      client.post('/v1/auth/register', { email, username, password }).then(r => r.data),
    guest: (email: string, username: string, password: string) =>
      client.post('/v1/auth/guest', { email, username, password }).then(r => r.data),
    me: () => client.get('/v1/users/me').then(r => r.data),
  },
  orders: {
    create: (data: any) => client.post('/v1/orders', data).then(r => r.data),
    list: (params?: any) => client.get('/v1/orders', { params }).then(r => r.data),
    get: (id: string) => client.get(`/v1/orders/${id}`).then(r => r.data),
    cancel: (id: string) => client.delete(`/v1/orders/${id}`).then(r => r.data),
  },
  portfolio: {
    get: () => client.get('/v1/portfolio').then(r => r.data),
    positions: () => client.get('/v1/positions').then(r => r.data),
  },
  accounts: {
    list: () => client.get('/v1/accounts').then(r => r.data),
    deposit: (currency: string, amount: number) =>
      client.post('/v1/accounts/deposit', { currency, amount }).then(r => r.data),
  },
  market: {
    tickers: () => client.get('/v1/market/tickers').then(r => r.data),
    ticker: (symbol: string) => client.get(`/v1/market/tickers/${symbol}`).then(r => r.data),
    orderbook: (symbol: string, depth?: number) =>
      client.get(`/v1/market/orderbook/${symbol}`, { params: { depth } }).then(r => r.data),
  },
  history: {
    orders: (params?: any) => client.get('/v1/history/orders', { params }).then(r => r.data),
    trades: (params?: any) => client.get('/v1/history/trades', { params }).then(r => r.data),
  },
};
```

### Главный компонент App

```tsx
function App() {
  const [user, setUser] = useState<{ email: string; username: string } | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      if (getToken()) {
        try {
          const res = await api.auth.me();
          setUser(res);
        } catch {
          clearToken();
        }
      } else {
        // Автоматический guest login
        try {
          const id = crypto.randomUUID().slice(0, 8);
          const email = `guest-${id}@guest.local`;
          const username = `guest_${id}`;
          const password = `guest_${id}_pass!`;
          const res = await api.auth.guest(email, username, password);
          setToken(res.access_token);
          const me = await api.auth.me();
          setUser(me);
        } catch {
          // Сервер недоступен
        }
      }
      setLoading(false);
    };
    init();
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={!user ? <LoginPage onLogin={handleLogin} /> : <Navigate to="/" />} />
        <Route path="/register" element={!user ? <RegisterPage onRegister={handleRegister} /> : <Navigate to="/" />} />
        <Route path="/" element={user ? <Layout user={user} onLogout={handleLogout} /> : <Navigate to="/" />}>
          <Route index element={<DashboardPage />} />
          <Route path="trade" element={<TradePage />} />
          <Route path="portfolio" element={<PortfolioPage />} />
          <Route path="history" element={<HistoryPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
```

### Trading Page — страница торговли

Самая сложная страница — страница торговли с книгой ордеров и графиком:

```tsx
export default function TradePage() {
  const [symbol, setSymbol] = useState('BTC/USDT');
  const [side, setSide] = useState<'BUY' | 'SELL'>('BUY');
  const [type, setType] = useState<'LIMIT' | 'MARKET'>('LIMIT');
  const [price, setPrice] = useState('');
  const [quantity, setQuantity] = useState('');
  const [orderBook, setOrderBook] = useState<any>(null);

  useEffect(() => {
    const load = () => {
      api.market.orderbook(symbol, 10).then((res: any) => setOrderBook(res.data)).catch(() => {});
    };
    load();
    const interval = setInterval(load, 5000);
    return () => clearInterval(interval);
  }, [symbol]);

  return (
    <div className="flex h-full">
      {/* Левая часть: график и статистика */}
      <div className="flex-1 flex flex-col">
        <div className="p-4 border-b flex items-center gap-4">
          {/* Выбор торговой пары */}
          {['BTC/USDT', 'ETH/USDT', 'SOL/USDT'].map(s => (
            <button key={s} onClick={() => setSymbol(s)}>{s}</button>
          ))}
        </div>
        <div className="flex-1 p-4">
          <CandlestickChart symbol={symbol} height={340} />
        </div>
        <div className="p-4 border-t">
          {/* 24h статистика */}
          <div className="grid grid-cols-4 gap-4">
            <div>24h High</div>
            <div>24h Low</div>
            <div>24h Volume</div>
            <div>Spread</div>
          </div>
        </div>
      </div>

      {/* Правая часть: форма ордера и книга */}
      <div className="w-72 flex flex-col border-l">
        <div className="p-4 border-b">
          {/* Выбор BUY/SELL и LIMIT/MARKET */}
          <button onClick={() => setSide('BUY')}>Buy</button>
          <button onClick={() => setSide('SELL')}>Sell</button>
          <button onClick={() => setType('LIMIT')}>LIMIT</button>
          <button onClick={() => setType('MARKET')}>MARKET</button>
        </div>

        <form onSubmit={handleSubmit} className="flex-1 p-4 flex flex-col">
          {type === 'LIMIT' && (
            <input type="number" value={price} onChange={e => setPrice(e.target.value)} placeholder="Price" />
          )}
          <input type="number" value={quantity} onChange={e => setQuantity(e.target.value)} placeholder="Quantity" />
          <button type="submit">{side} {symbol.split('/')[0]}</button>
        </form>

        <div className="p-4 border-t">
          {/* Книга ордеров */}
          <div>
            {orderBook?.asks?.slice(0, 6).reverse().map((ask: any) => (
              <div key={ask.price}>
                <span>{ask.price}</span>
                <span>{ask.quantity}</span>
              </div>
            ))}
          </div>
          <div>{orderBook?.asks?.[0]?.price || '--'}</div>
          <div>
            {orderBook?.bids?.slice(0, 6).map((bid: any) => (
              <div key={bid.price}>
                <span>{bid.price}</span>
                <span>{bid.quantity}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
```

### Candlestick Chart

Для графиков использовал lightweight-charts от TradingView:

```tsx
import { createChart, ColorType, CandlestickSeries } from 'lightweight-charts';

export default function CandlestickChart({ symbol = 'BTC/USDT', height = 400 }: Props) {
  const chartContainerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<ReturnType<typeof createChart> | null>(null);

  useEffect(() => {
    if (!chartContainerRef.current) return;

    const chart = createChart(chartContainerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: '#1e2329' },
        textColor: '#848e9c',
      },
      grid: {
        vertLines: { color: '#2b3139' },
        horzLines: { color: '#2b3139' },
      },
      crosshair: { mode: 0 },
      timeScale: {
        borderColor: '#2b3139',
        timeVisible: true,
        secondsVisible: false,
      },
      rightPriceScale: { borderColor: '#2b3139' },
      width: chartContainerRef.current.clientWidth,
      height,
    });

    const series = chart.addSeries(CandlestickSeries, {
      upColor: '#0ecb81',
      downColor: '#f6465d',
      borderDownColor: '#f6465d',
      borderUpColor: '#0ecb81',
      wickDownColor: '#f6465d',
      wickUpColor: '#0ecb81',
    });

    chartRef.current = chart;

    // Загрузка данных с Binance API
    const fetchKlines = async () => {
      const binanceSymbol = SYMBOL_MAP[symbol] || symbol.replace('/', '');
      const res = await fetch(
        `https://api.binance.com/api/v3/klines?symbol=${binanceSymbol}&interval=1h&limit=100`
      );
      const data = await res.json();

      const klines = data.map((k: any[]) => ({
        openTime: k[0],
        open: parseFloat(k[1]),
        high: parseFloat(k[2]),
        low: parseFloat(k[3]),
        close: parseFloat(k[4]),
        volume: parseFloat(k[5]),
      }));

      series.setData(
        klines.map(k => ({
          time: Math.floor(k.openTime / 1000) as any,
          open: k.open,
          high: k.high,
          low: k.low,
          close: k.close,
        }))
      );
    };

    fetchKlines();
    const interval = setInterval(fetchKlines, 30000);
    return () => clearInterval(interval);
  }, [symbol]);

  return <div ref={chartContainerRef} />;
}
```

### Dark Theme

Реализовал темную тему в стиле Binance:

```css
:root {
  --bg-primary: #0b0e11;
  --bg-secondary: #1e2329;
  --border: #2b3139;
  --text-primary: #eaecef;
  --text-secondary: #848e9c;
  --accent: #fcd535;
  --green: #0ecb81;
  --red: #f6465d;
}
```

---

## Глава 8: Тестирование

### Моки и интерфейсы

Я создал интерфейсы для всех репозиториев:

```go
// internal/order/repository/interfaces.go
type OrderRepositoryInterface interface {
    Create(ctx context.Context, order *models.Order) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
    List(ctx context.Context, filter ListFilter) ([]models.Order, int, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus, filled float64) error
}

type TradeRepositoryInterface interface {
    Create(ctx context.Context, trade *models.Trade) error
    ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Trade, int, error)
}
```

Затем написал моки с использованием `testify/mock`:

```go
// internal/order/handler/mock_test.go
type MockOrderRepository struct {
    mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *models.Order) error {
    args := m.Called(ctx, order)
    return args.Error(0)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Order), args.Error(1)
}
```

### Тесты обработчиков

```go
// internal/order/handler/handler_test.go
func TestCreateOrder_Success(t *testing.T) {
    mockRepo := new(MockOrderRepository)
    mockTradeRepo := new(MockTradeRepository)

    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    handler := New(mockRepo, mockTradeRepo)

    body := `{"symbol":"BTC/USDT","side":"BUY","type":"LIMIT","price":50000,"quantity":0.1}`
    req := httptest.NewRequest(http.MethodPost, "/v1/orders", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, uuid.New()))

    w := httptest.NewRecorder()
    handler.CreateOrder(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
    mockRepo.AssertExpectations(t)
}

func TestCreateOrder_InvalidSide(t *testing.T) {
    mockRepo := new(MockOrderRepository)
    mockTradeRepo := new(MockTradeRepository)

    handler := New(mockRepo, mockTradeRepo)

    body := `{"symbol":"BTC/USDT","side":"INVALID","type":"LIMIT","price":50000,"quantity":0.1}`
    req := httptest.NewRequest(http.MethodPost, "/v1/orders", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    handler.CreateOrder(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    assert.Contains(t, w.Body.String(), "side must be BUY or SELL")
}
```

### Тесты middleware

```go
// pkg/middleware/ratelimit_test.go
func TestRateLimiter_Allow(t *testing.T) {
    rl := NewRateLimiter(10, 5)

    // Первые 5 запросов должны проходить
    for i := 0; i < 5; i++ {
        assert.True(t, rl.allow("192.168.1.1"))
    }

    // 6-й запрос должен быть отклонен
    assert.False(t, rl.allow("192.168.1.1"))
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
    rl := NewRateLimiter(10, 5)

    // Разные IP независимы
    for i := 0; i < 5; i++ {
        assert.True(t, rl.allow("192.168.1.1"))
        assert.True(t, rl.allow("192.168.1.2"))
    }

    assert.False(t, rl.allow("192.168.1.1"))
    assert.False(t, rl.allow("192.168.1.2"))
}
```

### Запуск тестов

```bash
# Все тесты
go test ./...

# С race detection и coverage
go test -v -race -coverprofile=coverage.out ./...

# Просмотр отчета
go tool cover -html=coverage.out
```

---

## Глава 9: Продвинутые функции

### Rate Limiting

Реализовал token bucket алгоритм:

```go
type RateLimiter struct {
    mu       sync.Mutex
    visitors map[string]*visitor
    rate     float64
    burst    int
    cleanup  time.Duration
}

type visitor struct {
    tokens   float64
    lastSeen time.Time
    rate     float64
    burst    float64
}

func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
    rl := &RateLimiter{
        visitors: make(map[string]*visitor),
        rate:     requestsPerSecond,
        burst:    burst,
        cleanup:  30 * time.Second,
    }
    go rl.cleanupLoop()
    return rl
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    v, exists := rl.visitors[ip]
    if !exists {
        rl.visitors[ip] = &visitor{
            tokens:   float64(rl.burst),
            lastSeen: time.Now(),
            rate:     rl.rate,
            burst:    float64(rl.burst),
        }
        return rl.visitors[ip]
    }

    // Восстановление токенов
    elapsed := time.Since(v.lastSeen).Seconds()
    v.tokens += elapsed * v.rate
    if v.tokens > v.burst {
        v.tokens = v.burst
    }
    v.lastSeen = time.Now()
    return v
}

func (rl *RateLimiter) allow(ip string) bool {
    v := rl.getVisitor(ip)
    if v.tokens < 1 {
        return false
    }
    v.tokens--
    return true
}

func (rl *RateLimiter) cleanupLoop() {
    ticker := time.NewTicker(rl.cleanup)
    defer ticker.Stop()
    for range ticker.C {
        rl.mu.Lock()
        for ip, v := range rl.visitors {
            if time.Since(v.lastSeen) > rl.cleanup {
                delete(rl.visitors, ip)
            }
        }
        rl.mu.Unlock()
    }
}
```

### Structured Logging

Создал JSON-логгер с полями:

```go
type Entry struct {
    Time      string      `json:"time"`
    Level     Level       `json:"level"`
    Service   string      `json:"service"`
    Message   string      `json:"message"`
    RequestID string      `json:"request_id,omitempty"`
    Method    string      `json:"method,omitempty"`
    Path      string      `json:"path,omitempty"`
    Status    int         `json:"status,omitempty"`
    Duration  string      `json:"duration,omitempty"`
    Extra     interface{} `json:"extra,omitempty"`
}

type Logger struct {
    service string
}

func New(service string) *Logger {
    return &Logger{service: service}
}

func (l *Logger) log(level Level, msg string, fields map[string]interface{}) {
    entry := Entry{
        Time:    time.Now().UTC().Format(time.RFC3339Nano),
        Level:   level,
        Service: l.service,
        Message: msg,
    }
    if fields != nil {
        entry.Extra = fields
    }
    json.NewEncoder(os.Stdout).Encode(entry)
}

func (l *Logger) Request(method, path string, status int, duration time.Duration, requestID string) {
    l.log(LevelInfo, "request", map[string]interface{}{
        "method":     method,
        "path":       path,
        "status":     status,
        "duration":   duration.String(),
        "request_id": requestID,
    })
}
```

Пример лога:

```json
{
  "time": "2026-06-17T16:23:52.123456Z",
  "level": "INFO",
  "service": "api-gateway",
  "message": "request",
  "method": "POST",
  "path": "/v1/orders",
  "status": 201,
  "duration": "12ms",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Health Checks

Реализовал `/health` эндпоинт с проверкой подключения к БД:

```go
func HealthCheck(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()

        if err := db.PingContext(ctx); err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusServiceUnavailable)
            w.Write([]byte(`{"status":"unhealthy","error":"database connection failed"}`))
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"healthy","service":"` + serviceName + `"}`))
    }
}
```

### Graceful Shutdown

Создал паттерн graceful shutdown для всех сервисов:

```go
func Run(name string, handler http.Handler, cfg Config, log *logger.Logger) {
    srv := &http.Server{
        Addr:         ":" + os.Getenv("PORT"),
        Handler:      handler,
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
        IdleTimeout:  cfg.IdleTimeout,
    }

    go func() {
        log.Info(name + " starting", map[string]interface{}{
            "addr": srv.Addr,
        })
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(name+" failed to start", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    sig := <-quit

    log.Info(name+" received signal, shutting down", map[string]interface{}{
        "signal": sig.String(),
    })

    ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownWait)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Error(name+" shutdown error", map[string]interface{}{
            "error": err.Error(),
        })
    } else {
        log.Info(name+" stopped gracefully")
    }
}
```

### Fee Calculation

Создал модуль расчета комиссий:

```go
// pkg/trading/fees.go
package trading

func CalculateFee(amount float64, isMaker bool, makerBPS, takerBPS int) float64 {
    bps := takerBPS
    if isMaker {
        bps = makerBPS
    }
    return amount * float64(bps) / 10000
}

func IsMakerOrder(orderPrice, bestBid, bestAsk float64, side string) bool {
    if side == "BUY" {
        return orderPrice < bestAsk
    }
    return orderPrice > bestBid
}
```

Тесты:

```go
// pkg/trading/fees_test.go
func TestCalculateFee_Maker(t *testing.T) {
    fee := CalculateFee(1000, true, 10, 20)
    assert.Equal(t, 1.0, fee) // 0.10%
}

func TestCalculateFee_Taker(t *testing.T) {
    fee := CalculateFee(1000, false, 10, 20)
    assert.Equal(t, 2.0, fee) // 0.20%
}

func TestIsMakerOrder_Buy(t *testing.T) {
    // BUY ордер ниже best Ask — это maker
    assert.True(t, IsMakerOrder(49000, 49500, 50000, "BUY"))
    // BUY ордер выше best Ask — это taker
    assert.False(t, IsMakerOrder(50500, 49500, 50000, "BUY"))
}
```

### Market Feeder

Добавил генератор рыночных данных с Binance API:

```go
// internal/market/feeder/feeder.go
type Feeder struct {
    repo   *repository.MarketRepository
    client *http.Client
    pairs  []pair
}

type pair struct {
    binanceSymbol string
    dbSymbol      string
}

var defaultPairs = []pair{
    {binanceSymbol: "BTCUSDT", dbSymbol: "BTC/USDT"},
    {binanceSymbol: "ETHUSDT", dbSymbol: "ETH/USDT"},
    {binanceSymbol: "SOLUSDT", dbSymbol: "SOL/USDT"},
}

func (f *Feeder) Start(ctx context.Context) {
    log.Println("[feeder] starting Binance price feeder")
    go f.loop(ctx)
}

func (f *Feeder) loop(ctx context.Context) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    f.fetch(ctx)

    for {
        select {
        case <-ctx.Done():
            log.Println("[feeder] shutting down")
            return
        case <-ticker.C:
            f.fetch(ctx)
        }
    }
}

func (f *Feeder) fetchPair(ctx context.Context, p pair) error {
    url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/24hr?symbol=%s", p.binanceSymbol)
    resp, err := f.client.Get(url)
    if err != nil {
        return fmt.Errorf("HTTP request failed: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Binance API returned status %d: %s", resp.StatusCode, string(body))
    }

    var tickerResp binanceTicker
    if err := json.Unmarshal(body, &tickerResp); err != nil {
        return fmt.Errorf("failed to unmarshal response: %w", err)
    }

    ticker := &models.MarketTicker{
        ID:           uuid.New(),
        Symbol:       p.dbSymbol,
        LastPrice:    parseFloat(tickerResp.LastPrice),
        BestBid:      parseFloat(tickerResp.BidPrice),
        BestAsk:      parseFloat(tickerResp.AskPrice),
        Volume24h:    parseFloat(tickerResp.Volume),
        High24h:      parseFloat(tickerResp.HighPrice),
        Low24h:       parseFloat(tickerResp.LowPrice),
        Change24h:    parseFloat(tickerResp.PriceChange),
        ChangePct24h: parseFloat(tickerResp.PriceChangePercent),
        UpdatedAt:    time.Now().UTC(),
    }

    if err := f.repo.UpsertTicker(ctx, ticker); err != nil {
        return fmt.Errorf("upsert ticker failed: %w", err)
    }

    log.Printf("[feeder] updated %s: price=%.2f", p.dbSymbol, ticker.LastPrice)
    return nil
}
```

---

## Глава 10: Проблемы и решения

### Проблема 1: Версия Go в Dockerfile

**Проблема:** Dockerfile указывал Go 1.22, а проект использовал Go 1.25.

**Решение:** Обновил Dockerfile до Go 1.25:

```dockerfile
FROM golang:1.25-alpine AS builder
```

### Проблема 2: golangci-lint несовместимость

**Проблема:** golangci-lint требовал Go 1.22, проект на Go 1.25.

**Решение:** Заменил на go vet + staticcheck:

```yaml
lint:
  steps:
    - run: go vet ./...
    - run: go install honnef.co/go/tools/cmd/staticcheck@latest
    - run: staticcheck ./...
```

### Проблема 3: Гонки данных в Matching Engine

**Проблема:** Одновременный доступ к книге заявок вызывал гонки данных.

**Решение:** Добавил мьютексы на уровне книги заявок и движка:

```cpp
class OrderBook {
    mutable std::mutex mu_;
    // ...
};

double OrderBook::best_bid() const {
    std::lock_guard<std::mutex> lock(mu_);
    return bids_.empty() ? 0 : bids_.begin()->first;
}
```

### Проблема 4: JWT Token Refresh

**Проблема:** Access tokens быстро истекали, нужно было обновление.

**Решение:** Реализовал refresh token flow с 7-дневным TTL:

```go
accessToken, _ := h.generateToken(user.ID, time.Duration(h.jwtExpiry)*time.Hour)
refreshToken, _ := h.generateToken(user.ID, 7*24*time.Hour)
```

### Проблема 5: WebSocket Auth

**Проблема:** WebSocket не поддерживает заголовки Authorization.

**Решение:** Передаем JWT в query parameter:

```go
tokenStr := r.URL.Query().Get("token")
secret := r.URL.Query().Get("secret")
```

### Проблема 6: Cleanup WebSocket клиентов

**Проблема:** Устаревшие клиенты накапливались в памяти.

**Решение:** Добавил cleanup stale клиентов при broadcast:

```go
func (h *Hub) broadcastMessage(payload []byte, symbol string) {
    h.mu.RLock()
    var stale []*Client
    for client := range h.clients {
        select {
        case client.send <- payload:
        default:
            stale = append(stale, client)
        }
    }
    h.mu.RUnlock()

    if len(stale) > 0 {
        h.mu.Lock()
        for _, client := range stale {
            delete(h.clients, client)
            close(client.send)
        }
        h.mu.Unlock()
    }
}
```

---

## Глава 11: Итоги и выводы

### Что получилось

1. **Полноценная биржевая платформа** с реальным matching engine
2. **Микросервисная архитектура** с 7 независимыми сервисами
3. **React-интерфейс** в стиле Binance с графиками и книгой ордеров
4. **CI/CD пайплайн** с автоматическим тестированием
5. **Docker Compose** для запуска всей системы одной командой

### Цифры проекта

| Метрика | Значение |
|---------|----------|
| Строк Go | ~10,000 |
| Строк C++ | ~2,000 |
| Строк TypeScript | ~3,000 |
| Файлов | 94 |
| Коммитов | 20 |
| Контейнеров | 9 |
| Сервисов | 7 |
| Недель разработки | 3 |

### Ключевые уроки

1. **Архитектура важнее кода** — правильное разделение сервисов сэкономило часы переделок
2. **Тесты с первого дня** — моки и интерфейсы позволили тестировать изолированно
3. **CI/CD автоматизация** — каждый коммит автоматически проверяется
4. **Документация — часть продукта** — README, OpenAPI, архитектурные диаграммы
5. **C++ для критических компонентов** — matching engine на C++ обеспечивает нужную производительность

### Что можно улучшить

1. **Kubernetes** — текущий Docker Compose хорош для разработки, но для production нужен K8s
2. **Event-Driven Architecture** — заменить синхронные вызовы на события через Kafka/RabbitMQ
3. **Monitoring** — добавить Prometheus + Grafana
4. **Rate Limiting в Redis** — текущая реализация in-memory, для production нужен Redis
5. **HTTPS** — добавить TLS termination в API Gateway

---

## Заключение

Этот проект стал для меня отличным опытом создания production-grade системы с нуля. Я узнал, как проектировать микросервисную архитектуру, как писать высокопроизводительный matching engine на C++, как настраивать CI/CD и как создавать удобный пользовательский интерфейс для торговли.

Самое важное — это не просто учебный проект, а полноценная система, которую можно масштабировать и развивать. Каждый компонент продуман и протестирован, а документация позволяет другим разработчикам быстро включиться в проект.

Если вы думаете о создании подобного проекта — начинайте с архитектуры, пишите тесты с первого дня, и не бойтесь использовать разные языки программирования для разных компонентов. Это того стоит.
