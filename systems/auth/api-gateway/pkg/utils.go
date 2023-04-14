package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
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
	expStr := claims["expires_at"]
	exp, _ := time.Parse(time.RFC1123, expStr.(string))
	tUnix := exp.Local().Unix()
	if tUnix < time.Now().Local().Unix() {
		return errors.New("token expired")
	} else {
		return nil
	}
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

func SessionType(c *gin.Context, sessionKey string) (string, error) {
	cookies := map[string]string{}
	for _, cookie := range c.Request.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	if cookies[sessionKey] != "" {
		return "cookie", nil
	} else if c.Request.Header.Get("X-Session-Token") != "" {
		return "header", nil
	}
	return "", fmt.Errorf("no cookie/token found")

}

func GetCookieStr(c *gin.Context, sessionKey string) string {
	cookies := map[string]string{}
	for _, cookie := range c.Request.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	if cookies[sessionKey] != "" {
		return cookies[sessionKey]
	}
	return ""
}

func GetTokenStr(c *gin.Context) string {
	token := c.Request.Header.Get("X-Session-Token")
	return token
}
