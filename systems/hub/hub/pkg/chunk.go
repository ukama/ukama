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

	"github.com/ukama/ukama/systems/common/errors"

	"github.com/Masterminds/semver/v3"

	log "github.com/sirupsen/logrus"
)

const ChunksPath = "/v1/chunks"

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
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(ch.conf.TimeoutSecond)*time.Second)
	defer cancel()

	client := &http.Client{}
	fPath := fmt.Sprintf("%s/%s/%s", name, ver.String(), TarGzExtension)

	json, err := json.Marshal(chunkRequest{
		Store: "s3+" + strings.TrimSuffix(fileStorageUrl, fPath),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("%s%s/%s/%s", ch.conf.Host, ChunksPath, name, ver.String()), bytes.NewBuffer(json))
	if err != nil {
		return errors.Wrap(err, "failed create chunk request")
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send chunk request")
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("chunk server returned status code %d", resp.StatusCode)

		// print response body
		b, errB := io.ReadAll(resp.Body)
		if errB != nil {
			log.Errorf("failed to read response body. Error: %+v", err)
		}

		log.Errorf("response body: %s", string(b))

		return fmt.Errorf("failed to chunk file")
	}

	_, err = ch.storage.PutFile(ctx, name, ver, ChunkIndexExtension, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to save index file")
	}

	return nil
}
