package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/junaid9001/zvynt/gateway/config"

	"github.com/junaid9001/zvynt/pkg/shared"
)

func ValidiateJWT(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		var (
			ok     bool
			err    error
			claims shared.Claims
			token  *jwt.Token
		)

		tokenStr, ok = strings.CutPrefix(tokenStr, "Bearer ")

		if tokenStr == "" || !ok {

			tokenStr, err = c.Cookie("access_token")

			if err != nil || tokenStr == "" {

				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing or empty token"})
				return

			}

		}

		token, err = jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, err

			}
			return []byte(cfg.JWT_SECRET), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: invalid or expired token",
			})
			return

		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.UserRole)

		c.Next()
	}

}
