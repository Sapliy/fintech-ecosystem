# 💳 Microservices Fintech Ecosystem

A robust, developer-first microservices platform for financial operations, featuring a secure API Gateway, distributed services, and a Stripe-inspired CLI for seamless local development.

## 🚀 Architecture

The system is built on a distributed microservices architecture, leveraging **Go**, **Redis**, **PostgreSQL**, and **Docker**.

```mermaid
graph TD
    User((User/Dev)) -->|HTTP/WS| Gateway[API Gateway :8080]
    Gateway -->|Auth Check| AuthService[Auth Service :8081]
    Gateway -->|Proxy| PaymentsService[Payments Service :8082]
    Gateway -->|Proxy| LedgerService[Ledger Service :8083]
    
    PaymentsService -->|Events| Redis[(Redis Pub/Sub)]
    Redis -->|Stream| Gateway
    Gateway -->|WebSocket| CLI[Micro CLI]
    
    AuthService --- AuthDB[(Postgres: Auth)]
    PaymentsService --- PaymentsDB[(Postgres: Payments)]
    LedgerService --- LedgerDB[(Postgres: Ledger)]
```

- **API Gateway**: Central hub for rate limiting (100 req/min), API key validation, and WebSocket-based webhook streaming.
- **Auth Service**: Manages user identity, JWT sessions, and hashed API keys.
- **Payments Service**: Handles payment orchestration with idempotency support.
- **Ledger Service**: Immutable double-entry accounting system for financial integrity.
- **Micro CLI**: Professional tool for local webhook relay and session management.

## 🛠️ Tech Stack

- **Languge**: Go 1.24+
- **Infrastructure**: Docker, Docker Compose
- **Databases**: PostgreSQL (Isolated per service)
- **Messaging/Cache**: Redis
- **CLI Framework**: Cobra & Viper

---

## 📦 Getting Started

### 1. Launch the Ecosystem

```bash
docker-compose up --build -d
```
The Gateway is exposed at `http://localhost:8080`.

### 2. Install Micro CLI

```bash
go build -o micro ./cmd/cli
```

### 3. Developer Onboarding

```bash
# Register a new developer account
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "dev@example.com", "password": "securepassword"}'

# Login via CLI
./micro login
```

---

## 📖 API Reference

### 🔐 Auth Service
Requires session JWT for key management.

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/auth/register` | `POST` | Create a new developer account |
| `/auth/login` | `POST` | Authenticate and receive JWT |
| `/auth/api_keys` | `POST` | Generate a new API key (`test` or `live`) |

**Example API Key Generation:**
```bash
curl -X POST http://localhost:8080/auth/api_keys \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{"environment": "test"}'
```

### 💸 Payments Service
Requires `Authorization: Bearer sk_<key>` header.

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/payments/payment_intents` | `POST` | Create a new payment intent |
| `/payments/payment_intents/{id}/confirm` | `POST` | Confirm a payment intent |

**Create Payment (with Idempotency):**
```bash
curl -X POST http://localhost:8080/payments/payment_intents \
  -H "Authorization: Bearer sk_your_key" \
  -H "Idempotency-Key: unique_req_123" \
  -d '{"amount": 5000, "currency": "usd"}'
```

### 📒 Ledger Service
Requires `Authorization: Bearer sk_<key>` header.

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/ledger/accounts` | `POST` | Create a financial account |
| `/ledger/accounts/{id}` | `GET` | View account balance and type |
| `/ledger/transactions` | `POST` | Record a double-entry transaction |

**Record Balanced Transaction:**
```json
{
  "reference_id": "pay_987",
  "description": "Payment for Service",
  "entries": [
    {"account_id": "acc_revenue", "amount": 100, "direction": "credit"},
    {"account_id": "acc_assets", "amount": -100, "direction": "debit"}
  ]
}
```

---

## ⚡ Real-time Webhooks

Test local integrations without tunnels:

1. **Listen for events**:
   ```bash
   ./micro listen --forward-to http://localhost:4242/webhook
   ```
2. **Trigger**: When a payment is confirmed via `/payments/payment_intents/{id}/confirm`, the Gateway pushes the event through the WebSocket connection to your CLI.

---

## 💻 Developer Guide

### Makefile Commands
- `make build`: Build all services locally.
- `make test`: Run the test suite.
- `make clean`: Remove binaries and temporary files.

### Database Access
Each service has its own database. Connect via:
- `psql postgres://user:password@localhost:5433/microservices` (Auth)
- `psql postgres://user:password@localhost:5434/payments` (Payments)
- `psql postgres://user:password@localhost:5435/ledger` (Ledger)

---

## 📜 License
MIT
