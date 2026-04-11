package handlers

// RegisterData is the data payload returned after a successful registration.
type RegisterData struct {
	UserID uint `json:"user_id" example:"1"`
}

// UserData is the data payload returned for user profile endpoints.
type UserData struct {
	ID    uint   `json:"id"    example:"1"`
	Name  string `json:"name"  example:"Jane Doe"`
	Email string `json:"email" example:"jane@example.com"`
}

// HealthData is the data payload returned by the health check endpoint.
type HealthData struct {
	Status string `json:"status" example:"ok"`
	Env    string `json:"env"    example:"development"`
}

// TokenRefreshRequest is the request body for POST /auth/refresh.
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"d4f...3e"`
}

// LogoutRequest is the request body for POST /auth/logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"d4f...3e"`
}
