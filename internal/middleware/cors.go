// Package middleware HTTP 中间件。
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 返回跨域中间件，允许所有来源的请求。OPTIONS 预检请求直接返回 204。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
