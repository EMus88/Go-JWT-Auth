package model

type SignUpRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
