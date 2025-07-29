package controllers

import (
	"net/http"
	"vapi/internal/utils"

	"github.com/labstack/echo/v4"
)

type OAuthTokenRequest struct {
	Username     string `form:"username"`
	Password     string `form:"password"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	GrantType    string `form:"grant_type"` // nên kiểm tra = "password"
}

type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// OAuthToken godoc
// @Summary OAuth2 Password Grant Token
// @Tags Auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param client_id formData string true "Client ID"
// @Param client_secret formData string true "Client Secret"
// @Param grant_type formData string true "Grant Type (must be 'password')"
// @Success 200 {object} OAuthTokenResponse
// @Failure 400 {object} map[string]string
// @Router /oauth/token [post]
func OAuthToken(c echo.Context) error {
	conatiner := utils.GetContainer(c)
	//"testuser", "testpassword"

	req := new(OAuthTokenRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if req.GrantType != "password" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unsupported grant_type"})
	}

	// 🔒 Thực hiện xác thực người dùng tại đây
	_, err := conatiner.AccountSvc.Get().Login("testuser", "testpassword")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	conatiner.GetContext()
	return c.JSON(http.StatusOK, OAuthTokenResponse{
		AccessToken:  conatiner.GetDb().GetDBName() + "-access-token",
		TokenType:    "bearer",
		ExpiresIn:    3600,
		RefreshToken: "fake-refresh-token",
	})

}
