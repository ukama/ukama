package bootstrap

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type authenticator struct {
	bootstrapAuth AuthConfig
	token         string
	lock          sync.Mutex
}

type Authenticator interface {
	GetToken() (token string, err error)
}

func NewAuthenticator(bootstrapAuth AuthConfig) *authenticator {
	return &authenticator{bootstrapAuth: bootstrapAuth}
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// GetToken returns jwt token from the server. It validates the token every time to make sure it is not expired.
func (b *authenticator) GetToken() (token string, err error) {
	// we don't want to run  multiple simultaneous requests in case when no token is available
	b.lock.Lock()
	defer b.lock.Unlock()

	if !isTokenValid(b.token) {
		logrus.Infoln("Token is missing or invalid")
		t, err := b.getTokenFromServer()
		if err != nil {
			return "", err
		}
		b.token = t
	}

	return b.token, nil
}

func (b *authenticator) getTokenFromServer() (token string, err error) {
	logrus.Infoln("Retrieving token from server: ", b.bootstrapAuth.Auth0Host)
	client := resty.New()
	payload := map[string]string{
		"client_id":     b.bootstrapAuth.ClientId,
		"client_secret": b.bootstrapAuth.ClientSecret,
		"audience":      b.bootstrapAuth.Audience,
		"grant_type":    b.bootstrapAuth.GrantType}
	authResp := AuthResponse{}

	authUrl := b.bootstrapAuth.Auth0Host
	if !strings.HasSuffix(authUrl, "/oauth/token") {
		authUrl = authUrl + "/oauth/token"
	}
	resp, err := client.R().SetHeader("content-type", "application/json").
		SetBody(payload).SetResult(&authResp).Post(authUrl)
	if err != nil {
		return "", err
	}

	if resp.IsSuccess() {
		return authResp.AccessToken, nil
	}
	logrus.Errorf("Error from auth. Status code: %d. Body: %s", resp.StatusCode(), resp.String())
	return "", fmt.Errorf("error getting auth token")
}

func isTokenValid(token string) bool {
	if len(token) == 0 {
		logrus.Infoln("Token is missing")
		return false
	}
	parser := jwt.Parser{}
	isValid := false
	// ignoring error from this method as parsing errors are propagated to the key function
	_, err := parser.Parse(token, func(t *jwt.Token) (interface{}, error) {
		err := t.Claims.Valid()
		if err != nil {
			logrus.Warningf("Token is not valid. Error: %v", err)
			isValid = false
		} else {
			isValid = true
		}
		return []byte(""), nil
	})
	if err != nil {
		logrus.Infoln("Error parsing token. This error is expected when token is expired. Error: ", err.Error())
	}

	return isValid
}
