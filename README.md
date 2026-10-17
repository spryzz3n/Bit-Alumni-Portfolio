# Alumni Portal

College alumni platform. See `alumni_portal_implementation_specification.docx` for scope and roadmap.
Current status: **Milestone 1 (Foundation)** — React ↔ Go ↔ MySQL connectivity only.

## Prerequisites
- Go 1.22+
- Node 20+
- MySQL 8.x

## Setup

### 1. Database
```sql
CREATE DATABASE alumni_portal CHARACTER SET utf8mb4;
CREATE USER 'alumni'@'localhost' IDENTIFIED BY 'change-me';
GRANT ALL PRIVILEGES ON alumni_portal.* TO 'alumni'@'localhost';
```

### 2. Server
The server reads configuration from environment variables (see `.env.example`).
PowerShell:
```powershell
cd server
$env:DB_USER="alumni"; $env:DB_PASSWORD="change-me"; $env:DB_NAME="alumni_portal"
go run ./cmd/migrate up
go run ./cmd/api
```
Bash:
```bash
cd server
export DB_USER=alumni DB_PASSWORD=change-me DB_NAME=alumni_portal
go run ./cmd/migrate up
go run ./cmd/api
```

Endpoints: `GET /api/v1/health` (no DB), `GET /api/v1/db-check` (queries MySQL).

### 3. Client
```bash
cd client
npm install
npm run dev
```
Open http://localhost:5173. Optionally set `VITE_API_URL` (see `client/.env.example`).

## Tests
```bash
cd server && go test ./...
cd client && npm run lint && npm run build
```

## Layout
- `client/` — Vite + React + TypeScript + Tailwind
- `server/` — Go + Chi API (`cmd/api`, `cmd/migrate`, `internal/*`)
- `migrations/` — versioned SQL migrations
