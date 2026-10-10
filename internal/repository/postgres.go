package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"portfolio-marketplace/backend/internal/model"
)

var ErrNotFound = errors.New("not found")

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
	if errors.Is(err, pgx.ErrNoRows) {
		return u, "", ErrNotFound
	}
	return u, hash, err
}
func (p *Postgres) Templates(ctx context.Context) ([]model.Template, error) {
	rows, err := p.DB.Query(ctx, `SELECT id,name,description,preview_url,is_premium FROM templates ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Template{}
	for rows.Next() {
		var t model.Template
		if err = rows.Scan(&t.ID, &t.Name, &t.Description, &t.PreviewURL, &t.IsPremium); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
func (p *Postgres) CreatePortfolio(ctx context.Context, userID string, in model.CreatePortfolioRequest) (model.Portfolio, error) {
	var v model.Portfolio
	err := p.DB.QueryRow(ctx, `INSERT INTO portfolios(user_id,template_id,slug,title,headline,bio,avatar_url,location,github_url,linkedin_url,published) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id::text,user_id::text,template_id,slug,title,headline,bio,avatar_url,location,github_url,linkedin_url,published,created_at,updated_at`, userID, in.TemplateID, strings.ToLower(in.Slug), in.Title, in.Headline, in.Bio, in.AvatarURL, in.Location, in.GithubURL, in.LinkedInURL, in.Published).Scan(&v.ID, &v.UserID, &v.TemplateID, &v.Slug, &v.Title, &v.Headline, &v.Bio, &v.AvatarURL, &v.Location, &v.GithubURL, &v.LinkedInURL, &v.Published, &v.CreatedAt, &v.UpdatedAt)
	return v, err
}
func (p *Postgres) ListMyPortfolios(ctx context.Context, userID string) ([]model.Portfolio, error) {
	rows, err := p.DB.Query(ctx, `SELECT id::text,user_id::text,template_id,slug,title,headline,bio,avatar_url,location,github_url,linkedin_url,published,created_at,updated_at FROM portfolios WHERE user_id=$1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Portfolio{}
	for rows.Next() {
		v, scanErr := scanPortfolio(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

type scanner interface{ Scan(dest ...any) error }

func scanPortfolio(s scanner) (v model.Portfolio, err error) {
	err = s.Scan(&v.ID, &v.UserID, &v.TemplateID, &v.Slug, &v.Title, &v.Headline, &v.Bio, &v.AvatarURL, &v.Location, &v.GithubURL, &v.LinkedInURL, &v.Published, &v.CreatedAt, &v.UpdatedAt)
	return
}
func (p *Postgres) GetOwnedPortfolio(ctx context.Context, userID, portfolioID string) (model.Portfolio, error) {
	v, err := scanPortfolio(p.DB.QueryRow(ctx, `SELECT id::text,user_id::text,template_id,slug,title,headline,bio,avatar_url,location,github_url,linkedin_url,published,created_at,updated_at FROM portfolios WHERE id=$1 AND user_id=$2`, portfolioID, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return v, ErrNotFound
	}
	return v, err
}
func (p *Postgres) UpdatePortfolio(ctx context.Context, userID, portfolioID string, in model.CreatePortfolioRequest) (model.Portfolio, error) {
	var v model.Portfolio
	err := p.DB.QueryRow(ctx, `UPDATE portfolios SET template_id=$3,slug=$4,title=$5,headline=$6,bio=$7,avatar_url=$8,location=$9,github_url=$10,linkedin_url=$11,published=$12,updated_at=now() WHERE id=$1 AND user_id=$2 RETURNING id::text,user_id::text,template_id,slug,title,headline,bio,avatar_url,location,github_url,linkedin_url,published,created_at,updated_at`, portfolioID, userID, in.TemplateID, strings.ToLower(in.Slug), in.Title, in.Headline, in.Bio, in.AvatarURL, in.Location, in.GithubURL, in.LinkedInURL, in.Published).Scan(&v.ID, &v.UserID, &v.TemplateID, &v.Slug, &v.Title, &v.Headline, &v.Bio, &v.AvatarURL, &v.Location, &v.GithubURL, &v.LinkedInURL, &v.Published, &v.CreatedAt, &v.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, ErrNotFound
	}
	return v, err
}
func (p *Postgres) PublicPortfolio(ctx context.Context, slug string) (model.Portfolio, error) {
	v, err := scanPortfolio(p.DB.QueryRow(ctx, `SELECT id::text,user_id::text,template_id,slug,title,headline,bio,avatar_url,location,github_url,linkedin_url,published,created_at,updated_at FROM portfolios WHERE slug=$1 AND published=TRUE`, strings.ToLower(slug)))
	if errors.Is(err, pgx.ErrNoRows) {
		return v, ErrNotFound
	}
	if err != nil {
		return v, err
	}
	rows, err := p.DB.Query(ctx, `SELECT id::text,portfolio_id::text,title,description,image_url,github_url,demo_url,technologies,sort_order FROM projects WHERE portfolio_id=$1 ORDER BY sort_order,title`, v.ID)
	if err != nil {
		return v, err
	}
	defer rows.Close()
	v.Projects = []model.Project{}
	for rows.Next() {
		var pr model.Project
		if err = rows.Scan(&pr.ID, &pr.PortfolioID, &pr.Title, &pr.Description, &pr.ImageURL, &pr.GithubURL, &pr.DemoURL, &pr.Technologies, &pr.SortOrder); err != nil {
			return v, err
		}
		v.Projects = append(v.Projects, pr)
	}
	return v, rows.Err()
}
func (p *Postgres) AddProject(ctx context.Context, userID, portfolioID string, in model.UpsertProjectRequest) (model.Project, error) {
	var v model.Project
	err := p.DB.QueryRow(ctx, `INSERT INTO projects(portfolio_id,title,description,image_url,github_url,demo_url,technologies,sort_order) SELECT id,$3,$4,$5,$6,$7,$8,$9 FROM portfolios WHERE id=$1 AND user_id=$2 RETURNING id::text,portfolio_id::text,title,description,image_url,github_url,demo_url,technologies,sort_order`, portfolioID, userID, in.Title, in.Description, in.ImageURL, in.GithubURL, in.DemoURL, in.Technologies, in.SortOrder).Scan(&v.ID, &v.PortfolioID, &v.Title, &v.Description, &v.ImageURL, &v.GithubURL, &v.DemoURL, &v.Technologies, &v.SortOrder)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, ErrNotFound
	}
	return v, err
}
