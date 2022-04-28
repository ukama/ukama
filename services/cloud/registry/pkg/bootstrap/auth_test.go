package bootstrap

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/services/cloud/registry/pkg"
)

const VALID_TOKEN = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2MzM5NTk5MzEsImV4cCI6MjYxMjI2NzEzMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoidGVzdEBlY2FtcGxlLmNvbSJ9.zXmh32OAVtgieCmVJV92SmTYcQCFFBHRHKF2te6QdP8"

func Test_authenticator_isTokenValid(t *testing.T) {

	tests := []struct {
		name  string
		token string
		want  bool
	}{
		{"valid token",
			VALID_TOKEN,
			true},

		{"expired token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c",
			false},

		{
			"malformed",
			"yJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkxqY1NHZEpzVGdWVS1YNk9MM0tPYSJ9.eyJpc3MiOiJodHRwczovL3VrYW1hLnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJKajB4QmxxWWg1WUdSeTNBNVlyTkVLQjVJZkZwYnpJUUBjbGllbnRzIiwiYXVkIjoiYm9vdHN0cmFwLmRldi51a2FtYS5jb20iLCJpYXQiOjE2MzM5NDY4MTQsImV4cCI6MTYzNDAzMzIxNCwiYXpwIjoiSmoweEJscVloNVlHUnkzQTVZck5FS0I1SWZGcGJ6SVEiLCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMifQ.SmZSk_HI0Dg3kU4X_t0NZ1fL1qShVapQ7DaCKc7mEz1iY3iNy-3ENWEEQULUypEe-cALqcgSXc7ZEs_9_XLG3I90rrVCVcAxGvvRvhV_dHPuJjb1muPYhmKTA6iQSjOCGDB6CuZCJfOswcftTciErZmNPHkVIdMZV98uU2DTLpdCVz-eFcjHqhVwse0eVVZ1ZPVoEV_KEWVavlRn4eLOGNh7WW2MOJGLNSUjRoMm5f8P1z9Lj47mysZ2OYXp_RCrhXPeZjwEJ88CdJvsBlwu1j3yXnphcNetw_MJ9UYGSokfV1wIreBj8hA6qHZ5l_AED9Hv5xRXJoKhMpu1vyqhH",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isTokenValid(tt.token))
		})
	}
}

func Test_authenticator_GetToken(t *testing.T) {
	requestCount := 0
	j := newJwt()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		assert.Contains(t, r.URL.Path, "/oauth/token")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w,
			`{  "access_token": "`+j+`", "token_type": "Bearer" }`)
	}))
	auth := authenticator{
		bootstrapAuth: pkg.Auth{
			Auth0Host: ts.URL,
		},
	}

	// run several requests in parallel
	for i := 0; i < 10; i++ {
		go func() {
			// get token from cache
			token, err := auth.GetToken()
			assert.NoError(t, err)
			assert.Equal(t, j, token)
			assert.Equal(t, 1, requestCount)
		}()
	}

	// wait for token to expire
	time.Sleep(2 * time.Second)
	token, err := auth.GetToken()
	assert.NoError(t, err)
	assert.Equal(t, j, token)
	assert.Equal(t, 2, requestCount)
}

func newJwt() string {
	mySigningKey := []byte("AllYourBase")

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(2 * time.Second),
		},
		Issuer: "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	return ss
}
