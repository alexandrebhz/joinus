# Docker Setup Guide

## Overview

Docker Compose files are now located at the root of the project for better organization. This allows managing both backend and frontend services from a single location.

## File Structure

```
joinus/
├── docker-compose.yml          # Production configuration
├── docker-compose.dev.yml      # Development configuration
├── backend/
│   └── docker/
│       └── Dockerfile          # Backend production Dockerfile
└── frontend/
    ├── Dockerfile              # Frontend production Dockerfile
    └── Dockerfile.dev          # Frontend development Dockerfile
```

## Development Setup

### Start all services:

```bash
# From the root directory (joinus/)
docker-compose -f docker-compose.dev.yml up -d
```

### Start specific services:

```bash
# Only backend services
docker-compose -f docker-compose.dev.yml up postgres minio api

# Only frontend
docker-compose -f docker-compose.dev.yml up frontend
```

### View logs:

```bash
# All services
docker-compose -f docker-compose.dev.yml logs -f

# Specific service
docker-compose -f docker-compose.dev.yml logs -f frontend
```

### Stop services:

```bash
docker-compose -f docker-compose.dev.yml down
```

### Rebuild services:

```bash
docker-compose -f docker-compose.dev.yml up -d --build
```

## Production Setup

### Start production services:

```bash
docker-compose up -d
```

### Environment Variables

Create a `.env` file at the root with your production settings:

```env
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=startup_board
JWT_SECRET=your-very-secure-secret-key
NEXT_PUBLIC_API_URL=http://your-api-domain.com
```

## Services

### Development (`docker-compose.dev.yml`)

- **postgres**: PostgreSQL 15 database (port 5432)
- **minio**: MinIO object storage (ports 9000, 9001)
- **api**: Backend API (port 8080)
- **frontend**: Next.js frontend (port 3000)

### Production (`docker-compose.yml`)

- **postgres**: PostgreSQL 15 database
- **minio**: MinIO object storage
- **api**: Backend API
- **frontend**: Next.js frontend (production build)

## Troubleshooting

### Frontend build fails with "package-lock.json not found"

The Dockerfiles now use `npm install` instead of `npm ci` to handle cases where package-lock.json doesn't exist. After the first successful build, a package-lock.json will be generated.

### Port conflicts

If ports are already in use, you can modify the port mappings in docker-compose files:

```yaml
ports:
  - "3001:3000"  # Use port 3001 instead of 3000
```

### Volume issues

To reset volumes (this will delete all data):

```bash
docker-compose -f docker-compose.dev.yml down -v
```

### Rebuild from scratch

```bash
# Stop and remove containers, networks, volumes
docker-compose -f docker-compose.dev.yml down -v

# Remove images
docker-compose -f docker-compose.dev.yml rm -f

# Rebuild
docker-compose -f docker-compose.dev.yml up -d --build
```

## Development Workflow

1. **Start services**: `docker-compose -f docker-compose.dev.yml up -d`
2. **Make code changes**: Files are mounted as volumes, so changes are reflected immediately
3. **View logs**: `docker-compose -f docker-compose.dev.yml logs -f`
4. **Stop services**: `docker-compose -f docker-compose.dev.yml down`

## Notes

- Development services use volume mounts for hot-reloading
- Production services use optimized multi-stage builds
- Frontend connects to API at `http://api:8080` in Docker network
- For local development outside Docker, use `http://localhost:8080`

