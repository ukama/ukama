package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/errors"
)

type Chunker interface {
	Chunk(name string, ver *semver.Version, fileStorageUrl string) error
}

type chunker struct {
	conf    *ChunkerConfig
	storage Storage
}

func NewChunker(conf *ChunkerConfig, storage Storage) Chunker {
	return &chunker{
		conf:    conf,
		storage: storage,
	}
}

type chunkRequest struct {
	Store string `json:"store"`
}

// Chunk sends request to chunk server to chunk the file and uploads chunk index file to storage
func (ch *chunker) Chunk(name string, ver *semver.Version, fileStorageUrl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ch.conf.TimeoutSecond)*time.Second)
	defer cancel()

	client := &http.Client{}
	fPath := fmt.Sprintf("%s/%s/%s", name, ver.String(), TarGzExtension)

	json, err := json.Marshal(chunkRequest{
		Store: "s3+" + strings.TrimSuffix(fileStorageUrl, fPath),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/chunk/%s/%s", ch.conf.Host, name, ver.String()), bytes.NewBuffer(json))
	if err != nil {
		return errors.Wrap(err, "failed create chunk request")
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send chunk request")
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("chunk server returned status code %d", resp.StatusCode)
		// print response body
		b, errB := io.ReadAll(resp.Body)
		if errB != nil {
			logrus.Errorf("failed to read response body. Error: %+v", err)
		}
		logrus.Errorf("response body: %s", string(b))
		return errors.Errorf("failed to chunk file")
	}

	_, err = ch.storage.PutFile(ctx, name, ver, ChunkIndexExtension, resp.Body)
	return errors.Wrap(err, "failed to save index file")
}
