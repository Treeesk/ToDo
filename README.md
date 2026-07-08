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
|   +-- Dockerfile                  # Multi-stage Go backend build
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
|   +-- Dockerfile                  # Multi-stage React/Vite build served by Nginx
|   +-- nginx.conf                  # Nginx static-file and API proxy config
|   +-- src                         # React application and API client
|   +-- test                        # Frontend tests
|   +-- vite.config.js              # Dev server and API proxy
+-- .env.example
+-- docker-compose.yml              # Production-style container orchestration
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

For Docker Compose, the same variables point to services inside the Docker network. Typical container values are:

```env
BASE_URL=:8080
DB_HOST=database
DB_PORT=5432
DB_USER=youruser
DB_PASSWORD=yourpassword
DB_NAME=yourdbname
JWT_SECRET=secret
```

`BASE_URL=:8080` makes the Go server listen on port `8080` inside the backend container. `DB_HOST=database` uses the Compose service name for PostgreSQL.

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

## Docker

The repository includes Dockerfiles for both application layers and a Compose file for running the full stack.

- `backend/Dockerfile` builds the Go backend in a `golang:1.25` stage and copies the compiled `server` binary into a smaller Debian runtime image.
- `frontend/Dockerfile` builds the Vite React app with Node and copies the generated `dist` files into an Nginx image.
- `frontend/nginx.conf` serves the SPA files and proxies `/api/...` requests to the backend service.
- `docker-compose.yml` starts Nginx, backend, PostgreSQL, and a one-shot migration container.

The current Compose file uses prebuilt images:

```yaml
backend:
  image: treeesk/projectgo-backend:1.1

nginx:
  image: treeesk/projectgo-frontend:1.2
```

If those images are already pushed to Docker Hub, start the stack with:

```bash
docker compose pull
docker compose up -d
```

For local image builds instead, replace the `image` entries with `build` entries or build and push new tags manually.

Build and push the backend image from the repository root:

```bash
docker build -t treeesk/projectgo-backend:1.1 -f backend/Dockerfile .
docker push treeesk/projectgo-backend:1.1
```

Build and push the frontend/Nginx image:

```bash
docker build -t treeesk/projectgo-frontend:1.2 -f frontend/Dockerfile ./frontend
docker push treeesk/projectgo-frontend:1.2
```

After changing an image tag in `docker-compose.yml`, pull and recreate the affected containers:

```bash
docker compose pull
docker compose up -d
```

### Compose Services

`nginx` is the only public service:

```yaml
ports:
  - "80:80"
```

Open the application locally at:

```text
http://localhost
```

The backend is not published to the host. It is available only inside the Compose network:

```text
http://backend:8080
```

Nginx reaches it through `proxy_pass`:

```nginx
location /api/ {
    proxy_pass http://backend:8080;
}
```

PostgreSQL is exposed only inside the Docker network as:

```text
database:5432
```

### Database Migrations

Migrations live in `backend/migrations`. Compose runs them through a one-shot `migrate` service after PostgreSQL becomes healthy:

```yaml
migrate:
  image: migrate/migrate
  volumes:
    - ./backend/migrations:/migrations:ro
  command:
    [
      "-path", "/migrations",
      "-database", "postgres://${DB_USER}:${DB_PASSWORD}@database:${DB_PORT}/${DB_NAME}?sslmode=disable",
      "up"
    ]
```

The bind mount gives the migration container read-only access to the SQL files. The container sees them at `/migrations`.

### Database Volume

PostgreSQL data is stored in a named Docker volume:

```yaml
volumes:
  - database-data:/var/lib/postgresql
```

The top-level declaration creates the named volume:

```yaml
volumes:
  database-data:
```

This keeps database files outside the container filesystem, so recreating the database container does not erase data.

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
