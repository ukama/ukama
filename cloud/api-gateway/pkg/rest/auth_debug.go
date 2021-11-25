package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type DebugAuthMiddleware struct {
}

func (r *DebugAuthMiddleware) IsAuthenticated(c *gin.Context) {
	authHeader := c.Request.Header.Get("authorization")
	token := c.Request.Header.Get("token")

	if len(token) > 0 || len(authHeader) > 0 {
		logrus.Info("authorization header: ", authHeader)
		logrus.Info("token header: ", authHeader)
		logrus.Info("Bypassing authentication because we are in a debug mode")
		c.Set(USER_ID_KEY, "11111111-1111-1111-1111-111111111111")
		return
	}
	c.AbortWithStatus(http.StatusUnauthorized)
}

func (r *DebugAuthMiddleware) IsAuthorized(c *gin.Context) {
	logrus.Infoln("Authorization passed")
}

func NewDebugAuthMiddleware() *DebugAuthMiddleware {
	return &DebugAuthMiddleware{}
}
