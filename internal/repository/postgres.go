package repository

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio-marketplace/backend/internal/model"
)

type Postgres struct{ DB *pgxpool.Pool }

func New(db *pgxpool.Pool) *Postgres { return &Postgres{DB: db} }

func (p *Postgres) CreateUser(ctx context.Context, email, hash, name string) (model.User, error) {
	var u model.User
	err := p.DB.QueryRow(ctx, `INSERT INTO users(email,password_hash,display_name) VALUES($1,$2,$3) RETURNING id::text,email,display_name`, strings.ToLower(email), hash, name).Scan(&u.ID, &u.Email, &u.DisplayName)
	return u, err
}
func (p *Postgres) UserByEmail(ctx context.Context, email string) (model.User, string, error) {
	var u model.User
	var hash string
	err := p.DB.QueryRow(ctx, `SELECT id::text,email,display_name,password_hash FROM users WHERE email=$1`, strings.ToLower(email)).Scan(&u.ID, &u.Email, &u.DisplayName, &hash)
	return u, hash, err
}
