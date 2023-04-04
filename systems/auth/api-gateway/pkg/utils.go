package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ukama/ukama/systems/common/rest"
)

type User struct {
	Id       string       `json:"id"`
	Identity UserIdentity `json:"identity"`
}

type UserIdentity struct {
	Id     string     `json:"id"`
	Traits UserTraits `json:"traits"`
}

type UserTraits struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func parseCookies(c *gin.Context) map[string]string {
	cookies := map[string]string{}
	for _, cookie := range c.Request.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	return cookies
}

func getSessionCookie(cookies map[string]string) (string, error) {
	if cookies["ukama_session"] != "" {
		return cookies["ukama_session"], nil
	}
	return "", fmt.Errorf("no session cookie found")
}

func GetUserBySession(c *gin.Context, r *rest.RestClient) (*User, error) {
	cookies := parseCookies(c)
	session, err := getSessionCookie(cookies)
	if err != nil {
		return nil, err
	}
	urlObj, _ := url.Parse(r.C.BaseURL)
	cookie := &http.Cookie{
		Name:  "ukama_session",
		Value: session,
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
		return nil, err
	}
	jar.SetCookies(urlObj, []*http.Cookie{cookie})

	errStatus := &rest.ErrorMessage{}
	resp, err := r.C.SetCookieJar(jar).R().
		SetError(errStatus).
		Get(r.C.BaseURL + "/sessions/whoami")

	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(strings.NewReader(string(resp.Body())))
	var data User
	err = decoder.Decode(&data)

	if err != nil {
		return nil, err
	}
	return &data, nil
}
