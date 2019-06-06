package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("path: ", c.FullPath(), "  Request Method:", c.GetHeader("Request Method"))
		t := time.Now()
		// before request
		c.Next()
		// after request
		latency := time.Since(t)
		// access the status we are sending
		status := c.Writer.Status()
		log.Printf("latency: %fms, status: %d\n\n", float64(latency)/1000000, status)
	}
}
