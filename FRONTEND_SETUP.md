# Frontend Setup Guide

## Overview

The JoinUs frontend is a Next.js 14 application built with Clean Architecture and DDD principles, featuring a modern UI inspired by top job boards like RemoteOK, Resend, and Gupy.

## Architecture

The frontend follows Clean Architecture with clear separation of concerns:

```
src/
├── domain/              # Pure business logic (entities, value objects)
├── application/         # Use cases, DTOs, ports
├── infrastructure/      # API client, state management
└── presentation/        # React components, pages
```

## Quick Start

### Local Development (Without Docker)

1. **Install dependencies:**
   ```bash
   cd frontend
   npm install
   ```

2. **Set up environment:**
   ```bash
   cp .env.example .env.local
   # Edit .env.local and set NEXT_PUBLIC_API_URL=http://localhost:8080
   ```

3. **Start development server:**
   ```bash
   npm run dev
   ```

   The frontend will be available at http://localhost:3000

### Docker Development

The frontend is integrated into the main docker-compose setup:

```bash
# From the backend directory
cd backend
docker-compose -f docker-compose.dev.yml up
```

This will start:
- PostgreSQL (port 5432)
- MinIO (ports 9000, 9001)
- Backend API (port 8080)
- Frontend (port 3000)

The frontend will automatically connect to the API at `http://api:8080` (internal Docker network).

## Features

### ✅ Implemented

- **Clean Architecture**: Proper layer separation with DDD principles
- **Authentication**: Login, register, JWT token management
- **Job Listings**: Browse and search jobs
- **Startup Directory**: Explore startup profiles
- **Responsive Design**: Mobile-first, works on all devices
- **Modern UI**: Inspired by RemoteOK, Resend, Apple, Gupy
- **Type Safety**: Full TypeScript coverage
- **Docker Support**: Development and production Dockerfiles

### 🚧 Future Enhancements

- User dashboard with job applications
- Startup profile management
- Job creation and editing
- Advanced filtering and search
- Real-time notifications
- Application tracking

## Project Structure

### Pages

- `/` - Homepage with featured content
- `/jobs` - Job listings with search
- `/jobs/[id]` - Job detail page
- `/startups` - Startup directory
- `/startups/[slug]` - Startup profile page
- `/login` - User authentication
- `/register` - User registration
- `/dashboard` - User dashboard (protected)

### Key Components

**UI Components** (`presentation/components/ui/`):
- `Button` - Versatile button component
- `Input` - Form input with validation
- `Card` - Card container with variants

**Layout Components** (`presentation/components/layout/`):
- `Header` - Navigation header
- `Footer` - Site footer

**Feature Components**:
- `JobCard` - Job listing card
- `StartupCard` - Startup profile card

## API Integration

The frontend uses a centralized API client (`infrastructure/api/api-client.ts`) that:

- Handles all HTTP requests
- Manages authentication tokens
- Automatically refreshes expired tokens
- Provides type-safe API methods
- Handles errors gracefully

### Example Usage

```typescript
import { apiClient } from '@/infrastructure/api/api-client'

// List jobs
const response = await apiClient.listJobs({ page: 1, page_size: 10 })

// Create job (requires authentication)
const job = await apiClient.createJob({
  startup_id: '...',
  title: 'Software Engineer',
  // ... other fields
})
```

## State Management

Authentication state is managed with Zustand:

```typescript
import { useAuthStore } from '@/infrastructure/store/auth.store'

const { user, isAuthenticated, setUser, logout } = useAuthStore()
```

## Styling

The project uses Tailwind CSS with a custom design system:

- **Primary Color**: Blue (#0066FF)
- **Secondary Colors**: Gray scale
- **Typography**: Inter font family
- **Spacing**: Consistent 4px base unit
- **Components**: Custom styled with hover effects

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://localhost:8080` |
| `NODE_ENV` | Environment mode | `development` |

## Building for Production

```bash
npm run build
npm start
```

For Docker production builds, use the `Dockerfile` (not `Dockerfile.dev`).

## Development Tips

1. **Hot Reload**: Next.js provides instant feedback during development
2. **Type Safety**: TypeScript catches errors at compile time
3. **API Mocking**: Consider using MSW for API mocking in tests
4. **Component Isolation**: Each component is self-contained and reusable

## Troubleshooting

### Frontend can't connect to API

- **Docker**: Ensure `NEXT_PUBLIC_API_URL=http://api:8080`
- **Local**: Ensure `NEXT_PUBLIC_API_URL=http://localhost:8080` and API is running

### Authentication not working

- Check that tokens are being stored in localStorage
- Verify API CORS settings allow your frontend origin
- Check browser console for errors

### Build errors

- Clear `.next` directory: `rm -rf .next`
- Reinstall dependencies: `rm -rf node_modules && npm install`
- Check TypeScript errors: `npm run type-check`

## Contributing

When adding new features:

1. Follow Clean Architecture principles
2. Keep domain layer pure (no React dependencies)
3. Use TypeScript for all new code
4. Write reusable components
5. Follow existing code style and patterns

## Resources

- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [TypeScript](https://www.typescriptlang.org/docs)
- [React Hook Form](https://react-hook-form.com)
- [Zod](https://zod.dev)

