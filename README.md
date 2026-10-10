# Portfolio Marketplace Backend (Go)

Go REST API for the Portfolio Marketplace. PostgreSQL stores users, templates, portfolios, and projects. The Flutter web and mobile clients connect to this service.

## Requirements

- Go 1.23+
- Docker Engine and Docker Compose

## Configure

```bash
cp .env.example .env
```

The current Go app uses environment variables from the process; it does not automatically load `.env`. For local development, the defaults in the application and Docker Compose are intended to work. Export any environment overrides in your shell. Never use development secrets in production.

## Run PostgreSQL

```bash
make db-up
```

The SQL migration is applied automatically when the PostgreSQL data volume is created for the first time. If you change the initial migration and want a fresh local database, run `make db-reset` (this deletes local database data).

## Run the API

```bash
make run
```

The API listens on `http://localhost:8080`; health check: `GET /health`.

## Useful commands

```bash
make help
make tidy
make test
make fmt
make db-logs
make db-down
```

## Security before deployment

Set a strong unique JWT secret, restrict CORS to your deployed clients, validate all inputs, use HTTPS, configure database credentials securely, and review authorization policies. This starter is an MVP and should not be exposed publicly without a security review.
# portfolio-BE
