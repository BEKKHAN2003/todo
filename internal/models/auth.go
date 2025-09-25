package models

const (
	UserCtxKey = "user_id"
)

type AuthRequest struct {
	Username string `json:"username" binding:"required" example:"testuser"`
	Password string `json:"password" binding:"required" example:"secret123"`
}
type AuthResponse struct {
	Token string `json:"token"`
}
