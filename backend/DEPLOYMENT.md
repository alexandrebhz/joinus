# Deployment Guide

## Environment Variables

The application supports both `DATABASE_URL` and individual database parameters:

### Using DATABASE_URL (Recommended for Production)

```bash
DATABASE_URL=postgresql://user:password@host:port/dbname?sslmode=require
```

### Using Individual Parameters

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=startup_board
DB_SSL_MODE=disable
```

**Note:** If `DATABASE_URL` is set, it takes precedence over individual parameters.

## Running with Docker

### Using Docker Compose (Recommended)

The docker-compose files automatically handle `DATABASE_URL`:

```bash
# Production
docker-compose up -d

# Development
docker-compose -f docker-compose.dev.yml up -d
```

### Running Container Directly

When running the container directly (e.g., on Railway, Render, Heroku), pass `DATABASE_URL`:

```bash
docker run -d \
  -p 8080:8080 \
  -e DATABASE_URL=postgresql://user:password@host:port/dbname?sslmode=require \
  -e JWT_SECRET=your-secret-key \
  -e RESEND_API_KEY=your-resend-key \
  -e APP_URL=https://yourdomain.com \
  startup-job-board-api
```

### Platform-Specific Deployment

Most platforms (Railway, Render, Heroku, etc.) automatically provide `DATABASE_URL` as an environment variable. Just ensure:

1. Your database service is provisioned
2. `DATABASE_URL` is set in your platform's environment variables
3. The container receives the environment variable (usually automatic)

The application will automatically use `DATABASE_URL` if available, otherwise it falls back to individual parameters.

