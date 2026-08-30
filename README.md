# Loyalty System (Gophermart)

Gophermart is a loyalty points accumulation system — a service that accepts users' order numbers,
communicates with an external accrual calculation system, and credits bonus points that can later
be spent on future purchases.

## What it does

- User registration and authentication
- Accepting order numbers for loyalty points accrual
- Polling the external accrual system for order status and points amount
- Crediting and withdrawing points for user orders
- Viewing order history and current points balance

## Stack

- **Go** — main language
- **PostgreSQL** — data storage
- **Docker / docker-compose** — local environment
- **Makefile** — development commands

## Quick start

### With Docker

```bash
make up
```

Spins up PostgreSQL and the application (`http://localhost:8080`).

### Locally

1. Start the database:
   ```bash
   make compose-up-db
   ```
2. Install dependencies and run the server:
   ```bash
   make run
   ```
3. Connect to the local database (if needed):
   ```bash
   make psql-local-db
   ```

### Environment variables

| Variable        | Description                        | Example                                                         |
|-----------------|-------------------------------------|-------------------------------------------------------------------|
| `ADDRESS`       | Server host and port                | `:8080`                                                         |
| `DATABASE_URI`  | PostgreSQL connection string        | `postgres://user:pass@localhost:5432/postgres?sslmode=disable`|

## Architecture

The project follows Clean Architecture principles:

```
cmd/
  gophermart/   — main service entry point
  accrual/      — accrual calculation service entry point
internal/
  app/          — application bootstrap, HTTP server configuration
  entity/       — domain entities (user, order, balance, etc.)
  usecase/      — business logic (use cases)
config/         — application configuration
migrations/     — SQL database migrations
pkg/            — general-purpose reusable packages
```

The layers interact as `cmd` → `app` → `usecase` → `entity`, keeping business logic isolated from
transport (HTTP) and storage (PostgreSQL) details.

## Additional info

Detailed technical specification is available in [SPECIFICATION.md](SPECIFICATION.md).

### Updating the template

To pull in autotest updates from the course template:

```bash
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
git fetch template && git checkout template/master .github
```
