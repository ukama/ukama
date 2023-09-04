package pkg_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ukama/ukama/systems/hub/hub/mocks"
	"github.com/ukama/ukama/systems/hub/hub/pkg"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_chunker_Chunk(t *testing.T) {
	storage := mocks.Storage{}
	appName := "test-app"
	v := semver.MustParse("1.2.3")
	storeBaseURL := "http://store.example.com/artifacts"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.Contains(t, string(b), `"s3+`+storeBaseURL+`"`)
	}))

	storage.On("PutFile", mock.Anything, appName, v, pkg.ChunkIndexExtension,
		mock.Anything).Return("", nil)
	ch := pkg.NewChunker(&pkg.ChunkerConfig{
		Host: s.URL,
	}, &storage)

	err := ch.Chunk("test-app", v, storeBaseURL)
	assert.NoError(t, err)
	storage.AssertExpectations(t)
}
