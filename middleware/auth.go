package middleware

import (
	"net/http"
	"strings"

	"demodiqit_api/config"
	"demodiqit_api/helpers/crypt"
	"demodiqit_api/helpers/respond"

	"github.com/gin-gonic/gin"
)

const CurrentUserKey = "currentUser"

// JWTAuthMiddleware validates the Bearer token in the Authorization header.
// On success, it stores the parsed JWT claims in the Gin context under "currentUser".
func JWTAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, respond.ErrorRespond{
				Code:    "AUTH-005",
				Message: "Authorization header is missing or malformed",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := crypt.ParseJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, respond.ErrorRespond{
				Code:    "AUTH-006",
				Message: "Invalid or expired token",
			})
			return
		}

		// Store claims in context for use by subsequent handlers/middleware
		c.Set(CurrentUserKey, claims)
		c.Next()
	}
}
