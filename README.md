# Microservices Fintech Ecosystem

A robust, developer-first microservices platform for financial operations, featuring a secure API Gateway, distributed services, and a Stripe-inspired CLI for seamless local development.

## 🚀 Architecture

The system is built on a distributed microservices architecture, leveraging **Docker**, **Go**, and **Redis** for high performance and scalability.

- **API Gateway**: The entry point. Handles rate limiting (Redis-backed), API key verification, and WebSocket-based webhook streaming.
- **Auth Service**: Manages user registration, JWT-based sessions, and secure API key management (bcrypt-hashed).
- **Payments Service**: Orchestrates payment logic with built-in idempotency and publishes real-time events via Redis Pub/Sub.
- **Ledger Service**: A double-entry accounting system preserving financial integrity with atomic transaction recording.
- **Micro CLI**: A professional developer tool for authentication and local webhook relay.

## 🛠️ Tech Stack

- **Backend**: Go 1.24+
- **Database**: PostgreSQL (Service-specific isolation)
- **Caching/PubSub**: Redis
- **Containerization**: Docker & Docker Compose
- **CLI**: Cobra & Viper

## 📦 Getting Started

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go 1.24+](https://go.dev/dl/)

### 1. Launch the Ecosystem

```bash
docker-compose up --build -d
```
All services will initialize, migrate their schemas, and become active on their respective ports. The Gateway is exposed at `:8080`.

### 2. Build the CLI

```bash
go build -o micro ./cmd/cli
```

### 3. Developer Authentication

Sign up and log in directly from your terminal:

```bash
# Register via Gateway
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "dev@example.com", "password": "securepassword"}'

# Authenticate CLI
./micro login
```
The CLI automatically handles session persistence and generates a development API key stored in `~/.micro.yaml`.

## ⚡ Real-time Webhook Testing

Test your local webhooks without complex tunnel setups:

1. **Start Listening**:
   ```bash
   ./micro listen --forward-to http://localhost:4242/webhook
   ```
2. **Trigger an Event**: Create and confirm a payment.
3. **Relay**: The system streams the `payment.succeeded` event from the cloud (Docker) to your terminal via WebSockets, then forwards it to your local server.

## 🔒 Security & Operations

- **API Key Security**: Gateway enforces hashing and service-level validation.
- **Rate Limiting**: Fixed-window rate limiting (100 req/min) per API key.
- **Idempotency**: Payments service ensures safe retries via `Idempotency-Key` headers.
- **Double-Entry Ledger**: Immutable financial logs ensuring debits and credits always balance.

## 📜 License
MIT
