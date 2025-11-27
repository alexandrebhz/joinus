import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const authResponseTime = new Trend('auth_response_time');
const startupResponseTime = new Trend('startup_response_time');
const jobResponseTime = new Trend('job_response_time');

// Configuration
export const options = {
  stages: [
    { duration: '30s', target: 5 },   // Ramp up to 5 users
    { duration: '1m', target: 5 },    // Stay at 5 users
    { duration: '30s', target: 10 },  // Ramp up to 10 users
    { duration: '1m', target: 10 },   // Stay at 10 users
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% of requests should be below 2s
    http_req_failed: ['rate<0.1'],     // Error rate should be less than 10%
    errors: ['rate<0.1'],
  },
};

// Base URL - can be set via environment variable or default to local
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// API Token for token-based routes (optional)
// Format: sb_<base64_token>
// To get a token: Create a startup via JWT auth, then retrieve the API token from the startup
const API_TOKEN = __ENV.API_TOKEN || null;

// Test data generators
function generateEmail() {
  return `test_${Date.now()}_${Math.random().toString(36).substring(7)}@example.com`;
}

function generateName() {
  return `Test User ${Math.random().toString(36).substring(7)}`;
}

function generateStartupName() {
  return `Test Startup ${Math.random().toString(36).substring(7)}`;
}

// Shared state across VUs
let sharedState = {
  accessToken: null,
  refreshToken: null,
  userId: null,
  startupId: null,
  startupSlug: null,
  jobId: null,
};

export default function () {
  // Test flow: Health -> Register -> Login -> Protected Routes -> Token Routes
  
  // 1. Health Check
  testHealthCheck();
  sleep(1);

  // 2. Authentication Flow
  const authResult = testAuthFlow();
  if (!authResult.success) {
    return; // Skip rest if auth fails
  }
  sleep(1);

  // 3. Get Current User
  testGetMe();
  sleep(1);

  // 4. Startup Operations
  const startupResult = testStartupOperations();
  if (!startupResult.success) {
    return;
  }
  sleep(1);

  // 5. Job Operations
  testJobOperations();
  sleep(1);

  // 6. File Upload
  testFileUpload();
  sleep(1);

  // 7. Refresh Token
  testRefreshToken();
  sleep(1);

  // 8. Public Routes (without auth)
  testPublicRoutes();
  sleep(1);

  // 9. Token-based Routes (if API token is provided)
  if (API_TOKEN) {
    testTokenRoutes();
    sleep(1);
  }
}

function testHealthCheck() {
  const url = `${BASE_URL}/api/v1/health`;
  const response = http.get(url);
  
  const success = check(response, {
    'health check status is 200': (r) => r.status === 200,
    'health check has status ok': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.status === 'ok';
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!success);
}

function testAuthFlow() {
  const email = generateEmail();
  const name = generateName();
  const password = 'TestPassword123!';

  // Register
  const registerUrl = `${BASE_URL}/api/v1/auth/register`;
  const registerPayload = JSON.stringify({
    email: email,
    password: password,
    name: name,
  });

  const registerResponse = http.post(registerUrl, registerPayload, {
    headers: { 'Content-Type': 'application/json' },
  });

  const registerSuccess = check(registerResponse, {
    'register status is 200 or 201': (r) => r.status === 200 || r.status === 201,
    'register returns user data': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && (body.data.user || body.data.id);
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!registerSuccess);
  authResponseTime.add(registerResponse.timings.duration);

  // Login
  const loginUrl = `${BASE_URL}/api/v1/auth/login`;
  const loginPayload = JSON.stringify({
    email: email,
    password: password,
  });

  const loginResponse = http.post(loginUrl, loginPayload, {
    headers: { 'Content-Type': 'application/json' },
  });

  const loginSuccess = check(loginResponse, {
    'login status is 200': (r) => r.status === 200,
    'login returns access token': (r) => {
      try {
        const body = JSON.parse(r.body);
        if (body.data && body.data.access_token) {
          sharedState.accessToken = body.data.access_token;
          sharedState.refreshToken = body.data.refresh_token;
          if (body.data.user) {
            sharedState.userId = body.data.user.id;
          }
          return true;
        }
        return false;
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!loginSuccess);
  authResponseTime.add(loginResponse.timings.duration);

  return { success: registerSuccess && loginSuccess };
}

function testGetMe() {
  if (!sharedState.accessToken) {
    return;
  }

  const url = `${BASE_URL}/api/v1/me`;
  const response = http.get(url, {
    headers: {
      'Authorization': `Bearer ${sharedState.accessToken}`,
    },
  });

  const success = check(response, {
    'get me status is 200': (r) => r.status === 200,
    'get me returns user data': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && (body.data.user_id || body.data.id);
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!success);
}

function testStartupOperations() {
  if (!sharedState.accessToken) {
    return { success: false };
  }

  // Create Startup
  const createUrl = `${BASE_URL}/api/v1/startups`;
  const createPayload = JSON.stringify({
    name: generateStartupName(),
    description: 'A test startup description that is long enough to pass validation. This startup is created for testing purposes.',
    website: 'https://example.com',
    founded_year: 2020,
    industry: 'Technology',
    company_size: '1-10',
    location: 'San Francisco, CA',
    allow_public_join: false,
  });

  const createResponse = http.post(createUrl, createPayload, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${sharedState.accessToken}`,
    },
  });

  const createSuccess = check(createResponse, {
    'create startup status is 200 or 201': (r) => r.status === 200 || r.status === 201,
    'create startup returns startup data': (r) => {
      try {
        const body = JSON.parse(r.body);
        if (body.data && body.data.id) {
          sharedState.startupId = body.data.id;
          sharedState.startupSlug = body.data.slug;
          return true;
        }
        return false;
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!createSuccess);
  startupResponseTime.add(createResponse.timings.duration);

  if (!createSuccess || !sharedState.startupId) {
    return { success: false };
  }

  // List Startups
  const listUrl = `${BASE_URL}/api/v1/startups?page=1&page_size=10`;
  const listResponse = http.get(listUrl);

  const listSuccess = check(listResponse, {
    'list startups status is 200': (r) => r.status === 200,
    'list startups returns array': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && Array.isArray(body.data);
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!listSuccess);
  startupResponseTime.add(listResponse.timings.duration);

  // Get Startup by ID
  if (sharedState.startupId) {
    const getUrl = `${BASE_URL}/api/v1/startups/${sharedState.startupId}`;
    const getResponse = http.get(getUrl, {
      headers: {
        'Authorization': `Bearer ${sharedState.accessToken}`,
      },
    });

    const getSuccess = check(getResponse, {
      'get startup by id status is 200': (r) => r.status === 200,
      'get startup by id returns startup': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && body.data.id === sharedState.startupId;
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!getSuccess);
    startupResponseTime.add(getResponse.timings.duration);
  }

  // Get Startup by Slug
  if (sharedState.startupSlug) {
    const getBySlugUrl = `${BASE_URL}/api/v1/startups/slug/${sharedState.startupSlug}`;
    const getBySlugResponse = http.get(getBySlugUrl);

    const getBySlugSuccess = check(getBySlugResponse, {
      'get startup by slug status is 200': (r) => r.status === 200,
      'get startup by slug returns startup': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && body.data.slug === sharedState.startupSlug;
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!getBySlugSuccess);
    startupResponseTime.add(getBySlugResponse.timings.duration);
  }

  // Update Startup
  if (sharedState.startupId) {
    const updateUrl = `${BASE_URL}/api/v1/startups/${sharedState.startupId}`;
    const updatePayload = JSON.stringify({
      description: 'Updated description for the test startup. This is a longer description that meets the minimum requirements.',
      website: 'https://updated-example.com',
    });

    const updateResponse = http.put(updateUrl, updatePayload, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${sharedState.accessToken}`,
      },
    });

    const updateSuccess = check(updateResponse, {
      'update startup status is 200': (r) => r.status === 200,
      'update startup returns updated data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && body.data.id === sharedState.startupId;
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!updateSuccess);
    startupResponseTime.add(updateResponse.timings.duration);
  }

  return { success: createSuccess };
}

function testJobOperations() {
  if (!sharedState.accessToken || !sharedState.startupId) {
    return;
  }

  // Create Job
  const createUrl = `${BASE_URL}/api/v1/jobs`;
  const createPayload = JSON.stringify({
    startup_id: sharedState.startupId,
    title: 'Senior Software Engineer',
    description: 'We are looking for an experienced software engineer to join our team. This is a great opportunity to work on cutting-edge technology.',
    requirements: '5+ years of experience in software development. Strong knowledge of Go, JavaScript, and cloud technologies.',
    job_type: 'full_time',
    location_type: 'remote',
    city: 'San Francisco',
    country: 'USA',
    salary_min: 100000,
    salary_max: 150000,
    currency: 'USD',
    application_email: 'jobs@example.com',
  });

  const createResponse = http.post(createUrl, createPayload, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${sharedState.accessToken}`,
    },
  });

  const createSuccess = check(createResponse, {
    'create job status is 200 or 201': (r) => r.status === 200 || r.status === 201,
    'create job returns job data': (r) => {
      try {
        const body = JSON.parse(r.body);
        if (body.data && body.data.id) {
          sharedState.jobId = body.data.id;
          return true;
        }
        return false;
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!createSuccess);
  jobResponseTime.add(createResponse.timings.duration);

  // List Jobs
  const listUrl = `${BASE_URL}/api/v1/jobs?page=1&page_size=10`;
  const listResponse = http.get(listUrl);

  const listSuccess = check(listResponse, {
    'list jobs status is 200': (r) => r.status === 200,
    'list jobs returns array': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && Array.isArray(body.data);
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!listSuccess);
  jobResponseTime.add(listResponse.timings.duration);

  // List Jobs with filters
  const filteredListUrl = `${BASE_URL}/api/v1/jobs?startup_id=${sharedState.startupId}&job_type=full_time&location_type=remote&page=1&page_size=10`;
  const filteredListResponse = http.get(filteredListUrl);

  check(filteredListResponse, {
    'list jobs with filters status is 200': (r) => r.status === 200,
  });

  jobResponseTime.add(filteredListResponse.timings.duration);

  // Update Job
  if (sharedState.jobId) {
    const updateUrl = `${BASE_URL}/api/v1/jobs/${sharedState.jobId}`;
    const updatePayload = JSON.stringify({
      title: 'Senior Software Engineer - Updated',
      salary_min: 110000,
      salary_max: 160000,
    });

    const updateResponse = http.put(updateUrl, updatePayload, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${sharedState.accessToken}`,
      },
    });

    const updateSuccess = check(updateResponse, {
      'update job status is 200': (r) => r.status === 200,
      'update job returns updated data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && body.data.id === sharedState.jobId;
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!updateSuccess);
    jobResponseTime.add(updateResponse.timings.duration);
  }

  // Delete Job
  if (sharedState.jobId) {
    const deleteUrl = `${BASE_URL}/api/v1/jobs/${sharedState.jobId}`;
    const deleteResponse = http.del(deleteUrl, null, {
      headers: {
        'Authorization': `Bearer ${sharedState.accessToken}`,
      },
    });

    const deleteSuccess = check(deleteResponse, {
      'delete job status is 200': (r) => r.status === 200,
    });

    errorRate.add(!deleteSuccess);
    jobResponseTime.add(deleteResponse.timings.duration);

    // Clear job ID after deletion
    if (deleteSuccess) {
      sharedState.jobId = null;
    }
  }
}

function testFileUpload() {
  if (!sharedState.accessToken) {
    return;
  }

  // Create a simple text file for upload
  const fileContent = 'This is a test file content for upload testing.';
  const fileName = 'test.txt';
  
  // k6 doesn't have native multipart/form-data support, so we'll construct it manually
  const boundary = '----WebKitFormBoundary' + Math.random().toString(36);
  const formData = [
    `--${boundary}`,
    `Content-Disposition: form-data; name="file"; filename="${fileName}"`,
    'Content-Type: text/plain',
    '',
    fileContent,
    `--${boundary}--`,
  ].join('\r\n');

  const uploadUrl = `${BASE_URL}/api/v1/upload`;
  const uploadResponse = http.post(uploadUrl, formData, {
    headers: {
      'Content-Type': `multipart/form-data; boundary=${boundary}`,
      'Authorization': `Bearer ${sharedState.accessToken}`,
    },
  });

  const uploadSuccess = check(uploadResponse, {
    'upload file status is 200': (r) => r.status === 200,
    'upload file returns file data': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && (body.data.url || body.data.id);
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!uploadSuccess);
}

function testRefreshToken() {
  if (!sharedState.refreshToken) {
    return;
  }

  const url = `${BASE_URL}/api/v1/auth/refresh`;
  const payload = JSON.stringify({
    refresh_token: sharedState.refreshToken,
  });

  const response = http.post(url, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  const success = check(response, {
    'refresh token status is 200': (r) => r.status === 200,
    'refresh token returns new access token': (r) => {
      try {
        const body = JSON.parse(r.body);
        if (body.data && body.data.access_token) {
          sharedState.accessToken = body.data.access_token;
          if (body.data.refresh_token) {
            sharedState.refreshToken = body.data.refresh_token;
          }
          return true;
        }
        return false;
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!success);
  authResponseTime.add(response.timings.duration);
}

function testPublicRoutes() {
  // Test public startup listing with various filters
  const listStartupsUrl = `${BASE_URL}/api/v1/startups?page=1&page_size=5&order_by=created_at&order_dir=DESC`;
  const listStartupsResponse = http.get(listStartupsUrl);

  check(listStartupsResponse, {
    'public list startups status is 200': (r) => r.status === 200,
  });

  // Test public job listing
  const listJobsUrl = `${BASE_URL}/api/v1/jobs?page=1&page_size=5`;
  const listJobsResponse = http.get(listJobsUrl);

  check(listJobsResponse, {
    'public list jobs status is 200': (r) => r.status === 200,
  });
}

function testTokenRoutes() {
  if (!API_TOKEN) {
    return;
  }

  // Get startup info from token
  const getStartupUrl = `${BASE_URL}/api/v1/token/startup`;
  const getStartupResponse = http.get(getStartupUrl, {
    headers: {
      'Authorization': `Bearer ${API_TOKEN}`,
    },
  });

  const getStartupSuccess = check(getStartupResponse, {
    'token get startup status is 200': (r) => r.status === 200,
    'token get startup returns startup_id': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && body.data.startup_id;
      } catch (e) {
        return false;
      }
    },
  });

  errorRate.add(!getStartupSuccess);

  if (!getStartupSuccess) {
    return; // Skip rest if token is invalid
  }

  let tokenStartupId = null;
  try {
    const body = JSON.parse(getStartupResponse.body);
    if (body.data && body.data.startup_id) {
      tokenStartupId = body.data.startup_id;
    }
  } catch (e) {
    // Ignore
  }

  // List jobs via token
  const listJobsUrl = `${BASE_URL}/api/v1/token/jobs?page=1&page_size=10`;
  const listJobsResponse = http.get(listJobsUrl, {
    headers: {
      'Authorization': `Bearer ${API_TOKEN}`,
    },
  });

  check(listJobsResponse, {
    'token list jobs status is 200': (r) => r.status === 200,
  });

  // Create job via token
  if (tokenStartupId) {
    const createJobUrl = `${BASE_URL}/api/v1/token/jobs`;
    const createJobPayload = JSON.stringify({
      startup_id: tokenStartupId,
      title: 'Token-Based Job Posting',
      description: 'This job was created using an API token. This is a test job posting created via the token-based authentication method.',
      requirements: 'Experience with API integrations and testing. Strong problem-solving skills.',
      job_type: 'full_time',
      location_type: 'remote',
      city: 'Remote',
      country: 'USA',
      salary_min: 80000,
      salary_max: 120000,
      currency: 'USD',
      application_email: 'token-jobs@example.com',
    });

    const createJobResponse = http.post(createJobUrl, createJobPayload, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${API_TOKEN}`,
      },
    });

    const createJobSuccess = check(createJobResponse, {
      'token create job status is 200 or 201': (r) => r.status === 200 || r.status === 201,
      'token create job returns job data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && body.data.id;
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!createJobSuccess);
    jobResponseTime.add(createJobResponse.timings.duration);

    // Update job via token (if job was created)
    if (createJobSuccess) {
      try {
        const body = JSON.parse(createJobResponse.body);
        const jobId = body.data && body.data.id;
        
        if (jobId) {
          const updateJobUrl = `${BASE_URL}/api/v1/token/jobs/${jobId}`;
          const updateJobPayload = JSON.stringify({
            title: 'Token-Based Job Posting - Updated',
            salary_min: 90000,
            salary_max: 130000,
          });

          const updateJobResponse = http.put(updateJobUrl, updateJobPayload, {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${API_TOKEN}`,
            },
          });

          check(updateJobResponse, {
            'token update job status is 200': (r) => r.status === 200,
          });

          jobResponseTime.add(updateJobResponse.timings.duration);

          // Delete job via token
          const deleteJobUrl = `${BASE_URL}/api/v1/token/jobs/${jobId}`;
          const deleteJobResponse = http.del(deleteJobUrl, null, {
            headers: {
              'Authorization': `Bearer ${API_TOKEN}`,
            },
          });

          check(deleteJobResponse, {
            'token delete job status is 200': (r) => r.status === 200,
          });

          jobResponseTime.add(deleteJobResponse.timings.duration);
        }
      } catch (e) {
        // Ignore errors in cleanup
      }
    }
  }
}

export function handleSummary(data) {
  const summary = {
    'stdout': `
========================================
K6 Test Summary
========================================
Base URL: ${BASE_URL}
Total Requests: ${data.metrics.http_reqs.values.count}
Failed Requests: ${data.metrics.http_req_failed.values.count}
Success Rate: ${((1 - data.metrics.http_req_failed.values.rate) * 100).toFixed(2)}%
Average Response Time: ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms
P95 Response Time: ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms
P99 Response Time: ${data.metrics.http_req_duration.values['p(99)'].toFixed(2)}ms
Error Rate: ${(data.metrics.errors.values.rate * 100).toFixed(2)}%
========================================
    `,
    'summary.json': JSON.stringify(data, null, 2),
  };
  return summary;
}

