# Architecture Overview

This document describes the architecture of the Job Crawler System, which follows **Domain-Driven Design (DDD)**, **Clean Architecture**, and **DRY** principles.

## Project Structure

```
crawler/
├── backend/          # Go backend API (Clean Architecture)
├── frontend/         # Next.js frontend (Clean Architecture)
├── docker-compose.yml
├── Makefile
└── README.md
```

## Backend Architecture

The backend follows Clean Architecture with clear layer separation:

```
backend/
├── cmd/
│   └── crawler/      # Application entry point
├── internal/
│   ├── domain/       # Domain Layer (Pure Business Logic)
│   │   ├── entity/   # Domain entities
│   │   └── repository/ # Repository interfaces
│   ├── application/  # Application Layer (Use Cases)
│   │   ├── dto/      # Data Transfer Objects
│   │   └── usecase/  # Business use cases
│   ├── infrastructure/ # Infrastructure Layer
│   │   ├── persistence/ # Database implementations
│   │   ├── crawler/    # Crawler engine
│   │   ├── scheduler/  # Cron scheduler
│   │   └── sync/      # Backend API sync
│   └── presentation/  # Presentation Layer
│       └── http/       # HTTP handlers & router
└── migrations/        # Database migrations
```

### Layer Responsibilities

1. **Domain Layer**: Pure business logic, no external dependencies
2. **Application Layer**: Orchestrates use cases, validates business rules
3. **Infrastructure Layer**: Implements external concerns (DB, HTTP, crawler)
4. **Presentation Layer**: HTTP API handlers

## Frontend Architecture

The frontend follows Clean Architecture principles adapted for React/Next.js:

```
frontend/
├── src/
│   ├── domain/              # Domain Layer
│   │   ├── entities/        # Domain entities (CrawlSite, CrawlResult)
│   │   ├── value-objects/   # Value objects (CreateSiteInput, UpdateSiteInput)
│   │   └── repositories/    # Repository interfaces
│   ├── application/         # Application Layer
│   │   └── use-cases/       # Business use cases
│   ├── infrastructure/      # Infrastructure Layer
│   │   ├── http/           # HTTP client (ApiClient)
│   │   └── di/             # Dependency injection container
│   └── presentation/        # Presentation Layer
│       ├── components/      # React components
│       └── hooks/          # Custom React hooks
└── app/                     # Next.js app directory (pages)
```

### Layer Responsibilities

1. **Domain Layer**: TypeScript interfaces and types representing business concepts
2. **Application Layer**: Use cases that orchestrate business logic
3. **Infrastructure Layer**: HTTP client implementing repository interfaces
4. **Presentation Layer**: React components and hooks for UI

## Key Principles

### Domain-Driven Design (DDD)

- **Entities**: Core business objects with identity (CrawlSite, CrawledJob)
- **Value Objects**: Immutable data structures (CreateSiteInput, UpdateSiteInput)
- **Repositories**: Interfaces defining data access contracts
- **Use Cases**: Business operations encapsulated in use case classes

### Clean Architecture

- **Dependency Rule**: Dependencies point inward
  - Domain ← Application ← Infrastructure ← Presentation
- **Separation of Concerns**: Each layer has a single responsibility
- **Testability**: Easy to mock dependencies and test in isolation
- **Framework Independence**: Domain logic independent of frameworks

### DRY (Don't Repeat Yourself)

- **Reusable Components**: Shared UI components (SiteCard, SiteForm)
- **Custom Hooks**: Encapsulated state management (useSites)
- **Use Cases**: Single source of truth for business operations
- **Repository Pattern**: Centralized data access logic

## Dependency Injection

Both backend and frontend use dependency injection:

- **Backend**: Manual DI in `main.go`
- **Frontend**: DI Container in `src/infrastructure/di/container.ts`

## Data Flow

### Backend Request Flow
```
HTTP Request → Router → Handler → Use Case → Repository → Database
```

### Frontend Data Flow
```
Component → Hook → Use Case → Repository → API Client → Backend API
```

## Benefits

1. **Maintainability**: Clear separation makes code easy to understand and modify
2. **Testability**: Each layer can be tested independently
3. **Scalability**: Easy to add new features without affecting existing code
4. **Flexibility**: Can swap implementations (e.g., different HTTP clients)
5. **No Spaghetti Code**: Clear dependencies and responsibilities

## Best Practices

- ✅ Single Responsibility Principle
- ✅ Dependency Inversion Principle
- ✅ Interface Segregation
- ✅ Don't Repeat Yourself (DRY)
- ✅ Type Safety (TypeScript/Go)
- ✅ Error Handling at appropriate layers
- ✅ Validation in use cases

