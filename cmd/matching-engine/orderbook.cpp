#include "orderbook.h"
#include <algorithm>
#include <atomic>

namespace matching {

static std::atomic<uint64_t> trade_counter{0};

OrderBook::OrderBook(const std::string& symbol) : symbol_(symbol) {}

void OrderBook::submit(Order order) {
    std::lock_guard<std::mutex> lock(mu_);
    order.status = OrderStatus::OPEN;
    order.created_at = std::chrono::steady_clock::now();

    std::vector<Trade> trades;
    if (order.type == OrderType::LIMIT) {
        trades = match_limit(order);
    } else {
        trades = match_market(order);
    }

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
    } else if (order.remaining() <= 0) {
        order.status = OrderStatus::FILLED;
    }

    orders_[order.id] = order;

    if (on_order_) on_order_(order);
}

bool OrderBook::cancel(const std::string& order_id) {
    std::lock_guard<std::mutex> lock(mu_);

    auto it = orders_.find(order_id);
    if (it == orders_.end()) return false;

    Order& order = it->second;
    if (order.status != OrderStatus::OPEN && order.status != OrderStatus::PARTIAL) {
        return false;
    }

    remove_from_book(order_id);
    order.status = OrderStatus::CANCELLED;

    if (on_order_) on_order_(order);
    return true;
}

std::vector<Trade> OrderBook::match(Order& incoming) {
    std::lock_guard<std::mutex> lock(mu_);
    incoming.status = OrderStatus::OPEN;
    incoming.created_at = std::chrono::steady_clock::now();

    auto trades = (incoming.type == OrderType::LIMIT) ? match_limit(incoming) : match_market(incoming);

    if (incoming.remaining() > 0) {
        if (incoming.type == OrderType::MARKET && incoming.tif == TimeInForce::FOK) {
            incoming.status = OrderStatus::CANCELLED;
        } else if (incoming.tif == TimeInForce::IOC) {
            incoming.status = OrderStatus::CANCELLED;
        } else {
            insert_to_book(incoming);
            orders_[incoming.id] = incoming;
            remaining_[incoming.id] = incoming.remaining();
        }
    } else {
        incoming.status = OrderStatus::FILLED;
    }

    orders_[incoming.id] = incoming;
    if (on_order_) on_order_(incoming);
    return trades;
}

double OrderBook::best_bid() const {
    std::lock_guard<std::mutex> lock(mu_);
    return bids_.empty() ? 0 : bids_.begin()->first;
}

double OrderBook::best_ask() const {
    std::lock_guard<std::mutex> lock(mu_);
    return asks_.empty() ? 0 : asks_.begin()->first;
}

double OrderBook::last_price() const {
    std::lock_guard<std::mutex> lock(mu_);
    return last_price_;
}

double OrderBook::spread() const {
    std::lock_guard<std::mutex> lock(mu_);
    if (bids_.empty() || asks_.empty()) return 0;
    return asks_.begin()->first - bids_.begin()->first;
}

std::vector<std::pair<double, double>> OrderBook::bids(int depth) const {
    std::lock_guard<std::mutex> lock(mu_);
    std::vector<std::pair<double, double>> result;
    int n = 0;
    for (auto& [price, level] : bids_) {
        if (n++ >= depth) break;
        result.emplace_back(price, level.total_qty);
    }
    return result;
}

std::vector<std::pair<double, double>> OrderBook::asks(int depth) const {
    std::lock_guard<std::mutex> lock(mu_);
    std::vector<std::pair<double, double>> result;
    int n = 0;
    for (auto& [price, level] : asks_) {
        if (n++ >= depth) break;
        result.emplace_back(price, level.total_qty);
    }
    return result;
}

std::vector<Trade> OrderBook::match_limit(Order& order) {
    return match_against(order, true);
}

std::vector<Trade> OrderBook::match_market(Order& order) {
    return match_against(order, false);
}

std::vector<Trade> OrderBook::match_against(Order& incoming, bool check_price) {
    std::vector<Trade> trades;
    auto& opposite = incoming.is_buy() ? asks_ : bids_;

    for (auto it = opposite.begin(); it != opposite.end() && incoming.remaining() > 0;) {
        auto& level = it->second;

        if (check_price) {
            if (incoming.is_buy() && incoming.price < level.price) break;
            if (!incoming.is_buy() && incoming.price > level.price) break;
        }

        for (auto oid = level.order_ids.begin(); oid != level.order_ids.end() && incoming.remaining() > 0;) {
            auto oit = orders_.find(*oid);
            if (oit == orders_.end()) {
                oid = level.order_ids.erase(oid);
                continue;
            }

            Order& counter = oit->second;
            double avail = remaining_[*oid];
            double qty = std::min(incoming.remaining(), avail);

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

            double amount = trade.price * trade.quantity;
            trade.buyer_fee = calc_fee(amount, true);
            trade.seller_fee = calc_fee(amount, false);
            trade.executed_at = std::chrono::steady_clock::now();

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
    return trades;
}

void OrderBook::insert_to_book(const Order& order) {
    auto& book = order.is_buy() ? bids_ : asks_;
    auto it = book.find(order.price);
    if (it == book.end()) {
        PriceLevel level(order.price);
        level.total_qty = order.remaining();
        level.order_ids.push_back(order.id);
        book.emplace(order.price, std::move(level));
    } else {
        it->second.total_qty += order.remaining();
        it->second.order_ids.push_back(order.id);
    }
}

void OrderBook::remove_from_book(const std::string& order_id) {
    auto oit = orders_.find(order_id);
    if (oit == orders_.end()) return;

    const Order& order = oit->second;
    auto& book = order.is_buy() ? bids_ : asks_;

    auto lit = book.find(order.price);
    if (lit != book.end()) {
        lit->second.total_qty -= remaining_[order_id];
        lit->second.order_ids.remove(order_id);
        if (lit->second.total_qty <= 0 || lit->second.order_ids.empty()) {
            book.erase(lit);
        }
    }

    orders_.erase(order_id);
    remaining_.erase(order_id);
}

double OrderBook::calc_fee(double amount, bool maker) {
    return amount * (maker ? 0.001 : 0.002);
}

}  // namespace matching
