package services

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Suponiendo que el token válido sea "Bearer mysecrettoken"
		if !strings.HasPrefix(token, "Bearer ") || token != "Bearer mysecrettoken" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Continuar con el siguiente handler si el token es válido
		c.Next()
	}
}
