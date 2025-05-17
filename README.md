Below is a **draft README.md** you can drop into the root of the repository.
It gives new contributors a quick mental model of what the ledger-system does, how the folders fit together, and the fastest way to spin it up locally.

# Ledger-System

A lightweight, event-driven ledger written in Go.  
It keeps double-entry account balances in **PostgreSQL**, streams every state-change out on **Kafka**, and exposes a small **REST + OpenAPI** layer for external clients.  
MongoDB is used for side-channel look-ups that don’t belong in the immutable ledger tables (e.g., id-maps, metadata).

---

## Features

- **Double-entry accounting**: Ensures accurate and consistent account balances.
- **Event-driven architecture**: Uses Kafka for asynchronous processing.
- **REST API**: Exposes endpoints for account operations with OpenAPI documentation.
- **Multi-database support**: PostgreSQL for core ledger operations and MongoDB for metadata.
- **Scalability**: Designed to handle high transaction volumes.

---

## High-level Boot Flow

```plaintext
main.go
 ├─ config.Initialize()            # Read .env / defaults
 ├─ pg.InitPostgres()              # Open & ping Postgres
 ├─ mongo.InitMongo()              # Open & ping Mongo
 ├─ kafka.InitKafka()              # Create producer / consumer
 ├─ kafka.CreateTopics()           # Idempotent topic bootstrap
 ├─ service.Initialize(ctx)        # Wire repositories + start async consumers
 └─ http.ListenAndServe()          # Expose REST API
```

---

## Folder Map

| Path                   | Description                                            | Typical Use Case          |
| ---------------------- | ------------------------------------------------------ | ------------------------- |
| **api**                | Gin/Chi handlers, request DTOs, validation glue        | Add a new endpoint        |
| **service**            | Business rules – create account, post entry, transfer  | Extend domain logic       |
| **pg**                 | `sqlc`‐ or hand-rolled queries, migrations, Tx helpers | Change schema             |
| **mongo**              | Thin wrapper around `mongo.Client`                     | Add secondary indexes     |
| **kafka**              | Producer / consumer helpers, topic constants           | Tune partitions / acks    |
| **config**             | Viper-style env loading, strongly-typed config         | Add new env var           |
| **cmd**                | One-off CLIs (back-fill, repair, etc.)                 | Run batch jobs            |
| **utils**              | Generic helpers (error types, UUID, logging)           | Shared helpers            |
| **api.yaml**           | OpenAPI 3 spec – source of truth for the HTTP API      | Regenerate client         |
| **docker-compose.yml** | Local dev stack: Postgres, Mongo, Redpanda/Kafka       | `docker-compose up`       |

---

## Quick Start (Local Development)

Follow these steps to set up the project locally:

### Prerequisites

- **Docker** and **Docker Compose** installed.
- **Go** (version 1.20 or higher) installed.

### Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/shreetheja/ledger-system.git && cd ledger-system
   ```

2. **Spin up the infrastructure**:
   ```bash
   docker-compose up -d    # Postgres, Mongo, Kafka, Swagger
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

4. **Access the API documentation**:
   Open [http://localhost:8080/docs](http://localhost:8080/docs) to view and test the API using Swagger UI.

---

## REST API Endpoints

### Base URL: `http://localhost:8080`

| Endpoint              | Method | Description                     |
| --------------------- | ------ | ------------------------------- |
| `/balance`            | GET    | Retrieve the balance of a user  |
| `/balance`            | POST   | Create a new account            |
| `/balance/add`        | POST   | Add funds to an account         |
| `/balance/deduct`     | POST   | Deduct funds from an account    |
| `/logs`               | GET    | Logs of particular account      |

Refer to the OpenAPI Specification (`api.yaml`) for detailed request/response formats.

---

## Environment Variables

The application uses the following environment variables (defined in `.env`):

| Variable              | Description                     | Default Value         |
| --------------------- | ------------------------------- | --------------------- |
| `PORT`                | Port for the REST API           | `8080`               |
| `POSTGRES_USER`       | PostgreSQL username             | `ledger_user`        |
| `POSTGRES_PASSWORD`   | PostgreSQL password             | `ledger_pass`        |
| `POSTGRES_DB`         | PostgreSQL database name        | `ledger_db`          |
| `POSTGRES_HOST`       | PostgreSQL host                 | `localhost`          |
| `POSTGRES_PORT`       | PostgreSQL port                 | `5432`               |
| `MONGO_HOST`          | MongoDB host                   | `localhost`          |
| `MONGO_PORT`          | MongoDB port                   | `27017`              |
| `MONGO_DB`            | MongoDB database name          | `ledger_tx_log`      |
| `KAFKA_BROKER`        | Kafka broker address           | `localhost:9092`     |
| `KAFKA_CLUSTER_ID`    | Kafka cluster ID               | `kraft-cluster-1234` |

---

## Testing

The project includes a comprehensive testing strategy:

1. **Unit Tests**: Test individual components and business logic.
2. **Integration Tests**: Validate interactions between PostgreSQL, MongoDB, and Kafka.
3. **Mocking**: Use mocks for external dependencies to ensure robust validation.

Run tests using:
```bash
go test ./...
```

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request with a detailed description of your changes.

---

## License

This project is licensed under the MIT License. See the LICENSE file for details.

---

## Contact

For questions or support, please reach out to [shreetheja](mailto:shreetheja@example.com).
