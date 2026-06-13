#pragma once

#include "orderbook.h"
#include <unordered_map>
#include <memory>
#include <mutex>

namespace matching {

class MatchingEngine {
public:
    MatchingEngine();

    std::string submit(const Order& order);
    bool cancel(const std::string& order_id);

    double best_bid(const std::string& symbol) const;
    double best_ask(const std::string& symbol) const;
    double last_price(const std::string& symbol) const;
    std::vector<std::pair<double, double>> bids(const std::string& symbol, int depth = 20) const;
    std::vector<std::pair<double, double>> asks(const std::string& symbol, int depth = 20) const;

    void on_trade(OnTrade cb) { on_trade_ = std::move(cb); }
    void on_order(OnOrderUpdate cb) { on_order_ = std::move(cb); }

private:
    std::unordered_map<std::string, std::unique_ptr<OrderBook>> books_;
    mutable std::mutex mu_;
    OnTrade on_trade_;
    OnOrderUpdate on_order_;

    OrderBook& get_or_create(const std::string& symbol);
};

}  // namespace matching
