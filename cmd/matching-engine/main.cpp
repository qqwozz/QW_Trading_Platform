#include "engine.h"
#include <iostream>
#include <signal.h>
#include <atomic>
#include <thread>

static std::atomic<bool> running{true};

static void signal_handler(int) { running = false; }

int main(int argc, char** argv) {
    signal(SIGINT, signal_handler);
    signal(SIGTERM, signal_handler);

    matching::MatchingEngine engine;

    engine.on_trade([](const matching::Trade& t) {
        std::cout << "[TRADE] " << t.symbol
                  << " price=" << t.price
                  << " qty=" << t.quantity
                  << " buyer=" << t.buyer_id
                  << " seller=" << t.seller_id << "\n";
    });

    engine.on_order([](const matching::Order& o) {
        std::cout << "[ORDER] " << o.id
                  << " " << (o.is_buy() ? "BUY" : "SELL")
                  << " " << o.symbol
                  << " status=" << static_cast<int>(o.status) << "\n";
    });

    std::string port = argc > 1 ? argv[1] : "50051";
    std::cout << "Matching Engine starting on port " << port << "\n";

    matching::Order buy;
    buy.id = "order-1";
    buy.user_id = "user-1";
    buy.symbol = "BTC/USDT";
    buy.side = matching::OrderSide::BUY;
    buy.type = matching::OrderType::LIMIT;
    buy.price = 50000;
    buy.quantity = 1.0;
    engine.submit(buy);

    matching::Order sell;
    sell.id = "order-2";
    sell.user_id = "user-2";
    sell.symbol = "BTC/USDT";
    sell.side = matching::OrderSide::SELL;
    sell.type = matching::OrderType::LIMIT;
    sell.price = 50000;
    sell.quantity = 0.5;
    engine.submit(sell);

    while (running) {
        std::this_thread::sleep_for(std::chrono::milliseconds(100));
    }

    std::cout << "Matching Engine stopped\n";
    return 0;
}
