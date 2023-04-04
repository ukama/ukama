package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/ukama/ukama/systems/common/rest"
)

var SESSION_KEY = "ukama_session"
var WHOAMI_PATH = "/sessions/whoami"

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

func GetUserBySession(cookieStr string, r *rest.RestClient) (*User, error) {
	urlObj, _ := url.Parse(r.C.BaseURL)
	cookie := &http.Cookie{
		Name:  SESSION_KEY,
		Value: cookieStr,
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
		Get(r.C.BaseURL + WHOAMI_PATH)

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
