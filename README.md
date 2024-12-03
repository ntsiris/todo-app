# Todo App

A lightweight and efficient To-Do application built in Go, designed to manage tasks with PostgreSQL integration.

## Features

- **Task Management**: Add, update, delete, and view tasks.
- **Environment Configuration**: Uses `dotenv` for environment variables.
- **PostgreSQL Database**: Implements `pgx` for efficient database interactions.
- **Dockerized Deployment**: Simplified containerization using Docker.

## Requirements

- Go 1.22 or later (for local development)
- PostgreSQL database
- [godotenv](https://github.com/joho/godotenv) for environment variable management
- [pgx](https://github.com/jackc/pgx) for database interactions
- Docker (optional, for containerized deployment)

## Installation

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/ntsiris/todo-app.git
   cd todo-app

2. Install dependencies
    ```bash
    go mod tidy
    ```

3. Set up your `.env` file (optional)
4. Run The application

### Using Docker


1. Build the docker image
    ```bash
    docker build -t todo-app . 
    ```
2. Start the service using docker compose

    ```bash
    docker-compose up
    ```
   
## Configuration

Ensure the `.env` file includes your database connection details:
```dotenv
    DB_HOST=db
    DB_PORT=5432
    DB_USER=your_user
    DB_PASSWORD=your_password
    DB_NAME=todo_db
```
