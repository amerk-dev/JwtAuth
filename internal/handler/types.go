package handler

// TokenRequest represents a request with GUID
type TokenRequest struct {
	GuId string `json:"gu_id"`
}

// RefreshRequest represents a request to refresh token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// TokenResponse represents the pair of access and refresh tokens
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
