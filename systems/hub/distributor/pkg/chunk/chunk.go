/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package chunk

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/hub/distributor/pkg"
	"github.com/ukama/ukama/systems/hub/distributor/pkg/archiver"

	"github.com/Masterminds/semver/v3"

	casync "github.com/folbricht/desync"
	log "github.com/sirupsen/logrus"
)

var S3Schemes = []string{"s3+http", "s3+https"}

type tarOptions struct {
	casync.LocalFSOptions
}

/* file name with extension */
func fileWithExt(name string, ext string) string {
	if strings.HasPrefix(ext, ".") {
		return name + ext
	} else {
		return name + "." + ext
	}
}

/* Convert Index file into string */
func IndexToString(idx casync.Index) (*string, error) {
	r, w := io.Pipe()

	go func() {
		defer w.Close()

		_, err := idx.WriteTo(w)
		if err != nil {
			return
		}
	}()

	buf := new(strings.Builder)

	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	str := buf.String()
	//log.Debugf("Length is %d Content is : \n %s", n, str)

	return &str, nil
}

/* CA sync store options */
func storeOptions(storeCfg *pkg.StoreConfig, chunkCfg *pkg.ChunkConfig) *casync.StoreOptions {
	var opt = &casync.StoreOptions{
		N:             chunkCfg.N,
		ClientCert:    storeCfg.CaCert,
		CACert:        storeCfg.ClientCert,
		ClientKey:     storeCfg.ClientKey,
		SkipVerify:    storeCfg.SkipVerify,
		TrustInsecure: storeCfg.TrustInsecure,
		ErrorRetry:    storeCfg.ErrorRetry,
		Uncompressed:  storeCfg.Uncompressed,
	}

	return opt
}

/* Store Index file */
func storeIndex(name string, idx casync.Index) error {
	log.Debugf("Writing index to a file %s", name)

	/* Write the index to file */
	i, err := os.Create(name)
	if err != nil {
		return err
	}
	defer i.Close()

	_, err = idx.WriteTo(i)

	return err
}

type Store interface {
	Read(ctx context.Context, fname string, aType string, fversion *semver.Version, fext string, fstore string, wp string) error
}

type s3Store struct {
	Scheme []string
}

type localStore struct {
	Scheme []string
}

func NewS3Store() *s3Store {
	return &s3Store{Scheme: S3Schemes}
}

/* Read file from local store*/
func (s3 *localStore) Read(ctx context.Context, fname string, aType string, fversion *semver.Version, fext string, fstore string, wp string) error {
	var (
		tgzFile string
		err     error
	)

	log.Debugf("Using local store %s to read file %s", fstore, fname)

	tgzFile = wp + fileWithExt(fname, fext)

	/* Read from local store */
	err = GetArtifactFromLocalStore(ctx, fname, aType, fversion, fstore, tgzFile)
	if err != nil {
		log.Errorf("Failed to read artifact %s from local store %s : %s ",
			fname, fstore, err.Error())

		return err
	}

	/* Extract file*/
	err = archiver.Unarchive(tgzFile, wp)
	if err != nil {
		log.Errorf("Error while extracting file %s to %s: %s", tgzFile, wp, err.Error())

		return err
	} else {
		log.Debugf("Extraction look good for file %s at location %s", tgzFile, wp)
	}

	return err
}

/* Read file from local store*/
func (s3 *s3Store) Read(ctx context.Context, fname string, aType string, fversion *semver.Version, fext string, fstore string, wp string) error {
	var (
		tgzFile string
		err     error
	)

	tgzFile = wp + fileWithExt(fname, fext)

	log.Debugf("Using S3 store %s to read file %s", fstore, fname)

	/* Read data from artifact store */
	err = GetArtifactFromS3(ctx, fname, aType, fversion, fstore, tgzFile)
	if err != nil {
		log.Errorf("Error while pulling data from artifact store: %s", err.Error())

		return err
	}

	/* Extract the file as hub always provides tar.gz */
	err = archiver.Unarchive(tgzFile, wp)
	if err != nil {
		log.Errorf("Error while extracting file %s to %s: %s", tgzFile, wp, err.Error())

		return err
	} else {
		log.Debugf("Extraction look good for file %s at location %s", tgzFile, wp)
	}

	return err
}

/* Read from store */
func ReadFromStore(ctx context.Context, fname string, aType string, fversion *semver.Version, fext string, fstore string, wp string) error {
	loc, err := url.Parse(fstore)
	if err != nil {
		return fmt.Errorf("unable to parse store location %s : %s", fstore, err)
	}

	var st Store = nil
	for _, s := range S3Schemes {
		if loc.Scheme == s {
			st = NewS3Store()
			break
		}
	}

	/* Choose local store if scheme dosen't match. */
	if st == nil {
		st = new(localStore)
	}

	err = st.Read(ctx, fname, aType, fversion, fext, fstore, wp)
	if err != nil {
		return err
	}

	return err
}

/* Remove the workplace */
func removeWorkplace(wp string) {
	os.RemoveAll(wp)
}

/* Prepare w workplace dir */
func prepareWorkplace() (string, error) {
	/* Create temp Work place for chunking */
	u1 := uuid.NewV4()

	wp := "/tmp/" + u1.String() + "/"
	err := os.Mkdir(wp, 0755)
	if err != nil {
		removeWorkplace(wp)

		return "", err
	}

	return wp, nil
}

/* Read contents to be chunked from remote or S3 server and store them on locally*/
func ReadRemoteContents(ctx context.Context, fname string, aType string, fversion *semver.Version, fext string, fstore string, wp string) (string, bool, error) {
	isDir := false

	/* Read file from store */
	err := ReadFromStore(ctx, fname, aType, fversion, fext, fstore, wp)
	if err != nil {
		return "", false, err
	}

	/* Getting the new extracted contents path */
	content := wp + fname + "." + fext

	/* Check if extracted content exist */
	fileInfo, err := os.Stat(content)
	if err != nil {
		return "", false, err
	}

	if fileInfo.IsDir() {
		isDir = true
	}

	return content, isDir, nil
}

type chunker interface {
	MakeChunk(ctx context.Context, content string, store string, wp string) (*casync.Index, error)
	GetIndexExtension() string
}

type blob struct {
	indexExt string `default:".caibx"`
	config   *pkg.ChunkConfig
	store    *pkg.StoreConfig
}

type archive struct {
	indexExt string `default:".caidx"`
	config   *pkg.ChunkConfig
	store    *pkg.StoreConfig
}

func NewArchiveChunker(c *pkg.ChunkConfig, s *pkg.StoreConfig) *archive {
	return &archive{
		config:   c,
		store:    s,
		indexExt: ".caidx",
	}
}

func NewBlobChunker(c *pkg.ChunkConfig, s *pkg.StoreConfig) *blob {
	return &blob{
		config:   c,
		store:    s,
		indexExt: ".caibx",
	}
}

func (a *archive) GetIndexExtension() string {
	return a.indexExt
}

/* Create chunks for given tar file. */
func (a *archive) MakeChunk(ctx context.Context, content string, store string, wp string) (*casync.Index, error) {
	var fs casync.FilesystemReader
	var err error

	var tarOpt tarOptions

	opt := storeOptions(a.store, a.config)

	log.Debugf("Starting chunking process for %s from store %s, opt %+v tarOpt %+v.",
		content, store, opt, tarOpt)
	if store == "" {
		return nil, fmt.Errorf("requires store location from where contents needs to be copied")
	}

	/* What to expect from content */
	switch a.config.InFormat {
	case "disk":
		local := casync.NewLocalFS(content, tarOpt.LocalFSOptions)
		fs = local
	case "tar":
		log.Debugf("Reading contents of %s", content)

		var r *os.File
		if content == "-" {
			r = os.Stdin
		} else {
			r, err = os.Open(content)
			if err != nil {
				return nil, err
			}
			defer r.Close()
		}

		var op casync.TarReaderOptions

		op.AddRoot = true
		fs = casync.NewTarReader(r, op)
	default:
		return nil, fmt.Errorf("invalid input format '%s'", a.config.InFormat)
	}

	r, w := io.Pipe()

	/*Open the target store */
	s, err := WritableStore(store, *opt)
	if err != nil {
		log.Errorf("Error while opening writable store %s", err.Error())

		return nil, err
	}
	defer s.Close()

	if s == nil {
		log.Errorf("Error Writable store %s not found", store)
		return nil, fmt.Errorf("store '%s' not found", store)
	}

	/* Get chunker */
	c, err := casync.NewChunker(r, a.config.MinChunkSize, a.config.AvgChunkSize, a.config.MaxChunkSize)
	if err != nil {
		log.Errorf("Error while getting new chunker %s", err.Error())
		return nil, err
	}

	/* Run the tar bit in a goroutine, writing to the pipe */
	var tarErr error
	go func() {
		tarErr = casync.Tar(ctx, w, fs)
		w.Close()
	}()

	/* Store chunks */
	index, err := casync.ChunkStream(ctx, c, s, opt.N)
	if err != nil {
		log.Errorf("Error while chunking %s", err.Error())
		return nil, err
	}

	index.Index.FeatureFlags |= casync.TarFeatureFlags

	/* Any issues with tar */
	if tarErr != nil {
		log.Errorf("Error realted to tar %s", tarErr.Error())
		return nil, tarErr
	}

	return &index, nil
}

func (b *blob) GetIndexExtension() string {
	return b.indexExt
}

/* Create chunks for the blobs */
func (b *blob) MakeChunk(ctx context.Context, content string, storeLoc string, wp string) (*casync.Index, error) {
	// Open the target store if one was given
	var s casync.WriteStore
	opt := storeOptions(b.store, b.config)

	s, err := WritableStore(storeLoc, *opt)
	if err != nil {
		log.Errorf("failed to get writable store: %v", err)
		return nil, err
	}
	defer s.Close()

	if s == nil {
		log.Errorf("Err:: No store avalibale.")
		return nil, fmt.Errorf("not able to find store")
	}

	/* Create a index file. */
	log.Debugf("Creating index for a  file %s.", content)
	index, stats, err := casync.IndexFromFile(ctx, content, b.config.N, b.config.MinChunkSize,
		b.config.AvgChunkSize, b.config.MaxChunkSize, casync.NullProgressBar{})
	if err != nil {
		return nil, err
	}

	log.Debugf("Stats:: Chunk Accepted: %d Chunk Produced: %d",
		stats.ChunksAccepted, stats.ChunksProduced)

	/* Create chunks for the file */
	log.Debugf("Creating chunks for a  file %s and storing to %s.", content, storeLoc)
	err = casync.ChopFile(ctx, content, index.Chunks, s, b.config.N, casync.NullProgressBar{})
	if err != nil {
		return nil, err
	}

	return &index, nil
}

/* Handler for creating chunks */
func CreateChunks(ctx context.Context, storeCfg *pkg.StoreConfig, chunkCfg *pkg.ChunkConfig, fname string, aType string, fversion *semver.Version, fstore string) (*casync.Index, error) {
	var (
		index     *casync.Index
		err       error
		indexFile string
		ch        chunker
	)

	/* Prepare workplace */
	wp, err := prepareWorkplace()
	if err != nil {
		log.Errorf("Failed to prepare workplace %s", err.Error())

		return nil, err
	}

	/* Only first store is considered */
	storeLoc := chunkCfg.Stores[0]

	/* Read contents */
	content, isFS, err := ReadRemoteContents(ctx, fname, aType, fversion, chunkCfg.Extension, fstore, wp)
	if err != nil {
		log.Errorf("Failed to read contents for chunking %s", err.Error())

		return nil, err
	}
	log.Debugf("Workplace %s, Contents %s FS: %t Store %s", wp, content, isFS, storeLoc)

	/* Start chunking process */
	if isFS {
		log.Debugf("Creating chunks for FS.")
		ch = NewArchiveChunker(chunkCfg, storeCfg)

	} else {
		log.Debugf("Creating chunks for Blob.")
		ch = NewBlobChunker(chunkCfg, storeCfg)
	}

	index, err = ch.MakeChunk(ctx, content, storeLoc, wp)
	if err != nil {
		log.Errorf("Error while creating chunks for %s from %s: %s", fname, storeLoc, err.Error())

		return nil, err
	}

	/* Store index file */
	indexFile = wp + ch.GetIndexExtension()
	err = storeIndex(indexFile, *index)
	if err != nil {

		log.Errorf("failed to write index file.")

		return nil, fmt.Errorf("failed to write index file")
	}

	/* Remove work place */
	//removeWorkplace(wp)

	return index, err
}
