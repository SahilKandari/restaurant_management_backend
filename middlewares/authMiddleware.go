package middlewares

import (
	"net/http"
	"strings"

	"restaurant-management/helpers"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			return
		}

		// 2. Must start with “Bearer ”
		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Invalid Authorization header format (want: Bearer <token>)"})
			return
		}

		// 3. Extract the token part and trim any stray spaces
		tokenString := strings.TrimSpace(authHeader[len(bearerPrefix):])
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Token not provided"})
			return
		}

		err := helpers.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}
