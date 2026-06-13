#include "engine.h"

namespace matching {

MatchingEngine::MatchingEngine() = default;

std::string MatchingEngine::submit(const Order& order) {
    std::lock_guard<std::mutex> lock(mu_);
    auto& book = get_or_create(order.symbol);
    Order mutable_order = order;
    book.submit(mutable_order);
    return mutable_order.id;
}

bool MatchingEngine::cancel(const std::string& order_id) {
    std::lock_guard<std::mutex> lock(mu_);
    for (auto& [sym, book] : books_) {
        if (book->cancel(order_id)) return true;
    }
    return false;
}

double MatchingEngine::best_bid(const std::string& symbol) const {
    std::lock_guard<std::mutex> lock(mu_);
    auto it = books_.find(symbol);
    return it != books_.end() ? it->second->best_bid() : 0;
}

double MatchingEngine::best_ask(const std::string& symbol) const {
    std::lock_guard<std::mutex> lock(mu_);
    auto it = books_.find(symbol);
    return it != books_.end() ? it->second->best_ask() : 0;
}

double MatchingEngine::last_price(const std::string& symbol) const {
    std::lock_guard<std::mutex> lock(mu_);
    auto it = books_.find(symbol);
    return it != books_.end() ? it->second->last_price() : 0;
}

std::vector<std::pair<double, double>> MatchingEngine::bids(const std::string& symbol, int depth) const {
    std::lock_guard<std::mutex> lock(mu_);
    auto it = books_.find(symbol);
    return it != books_.end() ? it->second->bids(depth) : std::vector<std::pair<double, double>>{};
}

std::vector<std::pair<double, double>> MatchingEngine::asks(const std::string& symbol, int depth) const {
    std::lock_guard<std::mutex> lock(mu_);
    auto it = books_.find(symbol);
    return it != books_.end() ? it->second->asks(depth) : std::vector<std::pair<double, double>>{};
}

OrderBook& MatchingEngine::get_or_create(const std::string& symbol) {
    auto it = books_.find(symbol);
    if (it != books_.end()) return *it->second;

    auto book = std::make_unique<OrderBook>(symbol);
    if (on_trade_) book->on_trade(on_trade_);
    if (on_order_) book->on_order(on_order_);

    auto [ins, _] = books_.emplace(symbol, std::move(book));
    return *ins->second;
}

}  // namespace matching
