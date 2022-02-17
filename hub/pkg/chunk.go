package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"net/http"
	"time"
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
	Storage string `json:"storage"`
}

// Chunk sends request to chunk server to chunk the file and uploads chunk index file to storage
func (ch *chunker) Chunk(name string, ver *semver.Version, fileStorageUrl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ch.conf.TimeoutSecond)*time.Second)
	defer cancel()

	client := &http.Client{}

	json, err := json.Marshal(chunkRequest{
		Storage: fileStorageUrl,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s/%s", ch.conf.Host, name, ver.String()), bytes.NewBuffer(json))
	if err != nil {
		return errors.Wrap(err, "failed to chunk")
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to save index file")
	}

	_, err = ch.storage.PutFile(ctx, name, ver, ChunkIndexExtension, resp.Body)
	return errors.Wrap(err, "failed to save index file")
}
