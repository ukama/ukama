package server

import (
	"bytes"

	"github.com/Masterminds/semver/v3"
	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukamaX/common/rest"
	"github.com/ukama/ukamaX/hub/mocks"
	"github.com/ukama/ukamaX/hub/pkg"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func init() {
	pkg.IsDebugMode = true
}

var defaultCongif = &rest.HttpConfig{
	Cors: cors.Config{
		AllowAllOrigins: true,
	},
}

func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	s := mocks.Storage{}

	r := NewRouter(defaultCongif, &s, time.Second).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_RouterPut(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	f := getFileContent(t)
	defer f.Close()
	req, _ := http.NewRequest("PUT", "/capps/test-app/1.2.3", f)
	s := mocks.Storage{}
	ver := semver.MustParse("1.2.3")
	s.On("PutFile", mock.Anything, "test-app", ver,
		mock.MatchedBy(func(r io.Reader) bool {
			b, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("failed to read body: %s", err)
			}
			st, _ := f.Stat()
			assert.Equal(t, st.Size(), int64(len(b)))
			return true
		})).Return(nil)
	r := NewRouter(defaultCongif, &s, time.Second).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 201, w.Code)
}

func Test_RouterPutNotAtTargzFile(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	f := getFileContent(t)
	defer f.Close()
	//
	token := make([]byte, 1024*10)
	if _, err := rand.Read(token); err != nil {
		assert.FailNowf(t, "failed to generate token", err.Error())
	}
	req, _ := http.NewRequest("PUT", "/capps/test-app/1.2.3", bytes.NewReader(token))
	s := mocks.Storage{}
	r := NewRouter(defaultCongif, &s, time.Second).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Not a tar.gz file")
}

func Test_RouterGet(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	f := getFileContent(t)
	cont, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read testfile: %s", err)
	}
	defer f.Close()
	req, _ := http.NewRequest("GET", "/capps/test-app/1.2.3", bytes.NewReader(cont))
	s := mocks.Storage{}
	ver := semver.MustParse("1.2.3")

	s.On("GetFile", mock.Anything, "test-app", ver).Return(io.NopCloser(bytes.NewReader(cont)), nil)
	r := NewRouter(defaultCongif, &s, time.Second).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename=test-app-1.2.3.tar.gz", w.Header().Get("Content-Disposition"))
	assert.Equal(t, len(cont), w.Body.Len())
	assert.Equal(t, cont, w.Body.Bytes())
}

func getFileContent(t *testing.T) *os.File {
	f, err := os.Open("testdata/metrics.tar.gz")
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	return f
}
