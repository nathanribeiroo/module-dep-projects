package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func addLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

func xItauCorrelationId() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.GetHeader("x-itau-correlation-id")

		if id == "" {
			id = uuid.New().String()
		}

		c.Writer.Header().Set("x-itau-correlation-id", id)

		c.Next()
	}
}
