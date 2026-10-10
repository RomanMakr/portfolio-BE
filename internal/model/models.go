package model

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
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
