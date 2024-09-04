# Scheduling Api

This application is a Go-based backend service for managing schedules, users, and groups. It provides functionality for user authentication, group management, and shift scheduling.

## Features

- User registration and authentication
- Group creation and management
- Shift scheduling and management
- User-group membership management
- JWT-based authentication
- PostgreSQL database for data persistence

## Prerequisites

- Docker and Docker Compose (for running with Docker)
- Go 1.16 or later (for running locally)
- PostgreSQL (for running locally)

## Getting Started

### Running with Docker

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/scheduling-app.git
   cd scheduling-app
   ```

2. Build and run the application using Docker Compose:
   ```
   docker-compose up --build
   ```

   This command will build the Docker image for the application, start the application container, and a PostgreSQL container.

3. The application will be available at `http://localhost:8080`

### Running Locally

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/scheduling-app.git
   cd scheduling-app
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set up the PostgreSQL database and update the connection details in `db/connection.go`

4. Run database migrations:
   ```
   go run cmd/server/main.go migrate
   ```

5. Start the application:
   ```
   go run cmd/server/main.go serve
   ```

6. The application will be available at `http://localhost:8080`

## Project Structure

```
├── api/
│   ├── errors/
│   ├── handlers/
│   ├── middleware/
│   └── routers/
├── cmd/
│   └── server/
├── db/
│   ├── migrations/
│   ├── models/
│   └── queries/
├── internals/
│   └── auth/
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```
