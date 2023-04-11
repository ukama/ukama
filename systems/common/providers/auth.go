package providers

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const AuthenticateEndpoint = "/v1/auth"

type AuthRestClient struct {
	R   *rest.RestClient
	Jar *cookiejar.Jar
}

func NewAuthClient(u string, debug bool) (*AuthRestClient, error) {
	f, err := rest.NewRestClient(u, debug)
	if err != nil {
		logrus.Errorf("Can't conncet to %s url. Error %s", u, err.Error())

		return nil, err
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		logrus.Errorf("Got error while creating cookie jar %s", err.Error())
		return nil, err
	}
	urlObj, _ := url.Parse(u)
	jar.SetCookies(urlObj, []*http.Cookie{})

	N := &AuthRestClient{
		R:   f,
		Jar: jar,
	}

	return N, nil
}

func (a *AuthRestClient) AuthenticateUser(c *gin.Context, u string) (*resty.Response, error) {
	errStatus := &rest.ErrorMessage{}
	urlObj, _ := url.Parse(u)
	a.Jar.SetCookies(urlObj, c.Request.Cookies())
	resp, err := a.R.C.SetCookieJar(a.Jar).R().
		SetError(errStatus).
		Get(u + AuthenticateEndpoint)
	return resp, err
}
