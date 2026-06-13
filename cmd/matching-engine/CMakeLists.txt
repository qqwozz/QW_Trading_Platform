cmake_minimum_required(VERSION 3.14)
project(matching_engine LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_executable(matching-engine
    main.cpp
    orderbook.cpp
    engine.cpp
)

target_compile_options(matching-engine PRIVATE -Wall -Wextra -O2)
