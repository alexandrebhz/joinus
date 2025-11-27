# JoinUs Frontend

A modern Next.js frontend for the JoinUs startup job board platform, built with Clean Architecture principles and DDD.

## Tech Stack

- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **Forms**: React Hook Form + Zod
- **HTTP Client**: Axios
- **Icons**: Lucide React

## Architecture

This project follows Clean Architecture and Domain-Driven Design principles:

```
src/
├── domain/              # Domain Layer (Entities & Value Objects)
│   ├── entities/
│   └── value-objects/
├── application/         # Application Layer (Use Cases & DTOs)
│   ├── dto/
│   └── ports/
├── infrastructure/      # Infrastructure Layer (External Services)
│   ├── api/
│   └── store/
└── presentation/        # Presentation Layer (UI Components)
    └── components/
```

### Layer Responsibilities

- **Domain**: Pure business logic, no external dependencies
- **Application**: Use cases, DTOs, and application interfaces
- **Infrastructure**: API clients, state management, external services
- **Presentation**: React components, pages, UI elements

## Getting Started

### Prerequisites

- Node.js 20+
- npm or yarn

### Installation

```bash
npm install
```

### Development

```bash
npm run dev
```

The application will be available at [http://localhost:3000](http://localhost:3000)

### Docker Development

The frontend is integrated into the main docker-compose setup:

```bash
# From the backend directory
docker-compose -f docker-compose.dev.yml up frontend
```

Or run the entire stack:

```bash
docker-compose -f docker-compose.dev.yml up
```

### Production Build

```bash
npm run build
npm start
```

## Environment Variables

Create a `.env.local` file (see `.env.example`):

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Project Structure

### Pages

- `/` - Homepage with featured jobs and startups
- `/jobs` - Job listings with search and filters
- `/startups` - Startup directory
- `/login` - User authentication
- `/register` - User registration

### Key Components

- **UI Components**: Reusable components in `presentation/components/ui/`
- **Layout Components**: Header, Footer in `presentation/components/layout/`
- **Feature Components**: JobCard, StartupCard in `presentation/components/`

### API Client

The API client (`infrastructure/api/api-client.ts`) handles:
- HTTP requests to the backend
- Authentication token management
- Automatic token refresh
- Error handling

## Features

- ✅ Clean Architecture & DDD
- ✅ Type-safe API client
- ✅ Authentication with JWT
- ✅ Responsive design
- ✅ Modern UI inspired by top job boards
- ✅ Docker support
- ✅ TypeScript throughout

## Development Guidelines

1. **Domain Layer**: Keep pure, no React or external dependencies
2. **Application Layer**: Define use cases and DTOs
3. **Infrastructure**: Implement external integrations
4. **Presentation**: React components only

## Future Enhancements

- Dashboard for authenticated users
- Job application flow
- Startup profile management
- Advanced filtering and search
- Real-time notifications

