package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type DebugAuthMiddleware struct {
}

func NewDebugAuthMiddleware() *DebugAuthMiddleware {
	return &DebugAuthMiddleware{}
}

func (r *DebugAuthMiddleware) IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("authorization")
		token := c.Request.Header.Get("token")

		if len(token) > 0 || len(authHeader) > 0 {
			logrus.Info("authorization header: ", authHeader)
			logrus.Info("token header: ", authHeader)
			logrus.Info("Bypassing authentication because we are in a debug mode")
			return
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
