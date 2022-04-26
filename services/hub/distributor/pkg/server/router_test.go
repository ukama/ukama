package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukamaX/hub/distributor/pkg"
)

func init() {
	pkg.IsDebugMode = true
}

func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	defconf := pkg.NewConfig()

	r := NewRouter(defconf).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_RouterPut(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	f := []byte(`{ "store":"./test/data/art" }`)

	req, _ := http.NewRequest("PUT", "/chunk/ukamaos/1.0.1", bytes.NewBuffer(f))

	defconf := pkg.NewConfig()
	defconf.Distribution.Chunk.Stores[0] = "./test/data/store"

	r := NewRouter(defconf).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

	putIndexFileContent(t, w.Body)

}

func Test_RouterPutNoStore(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", "/chunk/ukamaos/1.0.1", nil)

	defconf := pkg.NewConfig()

	r := NewRouter(defconf).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation for 'Store'")
}

func putIndexFileContent(tt *testing.T, br io.Reader) {
	file := "./test/data/index/index.caidx"
	f, err := os.Create(file)
	if err != nil {
		assert.FailNow(tt, err.Error())
	}
	defer f.Close()

	bytes, err := io.Copy(f, br)
	if err != nil {
		assert.FailNow(tt, err.Error())
	}

	if bytes <= 0 {
		assert.FailNow(tt, "expected file contents but looks like its empty")
	}
	logrus.Debugf("Index file %s created with %d bytes.", file, bytes)
}
