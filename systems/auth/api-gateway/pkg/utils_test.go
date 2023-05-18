package pkg

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	ory "github.com/ory/client-go"
	"github.com/stretchr/testify/assert"
)

func TestGetUserTraitsFromSession(t *testing.T) {
	identity := ory.Identity{
		Id: "1234",
		Traits: map[string]interface{}{
			"name":       "John Doe",
			"email":      "johndoe@example.com",
			"role":       "user",
			"firstVisit": true,
		},
	}

	s := &ory.Session{
		Identity: identity,
	}

	expectedUserTraits := &UserTraits{
		Id:         "1234",
		Name:       "John Doe",
		Email:      "johndoe@example.com",
		Role:       "member",
		FirstVisit: true,
	}

	userTraits, err := GetUserTraitsFromSession("abc-123", s)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if userTraits.Id != expectedUserTraits.Id {
		t.Errorf("Unexpected Id: got %s, want %s", userTraits.Id, expectedUserTraits.Id)
	}

	if userTraits.Name != expectedUserTraits.Name {
		t.Errorf("Unexpected Name: got %s, want %s", userTraits.Name, expectedUserTraits.Name)
	}

	if userTraits.Email != expectedUserTraits.Email {
		t.Errorf("Unexpected Email: got %s, want %s", userTraits.Email, expectedUserTraits.Email)
	}

	if userTraits.Role != expectedUserTraits.Role {
		t.Errorf("Unexpected Role: got %s, want %s", userTraits.Role, expectedUserTraits.Role)
	}

	if userTraits.FirstVisit != expectedUserTraits.FirstVisit {
		t.Errorf("Unexpected FirstVisit: got %v, want %v", userTraits.FirstVisit, expectedUserTraits.FirstVisit)
	}
}

func TestGenerateJWT(t *testing.T) {
	session := "user123"
	expiresAt := "2023-04-17T12:00:00Z"
	authenticatedAt := "2023-04-16T12:00:00Z"
	key := "secret"

	tokenString, err := GenerateJWT(&session, expiresAt, authenticatedAt, key)
	if err != nil {
		t.Fatalf("Unexpected error generating JWT: %v", err)
	}

	// Parse the token to verify its contents
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key used to sign the token
		return []byte(key), nil
	})

	if err != nil {
		t.Fatalf("Unexpected error parsing JWT: %v", err)
	}

	// Verify the claims in the token
	claims := token.Claims.(jwt.MapClaims)

	if session != claims["session"] {
		t.Errorf("Unexpected session value: got %s, want %s", claims["session"], session)
	}

	if expiresAt != claims["expires_at"] {
		t.Errorf("Unexpected expires_at value: got %s, want %s", claims["expires_at"], expiresAt)
	}

	if authenticatedAt != claims["authenticated_at"] {
		t.Errorf("Unexpected authenticated_at value: got %s, want %s", claims["authenticated_at"], authenticatedAt)
	}
}

func TestValidateToken(t *testing.T) {
	// Set up a mock HTTP response writer
	w := httptest.NewRecorder()

	// Generate a JWT token
	key := "secret"
	token := "user123"
	a := time.Now()
	tokenString, err := GenerateJWT(&token, a.Add(time.Second*2).Format(time.RFC1123), a.Format(time.RFC1123), key)
	if err != nil {
		t.Fatalf("Unexpected error generating JWT: %v", err)
	}

	// Call the ValidateToken function with the generated token and key
	err = ValidateToken(w, tokenString, key)
	assert.Nil(t, err)

	// Wait for the token to expire, then call the function again
	time.Sleep(3 * time.Second)
	err = ValidateToken(w, tokenString, key)
	assert.Error(t, err)
}

func TestGetSessionFromToken(t *testing.T) {
	// Set up a mock HTTP response writer
	w := httptest.NewRecorder()

	// Generate a JWT token
	key := "secret"
	token := "user123"
	a := time.Now()
	e := a.Add(time.Second * 2).Format(time.RFC1123)
	tokenString, err := GenerateJWT(&token, e, a.Format(time.RFC1123), key)
	if err != nil {
		t.Fatalf("Unexpected error generating JWT: %v", err)
	}

	// Call the GetSessionFromToken function with the generated token and key
	session, err := GetSessionFromToken(w, tokenString, key)

	// Validate that no error was returned
	if err != nil {
		t.Fatalf("Unexpected error getting session from token: %v", err)
	}

	// Validate that the session object was correctly parsed from the token
	if session.Session != token || session.ExpiresAt != e || session.AuthenticatedAt != a.Format(time.RFC1123) {
		t.Fatalf("Unexpected session object: %+v", session)
	}

	// Call the function with an invalid token
	invalidToken := "invalid-token"
	_, err = GetSessionFromToken(w, invalidToken, key)

	// Validate that the "invalid token" error was returned
	if err == nil || err.Error() != "token error" {
		t.Fatalf("Unexpected error getting session from token: %v", err)
	}
}

func TestSessionType(t *testing.T) {
	// Create a new gin context with a cookie and header for testing
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_key",
		Value: "test_cookie",
	})
	req.Header.Set("X-Session-Token", "test_token")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	gin.SetMode(gin.TestMode)
	c.Request = req

	// Test with valid cookie and header
	sessionType, err := SessionType(c, "session_key")
	assert.Nil(t, err)
	assert.Equal(t, "cookie", sessionType)
	req.Header.Del("Cookie")

	// Test with valid header
	req.Header.Set("X-Session-Token", "test_token")
	sessionType, err = SessionType(c, "session_key")
	assert.Nil(t, err)
	assert.Equal(t, "header", sessionType)
}

func TestGetCookieStr(t *testing.T) {
	// Create a new gin context
	gin.SetMode(gin.TestMode)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{Name: "session_key", Value: "abc123"}
	req.AddCookie(cookie)
	// w := httptest.NewRecorder()
	c := gin.Context{Request: req}

	// Test with existing session key
	cookieStr := GetCookieStr(&c, "session_key")
	expectedCookieStr := "abc123"
	if cookieStr != expectedCookieStr {
		t.Errorf("Expected cookie string to be %s, but got %s", expectedCookieStr, cookieStr)
	}

	// Test with non-existing session key
	cookieStr = GetCookieStr(&c, "non_existing_session_key")
	expectedCookieStr = ""
	if cookieStr != expectedCookieStr {
		t.Errorf("Expected cookie string to be %s, but got %s", expectedCookieStr, cookieStr)
	}
}
