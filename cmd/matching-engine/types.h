#pragma once

#include <cstdint>
#include <string>
#include <chrono>

namespace matching {

enum class OrderSide : uint8_t { BUY = 0, SELL = 1 };
enum class OrderType : uint8_t { LIMIT = 0, MARKET = 1 };
enum class OrderStatus : uint8_t { PENDING = 0, OPEN = 1, PARTIAL = 2, FILLED = 3, CANCELLED = 4 };
enum class TimeInForce : uint8_t { GTC = 0, IOC = 1, FOK = 2 };

struct Order {
    std::string id;
    std::string user_id;
    std::string symbol;
    OrderSide side;
    OrderType type;
    double price = 0;
    double quantity = 0;
    double filled = 0;
    OrderStatus status = OrderStatus::PENDING;
    TimeInForce tif = TimeInForce::GTC;
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
    double price = 0;
    double quantity = 0;
    double buyer_fee = 0;
    double seller_fee = 0;
    std::chrono::steady_clock::time_point executed_at;
};

}  // namespace matching
