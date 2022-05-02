package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ukama/ukama/services/cloud/api-gateway/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/services/cloud/api-gateway/pkg"
)

const KRATOS_USER_ID_HEADER = "X-Kratos-Authenticated-Identity-Id"
const userId = "expected-user-id"

func TestKratosAuthMiddleware_IsAuthenticated(t *testing.T) {

	kratosConf := &pkg.Kratos{
		Url: "https://kratos.test",
	}

	r := &KratosAuthMiddleware{
		kratosConf, true, &mocks.AuthorizationService{},
	}

	t.Run("validToken", func(t *testing.T) {
		// arrange
		recorder, testContext := newRecorderWithContex()
		testContext.Request.Header.Set("Authorization", "Bearer test")

		//
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Add(KRATOS_USER_ID_HEADER, userId)
			_, _ = res.Write([]byte("body"))
			res.WriteHeader(http.StatusOK)
		}))
		defer testServer.Close()
		kratosConf.Url = testServer.URL
		// Act
		r.IsAuthenticated(testContext)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, userId, testContext.GetString("UserId"))
	})

	t.Run("withCookie", func(t *testing.T) {
		// arrange
		recorder, testContext := newRecorderWithContex()
		testContext.Request.AddCookie(&http.Cookie{
			Name: CookieName, Value: "cookie with from kratos",
		})

		//
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if c, err := req.Cookie(CookieName); err == http.ErrNoCookie || len(c.Value) == 0 {
				assert.Fail(t, "cookie not found")
			}
			res.Header().Add(KRATOS_USER_ID_HEADER, userId)
			_, _ = res.Write([]byte("body"))
			res.WriteHeader(http.StatusOK)
		}))
		defer testServer.Close()
		kratosConf.Url = testServer.URL
		// Act
		r.IsAuthenticated(testContext)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, userId, testContext.GetString("UserId"))
	})

	t.Run("NotBearerToken", func(t *testing.T) {
		recorder, testContext := newRecorderWithContex()
		testContext.Request.Header.Set("Authorization", "test")

		r.IsAuthenticated(testContext)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		// arrange
		recorder, testContext := newRecorderWithContex()
		testContext.Request.Header.Set("Authorization", "Bearer test")

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Add(KRATOS_USER_ID_HEADER, userId)
			res.WriteHeader(http.StatusUnauthorized)
		}))
		defer testServer.Close()
		kratosConf.Url = testServer.URL
		// Act
		r.IsAuthenticated(testContext)
		// Assert
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("NoAuthorizationHeader", func(t *testing.T) {
		recorder, testContext := newRecorderWithContex()

		r.IsAuthenticated(testContext)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

}

func newRecorderWithContex() (*httptest.ResponseRecorder, *gin.Context) {
	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = &http.Request{}
	testContext.Request.Header = http.Header{}
	return recorder, testContext
}
