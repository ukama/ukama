package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ory/client-go"
	ory "github.com/ory/client-go"
)

var SESSION_KEY = "ukama_session"
var WHOAMI_PATH = "/sessions/whoami"

type Session struct {
	Session         string `json:"session"`
	ExpiresAt       string `json:"expires_at"`
	AuthenticatedAt string `json:"authenticated_at"`
}

type UserTraits struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	FirstVisit bool   `json:"firstVisit"`
}

func GetUserTraitsFromSession(s *ory.Session) (*UserTraits, error) {
	data, err := json.Marshal(s.Identity.Traits)
	if err != nil {
		return nil, err
	}

	var userTraits UserTraits
	if err := json.Unmarshal(data, &userTraits); err != nil {
		return nil, err
	}
	return &UserTraits{
		Id:         s.Identity.Id,
		Name:       userTraits.Name,
		Email:      userTraits.Email,
		Role:       userTraits.Role,
		FirstVisit: userTraits.FirstVisit,
	}, nil
}

func GenerateJWT(s *string, e string, a string, k string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["session"] = s
	claims["expires_at"] = e
	claims["authenticated_at"] = a

	tokenString, err := token.SignedString([]byte(k))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(w http.ResponseWriter, t string, k string) (err error) {

	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return []byte(k), nil
	})

	if err != nil {
		return err
	}

	if token == nil {
		return errors.New("token is nil")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("token error")
	}

	exp := claims["expires_at"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return errors.New("token expired")
	}

	return nil
}

func GetSessionFromToken(w http.ResponseWriter, t string, k string) (*Session, error) {

	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return []byte(k), nil
	})

	if token == nil {
		fmt.Fprintf(w, "invalid token")
		return nil, errors.New("token error")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Fprintf(w, "couldn't parse token")
		return nil, errors.New("token error")
	}

	return &Session{
		Session:         claims["session"].(string),
		ExpiresAt:       claims["expires_at"].(string),
		AuthenticatedAt: claims["authenticated_at"].(string),
	}, nil
}

func GetSessionFromCookie(c *gin.Context, sessionKey string) (string, error) {
	cookies := map[string]string{}
	for _, cookie := range c.Request.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	if cookies[sessionKey] != "" {
		return cookies[sessionKey], nil
	}
	return "", fmt.Errorf("no session cookie found")
}

func ValidateSession(ss string, o *ory.APIClient) (*client.Session, error) {
	urlObj, _ := url.Parse(o.GetConfig().Servers[0].URL)
	cookie := &http.Cookie{
		Name:  SESSION_KEY,
		Value: ss,
	}
	o.GetConfig().HTTPClient.Jar.SetCookies(urlObj, []*http.Cookie{cookie})
	resp, r, err := o.FrontendApi.ToSession(context.Background()).Execute()
	if err != nil {
		return nil, err
	}
	if r.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("no valid session cookie found")
	}
	return resp, nil
}

func LoginUser(email string, password string, o *ory.APIClient) (*client.SuccessfulNativeLogin, error) {
	flow, _, err := o.FrontendApi.CreateNativeLoginFlow(context.Background()).Execute()
	if err != nil {
		return nil, err
	}
	b := client.UpdateLoginFlowWithPasswordMethod{
		Password:           password,
		Method:             "password",
		Identifier:         email,
		PasswordIdentifier: &email,
	}
	body := client.UpdateLoginFlowBody{
		UpdateLoginFlowWithPasswordMethod: &b,
	}
	flow1, _, err := o.FrontendApi.UpdateLoginFlow(context.Background()).Flow(flow.Id).UpdateLoginFlowBody(body).Execute()
	if err != nil {
		return nil, err
	}
	return flow1, nil
}

func CheckSession(sessionToken string, o *ory.APIClient) (session *client.Session, err error) {
	session, _, err = o.FrontendApi.ToSession(context.Background()).
		XSessionToken(sessionToken).
		Execute()
	if err != nil {
		return nil, err
	}
	return session, nil
}
