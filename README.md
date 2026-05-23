# Social Network

A full-stack social networking platform featuring real-time messaging, posts and comments, follow relationships, groups with events, and live notifications.

Built with a **Go** backend (Gorilla Mux, WebSockets, SQLite) and a modern **Next.js 16 / React 19** frontend in TypeScript with Tailwind CSS.


---

## Table of Contents

- [Getting Started](#getting-started)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Available Scripts](#available-scripts)
- [API Overview](#api-overview)


---
## Getting Started

The repo ships with a multi-stage `Dockerfile` for each service and a root `docker-compose.yml` that wires them together. SQLite data and uploads are persisted in a named volume (`backend-data`), so they survive container restarts.

### Prerequisites

- **Docker** ≥ 24
- **Docker Compose** v2 (bundled with modern Docker Desktop / Docker Engine)

### 1. Clone

```bash
git clone https://platform.zone01.gr/git/cgkaldan/social-network
cd social-network
```

### 2. Start everything

From the repo root:

```bash
docker compose up --build
```

- Frontend: [http://localhost:3000](http://localhost:3000)
- Backend:  [http://localhost:8000](http://localhost:8000)

Database migrations are applied automatically on backend startup.

Run it detached with `docker compose up -d --build`, stop with `docker compose down`, and stop **and** wipe the SQLite volume with `docker compose down -v`.

---

## Features

### Identity & Profiles
- Email + password registration with bcrypt-hashed credentials
- Session-based authentication
- Public / private profile visibility
- Avatar upload, bio, and personal details
- Online presence indicators driven by WebSocket activity

### Posts & Comments
- Create posts with categories and rich content
- Per-post privacy levels: `public`, `almost_private`, `private`
- Image / media attachments
- Threaded comments with media support
- Like counts and comment counts
- Full CRUD: create, read, update, delete on posts and comments

### Social Graph
- Follow / unfollow users
- Follow requests for private accounts
- Followers and following lists
- Tailored feed based on connections

### Groups
- Create groups and invite members
- Group posts and discussions
- Group events with RSVP
- Update, leave, and delete groups
- Full CRUD on group events

### Real-time Messaging
- Private and group conversations
- WebSocket-powered live message delivery
- Read receipts and unread counts
- Conversation list with last-message preview
- Infinite-scroll message history
- Emoji picker

### Notifications
- Real-time notification stream
- Notification count badge
- Per-event types: follow, message, group invite, event, etc.

### UX
- Responsive layout
- Light / dark theme toggle

---

## Tech Stack

### Backend
| Layer | Choice |
|------|--------|
| Language | Go 1.24 |
| HTTP router | [`gorilla/mux`](https://github.com/gorilla/mux) |
| WebSockets | [`gorilla/websocket`](https://github.com/gorilla/websocket) |
| Database | SQLite via [`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) |
| Migrations | [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate) |
| Auth crypto | `golang.org/x/crypto/bcrypt` |

### Frontend
| Layer | Choice |
|------|--------|
| Framework | Next.js 16 (App Router) |
| UI runtime | React 19 |
| Language | TypeScript |
| Styling | Tailwind CSS 4 |
| Lint | ESLint (flat config) |

---

## Architecture

The backend follows a clean, layered architecture:

```
HTTP / WebSocket
       │
       ▼
┌────────────────┐
│   Handlers     │  (internal/api/handlers)
└────────────────┘
       │
       ▼
┌────────────────┐
│   Services     │  (internal/service)  — business logic
└────────────────┘
       │
       ▼
┌────────────────┐
│ Repositories   │  (internal/repository) — persistence
└────────────────┘
       │
       ▼
┌────────────────┐
│   SQLite DB    │  (managed via migrations)
└────────────────┘
```

- **Domain models** live in `internal/domain` and are shared across layers.
- **Events** are dispatched through `internal/event` and consumed by `internal/consumer` for fan-out work like real-time notifications.
- **WebSocket hub** in `internal/websocket` brokers presence and live messages.
- **Middleware** in `internal/api/middleware` covers auth, CORS, rate limiting, and request logging.

The frontend uses the Next.js App Router with route groups:

- `(auth)` — login / register pages without the authenticated shell.
- `(app)` — protected routes wrapped in `AppShell`, including feed, profile, groups, messages, notifications, follow, posts, and users pages.
- Cross-cutting providers (`AuthProvider`, `ThemeProvider`, `WebSocketProvider`, `NotificationCountProvider`, `MessageCountProvider`) supply state to the tree.

---

## Project Structure

```
social-network/
├── backend/
│   ├── cmd/
│   │   └── server/              # main.go entry point
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/        # HTTP & WS handlers
│   │   │   ├── middleware/      # auth, CORS, rate-limit
│   │   │   └── router/          # route wiring
│   │   ├── config/              # env loading
│   │   ├── consumer/            # event consumers
│   │   ├── database/
│   │   │   └── migrations/      # SQL migrations
│   │   ├── domain/              # core types
│   │   ├── event/               # event bus
│   │   ├── repository/          # SQLite data access
│   │   ├── service/             # business logic
│   │   └── websocket/           # WS hub
│   ├── packages/
│   │   └── logger/              # structured logging
│   ├── data/                    # SQLite db & uploads (gitignored)
│   ├── go.mod
│   └── Makefile
│
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── (auth)/          # login, register
│   │   │   └── (app)/           # feed, groups, messages, profile, ...
│   │   ├── components/          # shared UI
│   │   ├── hooks/               # custom React hooks
│   │   ├── lib/                 # api client, ws client, auth, helpers
│   │   └── types/               # shared TS types
│   ├── package.json
│   └── next.config.ts
│
├── start.md                     # quickstart reference
└── README.md
```

---

## Configuration

### Backend (`backend/.env`)

| Variable | Default | Description |
|---|---|---|
| `ENVIRONMENT` | `development` | Runtime environment |
| `SERVER_PORT` | `8000` | HTTP port |
| `SERVER_READ_TIMEOUT` | `5s` | HTTP read timeout |
| `SERVER_WRITE_TIMEOUT` | `15s` | HTTP write timeout |
| `SERVER_IDLE_TIMEOUT` | `60s` | HTTP idle timeout |
| `DATABASE_PATH` | `./data/database/.db` | SQLite file path |
| `SESSION_DURATION` | `24h` | Session lifetime |
| `RATE_LIMIT_ENABLED` | `true` | Toggle rate limiting |
| `RATE_LIMIT_RPM` | `100` | Requests per minute per client |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:3000` | Comma-separated origins |
| `WS_READ_BUFFER_SIZE` | `1024` | WebSocket read buffer |
| `WS_WRITE_BUFFER_SIZE` | `1024` | WebSocket write buffer |
| `WS_PING_PERIOD` | `54s` | WS ping interval |
| `WS_PONG_WAIT` | `60s` | WS pong timeout |
| `WS_WRITE_WAIT` | `10s` | WS write timeout |
| `UPLOAD_PATH` | `./data/uploads` | Upload destination |
| `MAX_FILE_SIZE` | `5242880` | Max upload size (bytes) |

### Frontend (`frontend/.env`)

| Variable | Default | Description |
|---|---|---|
| `NEXT_PUBLIC_API_URL` | `http://localhost:8000` | Backend HTTP base URL |
| `NEXT_PUBLIC_WS_URL` | `ws://localhost:8000/ws` | Backend WebSocket URL |

---

## Available Scripts

### Backend (`backend/Makefile`)

| Command | Description |
|---|---|
| `make run` | Run the server |
| `make build` | Build a static binary into `bin/forum-backend` |
| `make test` | Run all tests with race detector and coverage report |
| `make migrate` | Apply pending migrations |
| `make migrate-down` | Roll back the last migration |
| `make lint` | Run linter |
| `make fmt` | Format code |
| `make clean` | Remove build artifacts |
| `make docker-build` | Build a Docker image |
| `make docker-run` | Run the Docker container on port 8000 |

### Frontend (`frontend/package.json`)

| Command | Description |
|---|---|
| `npm run dev` | Start the Next.js dev server |
| `npm run build` | Build the production bundle |
| `npm run start` | Run the production server |
| `npm run lint` | Lint the codebase |

---

## API Overview

The backend exposes a versioned REST API plus a single WebSocket endpoint. Authenticated endpoints require a valid session cookie.

| Domain | Example endpoints |
|---|---|
| Auth | `POST /api/auth/register`, `POST /api/auth/login`, `POST /api/auth/logout`, `GET /api/auth/me` |
| Users | `GET /api/users`, `GET /api/users/{id}`, `PATCH /api/users/{id}` |
| Posts | `GET /api/posts`, `POST /api/posts`, `PATCH /api/posts/{id}`, `DELETE /api/posts/{id}` |
| Comments | `GET /api/posts/{id}/comments`, `POST /api/posts/{id}/comments` |
| Follows | `POST /api/follow/{id}`, `DELETE /api/follow/{id}`, `GET /api/follow/requests` |
| Groups | `GET /api/groups`, `POST /api/groups`, `POST /api/groups/{id}/events` |
| Messages | `GET /api/conversations`, `GET /api/conversations/{id}/messages` |
| Notifications | `GET /api/notifications` |
| Uploads | `POST /api/uploads` |
| Health | `GET /health` |
| WebSocket | `GET /ws` (presence + live messages + notifications) |

> Routes are defined in [`backend/internal/api/router`](backend/internal/api/router) and handlers in [`backend/internal/api/handlers`](backend/internal/api/handlers).

## Authors
- Christos Gkaldanidis
- Christos Markos

Creators and primary Developers