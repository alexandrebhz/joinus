# Frontend - Job Crawler System

A Next.js frontend application built with Clean Architecture, DDD principles, and DRY practices.

## Architecture

The frontend follows Clean Architecture principles with clear separation of concerns:

```
frontend/
├── src/
│   ├── domain/              # Domain Layer (Pure Business Logic)
│   │   ├── entities/        # Domain entities
│   │   ├── value-objects/   # Value objects
│   │   └── repositories/    # Repository interfaces
│   ├── application/         # Application Layer (Use Cases)
│   │   └── use-cases/       # Business use cases
│   ├── infrastructure/      # Infrastructure Layer
│   │   ├── http/           # HTTP client implementations
│   │   └── di/             # Dependency injection container
│   └── presentation/        # Presentation Layer
│       ├── components/      # React components
│       └── hooks/          # Custom React hooks
└── app/                     # Next.js app directory
```

## Key Principles

### Domain-Driven Design (DDD)
- **Entities**: Core business objects (`CrawlSite`, `CrawlResult`)
- **Value Objects**: Immutable data structures (`CreateSiteInput`, `UpdateSiteInput`)
- **Repositories**: Interfaces defining data access contracts

### Clean Architecture
- **Dependency Rule**: Dependencies point inward (Domain ← Application ← Infrastructure ← Presentation)
- **Separation of Concerns**: Each layer has a single responsibility
- **Testability**: Easy to mock dependencies and test in isolation

### DRY (Don't Repeat Yourself)
- **Reusable Components**: `SiteCard`, `SiteForm` can be used across pages
- **Custom Hooks**: `useSites` encapsulates site management logic
- **Shared Utilities**: Common logic extracted to use cases

## Project Structure

### Domain Layer
- **Entities**: Pure TypeScript interfaces representing domain concepts
- **Value Objects**: Input/output DTOs with validation
- **Repositories**: Interfaces for data access (no implementation details)

### Application Layer
- **Use Cases**: Business logic orchestration
- Each use case handles a single business operation
- Validates domain rules before delegating to repositories

### Infrastructure Layer
- **HTTP Client**: Implements repository interfaces using fetch API
- **DI Container**: Manages dependency injection and singleton instances

### Presentation Layer
- **Components**: Reusable UI components
- **Hooks**: Custom React hooks for state management
- **Pages**: Next.js pages that compose components and hooks

## Usage

### Development

```bash
npm run dev
```

### Build

```bash
npm run build
npm start
```

### Docker

```bash
docker-compose up frontend
```

## Environment Variables

Create a `.env.local` file:

```env
NEXT_PUBLIC_API_URL=http://localhost:8081
```

## Key Features

- ✅ Clean Architecture with proper layer separation
- ✅ DDD principles with entities and value objects
- ✅ Dependency Injection container
- ✅ Reusable components and hooks
- ✅ Type-safe with TypeScript
- ✅ Responsive UI with Tailwind CSS
