# K6 Load Testing Guide

This document describes how to run k6 load tests against the JoinUs backend API.

## Prerequisites

1. Install k6: https://k6.io/docs/getting-started/installation/
2. Ensure the backend is running (locally or in production)

## Test File

The test file `k6_test.js` covers all available endpoints in the backend:

### Public Routes
- `GET /api/v1/health` - Health check
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/startups` - List startups (with filters)
- `GET /api/v1/startups/slug/:slug` - Get startup by slug
- `GET /api/v1/jobs` - List jobs (with filters)

### Protected Routes (JWT)
- `GET /api/v1/me` - Get current user
- `POST /api/v1/startups` - Create startup
- `PUT /api/v1/startups/:id` - Update startup
- `GET /api/v1/startups/:id` - Get startup by ID
- `POST /api/v1/jobs` - Create job
- `PUT /api/v1/jobs/:id` - Update job
- `DELETE /api/v1/jobs/:id` - Delete job
- `POST /api/v1/upload` - Upload file

### Token Routes (API Token)
- `GET /api/v1/token/startup` - Get startup info from token
- `POST /api/v1/token/jobs` - Create job (token-based)
- `PUT /api/v1/token/jobs/:id` - Update job (token-based)
- `DELETE /api/v1/token/jobs/:id` - Delete job (token-based)
- `GET /api/v1/token/jobs` - List jobs (token-based)

**Note:** Token-based routes are not currently tested in the script as they require an API token that needs to be obtained from a startup entity. To test these routes, you would need to:
1. Create a startup via JWT auth
2. Retrieve the API token from the startup
3. Use that token to test the `/api/v1/token/*` routes

## Running Tests

### Test Against Local Environment

```bash
# Default (localhost:8080)
k6 run k6_test.js

# Or explicitly set the base URL
BASE_URL=http://localhost:8080 k6 run k6_test.js
```

### Test Against Production Environment

```bash
BASE_URL=https://api.yourdomain.com k6 run k6_test.js
```

### Test Against Staging Environment

```bash
BASE_URL=https://staging-api.yourdomain.com k6 run k6_test.js
```

### Test Token-Based Routes

To test token-based routes, you need to provide an API token:

```bash
# Test with API token (format: sb_<base64_token>)
BASE_URL=http://localhost:8080 API_TOKEN=sb_your_token_here k6 run k6_test.js
```

**How to get an API token:**
1. Create a startup via JWT authentication (the test does this automatically)
2. Retrieve the API token from the startup entity (check your database or API response)
3. Use that token with the `API_TOKEN` environment variable

**Note:** Token-based routes are optional. If no `API_TOKEN` is provided, the test will skip these routes.

## Test Configuration

The test is configured with the following load pattern:

- **Stage 1**: Ramp up to 5 users over 30 seconds
- **Stage 2**: Stay at 5 users for 1 minute
- **Stage 3**: Ramp up to 10 users over 30 seconds
- **Stage 4**: Stay at 10 users for 1 minute
- **Stage 5**: Ramp down to 0 users over 30 seconds

### Thresholds

The test has the following performance thresholds:
- 95% of requests should complete in under 2 seconds
- Error rate should be less than 10%

### Custom Metrics

The test tracks custom metrics:
- `auth_response_time` - Response times for authentication endpoints
- `startup_response_time` - Response times for startup endpoints
- `job_response_time` - Response times for job endpoints
- `errors` - Overall error rate

## Test Flow

The test follows this flow for each virtual user:

1. **Health Check** - Verify API is accessible
2. **Authentication Flow** - Register a new user and login
3. **Get Current User** - Retrieve authenticated user info
4. **Startup Operations**:
   - Create a startup
   - List startups (public)
   - Get startup by ID
   - Get startup by slug
   - Update startup
5. **Job Operations**:
   - Create a job
   - List jobs (public)
   - List jobs with filters
   - Update job
   - Delete job
6. **File Upload** - Upload a test file
7. **Refresh Token** - Test token refresh
8. **Public Routes** - Test public endpoints without authentication

## Customizing Tests

### Adjust Load Pattern

Edit the `stages` in the `options` object:

```javascript
export const options = {
  stages: [
    { duration: '1m', target: 20 },  // Ramp to 20 users
    { duration: '3m', target: 20 },  // Stay at 20 users
    { duration: '1m', target: 0 },   // Ramp down
  ],
};
```

### Adjust Thresholds

Modify the `thresholds` in the `options` object:

```javascript
thresholds: {
  http_req_duration: ['p(95)<1000'], // 95% under 1 second
  http_req_failed: ['rate<0.05'],   // Less than 5% errors
},
```

### Test Specific Endpoints Only

Comment out or remove function calls in the `default` function to test only specific flows.

## Output

The test generates:
- Real-time console output showing test progress
- Summary statistics at the end
- A `summary.json` file with detailed metrics

## Troubleshooting

### Connection Refused

If you see connection errors, ensure:
- The backend is running
- The BASE_URL is correct
- The port is accessible

### Authentication Failures

If authentication fails:
- Check that the database is properly set up
- Verify that user registration is working
- Check backend logs for errors

### High Error Rates

If error rates are high:
- Check backend logs
- Verify database connectivity
- Check if rate limiting is enabled
- Ensure sufficient resources are available

## Example Output

```
          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: k6_test.js
     output: -

  scenarios: (100.00%) 1 scenario, 10 max VUs, 3m30s max duration
           ✓ setup=0s
           ✓ teardown=0s

     ✓ health check status is 200
     ✓ register status is 200 or 201
     ✓ login status is 200
     ✓ get me status is 200
     ...

     checks.........................: 100.00% ✓ 150  ✗ 0
     data_received..................: 2.5 MB  12 kB/s
     data_sent......................: 150 kB  700 B/s
     http_req_duration..............: avg=250ms min=50ms med=200ms max=2000ms p(90)=500ms p(95)=800ms
     http_req_failed................: 0.00%   ✓ 0    ✗ 150
     http_reqs......................: 150     0.714/s
     iteration_duration.............: avg=3.5s min=2s max=5s
     iterations.....................: 10      0.048/s
     vus............................: 1       min=1 max=10
     vus_max........................: 10      min=10 max=10
```

