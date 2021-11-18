package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/ukama/ukamaX/cloud/api-gateway/pkg"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
	urest "github.com/ukama/ukamaX/common/rest"
)

const USER_ID_KEY = "UserId"
const CookieName = "ukama_session"

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
		userId, err := r.isTokenValid(c.Request)
		if err != nil {
			logrus.Warning("Error validating token. ", err.Error())
			if r.isDebugMode {
				urest.ThrowError(c, http.StatusUnauthorized, err.Error(), "", nil)
			} else {
				urest.ThrowError(c, http.StatusUnauthorized, "Unauthorized", "", nil)
			}
			c.Abort()
			return
		}
		c.Set("UserId", userId)
	}
}

func (r *KratosAuthMiddleware) isTokenValid(request *http.Request) (userId string, err error) {
	kratosUrl := r.kratosConf.Url + "/sessions/whoami"
	client := resty.New()

	var resp *resty.Response
	cookie, err := request.Cookie(CookieName)
	if err != nil {
		logrus.Warning("Can't read cookie: ", err)
	}
	if err == nil {
		resp, err = client.R().SetCookie(cookie).
			Get(kratosUrl)
	} else {
		authHeader := request.Header.Get("authorization")

		if len(authHeader) == 0 {
			return "", fmt.Errorf("no header")
		}

		if len(authHeader) < 6 {
			return "", fmt.Errorf("invalid authorization format")
		}
		token := authHeader[6:]
		token = strings.TrimSpace(token)

		resp, err = client.R().SetHeader("X-Session-Token", token).
			Get(kratosUrl)
	}

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("error getting session %v", resp.String())
	}

	userId = resp.Header().Get("X-Kratos-Authenticated-Identity-Id")
	return userId, nil
}
