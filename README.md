# MeGo (RestApi-Go)

A lightweight and efficient REST API backend built with Go, featuring user management, JWT authentication, and PostgreSQL integration.

## Features

- RESTful API architecture using the Chi router.
- Secure user registration and login with password hashing.
- JWT-based authentication using HTTP-only cookies.
- Type-safe database interactions with SQLC.
- PostgreSQL database integration.
- Environment variable management with godotenv.
- Graceful server shutdown and health checks.
- Comprehensive unit testing suite.

## Tech Stack

- Language: Go (1.25.3+)
- Router: Chi
- Database: PostgreSQL
- Query Generator: SQLC
- Authentication: JWT (github.com/golang-jwt/jwt/v4)
- Others: godotenv, uuid, crypto/bcrypt

## Prerequisites

- Go 1.25.3 or higher.
- PostgreSQL database instance.
- SQLC (optional, for regenerating database code).

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Blue-Onion/RestApi-Go.git
   cd MeGo
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Configuration

Create a `.env` file in the root directory and configure the following variables:

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
JWT_SECERT=your_super_secret_key
```

## Running the Application

You can start the server using the provided Makefile or directly with Go:

Using Makefile:
```bash
make run
```

Using Go:
```bash
go run cmd/main.go
```

The server will start listening on the port specified in your `.env` file (default: http://localhost:8080).

## API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|--------------|
| GET | /health | Server health check | No |
| POST | /api/users | Register a new user | No |
| POST | /api/login | User login and JWT issuance | No |
| POST | /api/logOut | User logout and cookie clear | Yes |

## Project Structure

- `cmd/`: Application entry point.
- `config/`: Configuration loading and database connection.
- `handler/`: HTTP handlers for various routes.
- `internal/database/`: Auto-generated database code by SQLC.
- `middleware/`: Authentication and other middlewares.
- `model/`: Data models and structures.
- `sql/`: SQL schema and queries.
- `test/`: Automated test suite.
- `utlis/`: Utility functions (JWT, hashing, etc.).

## Testing

To run the automated tests, use the following command:

```bash
go test ./test/...
```
