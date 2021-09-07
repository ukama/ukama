package rest

import (
	"fmt"
	"net/http"
	"strings"
	"ukamaX/cloud/api-gateway/pkg"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	urest "github.com/ukama/ukamaX/common/rest"
)

type KratosAuthMiddleware struct {
	kratosConf  *pkg.Kratos
	isDebugMode bool
}

func NewKratosAuthMiddleware(kratosConf *pkg.Kratos, isDebugMode bool) *KratosAuthMiddleware {
	return &KratosAuthMiddleware{
		kratosConf:  kratosConf,
		isDebugMode: isDebugMode,
	}
}

func (r *KratosAuthMiddleware) IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := r.isTokenValid(c.Request)
		if err != nil {
			logrus.Warning("Error validating token. ", err.Error())
			if r.isDebugMode {
				urest.ThrowError(c, http.StatusUnauthorized, err.Error(), "", nil)
			} else {
				urest.ThrowError(c, http.StatusUnauthorized, "Unauthorized", "", nil)
			}
			c.Abort()
		}
	}
}

func (r *KratosAuthMiddleware) isTokenValid(request *http.Request) error {
	authHeader := request.Header.Get("authorization")
	if len(authHeader) == 0 {
		return fmt.Errorf("no header")
	}

	kratosUrl := r.kratosConf.Url

	token := authHeader[6:]
	token = strings.TrimSpace(token)

	client := resty.New()

	resp, err := client.R().
		EnableTrace().SetHeader("X-Session-Token", token).
		Get(kratosUrl + "/sessions/whoami")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("error getting session %v", resp.String())
	}

	return nil
}
