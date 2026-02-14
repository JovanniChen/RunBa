# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

- `npm run dev` - Start development server on port 3000
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint
- `npm install` - Install dependencies

## Architecture Overview

This is a Next.js 15 App Router application for Steam account management with modern React patterns:

### Frontend Stack
- **Next.js 15** with App Router and TypeScript
- **TailwindCSS + Shadcn/ui** for styling and components
- **React Hook Form + Zod** for form handling and validation
- **Axios** for HTTP requests with Steam ID transformation
- **React Context** for authentication state management

### Key Architecture Patterns

#### API Layer (`src/lib/api.ts`)
- Axios instances with JWT auth interceptors and Steam ID transformation
- Comprehensive API modules organized by domain (auth, user, points, activation, proxy, config)
- Custom JSON transformer handles large Steam IDs (15+ digits) by converting to strings
- Separate `publicRequest` instance for unauthenticated Steam CDK endpoints

#### Authentication (`src/lib/auth.tsx`)
- React Context provider for global auth state
- localStorage persistence with automatic initialization
- Route guard components: `ProtectedRoute` and `PublicRoute`
- Automatic redirect handling for auth state changes

#### App Router Structure (`src/app/`)
- `/login`, `/register` - Public authentication pages
- `/steam-cdk` - Public Steam CDK exchange interface
- Protected management routes: `/users`, `/certificates`, `/activation-record`, etc.
- All protected routes use `ProtectedRoute` wrapper and `DashboardLayout`

#### Component Architecture
- `DashboardLayout` - Main layout with responsive sidebar navigation
- Shadcn/ui components for consistent design system
- Form components with React Hook Form + Zod validation
- Table components with sorting, filtering, and actions

### Backend Integration
- Backend runs on `http://localhost:8080/api/v1`
- All authenticated requests include `Bearer` token in Authorization header
- 401 responses trigger automatic logout and redirect
- Steam ID handling prevents JavaScript precision loss

### Data Handling
- Steam IDs transformed from numbers to strings in response interceptors
- File exports use blob responses with proper Authorization headers
- Form validation with Zod schemas for type safety
- Loading states and error handling throughout UI

### Authentication Flow
- JWT tokens stored in localStorage with Context state sync
- Route guards prevent unauthorized access
- Automatic cleanup on 401 responses
- Seamless redirect handling between public/private routes

## Important Implementation Notes

### Steam ID Handling
The app includes comprehensive Steam ID transformation logic in both `request` and `publicRequest` instances. When working with Steam-related APIs, ensure large IDs are handled as strings to prevent precision loss.

### Route Structure
- Public routes: `/login`, `/register`, `/steam-cdk` (no auth required)
- Protected routes: All management interfaces under dashboard layout
- Root `/` redirects based on auth status

### Form Patterns
All forms use React Hook Form with Zod validation:
- Schema definition with TypeScript inference
- Consistent error handling and loading states
- Form state managed with `useForm` hook

### API Organization
APIs organized by domain with consistent patterns:
- CRUD operations follow RESTful conventions
- Status updates use PATCH with boolean payloads
- Search/filtering via query parameters
- File operations handle blob responses

### Component Patterns
- Use `'use client'` directive for interactive components
- Shadcn/ui components provide consistent styling
- Loading and error states handled at component level
- Responsive design with mobile-first approach

### Development Workflow
- Hot reload enabled with Next.js dev server
- TypeScript provides compile-time type checking
- ESLint for code quality
- No testing framework configured currently

When adding new features, follow the established patterns for API integration, form handling, authentication, and component architecture. Maintain consistency with existing code style and organization.
