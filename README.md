# Subscription Service

REST API service for managing user subscriptions with PostgreSQL database.

## Features

- CRUDL operations for subscriptions
- Aggregation of subscription costs by period
- PostgreSQL database with migrations
- Swagger documentation
- Docker Compose support
- Logging with logrus
- Configuration via .env file

## Requirements

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL 15+

## Installation

### Using Docker Compose (Recommended)

1. Clone the repository
2. Run `docker-compose up --build`
3. The service will be available at `http://localhost:8080`
4. Swagger documentation at `http://localhost:8080/swagger/index.html`

### Local Development

1. Install dependencies: `go mod download`
2. Set up PostgreSQL database
3. Create `.env` file with your configuration
4. Run migrations manually if needed
5. Start the server: `go run cmd/server/main.go`

## API Endpoints

- `POST /subscriptions` - Create a new subscription
- `GET /subscriptions/{id}` - Get subscription by ID
- `GET /subscriptions` - List all subscriptions
- `PUT /subscriptions/{id}` - Update subscription
- `DELETE /subscriptions/{id}` - Delete subscription
- `GET /subscriptions/aggregate` - Aggregate total price by filters

## Example Request

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
```

## Configuration

Edit `.env` file for configuration:

```env
SERVER_PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscription_db
DB_SSL_MODE=disable
LOG_LEVEL=debug
```

## Database Migrations

Migrations are located in the `migrations/` directory and are automatically applied when using Docker Compose.

## Swagger Documentation

Swagger documentation is available at `/swagger/index.html` when the service is running.

## Project Structure

```
.
├── cmd/
│   └── server/          # Main server application
├── internal/
│   ├── config/          # Configuration
│   ├── database/        # Database connection
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   ├── repositories/    # Database repositories
│   └── services/        # Business logic
├── migrations/          # Database migrations
├── docs/                # Swagger documentation
├── .env                 # Environment variables
└── docker-compose.yml   # Docker Compose configuration
```

## Testing

The service includes comprehensive logging and error handling. Check logs for debugging information.

## License

MIT