package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

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
