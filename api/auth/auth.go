package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug("Entered API key middleware")
		err := IsKeyValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

func IsKeyValid(c *gin.Context) error {
	key := c.GetHeader("Authorization")
	log.Debug("API key: ", key)
	if key == "testkey" {
		return nil
	}
	return fmt.Errorf("Unauthorized")
}