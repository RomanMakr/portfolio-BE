# Portfolio Backend (Go)

Go REST API with PostgreSQL-backed user registration and sign-in. Successful registration and sign-in return a JWT and the user profile.

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

The initial SQL migration creates the users table when the PostgreSQL data volume is created.

To apply a changed initial schema to your local database, run:

```bash
make db-reset
```

This deletes the local PostgreSQL volume and all its data, including registered accounts, then recreates the database. You will need to sign up again afterward. Do not use this command if you need to keep the current database data.

## Authentication API

- `POST /api/v1/auth/register` — JSON body: `{"email":"you@example.com","password":"at-least-8-chars","display_name":"Your Name"}`
- `POST /api/v1/auth/login` — JSON body: `{"email":"you@example.com","password":"your-password"}`
- `GET /health` — API health check

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
