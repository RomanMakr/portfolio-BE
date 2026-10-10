package config

import "os"

type Config struct{ Port, DatabaseURL, JWTSecret, FrontendOrigin string }

func Load() Config {
	return Config{
		Port:           env("PORT", "8080"),
		DatabaseURL:    env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/portfolio_marketplace?sslmode=disable"),
		JWTSecret:      env("JWT_SECRET", "dev-only-change-me"),
		FrontendOrigin: env("FRONTEND_ORIGIN", "http://localhost:3000"),
	}
}
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
