# K6 Test Coverage Report

This document provides a comprehensive overview of all endpoints and test scenarios covered by the k6 load test suite.

## Test Execution Flow

Each virtual user (VU) executes the following test sequence:

```
1. Health Check
   ↓
2. Authentication Flow (Register + Login)
   ↓
3. Get Current User (JWT Protected)
   ↓
4. Startup Operations (Create, List, Get, Update)
   ↓
5. Job Operations (Create, List, Update, Delete)
   ↓
6. File Upload
   ↓
7. Refresh Token
   ↓
8. Public Routes (without authentication)
   ↓
9. Token-Based Routes (optional, if API_TOKEN provided)
```

---

## Complete Endpoint Coverage

### ✅ 1. Health Check

**Endpoint:** `GET /api/v1/health`

**Test Function:** `testHealthCheck()`

**Validations:**
- ✅ Status code is 200
- ✅ Response body contains `{"status": "ok"}`

**Authentication:** None (Public)

**Test Data:** N/A

---

### ✅ 2. Authentication Endpoints

#### 2.1 User Registration

**Endpoint:** `POST /api/v1/auth/register`

**Test Function:** `testAuthFlow()` → Register

**Validations:**
- ✅ Status code is 200 or 201
- ✅ Response contains user data (user object or user ID)

**Authentication:** None (Public)

**Test Data:**
```json
{
  "email": "test_{timestamp}_{random}@example.com",
  "password": "TestPassword123!",
  "name": "Test User {random}"
}
```

**Metrics Tracked:** `auth_response_time`

---

#### 2.2 User Login

**Endpoint:** `POST /api/v1/auth/login`

**Test Function:** `testAuthFlow()` → Login

**Validations:**
- ✅ Status code is 200
- ✅ Response contains `access_token` and `refresh_token`
- ✅ Tokens are stored in shared state for subsequent requests

**Authentication:** None (Public)

**Test Data:**
```json
{
  "email": "{email_from_registration}",
  "password": "TestPassword123!"
}
```

**Metrics Tracked:** `auth_response_time`

**State Updates:**
- `sharedState.accessToken` ← access_token
- `sharedState.refreshToken` ← refresh_token
- `sharedState.userId` ← user.id

---

#### 2.3 Refresh Token

**Endpoint:** `POST /api/v1/auth/refresh`

**Test Function:** `testRefreshToken()`

**Validations:**
- ✅ Status code is 200
- ✅ Response contains new `access_token`
- ✅ New tokens are stored in shared state

**Authentication:** None (Public)

**Test Data:**
```json
{
  "refresh_token": "{refresh_token_from_login}"
}
```

**Metrics Tracked:** `auth_response_time`

**State Updates:**
- `sharedState.accessToken` ← new access_token
- `sharedState.refreshToken` ← new refresh_token (if provided)

---

### ✅ 3. User Endpoints

#### 3.1 Get Current User

**Endpoint:** `GET /api/v1/me`

**Test Function:** `testGetMe()`

**Validations:**
- ✅ Status code is 200
- ✅ Response contains user data (user_id or id)

**Authentication:** JWT Bearer Token (Required)

**Headers:**
```
Authorization: Bearer {access_token}
```

**Metrics Tracked:** N/A

---

### ✅ 4. Startup Endpoints

#### 4.1 Create Startup

**Endpoint:** `POST /api/v1/startups`

**Test Function:** `testStartupOperations()` → Create

**Validations:**
- ✅ Status code is 200 or 201
- ✅ Response contains startup data with `id` and `slug`
- ✅ Startup ID and slug stored in shared state

**Authentication:** JWT Bearer Token (Required)

**Test Data:**
```json
{
  "name": "Test Startup {random}",
  "description": "A test startup description that is long enough to pass validation...",
  "website": "https://example.com",
  "founded_year": 2020,
  "industry": "Technology",
  "company_size": "1-10",
  "location": "San Francisco, CA",
  "allow_public_join": false
}
```

**Metrics Tracked:** `startup_response_time`

**State Updates:**
- `sharedState.startupId` ← startup.id
- `sharedState.startupSlug` ← startup.slug

---

#### 4.2 List Startups (Public)

**Endpoint:** `GET /api/v1/startups?page=1&page_size=10`

**Test Function:** `testStartupOperations()` → List

**Validations:**
- ✅ Status code is 200
- ✅ Response contains array of startups in `data` field

**Authentication:** None (Public)

**Query Parameters:**
- `page=1`
- `page_size=10`

**Metrics Tracked:** `startup_response_time`

---

#### 4.3 Get Startup by ID

**Endpoint:** `GET /api/v1/startups/:id`

**Test Function:** `testStartupOperations()` → Get by ID

**Validations:**
- ✅ Status code is 200
- ✅ Response startup ID matches requested ID

**Authentication:** JWT Bearer Token (Required)

**Path Parameters:**
- `id`: `{startupId_from_create}`

**Metrics Tracked:** `startup_response_time`

---

#### 4.4 Get Startup by Slug

**Endpoint:** `GET /api/v1/startups/slug/:slug`

**Test Function:** `testStartupOperations()` → Get by Slug

**Validations:**
- ✅ Status code is 200
- ✅ Response startup slug matches requested slug

**Authentication:** None (Public)

**Path Parameters:**
- `slug`: `{startupSlug_from_create}`

**Metrics Tracked:** `startup_response_time`

---

#### 4.5 Update Startup

**Endpoint:** `PUT /api/v1/startups/:id`

**Test Function:** `testStartupOperations()` → Update

**Validations:**
- ✅ Status code is 200
- ✅ Response contains updated startup data
- ✅ Startup ID matches requested ID

**Authentication:** JWT Bearer Token (Required)

**Path Parameters:**
- `id`: `{startupId_from_create}`

**Test Data:**
```json
{
  "description": "Updated description for the test startup...",
  "website": "https://updated-example.com"
}
```

**Metrics Tracked:** `startup_response_time`

---

### ✅ 5. Job Endpoints

#### 5.1 Create Job

**Endpoint:** `POST /api/v1/jobs`

**Test Function:** `testJobOperations()` → Create

**Validations:**
- ✅ Status code is 200 or 201
- ✅ Response contains job data with `id`
- ✅ Job ID stored in shared state

**Authentication:** JWT Bearer Token (Required)

**Test Data:**
```json
{
  "startup_id": "{startupId_from_create}",
  "title": "Senior Software Engineer",
  "description": "We are looking for an experienced software engineer...",
  "requirements": "5+ years of experience in software development...",
  "job_type": "full_time",
  "location_type": "remote",
  "city": "San Francisco",
  "country": "USA",
  "salary_min": 100000,
  "salary_max": 150000,
  "currency": "USD",
  "application_email": "jobs@example.com"
}
```

**Metrics Tracked:** `job_response_time`

**State Updates:**
- `sharedState.jobId` ← job.id

---

#### 5.2 List Jobs (Public)

**Endpoint:** `GET /api/v1/jobs?page=1&page_size=10`

**Test Function:** `testJobOperations()` → List

**Validations:**
- ✅ Status code is 200
- ✅ Response contains array of jobs in `data` field

**Authentication:** None (Public)

**Query Parameters:**
- `page=1`
- `page_size=10`

**Metrics Tracked:** `job_response_time`

---

#### 5.3 List Jobs with Filters

**Endpoint:** `GET /api/v1/jobs?startup_id={id}&job_type=full_time&location_type=remote&page=1&page_size=10`

**Test Function:** `testJobOperations()` → List with Filters

**Validations:**
- ✅ Status code is 200

**Authentication:** None (Public)

**Query Parameters:**
- `startup_id`: `{startupId_from_create}`
- `job_type`: `full_time`
- `location_type`: `remote`
- `page=1`
- `page_size=10`

**Metrics Tracked:** `job_response_time`

---

#### 5.4 Update Job

**Endpoint:** `PUT /api/v1/jobs/:id`

**Test Function:** `testJobOperations()` → Update

**Validations:**
- ✅ Status code is 200
- ✅ Response contains updated job data
- ✅ Job ID matches requested ID

**Authentication:** JWT Bearer Token (Required)

**Path Parameters:**
- `id`: `{jobId_from_create}`

**Test Data:**
```json
{
  "title": "Senior Software Engineer - Updated",
  "salary_min": 110000,
  "salary_max": 160000
}
```

**Metrics Tracked:** `job_response_time`

---

#### 5.5 Delete Job

**Endpoint:** `DELETE /api/v1/jobs/:id`

**Test Function:** `testJobOperations()` → Delete

**Validations:**
- ✅ Status code is 200, 204, or 404 (404 acceptable if already deleted)
- ✅ Response is valid (if status 200, contains data or message)

**Authentication:** JWT Bearer Token (Required)

**Path Parameters:**
- `id`: `{jobId_from_create}`

**Metrics Tracked:** `job_response_time`

**State Updates:**
- `sharedState.jobId` ← null (cleared after deletion)

**Note:** Accepts 404 as valid since job might have been deleted in a previous iteration.

---

### ✅ 6. File Upload Endpoint

#### 6.1 Upload File

**Endpoint:** `POST /api/v1/upload`

**Test Function:** `testFileUpload()`

**Validations:**
- ✅ Status code is 200 (success) or 400 (expected failure - storage not configured)
- ✅ Response contains file data (if status 200) or error message (if status 400)

**Authentication:** JWT Bearer Token (Required)

**Content-Type:** `multipart/form-data`

**Test Data:**
- File: Minimal 1x1 transparent PNG (43 bytes)
- Filename: `test.png`
- Content-Type: `image/png`

**Metrics Tracked:** N/A (errors only tracked for 500+ status codes)

**Note:** 
- 400 errors are expected if storage service is not configured
- Binary file upload in k6 has limitations, so this test verifies endpoint accessibility
- Only 500+ status codes are counted as errors

---

### ✅ 7. Public Routes (Without Authentication)

#### 7.1 List Startups (Public with Filters)

**Endpoint:** `GET /api/v1/startups?page=1&page_size=5&order_by=created_at&order_dir=DESC`

**Test Function:** `testPublicRoutes()`

**Validations:**
- ✅ Status code is 200

**Authentication:** None (Public)

**Query Parameters:**
- `page=1`
- `page_size=5`
- `order_by=created_at`
- `order_dir=DESC`

---

#### 7.2 List Jobs (Public)

**Endpoint:** `GET /api/v1/jobs?page=1&page_size=5`

**Test Function:** `testPublicRoutes()`

**Validations:**
- ✅ Status code is 200

**Authentication:** None (Public)

**Query Parameters:**
- `page=1`
- `page_size=5`

---

### ✅ 8. Token-Based Routes (Optional)

**Note:** These routes are only tested if `API_TOKEN` environment variable is provided.

#### 8.1 Get Startup Info from Token

**Endpoint:** `GET /api/v1/token/startup`

**Test Function:** `testTokenRoutes()` → Get Startup

**Validations:**
- ✅ Status code is 200
- ✅ Response contains `startup_id`

**Authentication:** API Token (Bearer format)

**Headers:**
```
Authorization: Bearer {API_TOKEN}
```

**State Updates:**
- `tokenStartupId` ← startup_id (used for subsequent token-based job operations)

---

#### 8.2 List Jobs via Token

**Endpoint:** `GET /api/v1/token/jobs?page=1&page_size=10`

**Test Function:** `testTokenRoutes()` → List Jobs

**Validations:**
- ✅ Status code is 200

**Authentication:** API Token (Bearer format)

**Query Parameters:**
- `page=1`
- `page_size=10`

---

#### 8.3 Create Job via Token

**Endpoint:** `POST /api/v1/token/jobs`

**Test Function:** `testTokenRoutes()` → Create Job

**Validations:**
- ✅ Status code is 200 or 201
- ✅ Response contains job data with `id`

**Authentication:** API Token (Bearer format)

**Test Data:**
```json
{
  "startup_id": "{tokenStartupId}",
  "title": "Token-Based Job Posting",
  "description": "This job was created using an API token...",
  "requirements": "Experience with API integrations...",
  "job_type": "full_time",
  "location_type": "remote",
  "city": "Remote",
  "country": "USA",
  "salary_min": 80000,
  "salary_max": 120000,
  "currency": "USD",
  "application_email": "token-jobs@example.com"
}
```

**Metrics Tracked:** `job_response_time`

---

#### 8.4 Update Job via Token

**Endpoint:** `PUT /api/v1/token/jobs/:id`

**Test Function:** `testTokenRoutes()` → Update Job

**Validations:**
- ✅ Status code is 200

**Authentication:** API Token (Bearer format)

**Path Parameters:**
- `id`: `{jobId_from_token_create}`

**Test Data:**
```json
{
  "title": "Token-Based Job Posting - Updated",
  "salary_min": 90000,
  "salary_max": 130000
}
```

**Metrics Tracked:** `job_response_time`

---

#### 8.5 Delete Job via Token

**Endpoint:** `DELETE /api/v1/token/jobs/:id`

**Test Function:** `testTokenRoutes()` → Delete Job

**Validations:**
- ✅ Status code is 200

**Authentication:** API Token (Bearer format)

**Path Parameters:**
- `id`: `{jobId_from_token_create}`

**Metrics Tracked:** `job_response_time`

---

## Test Coverage Summary

### Endpoints by Category

| Category | Total Endpoints | Tested | Coverage |
|----------|----------------|--------|----------|
| **Health** | 1 | 1 | ✅ 100% |
| **Authentication** | 3 | 3 | ✅ 100% |
| **User** | 1 | 1 | ✅ 100% |
| **Startups** | 5 | 5 | ✅ 100% |
| **Jobs (JWT)** | 5 | 5 | ✅ 100% |
| **File Upload** | 1 | 1 | ✅ 100% |
| **Public Routes** | 2 | 2 | ✅ 100% |
| **Token Routes** | 5 | 5* | ✅ 100%* |
| **TOTAL** | **23** | **23** | **✅ 100%** |

*Token routes are optional and only tested if API_TOKEN is provided

### HTTP Methods Coverage

| Method | Count | Endpoints |
|--------|-------|-----------|
| **GET** | 11 | Health, Get Me, List Startups (2x), Get Startup by ID, Get Startup by Slug, List Jobs (3x), Token Get Startup, Token List Jobs |
| **POST** | 7 | Register, Login, Refresh Token, Create Startup, Create Job, Upload File, Token Create Job |
| **PUT** | 4 | Update Startup, Update Job, Token Update Job |
| **DELETE** | 2 | Delete Job, Token Delete Job |

### Authentication Methods Tested

| Method | Endpoints | Count |
|--------|-----------|-------|
| **None (Public)** | Health, Register, Login, Register, Refresh Token, List Startups, Get Startup by Slug, List Jobs | 8 |
| **JWT Bearer Token** | Get Me, Create/Update/Get Startup, Create/Update/Delete Job, Upload File | 8 |
| **API Token** | Token routes (Get Startup, List/Create/Update/Delete Jobs) | 5 |

---

## Test Validations Performed

### Response Validations

1. **Status Code Validation**
   - Expected status codes for each endpoint
   - Acceptable alternative status codes (e.g., 404 for delete operations)

2. **Response Body Validation**
   - JSON structure validation
   - Required fields presence
   - Data type validation
   - Value matching (IDs, slugs, etc.)

3. **State Management**
   - Token storage and reuse
   - Resource ID tracking (startup ID, job ID)
   - State cleanup after operations

### Error Handling

- **Expected Errors:** 400 status codes are accepted for file upload (storage not configured)
- **Unexpected Errors:** 500+ status codes are tracked as errors
- **Error Rate Tracking:** Custom `errors` metric tracks all failures

---

## Custom Metrics Tracked

| Metric | Description | Tracked For |
|--------|-------------|-------------|
| `auth_response_time` | Response times for auth endpoints | Register, Login, Refresh Token |
| `startup_response_time` | Response times for startup endpoints | All startup operations |
| `job_response_time` | Response times for job endpoints | All job operations |
| `errors` | Overall error rate | All endpoints |

---

## Test Data Generation

### Dynamic Test Data

- **Emails:** `test_{timestamp}_{random}@example.com`
- **Names:** `Test User {random}`
- **Startup Names:** `Test Startup {random}`

### Static Test Data

- **Password:** `TestPassword123!`
- **Startup Details:** Predefined valid startup data
- **Job Details:** Predefined valid job posting data
- **File:** Minimal 1x1 PNG image

---

## Test Execution Configuration

### Load Pattern

```
Stage 1: Ramp up to 5 users over 30s
Stage 2: Stay at 5 users for 1m
Stage 3: Ramp up to 10 users over 30s
Stage 4: Stay at 10 users for 1m
Stage 5: Ramp down to 0 users over 30s
```

### Performance Thresholds

- **Response Time:** 95% of requests < 2000ms
- **Error Rate:** < 10%
- **Custom Error Rate:** < 10%

---

## Notes and Limitations

1. **File Upload:** Binary file upload in k6 has limitations. The test uses a minimal PNG but may fail if storage service is not configured (400 is acceptable).

2. **Token Routes:** Token-based routes require an API token to be provided via `API_TOKEN` environment variable. These routes are skipped if no token is provided.

3. **State Management:** Each VU maintains its own state. Resources created in one iteration may not be available in subsequent iterations.

4. **Delete Operations:** Delete operations accept 404 as valid since resources may have been deleted in previous iterations.

5. **Concurrent Execution:** Multiple VUs run concurrently, which may cause database conflicts (e.g., duplicate email registration). The test uses unique email generation to mitigate this.

---

## Running Tests with Coverage Report

To see what's being tested, refer to this document or run:

```bash
# Run tests and view summary
k6 run k6_test.js

# Run with specific base URL
BASE_URL=http://localhost:8080 k6 run k6_test.js

# Run with API token for token route testing
BASE_URL=http://localhost:8080 API_TOKEN=sb_your_token k6 run k6_test.js
```

The test output will show:
- ✅ Check marks for passed validations
- ✗ Marks for failed validations
- Custom metrics for each endpoint category
- Overall performance statistics

