#pragma once

#include "types.h"
#include <map>
#include <set>
#include <unordered_map>
#include <vector>
#include <mutex>
#include <functional>
#include <list>

namespace matching {

struct PriceLevel {
    double price;
    double total_qty;
    std::list<std::string> order_ids;

    explicit PriceLevel(double p) : price(p), total_qty(0) {}
};

using OnTrade = std::function<void(const Trade&)>;
using OnOrderUpdate = std::function<void(const Order&)>;

class OrderBook {
public:
    explicit OrderBook(const std::string& symbol);

    void submit(Order order);
    bool cancel(const std::string& order_id);

    std::vector<Trade> match(Order& incoming);

    double best_bid() const;
    double best_ask() const;
    double last_price() const;
    double spread() const;

    std::vector<std::pair<double, double>> bids(int depth = 20) const;
    std::vector<std::pair<double, double>> asks(int depth = 20) const;

    void on_trade(OnTrade cb) { on_trade_ = std::move(cb); }
    void on_order(OnOrderUpdate cb) { on_order_ = std::move(cb); }

private:
    std::string symbol_;
    std::map<double, PriceLevel, std::greater<double>> bids_;
    std::map<double, PriceLevel, std::less<double>> asks_;
    std::unordered_map<std::string, Order> orders_;
    std::unordered_map<std::string, double> remaining_;
    double last_price_ = 0;
    mutable std::mutex mu_;

    OnTrade on_trade_;
    OnOrderUpdate on_order_;

    std::vector<Trade> match_limit(Order& order);
    std::vector<Trade> match_market(Order& order);
    std::vector<Trade> match_against(Order& incoming, bool check_price);

    void insert_to_book(const Order& order);
    void remove_from_book(const std::string& order_id);
    static double calc_fee(double amount, bool maker);
};

}  // namespace matching
