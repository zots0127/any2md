package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func RateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	type client struct {
		count     int
		lastReset time.Time
	}
	
	clients := make(map[string]*client)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()
		
		if cl, exists := clients[clientIP]; exists {
			if now.Sub(cl.lastReset) > window {
				cl.count = 1
				cl.lastReset = now
			} else {
				cl.count++
				if cl.count > maxRequests {
					c.JSON(429, gin.H{
						"error": gin.H{
							"code":    "RATE_LIMIT_EXCEEDED",
							"message": "Too many requests",
						},
					})
					c.Abort()
					return
				}
			}
		} else {
			clients[clientIP] = &client{
				count:     1,
				lastReset: now,
			}
		}
		
		c.Next()
	}
}

func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	})
}