package version

import "github.com/gin-gonic/gin"

func Middleware(ver string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-ACCOUNT-VERSION", ver)
		c.Next()
	}
}