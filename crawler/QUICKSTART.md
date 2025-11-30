# Quick Start Guide

## Prerequisites

- Go 1.23+
- PostgreSQL 15+
- Docker and Docker Compose (optional)

## Local Development Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Set Up Environment Variables

Create a `.env` file:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/crawler?sslmode=disable
BACKEND_URL=http://localhost:8080
BACKEND_TOKEN=your-backend-api-token
ENVIRONMENT=development
```

### 3. Set Up Database

#### Option A: Using Docker Compose

```bash
docker-compose up -d postgres
```

#### Option B: Local PostgreSQL

Create a database:

```sql
CREATE DATABASE crawler;
```

Run migrations:

```bash
make migrate
# or manually:
psql $DATABASE_URL -f migrations/001_initial_schema.sql
```

### 4. Run the Application

```bash
make run
# or
go run ./cmd/crawler
```

The service will start on `http://localhost:8080`

## Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f crawler

# Stop services
docker-compose down
```

## API Endpoints

### Health Check
```bash
GET /health
```

### Site Management

#### Create a Crawl Site
```bash
POST /api/v1/sites
Content-Type: application/json

{
  "name": "Example Job Board",
  "base_url": "https://example.com/jobs",
  "backend_startup_id": "uuid-here",
  "schedule": "0 0 * * *",
  "crawl_interval": "daily",
  "pagination_config": {
    "type": "query_param",
    "param_name": "page",
    "start_page": 1,
    "increment": 1,
    "max_pages": 10
  },
  "extraction_rules": {
    "job_list_selector": ".job-item",
    "job_detail_url": {
      "type": "relative",
      "selector": "a.job-link",
      "attribute": "href",
      "base_url": "https://example.com"
    },
    "fields": {
      "title": {
        "selector": "h2.job-title",
        "type": "text",
        "required": true
      },
      "description": {
        "selector": ".job-description",
        "type": "html",
        "required": true
      }
    }
  },
  "deduplication_key": "url",
  "request_delay": 2
}
```

#### List Sites
```bash
GET /api/v1/sites
```

#### Get Site
```bash
GET /api/v1/sites/:id
```

#### Update Site
```bash
PUT /api/v1/sites/:id
Content-Type: application/json

{
  "active": false
}
```

#### Delete Site
```bash
DELETE /api/v1/sites/:id
```

#### Execute Manual Crawl
```bash
POST /api/v1/sites/:id/crawl
```

## Example Site Configuration

### Query Parameter Pagination
```json
{
  "type": "query_param",
  "param_name": "page",
  "start_page": 1,
  "increment": 1,
  "max_pages": 50
}
```

### URL Pattern Pagination
```json
{
  "type": "url_pattern",
  "pattern": "/jobs/page/{page}",
  "start_page": 1,
  "increment": 1,
  "max_pages": 50
}
```

### Link Following Pagination
```json
{
  "type": "link_follow",
  "next_page_selector": "a.pagination-next",
  "max_depth": 10
}
```

## Scheduling

Sites are automatically crawled based on their cron schedule. The scheduler runs in the background and executes crawls according to each site's configuration.

Example schedules:
- Daily at midnight: `0 0 * * *`
- Every 6 hours: `0 */6 * * *`
- Weekly (Monday 9 AM): `0 9 * * 1`
- Every 30 minutes: `*/30 * * * *`

## Synchronization

Crawled jobs are automatically synced to the backend API every 5 minutes. The sync service:
1. Fetches unsynced jobs from the local database
2. Maps them to backend job DTOs
3. Sends them to the backend API
4. Marks jobs as synced upon success

## Troubleshooting

### Database Connection Issues
- Verify PostgreSQL is running
- Check `DATABASE_URL` in `.env`
- Ensure database exists

### Crawl Failures
- Check site configuration (selectors, pagination)
- Verify site is accessible
- Check logs for specific errors

### Sync Failures
- Verify `BACKEND_URL` and `BACKEND_TOKEN` are correct
- Check backend API is accessible
- Verify backend startup ID mapping

