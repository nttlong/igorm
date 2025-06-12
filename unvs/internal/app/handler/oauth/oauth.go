package oauth

// OAuthHandler là struct chứa dependency đến AccountService.
type OAuthHandler struct {
}

// ErrorResponse struct for consistent error messages
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
