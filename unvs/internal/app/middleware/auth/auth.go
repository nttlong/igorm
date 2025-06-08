// internal/app/middleware/auth_middleware.go
package auth

import (
	"fmt" // Thêm import fmt
	"net/http"
	"strings"
	"time"

	jwtUtils "unvs/pkg/jwt_utils" // Import package jwt_utils

	"github.com/labstack/echo/v4"
)

// JWTAuthMiddleware xác thực JWT token từ header Authorization (Bearer).
// Nếu token hợp lệ, nó sẽ gắn UserClaims vào context của Echo.
func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"code":    "MISSING_AUTH_HEADER",
				"message": "Thiếu tiêu đề Authorization (Bearer Token).",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"code":    "INVALID_AUTH_FORMAT",
				"message": "Định dạng tiêu đề Authorization không hợp lệ. Phải là 'Bearer <token>'.",
			})
		}

		tokenString := parts[1]
		start := time.Now()                            // Thêm code để đo thời gian chạy hàm
		claims, err := jwtUtils.DecodeJWT(tokenString) // Sử dụng hàm DecodeJWT từ jwt_utils
		n := time.Since(start).Milliseconds()
		fmt.Println("Elapse time in ms ", n)
		if err != nil {
			// fmt.Printf("DEBUG: Lỗi khi decode JWT: %v\n", err) // Để debug nếu cần
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"code":    "INVALID_TOKEN",
				"message": fmt.Sprintf("Token không hợp lệ hoặc đã hết hạn: %v", err),
			})
		}

		// Lưu claims vào context của Echo để các handler tiếp theo có thể truy cập
		c.Set(string(jwtUtils.UserClaimsContextKey), claims)

		return next(c) // Chuyển request đến handler tiếp theo
	}
}
