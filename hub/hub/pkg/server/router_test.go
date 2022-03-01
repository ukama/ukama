package server

import (
	"bytes"
	"github.com/minio/minio-go/v7"
	"reflect"

	"github.com/Masterminds/semver/v3"
	"github.com/gin-contrib/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukamaX/common/rest"
	"github.com/ukama/ukamaX/hub/hub/mocks"
	"github.com/ukama/ukamaX/hub/hub/pkg"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var emptyChunker = &mocks.Chunker{}

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
	r := NewRouter(defaultCongif, &s, emptyChunker, time.Second).fizz.Engine()

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
	ch := mocks.Chunker{}
	ch.On("Chunk", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ver := semver.MustParse("1.2.3")
	s.On("PutFile", mock.Anything, "test-app", ver, pkg.TarGzExtension,
		mock.MatchedBy(func(r io.Reader) bool {
			b, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("failed to read body: %s", err)
			}
			st, _ := f.Stat()
			assert.Equal(t, st.Size(), int64(len(b)))
			return true
		})).Return("", nil)
	r := NewRouter(defaultCongif, &s, &ch, time.Second).fizz.Engine()

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

	r := NewRouter(defaultCongif, &s, emptyChunker, time.Second).fizz.Engine()

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
	req, _ := http.NewRequest("GET", "/capps/test-app/1.2.3.tar.gz", bytes.NewReader(cont))
	s := mocks.Storage{}
	ver := semver.MustParse("1.2.3")

	s.On("GetFile", mock.Anything, "test-app", ver, pkg.TarGzExtension).Return(io.NopCloser(bytes.NewReader(cont)), nil)
	r := NewRouter(defaultCongif, &s, emptyChunker, time.Second).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename=test-app-1.2.3.tar.gz", w.Header().Get("Content-Disposition"))
	assert.Equal(t, len(cont), w.Body.Len())
	if !bytes.Equal(cont, w.Body.Bytes()) {
		assert.Fail(t, "actual content is not equal to expected")
	}
}

type FakeReader struct {
}

func (f FakeReader) Read(p []byte) (n int, err error) {
	return 0, minio.ErrorResponse{
		Code: "NoSuchKey",
	}
}

func (f FakeReader) Close() error {
	return nil
}

func Test_RouterGetReturnError(t *testing.T) {

	noContentDistrValidator := func(w *httptest.ResponseRecorder) {
		assert.Equal(t, "", w.Header().Get("Content-Disposition"))
	}
	tests := []struct {
		name            string
		request         string
		storageMockFunc func() pkg.Storage
		expectedCode    int
		validateRequst  func(*httptest.ResponseRecorder)
	}{
		{
			name:    "NotFoundInBucket",
			request: "/capps/test-app/1.2.3.tar.gz",
			storageMockFunc: func() pkg.Storage {
				s := mocks.Storage{}
				s.On("GetFile", mock.Anything, "test-app", semver.MustParse("1.2.3"), pkg.TarGzExtension).Return(&FakeReader{}, nil)
				return &s
			},
			expectedCode:   404,
			validateRequst: noContentDistrValidator,
		},
		{
			name:    "BadExtension",
			request: "/capps/test-app/1.2.3.bad-extension",
			storageMockFunc: func() pkg.Storage {
				s := mocks.Storage{}
				return &s
			},
			expectedCode:   404,
			validateRequst: noContentDistrValidator,
		},
		{
			name:    "NoExtension",
			request: "/capps/test-app/1.2.3",
			storageMockFunc: func() pkg.Storage {
				s := mocks.Storage{}
				return &s
			},
			expectedCode:   404,
			validateRequst: noContentDistrValidator,
		},
		{
			name:    "BadVersion",
			request: "/capps/test-app/1.this-is-bad.3",
			storageMockFunc: func() pkg.Storage {
				s := mocks.Storage{}
				return &s
			},
			expectedCode:   404,
			validateRequst: noContentDistrValidator,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", tt.request, nil)

			r := NewRouter(defaultCongif, tt.storageMockFunc(), emptyChunker, time.Second).fizz.Engine()

			// act
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.validateRequst != nil {
				tt.validateRequst(w)
			}

		})
	}

}

func TestListApps(t *testing.T) {
	tests := []struct {
		name             string
		artifacts        *[]pkg.AritfactInfo
		wantCode         int
		wantBodyContains []string
	}{
		{
			name: "ReturnsList",
			artifacts: &[]pkg.AritfactInfo{
				{
					Url:       "capps/test-app",
					Version:   "1.2.3",
					CreatedAt: time.Now().Add(-5 * time.Hour),
				},
				{
					Url:       "capps/test-app",
					Version:   "1.2.4",
					CreatedAt: time.Now().Add(-4 * time.Hour),
				},
			},
			wantBodyContains: []string{"capps/test-app", "1.2.4", "1.2.3"},
			wantCode:         200,
		},

		{
			name:             "ReturnsList",
			artifacts:        &[]pkg.AritfactInfo{},
			wantBodyContains: []string{},
			wantCode:         404,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// arrange
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/capps/test-app", nil)
			s := mocks.Storage{}

			s.On("ListVersions", mock.Anything, "test-app").Return(test.artifacts, nil)
			r := NewRouter(defaultCongif, &s, emptyChunker, time.Second).fizz.Engine()

			// act
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, test.wantCode, w.Code)
			for _, c := range test.wantBodyContains {
				assert.Contains(t, w.Body.String(), c)
			}

		})
	}

}

func getFileContent(t *testing.T) *os.File {
	f, err := os.Open("testdata/metrics.tar.gz")
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	return f
}

func Test_parseArtifactName(t *testing.T) {
	tests := []struct {
		name         string
		artifactName string
		wantVer      *semver.Version
		wantExt      string
		wantErr      bool
	}{
		{
			name:         "valid_tar.gz",
			artifactName: "1.2.3.tar.gz",
			wantVer:      semver.MustParse("1.2.3"),
			wantExt:      ".tar.gz",
		},
		{
			name:         "valid_fancy_version",
			artifactName: "1.2.3-debug.tar.gz",
			wantVer:      semver.MustParse("1.2.3-debug"),
			wantExt:      ".tar.gz",
		},
		{
			name:         "valid_chunkindex",
			artifactName: "1.2.3.caidx",
			wantVer:      semver.MustParse("1.2.3"),
			wantExt:      ".caidx",
		},
		{
			name:         "invalid_no_extension",
			artifactName: "test-app-1.2.3",
			wantErr:      true,
		},
		{
			name:         "invalid_bad_version",
			artifactName: "test-app-1.s.3.tar.gz",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVer, gotExt, err := parseArtifactName(tt.artifactName)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArtifactName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVer, tt.wantVer) {
				t.Errorf("parseArtifactName() gotVer = %v, want %v", gotVer, tt.wantVer)
			}
			if gotExt != tt.wantExt {
				t.Errorf("parseArtifactName() gotExt = %v, want %v", gotExt, tt.wantExt)
			}
		})
	}
}
