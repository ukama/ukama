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
	kratosConf           *pkg.Kratos
	isDebugMode          bool
	authorizationService AuthorizationService
}

type AuthorizationService interface {
	// checks if user with userId is authorized to access org
	IsAuthorized(userId string, org string) (bool, error)
}

func NewKratosAuthMiddleware(kratosConf *pkg.Kratos, authorizationSvc AuthorizationService, isDebugMode bool) *KratosAuthMiddleware {
	return &KratosAuthMiddleware{
		kratosConf:           kratosConf,
		isDebugMode:          isDebugMode,
		authorizationService: authorizationSvc,
	}
}

func (r *KratosAuthMiddleware) IsAuthenticated(c *gin.Context) {
	userId, err := r.isTokenValid(c.Request)
	if err != nil {
		r.processErr(c, err)
		return
	}
	c.Set(USER_ID_KEY, userId)
}

func (r *KratosAuthMiddleware) processErr(c *gin.Context, err error) {
	logrus.Warning("Error validating token. ", err.Error())
	if r.isDebugMode {
		urest.ThrowError(c, http.StatusUnauthorized, err.Error(), "", nil)
	} else {
		urest.ThrowError(c, http.StatusUnauthorized, "Unauthorized", "", nil)
	}
	c.Abort()
}

func (r *KratosAuthMiddleware) IsAuthorized(c *gin.Context) {
	userId := c.GetString(USER_ID_KEY)
	org := c.Param(ORG_URL_PARAMETER)
	res, err := r.authorizationService.IsAuthorized(userId, org)
	if err != nil || !res {
		if err != nil {
			logrus.Warning("error checking auhorization")
		}
		logrus.Infof("Access denied for user %s to organization %s", userId, org)
		urest.ThrowError(c, http.StatusNotFound, "Organization not found", "", nil)
		c.Abort()
	}
}

func (r *KratosAuthMiddleware) isTokenValid(request *http.Request) (userId string, err error) {
	kratosUrl := r.kratosConf.Url + "/sessions/whoami"
	client := resty.New()

	var resp *resty.Response
	cookie, err := request.Cookie(CookieName)
	if err != nil {
		logrus.Infoln("Cannot read cookie: ", err, " falling back to session token")
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
