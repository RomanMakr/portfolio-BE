package model

import "time"

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}
type Template struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PreviewURL  string `json:"preview_url"`
	IsPremium   bool   `json:"is_premium"`
}
type Project struct {
	ID           string   `json:"id"`
	PortfolioID  string   `json:"portfolio_id,omitempty"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ImageURL     string   `json:"image_url"`
	GithubURL    string   `json:"github_url"`
	DemoURL      string   `json:"demo_url"`
	Technologies []string `json:"technologies"`
	SortOrder    int      `json:"sort_order"`
}
type Portfolio struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id,omitempty"`
	TemplateID  string    `json:"template_id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Headline    string    `json:"headline"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	Location    string    `json:"location"`
	GithubURL   string    `json:"github_url"`
	LinkedInURL string    `json:"linkedin_url"`
	Published   bool      `json:"published"`
	Projects    []Project `json:"projects,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type CreatePortfolioRequest struct {
	TemplateID  string `json:"template_id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Headline    string `json:"headline"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	Location    string `json:"location"`
	GithubURL   string `json:"github_url"`
	LinkedInURL string `json:"linkedin_url"`
	Published   bool   `json:"published"`
}
type UpsertProjectRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ImageURL     string   `json:"image_url"`
	GithubURL    string   `json:"github_url"`
	DemoURL      string   `json:"demo_url"`
	Technologies []string `json:"technologies"`
	SortOrder    int      `json:"sort_order"`
}
