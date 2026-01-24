# Contributing to the Finitech Ecosystem

First off, thank you for considering contributing to this project! It's people like you that make the open source community such an amazing place to learn, inspire, and create.

## 🤝 How to Contribute

### 1. Reporting Bugs
- **Ensure the bug was not already reported** by searching on GitHub under [Issues].
- If you're unable to find an open issue addressing the problem, [open a new one]. Be sure to include a **title and clear description**, as well as as much relevant information as possible, and a **code sample** or an **executable test case** demonstrating the expected behavior that is not occurring.

### 2. Suggesting Enhancements
- Open a new issue with the label `enhancement`.
- Explain why this enhancement would be useful to most users.

### 3. Pull Requests
1.  Fork the repo and create your branch from `main`.
2.  If you've added code that should be tested, add tests.
3.  If you've changed APIs, update the documentation.
4.  Ensure the test suite passes (`make test`).
5.  Make sure your code lints (`golangci-lint run`).
6.  Issue that pull request!

## 🛠️ Development Guide

### Prerequisites
- **Go** 1.24+
- **Docker** & **Docker Compose**
- **Make**
- **Protoc** (if modifying `.proto` files)

### Local Setup
```bash
# Clone the repo
git clone https://github.com/your-org/microservices.git
cd microservices

# Start infrastructure (Postgres, Redis, Kafka, RabbitMQ)
docker-compose up -d

# Build and run services
make build
./bin/gateway
```

### Protocol Buffers
If you modify any file in `proto/`, you must regenerate the Go code:
```bash
make proto
```

## 🎨 Coding Standards

- We follow the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).
- All exported functions and methods should have comments.
- Run `go fmt` before committing.

## ⚖️ License
By contributing, you agree that your contributions will be licensed under its MIT License.
