zvynt
High-throughput stock trading backend in Go. WebSocket price feeds, order matching engine, async settlement via Kafka. Deployed on AWS EKS with Terraform and GitHub Actions CI/CD.

🚧 Under active development

Architecture

                         Internet
                            |
                       [AWS ALB]
                       /        \
                   /api/*       /ws
                     |            |
               [API Gateway]  [Market Data] ←── Kafka ←── [Simulator]
                     |            |
              ┌──────┼──────┐   gRPC
              |      |      |    |
           [Auth] [Orders] [Execution]
                     |      |
                     └──┬───┘
                        |
                   [PostgreSQL]
              Redis: LTP cache, order book, locks, rate limits
Services
Service	What it does
gateway	JWT validation, rate limiting, routes REST traffic
auth	Registration, login, token management
orders	Place/cancel orders, wallet, portfolio
marketdata	WebSocket price streaming, LTP cache, OHLCV candles
execution	Order book (Redis ZSET), matching with partial fills, settlement
simulator	Publishes realistic stock ticks to Kafka every 500ms
How It Works
Prices flow from the simulator through Kafka into the Market Data Service, which caches the latest price in Redis and fans it out to thousands of WebSocket connections using a goroutine-per-connection model.

Orders go through the gateway into the Orders Service. Funds are blocked using an atomic PostgreSQL update (row-level locking handles concurrency, no distributed lock needed here). The order is then published to Kafka for async processing.

Matching happens in the Execution Service. It maintains an order book per stock using Redis sorted sets. When a match is found, it runs a fill loop, popping the best price from the opposite side, settling each partial fill as an independent database transaction, and continuing until the order is fully filled or no matches remain. Distributed locks (Redis + TTL) protect against concurrent settlements for the same user.

Tech Stack
Backend	Go, Gin, gorilla/websocket, gRPC, Kafka, Redis, PostgreSQL
Auth	JWT (access + refresh tokens), Redis blacklisting
Infra	Docker, AWS EKS, ECR, RDS, Terraform
CI/CD	GitHub Actions (service-scoped pipelines)
Monitoring	AWS CloudWatch
Status
Service	Status
Auth	🔲 Not started
Gateway	🔲 Not started
Market Data	🔲 Not started
Orders	🔲 Not started
Execution	🔲 Not started
Simulator	🔲 Not started
Infra	🔲 Not started
CI/CD	🔲 Not started
Running Locally
bash

docker-compose up -d
make migrate-up
make seed
make run-simulator