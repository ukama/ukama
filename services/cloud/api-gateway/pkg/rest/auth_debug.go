package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DebugAuthMiddleware struct {
}

func (r *DebugAuthMiddleware) IsAuthenticated(c *gin.Context) {
	authHeader := c.Request.Header.Get("authorization")
	token := c.Request.Header.Get("token")
	if len(token) > 0 || len(authHeader) > 0 {
		logrus.Info("cookie header: ", authHeader)
		logrus.Info("token header: ", authHeader)
		logrus.Info("Bypassing authentication because we are in a debug mode")
		c.Set(USER_ID_KEY, "11111111-1111-1111-1111-111111111111")
		return
	}

	_, err := c.Request.Cookie("ukama_session")
	if err == nil {

		return
	}
	if err != nil {
		logrus.Error("Cookie not found. Error: ", err)
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func (r *DebugAuthMiddleware) IsAuthorized(c *gin.Context) {
	logrus.Infoln("Authorization passed")
}

func NewDebugAuthMiddleware() *DebugAuthMiddleware {
	return &DebugAuthMiddleware{}
}
