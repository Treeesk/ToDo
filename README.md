# ProjectGo

ProjectGo is a notes and to-do web application with a Go backend and a Vite React frontend. Users can register, log in, and manage their own private notes through cookie-based authentication.

## Screenshots

Add product screenshots in `docs/screenshots/` and update the image paths below if needed.

| Landing page | Auth screen |
| --- | --- |
| ![Landing page screenshot placeholder](docs/screenshots/landing.png) | ![Auth screen screenshot placeholder](docs/screenshots/auth.png) |

| Notes workspace | Mobile view |
| --- | --- |
| ![Notes workspace screenshot placeholder](docs/screenshots/notes.png) | ![Mobile screenshot placeholder](docs/screenshots/mobile.png) |

## Tech Stack

- Backend: Go `1.25.0`, standard `net/http` router, `pgxpool` for PostgreSQL connections.
- Database: PostgreSQL.
- Authentication: JWT access tokens, secure random refresh tokens, bcrypt password hashes.
- Frontend: React `19`, Vite `7`, `lucide-react` icons.
- Tests: Go tests for backend packages and Node's built-in test runner for frontend API tests.

## Project Structure

```text
.
+-- backend
|   +-- cmd/main.go                 # Backend entry point
|   +-- docs/api.md                 # Short API notes
|   +-- internal
|   |   +-- config                  # Environment configuration
|   |   +-- customerrors            # Domain error types
|   |   +-- entity                  # Response entities
|   |   +-- handlers                # HTTP handlers and JSON errors
|   |   +-- repos                   # PostgreSQL queries
|   |   +-- services                # Business logic and auth logic
|   |   +-- transport               # Route registration
|   +-- migrations                  # SQL migrations
+-- frontend
|   +-- src                         # React application and API client
|   +-- test                        # Frontend tests
|   +-- vite.config.js              # Dev server and API proxy
+-- .env.example
+-- go.mod
+-- README.md
```

## Backend Overview

The backend starts from `backend/cmd/main.go`. It loads `.env`, builds the application config, opens a PostgreSQL connection pool, creates note and auth services, registers HTTP handlers, and starts `http.ListenAndServe`.

The request flow is:

1. `transport.Setuprouter` registers `/api/...` routes with `net/http`.
2. `handlers` decode JSON requests, validate required fields, read auth cookies, and return JSON errors.
3. `services` hold application logic for notes and authentication.
4. `repos` execute SQL queries against PostgreSQL with `pgxpool`.

Protected notes endpoints read the `access-token` cookie, verify the JWT, extract `user_id`, and only operate on notes that belong to that user.

## Database

ProjectGo uses PostgreSQL. Migrations are stored in `backend/migrations` and should be applied in numeric order.

Current tables:

- `users`
  - `id SERIAL PRIMARY KEY`
  - `user_login TEXT`
  - `user_password TEXT`
  - `unique_user_login` constraint on `user_login`
- `notes`
  - `id SERIAL PRIMARY KEY`
  - `user_id INT NOT NULL`
  - `note TEXT NOT NULL`
  - foreign key `user_id -> users(id)` with `ON DELETE CASCADE`
- `refresh_tokens`
  - `id BIGSERIAL PRIMARY KEY`
  - `user_id INTEGER NOT NULL REFERENCES users(id)`
  - `token_hash BYTEA NOT NULL`
  - `expires_at TIMESTAMPTZ NOT NULL`
  - `created_at TIMESTAMPTZ NOT NULL`
  - unique constraint on `token_hash`

Passwords are stored as bcrypt hashes. Refresh tokens are generated as random URL-safe strings, hashed with SHA-256, and only the hash is stored in the database.

## Authentication

Registration and login create two cookies:

- `access-token`: JWT signed with `JWT_SECRET`, valid for 15 minutes.
- `refresh-token`: random token, valid for 30 days and persisted as a SHA-256 hash.

Both cookies are `HttpOnly`, `Secure`, and `SameSite=Lax` in backend responses. In local Vite development, the proxy removes the `Secure` flag from proxied `Set-Cookie` headers so the app can work on `http://localhost`.

When an access token expires, the frontend API client retries protected requests once after calling `/api/refresh/`. Refresh rotates the refresh token by deleting the old token hash and inserting a new one.

## API

All request and response bodies are JSON unless the endpoint returns an empty success body.

| Method | Path | Auth | Body | Success |
| --- | --- | --- | --- | --- |
| `GET` | `/api/` | Required | none | `200`, array of notes |
| `POST` | `/api/add/` | Required | `{ "text": "..." }` | `201` |
| `DELETE` | `/api/del/` | Required | `{ "id": 1 }` | `200` |
| `PUT` or `PATCH` | `/api/edit/` | Required | `{ "id": 1, "text": "..." }` | `200` |
| `POST` | `/api/register/` | No | `{ "login": "...", "password": "..." }` | `201`, auth cookies |
| `POST` | `/api/login/` | No | `{ "login": "...", "password": "..." }` | `200`, auth cookies |
| `POST` | `/api/logout/` | Refresh cookie | none | `200`, clears cookies |
| `POST` | `/api/refresh/` | Refresh cookie | none | `200`, rotated auth cookies |

Note response shape:

```json
{
  "id": 1,
  "user_id": 2,
  "text": "Buy milk"
}
```

Error response shape:

```json
{
  "message": "Error message"
}
```

## Environment

Copy the example file and fill in your local database values:

```bash
cp .env.example .env
```

Available variables:

```env
BASE_URL=localhost:8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=youruser
DB_PASSWORD=yourpassword
DB_NAME=yourdbname
JWT_SECRET=secret
```

`BASE_URL` is passed directly to `http.ListenAndServe`, so it should be a host and port, not a URL with `http://`.

## Local Development

Prerequisites:

- Go installed.
- Node.js and npm installed.
- PostgreSQL running locally.
- A database created for the values in `.env`.
- SQL migrations applied from `backend/migrations`.

Install frontend dependencies:

```bash
cd frontend
npm install
```

Run the backend from the repository root:

```bash
go run ./backend/cmd
```

Run the frontend from `frontend/`:

```bash
npm run dev
```

By default, Vite serves the frontend on `http://127.0.0.1:5173` and proxies `/api` requests to `http://localhost:8080`. You can override the backend origin for frontend development:

```bash
BACKEND_ORIGIN=http://localhost:8080 npm run dev
```

## Testing

Run backend tests from the repository root:

```bash
go test ./...
```

Some backend tests connect to the configured PostgreSQL database, so PostgreSQL must be available before running them.

Run frontend tests from `frontend/`:

```bash
npm test
```

Build the frontend:

```bash
npm run build
```

## Frontend Behavior

The React app has four main routes managed in the browser:

- `/`: public landing page.
- `/login`: login form.
- `/register`: registration form.
- `/notes`: authenticated notes workspace.

The API client in `frontend/src/api.js` always sends `credentials: "include"` because the backend stores auth in `HttpOnly` cookies. The frontend cannot read those cookies directly, which keeps token handling out of browser JavaScript.

## Production Notes

- Serve the app over HTTPS because backend auth cookies use the `Secure` flag.
- Use a strong, private `JWT_SECRET`.
- Keep `.env` out of version control.
- Apply database migrations before starting a new backend deployment.
- If the frontend and backend are deployed on different origins, add the required CORS and cookie settings intentionally on the backend.
