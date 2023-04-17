package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

type TRole struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organizationId"`
}

type UserTraits struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	FirstVisit bool   `json:"firstVisit"`
}

func GetUserTraitsFromSession(orgId string, s *ory.Session) (*UserTraits, error) {
	data, err := json.Marshal(s.Identity.Traits)
	if err != nil {
		return nil, err
	}

	rdata, err := json.Marshal(s.Identity.MetadataPublic["roles"])
	if err != nil {
		return nil, err
	}

	var userTraits UserTraits
	if err := json.Unmarshal(data, &userTraits); err != nil {
		return nil, err
	}
	var roles []TRole
	if err := json.Unmarshal(rdata, &roles); err != nil {
		return nil, err
	}
	var role string = ""
	for _, r := range roles {
		if r.OrganizationId == orgId {
			role = r.Name
			break
		}
	}
	return &UserTraits{
		Id:         s.Identity.Id,
		Name:       userTraits.Name,
		Email:      userTraits.Email,
		Role:       role,
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

func GetMemberDetails(c *gin.Context) (string, string) {
	userId := c.Request.Header.Get("User-id")
	orgId := c.Request.Header.Get("Org-id")

	return userId, orgId
}

func GetMetaHeaderValues(s string) (string, string, string, error) {
	parts := strings.Split(s, ",")
	if len(parts) < 3 {
		return "", "", "", errors.New("meta header not provider")
	}

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]), nil
}
