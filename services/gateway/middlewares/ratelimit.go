package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.RemoteIP()
		key := "ratelimit:" + clientIP

		val, err := rdb.Incr(c.Request.Context(), key).Result()
		if err != nil {

			log.Printf("redis rate limit error: %v", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})

			return

		}

		if val == 1 {
			err = rdb.Expire(c.Request.Context(), key, 60*time.Second).Err()

			if err != nil {
				log.Printf("redis expire error: %v", err)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})

				return
			}
		}

		if val > 60 {

			log.Printf("rate limit exceeded for ip: %s", clientIP)

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})

			return
		}

		c.Next()
	}
}
