package middleware

import (
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/PPIO/pi-cloud-monitor-backend/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWTMiddleware(publickKey *rsa.PublicKey, skippedPrefixPath map[string]bool) echo.MiddlewareFunc {
	var jwtMiddleware = middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			for prefix := range skippedPrefixPath {
				if strings.HasPrefix(path, prefix) {
					return true
				}
			}
			return false
		},
		SigningKey:    publickKey,
		SigningMethod: "RS512",
		Claims:        &common.JWTCustomClaims{},
		ErrorHandler: func(err error) error {
			return &echo.HTTPError{
				Code: http.StatusUnauthorized,
			}
		},
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*common.JWTCustomClaims)
			c.Set("uid", claims.UID)
		},
	})
	return jwtMiddleware
}
