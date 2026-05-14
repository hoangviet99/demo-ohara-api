package middleware

import (
	"net/http"

	"demodiqit_api/helpers/crypt"
	"demodiqit_api/helpers/respond"

	"github.com/gin-gonic/gin"
)

// AdminOnlyMiddleware restricts access to users with the "admin" role.
// Must be used AFTER JWTAuthMiddleware, which sets the "currentUser" context key.
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(CurrentUserKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, respond.ErrorRespond{
				Code:    "AUTH-005",
				Message: "Unauthorized",
			})
			return
		}

		claims, ok := val.(*crypt.JWTClaim)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, respond.ErrorRespond{
				Code:    "AUTH-006",
				Message: "Invalid token claims",
			})
			return
		}

		for _, role := range claims.Roles {
			if role == "admin" {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, respond.ErrorRespond{
			Code:    "AUTH-007",
			Message: "Forbidden: admin access required",
		})
	}
}
