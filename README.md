# Startup Job Board Backend - DDD Clean Architecture Specification

## Tech Stack
- **Language**: Golang 1.21+
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL 15+
- **Authentication**: JWT + API Tokens
- **Storage**: MinIO (local/dev) / AWS S3 (production)
- **Containerization**: Docker + Docker Compose

## Architecture Principles

### Domain-Driven Design (DDD)
- **Domain Layer**: Pure business logic, no external dependencies
- **Application Layer**: Use cases and application services
- **Infrastructure Layer**: External concerns (DB, HTTP, storage)
- **Presentation Layer**: API handlers and DTOs

### Clean Architecture
- Dependencies point inward (Domain ← Application ← Infrastructure ← Presentation)
- Domain entities are independent of frameworks
- Use interfaces for external dependencies
- Dependency injection throughout

### DRY (Don't Repeat Yourself)
- Shared utilities and helpers
- Generic repository patterns
- Reusable middleware
- Common response structures

---

## Project Structure

```
startup-job-board/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
├── internal/
│   ├── domain/                        # Domain Layer (Core Business Logic)
│   │   ├── entity/                    # Domain entities
│   │   │   ├── startup.go
│   │   │   ├── job.go
│   │   │   ├── user.go
│   │   │   └── file.go
│   │   ├── repository/                # Repository interfaces
│   │   │   ├── startup_repository.go
│   │   │   ├── job_repository.go
│   │   │   ├── user_repository.go
│   │   │   └── file_repository.go
│   │   ├── service/                   # Domain services (business rules)
│   │   │   └── authorization_service.go
│   │   └── valueobject/               # Value objects
│   │       ├── email.go
│   │       ├── api_token.go
│   │       └── job_type.go
│   ├── application/                   # Application Layer (Use Cases)
│   │   ├── usecase/
│   │   │   ├── startup/
│   │   │   │   ├── create_startup.go
│   │   │   │   ├── update_startup.go
│   │   │   │   ├── get_startup.go
│   │   │   │   └── list_startups.go
│   │   │   ├── job/
│   │   │   │   ├── create_job.go
│   │   │   │   ├── update_job.go
│   │   │   │   ├── list_jobs.go
│   │   │   │   └── delete_job.go
│   │   │   ├── auth/
│   │   │   │   ├── register.go
│   │   │   │   ├── login.go
│   │   │   │   └── refresh_token.go
│   │   │   └── file/
│   │   │       └── upload_file.go
│   │   ├── dto/                       # Data Transfer Objects
│   │   │   ├── startup_dto.go
│   │   │   ├── job_dto.go
│   │   │   └── auth_dto.go
│   │   └── port/                      # Application interfaces
│   │       ├── storage_service.go
│   │       └── token_service.go
│   ├── infrastructure/                # Infrastructure Layer
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   │   ├── startup_repository_impl.go
│   │   │   │   ├── job_repository_impl.go
│   │   │   │   ├── user_repository_impl.go
│   │   │   │   └── file_repository_impl.go
│   │   │   └── gorm_model/           # GORM models (separate from domain)
│   │   │       ├── startup_model.go
│   │   │       ├── job_model.go
│   │   │       └── user_model.go
│   │   ├── storage/
│   │   │   ├── minio_storage.go
│   │   │   ├── s3_storage.go
│   │   │   └── storage_factory.go
│   │   ├── auth/
│   │   │   ├── jwt_service.go
│   │   │   └── token_generator.go
│   │   └── config/
│   │       ├── config.go
│   │       └── database.go
│   └── presentation/                  # Presentation Layer
│       ├── http/
│       │   ├── handler/
│       │   │   ├── startup_handler.go
│       │   │   ├── job_handler.go
│       │   │   ├── auth_handler.go
│       │   │   └── file_handler.go
│       │   ├── middleware/
│       │   │   ├── auth_middleware.go
│       │   │   ├── api_token_middleware.go
│       │   │   ├── cors_middleware.go
│       │   │   ├── rate_limit_middleware.go
│       │   │   └── logger_middleware.go
│       │   ├── router/
│       │   │   └── router.go
│       │   └── response/
│       │       ├── success_response.go
│       │       └── error_response.go
│       └── validator/
│           └── request_validator.go
├── pkg/                               # Shared packages
│   ├── logger/
│   │   └── logger.go
│   ├── errors/
│   │   └── app_error.go
│   └── utils/
│       ├── hash.go
│       ├── slug.go
│       └── pagination.go
├── migrations/
│   └── 001_initial_schema.sql
├── docker/
│   ├── Dockerfile
│   └── Dockerfile.dev
├── scripts/
│   └── init-db.sh
├── docker-compose.yml
├── docker-compose.dev.yml
├── .env.example
├── .dockerignore
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Domain Layer Specifications

### Entities

**Startup Entity** (`internal/domain/entity/startup.go`)
```go
type Startup struct {
    ID          string
    Name        string
    Slug        string
    Description string
    LogoURL     *string
    Website     string
    FoundedYear int
    Industry    string
    CompanySize string
    Location    string
    SocialLinks SocialLinks
    APIToken    string
    Status      StartupStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   *time.Time
}

type SocialLinks struct {
    LinkedIn string
    Twitter  string
    GitHub   string
}

type StartupStatus string
const (
    StartupStatusActive   StartupStatus = "active"
    StartupStatusInactive StartupStatus = "inactive"
)
```

**Job Entity** (`internal/domain/entity/job.go`)
```go
type Job struct {
    ID            string
    StartupID     string
    Title         string
    Description   string
    Requirements  string
    JobType       JobType
    LocationType  LocationType
    City          string
    Country       string
    SalaryMin     *int
    SalaryMax     *int
    Currency      string
    ApplicationURL *string
    ApplicationEmail *string
    Status        JobStatus
    ExpiresAt     *time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type JobType string
const (
    JobTypeFullTime  JobType = "full_time"
    JobTypePartTime  JobType = "part_time"
    JobTypeContract  JobType = "contract"
    JobTypeInternship JobType = "internship"
)

type LocationType string
const (
    LocationRemote LocationType = "remote"
    LocationHybrid LocationType = "hybrid"
    LocationOnsite LocationType = "onsite"
)

type JobStatus string
const (
    JobStatusActive JobStatus = "active"
    JobStatusClosed JobStatus = "closed"
)
```

**User Entity** (`internal/domain/entity/user.go`)
```go
type User struct {
    ID        string
    Email     string
    Password  string // Hashed
    Name      string
    Role      UserRole
    StartupID *string
    Status    UserStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}

type UserRole string
const (
    UserRoleAdmin        UserRole = "admin"
    UserRoleStartupOwner UserRole = "startup_owner"
    UserRoleMember       UserRole = "member"
)

type UserStatus string
const (
    UserStatusActive  UserStatus = "active"
    UserStatusPending UserStatus = "pending"
    UserStatusInactive UserStatus = "inactive"
)
```

**Startup Member Entity** (`internal/domain/entity/startup_member.go`)
```go
type StartupMember struct {
    ID         string
    StartupID  string
    UserID     string
    Role       MemberRole
    Status     MemberStatus
    InvitedBy  *string
    InvitedAt  time.Time
    JoinedAt   *time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type MemberRole string
const (
    MemberRoleOwner      MemberRole = "owner"
    MemberRoleAdmin      MemberRole = "admin"
    MemberRoleMember     MemberRole = "member"
    MemberRoleRecruiter  MemberRole = "recruiter"  // Can only manage jobs
)

type MemberStatus string
const (
    MemberStatusActive  MemberStatus = "active"
    MemberStatusPending MemberStatus = "pending"
    MemberStatusRemoved MemberStatus = "removed"
)

// Domain methods
func (m *StartupMember) CanManageJobs() bool {
    return m.Status == MemberStatusActive && 
           (m.Role == MemberRoleOwner || m.Role == MemberRoleAdmin || m.Role == MemberRoleRecruiter)
}

func (m *StartupMember) CanManageMembers() bool {
    return m.Status == MemberStatusActive && 
           (m.Role == MemberRoleOwner || m.Role == MemberRoleAdmin)
}

func (m *StartupMember) CanManageStartup() bool {
    return m.Status == MemberStatusActive && 
           (m.Role == MemberRoleOwner || m.Role == MemberRoleAdmin)
}
```

**Invitation Entity** (`internal/domain/entity/invitation.go`)
```go
type Invitation struct {
    ID         string
    StartupID  string
    Email      string
    Token      string // Unique invitation token
    Role       MemberRole
    InvitedBy  string // UserID
    ExpiresAt  time.Time
    AcceptedAt *time.Time
    Status     InvitationStatus
    CreatedAt  time.Time
}

type InvitationStatus string
const (
    InvitationStatusPending  InvitationStatus = "pending"
    InvitationStatusAccepted InvitationStatus = "accepted"
    InvitationStatusExpired  InvitationStatus = "expired"
    InvitationStatusRevoked  InvitationStatus = "revoked"
)

func (i *Invitation) IsValid() bool {
    return i.Status == InvitationStatusPending && 
           i.ExpiresAt.After(time.Now())
}
```

**File Entity** (`internal/domain/entity/file.go`)
```go
type File struct {
    ID        string
    FileName  string
    FileSize  int64
    MimeType  string
    StorageKey string
    URL       string
    UploadedBy string
    CreatedAt time.Time
}
```

### Repository Interfaces

**Startup Repository** (`internal/domain/repository/startup_repository.go`)
```go
type StartupRepository interface {
    Create(ctx context.Context, startup *entity.Startup) error
    Update(ctx context.Context, startup *entity.Startup) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.Startup, error)
    FindBySlug(ctx context.Context, slug string) (*entity.Startup, error)
    FindByAPIToken(ctx context.Context, token string) (*entity.Startup, error)
    List(ctx context.Context, filter StartupFilter) ([]*entity.Startup, int64, error)
}

type StartupFilter struct {
    Industry string
    Status   entity.StartupStatus
    Search   string
    Pagination
}
```

**Job Repository** (`internal/domain/repository/job_repository.go`)
```go
type JobRepository interface {
    Create(ctx context.Context, job *entity.Job) error
    Update(ctx context.Context, job *entity.Job) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.Job, error)
    List(ctx context.Context, filter JobFilter) ([]*entity.Job, int64, error)
    FindByStartupID(ctx context.Context, startupID string, limit int) ([]*entity.Job, error)
}

type JobFilter struct {
    StartupID    string
    JobType      entity.JobType
    LocationType entity.LocationType
    Status       entity.JobStatus
    Search       string
    Pagination
}

type Pagination struct {
    Page     int
    PageSize int
    OrderBy  string
    OrderDir string // ASC, DESC
}
```

**Startup Member Repository** (`internal/domain/repository/startup_member_repository.go`)
```go
type StartupMemberRepository interface {
    Create(ctx context.Context, member *entity.StartupMember) error
    Update(ctx context.Context, member *entity.StartupMember) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.StartupMember, error)
    FindByUserAndStartup(ctx context.Context, userID, startupID string) (*entity.StartupMember, error)
    FindByStartupID(ctx context.Context, startupID string) ([]*entity.StartupMember, error)
    FindByUserID(ctx context.Context, userID string) ([]*entity.StartupMember, error)
}
```

**Invitation Repository** (`internal/domain/repository/invitation_repository.go`)
```go
type InvitationRepository interface {
    Create(ctx context.Context, invitation *entity.Invitation) error
    Update(ctx context.Context, invitation *entity.Invitation) error
    FindByID(ctx context.Context, id string) (*entity.Invitation, error)
    FindByToken(ctx context.Context, token string) (*entity.Invitation, error)
    FindByEmail(ctx context.Context, email string) ([]*entity.Invitation, error)
    FindPendingByStartup(ctx context.Context, startupID string) ([]*entity.Invitation, error)
    ExpireOldInvitations(ctx context.Context) error
}
```

---

## Application Layer Specifications

### Use Cases

**Create Startup Use Case** (`internal/application/usecase/startup/create_startup.go`)
```go
type CreateStartupUseCase struct {
    startupRepo       repository.StartupRepository
    memberRepo        repository.StartupMemberRepository
    tokenGen          port.TokenService
    logger            logger.Logger
}

func (uc *CreateStartupUseCase) Execute(ctx context.Context, input dto.CreateStartupInput, userID string) (*dto.StartupOutput, error)
// Creates startup and automatically adds creator as owner
```

**Invite Member Use Case** (`internal/application/usecase/member/invite_member.go`)
```go
type InviteMemberUseCase struct {
    invitationRepo repository.InvitationRepository
    memberRepo     repository.StartupMemberRepository
    startupRepo    repository.StartupRepository
    userRepo       repository.UserRepository
    emailService   port.EmailService
    tokenGen       port.TokenService
}

func (uc *InviteMemberUseCase) Execute(ctx context.Context, input dto.InviteMemberInput, inviterID string) (*dto.InvitationOutput, error)
// Validates inviter has permission, checks if user exists, sends invitation email
```

**Accept Invitation Use Case** (`internal/application/usecase/member/accept_invitation.go`)
```go
type AcceptInvitationUseCase struct {
    invitationRepo repository.InvitationRepository
    memberRepo     repository.StartupMemberRepository
    userRepo       repository.UserRepository
}

func (uc *AcceptInvitationUseCase) Execute(ctx context.Context, token string, userID string) (*dto.MemberOutput, error)
// Validates invitation, creates member record, updates invitation status
```

**Join via Public Link Use Case** (`internal/application/usecase/member/join_startup.go`)
```go
type JoinStartupUseCase struct {
    startupRepo repository.StartupRepository
    memberRepo  repository.StartupMemberRepository
    userRepo    repository.UserRepository
}

func (uc *JoinStartupUseCase) Execute(ctx context.Context, startupSlug string, joinCode string, userID string) (*dto.MemberOutput, error)
// Validates join code (if startup uses it), creates member with default role
```

**Request to Join Use Case** (`internal/application/usecase/member/request_join.go`)
```go
type RequestJoinUseCase struct {
    startupRepo    repository.StartupRepository
    memberRepo     repository.StartupMemberRepository
    userRepo       repository.UserRepository
    notifyService  port.NotificationService
}

func (uc *RequestJoinUseCase) Execute(ctx context.Context, input dto.JoinRequestInput, userID string) (*dto.MemberOutput, error)
// Creates pending member, notifies startup admins/owners
```

**Create Job Use Case** (`internal/application/usecase/job/create_job.go`)
```go
type CreateJobUseCase struct {
    jobRepo     repository.JobRepository
    startupRepo repository.StartupRepository
    memberRepo  repository.StartupMemberRepository
    authService domain.AuthorizationService
}

func (uc *CreateJobUseCase) Execute(ctx context.Context, input dto.CreateJobInput, userID string) (*dto.JobOutput, error)
// Checks if user has permission via member record
```

**Upload File Use Case** (`internal/application/usecase/file/upload_file.go`)
```go
type UploadFileUseCase struct {
    fileRepo       repository.FileRepository
    storageService port.StorageService
}

func (uc *UploadFileUseCase) Execute(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID string) (*dto.FileOutput, error)
```

### DTOs

**Startup DTOs** (`internal/application/dto/startup_dto.go`)
```go
type CreateStartupInput struct {
    Name        string
    Description string
    Website     string
    FoundedYear int
    Industry    string
    CompanySize string
    Location    string
    AllowPublicJoin bool  // Enable public join link
}

type UpdateStartupInput struct {
    ID          string
    Name        *string
    Description *string
    LogoURL     *string
    Website     *string
    AllowPublicJoin *bool
    JoinCode    *string  // Optional code for joining
    // ... other fields
}

type StartupOutput struct {
    ID              string
    Name            string
    Slug            string
    Description     string
    AllowPublicJoin bool
    MemberCount     int
    JobCount        int
    // ... all public fields
}
```

**Member DTOs** (`internal/application/dto/member_dto.go`)
```go
type InviteMemberInput struct {
    StartupID string
    Email     string
    Role      string // "admin", "member", "recruiter"
    Message   string // Optional personal message
}

type InvitationOutput struct {
    ID        string
    Email     string
    Role      string
    Token     string
    ExpiresAt time.Time
    InviteURL string // Full URL to accept invitation
}

type AcceptInvitationInput struct {
    Token string
}

type JoinRequestInput struct {
    StartupSlug string
    JoinCode    string // Optional if startup requires it
    Message     string // Why they want to join
}

type MemberOutput struct {
    ID        string
    UserID    string
    StartupID string
    Role      string
    Status    string
    User      UserSummary
    JoinedAt  *time.Time
}

type UserSummary struct {
    ID    string
    Name  string
    Email string
}

type UpdateMemberRoleInput struct {
    MemberID string
    Role     string
}
```
```

---

## Infrastructure Layer Specifications

### Storage Service

**Storage Service Interface** (`internal/application/port/storage_service.go`)
```go
type StorageService interface {
    Upload(ctx context.Context, file io.Reader, key string, contentType string) (string, error)
    Delete(ctx context.Context, key string) error
    GetURL(ctx context.Context, key string) (string, error)
}
```

**MinIO Implementation** (`internal/infrastructure/storage/minio_storage.go`)
```go
type MinIOStorage struct {
    client     *minio.Client
    bucketName string
}

func (s *MinIOStorage) Upload(ctx context.Context, file io.Reader, key string, contentType string) (string, error)
```

**S3 Implementation** (`internal/infrastructure/storage/s3_storage.go`)
```go
type S3Storage struct {
    client     *s3.Client
    bucketName string
    region     string
}

func (s *S3Storage) Upload(ctx context.Context, file io.Reader, key string, contentType string) (string, error)
```

**Storage Factory** (`internal/infrastructure/storage/storage_factory.go`)
```go
func NewStorageService(config *config.Config) (port.StorageService, error) {
    if config.Environment == "production" {
        return NewS3Storage(config.S3Config)
    }
    return NewMinIOStorage(config.MinIOConfig)
}
```

### Repository Implementation

**Startup Repository Implementation** (`internal/infrastructure/persistence/postgres/startup_repository_impl.go`)
```go
type StartupRepositoryImpl struct {
    db *gorm.DB
}

// Convert between GORM models and domain entities
func (r *StartupRepositoryImpl) toDomain(model *gorm_model.Startup) *entity.Startup
func (r *StartupRepositoryImpl) toModel(startup *entity.Startup) *gorm_model.Startup
```

---

## Presentation Layer Specifications

### HTTP Handlers

**Startup Handler** (`internal/presentation/http/handler/startup_handler.go`)
```go
type StartupHandler struct {
    createUseCase  *usecase.CreateStartupUseCase
    updateUseCase  *usecase.UpdateStartupUseCase
    getUseCase     *usecase.GetStartupUseCase
    listUseCase    *usecase.ListStartupsUseCase
    validator      *validator.Validator
}

func (h *StartupHandler) Create(c *gin.Context)
func (h *StartupHandler) Update(c *gin.Context)
func (h *StartupHandler) Get(c *gin.Context)
func (h *StartupHandler) List(c *gin.Context)
```

### Middleware

**JWT Middleware** (`internal/presentation/http/middleware/auth_middleware.go`)
```go
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc
```

**API Token Middleware** (`internal/presentation/http/middleware/api_token_middleware.go`)
```go
func APITokenMiddleware(startupRepo repository.StartupRepository) gin.HandlerFunc
```

### Response Helpers

**Success Response** (`internal/presentation/http/response/success_response.go`)
```go
type SuccessResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    TotalPages int   `json:"total_pages"`
    TotalCount int64 `json:"total_count"`
}

func Success(c *gin.Context, data interface{})
func SuccessWithMeta(c *gin.Context, data interface{}, meta *Meta)
```

**Error Response** (`internal/presentation/http/response/error_response.go`)
```go
type ErrorResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
    Code    string `json:"code,omitempty"`
}

func Error(c *gin.Context, statusCode int, err error)
func BadRequest(c *gin.Context, message string)
func Unauthorized(c *gin.Context)
func Forbidden(c *gin.Context)
func NotFound(c *gin.Context, resource string)
func InternalError(c *gin.Context)
```

---

## API Endpoints

### Public Endpoints
```
GET    /api/v1/health                 # Health check
GET    /api/v1/jobs                   # List jobs (public)
GET    /api/v1/jobs/:id               # Get job details
GET    /api/v1/startups               # List startups
GET    /api/v1/startups/:slug         # Get startup profile + jobs
```

### Authentication Endpoints
```
POST   /api/v1/auth/register          # Register user
POST   /api/v1/auth/login             # Login
POST   /api/v1/auth/refresh           # Refresh JWT token
```

### Protected Endpoints (JWT)
```
POST   /api/v1/startups               # Create startup (admin)
PUT    /api/v1/startups/:id           # Update startup (owner/admin)
DELETE /api/v1/startups/:id           # Delete startup (admin)
POST   /api/v1/startups/:id/regenerate-token  # Regenerate API token

# Member Management
POST   /api/v1/startups/:id/members/invite     # Invite member (owner/admin)
GET    /api/v1/startups/:id/members            # List members (owner/admin/member)
PUT    /api/v1/startups/:id/members/:memberId  # Update member role (owner/admin)
DELETE /api/v1/startups/:id/members/:memberId  # Remove member (owner/admin)

POST   /api/v1/invitations/:token/accept       # Accept invitation
GET    /api/v1/invitations/my                  # Get my pending invitations
POST   /api/v1/startups/:slug/join             # Request to join (public/with code)

# Jobs
POST   /api/v1/jobs                   # Create job (startup owner/admin/recruiter)
PUT    /api/v1/jobs/:id               # Update job
DELETE /api/v1/jobs/:id               # Delete job

# User
GET    /api/v1/me                     # Get current user info
GET    /api/v1/me/startups            # Get my startups
PUT    /api/v1/me                     # Update profile

POST   /api/v1/upload                 # Upload file (logo, etc.)
```

### API Token Endpoints (for integrations)
```
GET    /api/v1/token/startup          # Get own startup info
PUT    /api/v1/token/startup          # Update own startup
POST   /api/v1/token/jobs             # Create job
PUT    /api/v1/token/jobs/:id         # Update job
DELETE /api/v1/token/jobs/:id         # Delete job
GET    /api/v1/token/jobs             # List own jobs
```

---

## Docker Configuration

### Dockerfile (`docker/Dockerfile`)
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./server"]
```

### Docker Compose (`docker-compose.yml`)
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: startup_board_db
    environment:
      POSTGRES_USER: ${DB_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-postgres}
      POSTGRES_DB: ${DB_NAME:-startup_board}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  minio:
    image: minio/minio:latest
    container_name: startup_board_minio
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ROOT_USER:-minioadmin}
      MINIO_ROOT_PASSWORD: ${MINIO_ROOT_PASSWORD:-minioadmin}
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  api:
    build:
      context: .
      dockerfile: docker/Dockerfile
    container_name: startup_board_api
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER:-postgres}
      - DB_PASSWORD=${DB_PASSWORD:-postgres}
      - DB_NAME=${DB_NAME:-startup_board}
      - STORAGE_TYPE=minio
      - MINIO_ENDPOINT=minio:9000
      - MINIO_ACCESS_KEY=${MINIO_ROOT_USER:-minioadmin}
      - MINIO_SECRET_KEY=${MINIO_ROOT_PASSWORD:-minioadmin}
      - JWT_SECRET=${JWT_SECRET:-your-secret-key-change-in-production}
      - GIN_MODE=release
    depends_on:
      postgres:
        condition: service_healthy
      minio:
        condition: service_healthy
    restart: unless-stopped

volumes:
  postgres_data:
  minio_data:
```

### Development Docker Compose (`docker-compose.dev.yml`)
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: startup_board_dev
    ports:
      - "5432:5432"
    volumes:
      - ./scripts/init-db.sh:/docker-entrypoint-initdb.d/init-db.sh

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_dev_data:/data

volumes:
  minio_dev_data:
```

---

## Environment Variables (`.env.example`)

```env
# Server
PORT=8080
GIN_MODE=debug
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=startup_board
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-very-secure-secret-key-change-this
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h

# Storage (Development - MinIO)
STORAGE_TYPE=minio
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=startup-board-uploads
MINIO_USE_SSL=false

# Storage (Production - S3)
# STORAGE_TYPE=s3
# S3_BUCKET=startup-board-production
# S3_REGION=us-east-1
# AWS_ACCESS_KEY_ID=your-access-key
# AWS_SECRET_ACCESS_KEY=your-secret-key

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

---

## Makefile

```makefile
.PHONY: help build run test docker-up docker-down migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make docker-up   - Start Docker containers"
	@echo "  make docker-down - Stop Docker containers"
	@echo "  make migrate-up  - Run database migrations"

build:
	go build -o bin/server cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

docker-rebuild:
	docker-compose up -d --build

dev:
	docker-compose -f docker-compose.dev.yml up -d
	air

migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/startup_board?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/startup_board?sslmode=disable" down

lint:
	golangci-lint run

.DEFAULT_GOAL := help
```

---

## Key Implementation Notes

### Team/Member Management Flows

**Flow 1: Invitation Flow (Most Secure)**
1. Startup owner/admin invites via email
2. System generates unique invitation token (expires in 7 days)
3. Email sent with invitation link: `https://app.com/invitations/{token}/accept`
4. User clicks link (must be logged in or creates account)
5. System validates token, creates member record with specified role
6. User now has access to startup dashboard

**Flow 2: Public Join Link (Easiest)**
1. Startup enables "Allow Public Join" in settings
2. System generates join link: `https://app.com/startups/{slug}/join`
3. Startup shares link on their website/social media
4. User clicks link, logs in, clicks "Join Team"
5. System creates member record with default role (usually "member")
6. Optional: Startup can set a join code for semi-private access

**Flow 3: Join Request Flow (Moderated)**
1. User finds startup page, clicks "Request to Join"
2. User submits request with optional message
3. System creates pending member record
4. Startup owners/admins receive notification
5. Owner/admin approves/rejects from dashboard
6. If approved, member status changes to active

**Flow 4: Direct Add (Admin)**
1. Platform admin can directly add users to any startup
2. Useful for onboarding or support scenarios

### Authorization Service Implementation

```go
// internal/domain/service/authorization_service.go
type AuthorizationService struct {
    memberRepo repository.StartupMemberRepository
}

func (s *AuthorizationService) CanManageStartup(ctx context.Context, userID, startupID string) (bool, error) {
    member, err := s.memberRepo.FindByUserAndStartup(ctx, userID, startupID)
    if err != nil {
        return false, err
    }
    return member != nil && member.CanManageStartup(), nil
}

func (s *AuthorizationService) CanManageJobs(ctx context.Context, userID, startupID string) (bool, error) {
    member, err := s.memberRepo.FindByUserAndStartup(ctx, userID, startupID)
    if err != nil {
        return false, err
    }
    return member != nil && member.CanManageJobs(), nil
}

func (s *AuthorizationService) CanManageMembers(ctx context.Context, userID, startupID string) (bool, error) {
    member, err := s.memberRepo.FindByUserAndStartup(ctx, userID, startupID)
    if err != nil {
        return false, err
    }
    return member != nil && member.CanManageMembers(), nil
}

func (s *AuthorizationService) GetUserStartups(ctx context.Context, userID string) ([]*entity.StartupMember, error) {
    return s.memberRepo.FindByUserID(ctx, userID)
}
```

### Email Service Interface

```go
// internal/application/port/email_service.go
type EmailService interface {
    SendInvitationEmail(ctx context.Context, invitation *entity.Invitation, startup *entity.Startup) error
    SendJoinRequestNotification(ctx context.Context, member *entity.StartupMember, startup *entity.Startup) error
    SendMemberApprovedEmail(ctx context.Context, member *entity.StartupMember, startup *entity.Startup) error
}
```

### Database Schema Updates

```sql
-- Add to migrations
CREATE TABLE startup_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    startup_id UUID NOT NULL REFERENCES startups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    invited_by UUID REFERENCES users(id),
    invited_at TIMESTAMP NOT NULL DEFAULT NOW(),
    joined_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(startup_id, user_id),
    CHECK (role IN ('owner', 'admin', 'member', 'recruiter')),
    CHECK (status IN ('active', 'pending', 'removed'))
);

CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    startup_id UUID NOT NULL REFERENCES startups(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    invited_by UUID NOT NULL REFERENCES users(id),
    expires_at TIMESTAMP NOT NULL,
    accepted_at TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (role IN ('admin', 'member', 'recruiter')),
    CHECK (status IN ('pending', 'accepted', 'expired', 'revoked'))
);

-- Add to startups table
ALTER TABLE startups ADD COLUMN allow_public_join BOOLEAN DEFAULT false;
ALTER TABLE startups ADD COLUMN join_code VARCHAR(50);

-- Indexes
CREATE INDEX idx_startup_members_user_id ON startup_members(user_id);
CREATE INDEX idx_startup_members_startup_id ON startup_members(startup_id);
CREATE INDEX idx_invitations_token ON invitations(token);
CREATE INDEX idx_invitations_email ON invitations(email);
CREATE INDEX idx_invitations_startup_id ON invitations(startup_id);
```

### Dependency Injection
Use a DI container or manual injection in `main.go`:
```go
// Wire up dependencies
db := infrastructure.NewDatabase(config)
startupRepo := postgres.NewStartupRepository(db)
jobRepo := postgres.NewJobRepository(db)
storageService := storage.NewStorageService(config)

createStartupUC := usecase.NewCreateStartupUseCase(startupRepo, tokenGen)
startupHandler := handler.NewStartupHandler(createStartupUC, ...)

router := router.NewRouter(startupHandler, jobHandler, ...)
```

### Error Handling
Create domain-specific errors in `pkg/errors/app_error.go`:
```go
type AppError struct {
    Code    string
    Message string
    Err     error
}

var (
    ErrNotFound      = &AppError{Code: "NOT_FOUND", Message: "Resource not found"}
    ErrUnauthorized  = &AppError{Code: "UNAUTHORIZED", Message: "Unauthorized"}
    ErrForbidden     = &AppError{Code: "FORBIDDEN", Message: "Forbidden"}
)
```

### Validation
Use `go-playground/validator` with custom validators:
```go
type CreateJobInput struct {
    Title       string `json:"title" validate:"required,min=5,max=100"`
    JobType     string `json:"job_type" validate:"required,oneof=full_time part_time contract internship"`
    SalaryMin   *int   `json:"salary_min" validate:"omitempty,min=0"`
}
```

### Database Migrations
Include SQL migrations in `migrations/` directory for version control and rollback capability.

### Testing Strategy
- Unit tests for domain entities and use cases
- Integration tests for repositories
- E2E tests for API endpoints
- Mock interfaces for external dependencies

This architecture ensures:
- ✅ Clear separation of concerns
- ✅ Easy to test and maintain
- ✅ Scalable and extensible for future ATS features
- ✅ Infrastructure-agnostic domain logic
- ✅ Consistent error handling and responses