package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)
		c.Next()
		log.Printf("%s %s %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() >= http.StatusBadRequest {
			return
		}
	}
}
