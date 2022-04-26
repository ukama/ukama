package chunk

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"

	casync "github.com/folbricht/desync"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/hub/distributor/pkg"
	"github.com/ukama/ukamaX/hub/distributor/pkg/archiver"
)

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
	//logrus.Debugf("Length is %d Content is : \n %s", n, str)

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

	logrus.Debugf("Writing index to a file %s", name)

	/* Write the index to file */
	i, err := os.Create(name)
	if err != nil {
		return err
	}
	defer i.Close()
	_, err = idx.WriteTo(i)
	return err
}

/* Read file from local store*/
func ReadFromLocalStore(ctx context.Context, fname string, fversion *semver.Version, fext string, fstore string, wp string) error {

	var (
		tgzFile string
		err     error
	)

	logrus.Debugf("Using local store %s to read file %s", fstore, fname)

	tgzFile = wp + fileWithExt(fname, fext)

	/* Read from local store */
	err = GetArtifactFromLocalStore(ctx, fname, fversion, fstore, tgzFile)
	if err != nil {
		logrus.Errorf("Failed to read artifact %s from local store %s : %s ", fname, fstore, err.Error())
		return err
	}

	/* Extract file*/
	err = archiver.Unarchive(tgzFile, wp)
	if err != nil {
		logrus.Errorf("Error while extracting file %s to %s: %s", tgzFile, wp, err.Error())
		return err
	} else {
		logrus.Debugf("Extraction look good for file %s at location %s", tgzFile, wp)
	}

	return err
}

/* Read file from local store*/
func ReadFromS3Store(ctx context.Context, fname string, fversion *semver.Version, fext string, fstore string, wp string) error {

	var (
		tgzFile string
		err     error
	)

	tgzFile = wp + fileWithExt(fname, fext)

	logrus.Debugf("Using S3 store %s to read file %s", fstore, fname)

	/* Read data from artifact store */
	err = GetArtifactFromS3(ctx, fname, fversion, fstore, tgzFile)
	if err != nil {
		logrus.Errorf("Error while pulling data from artifact store: %s", err.Error())
		return err
	}

	/* Extract the file as hub always provides tar.gz */
	err = archiver.Unarchive(tgzFile, wp)
	if err != nil {
		logrus.Errorf("Error while extracting file %s to %s: %s", tgzFile, wp, err.Error())
		return err
	} else {
		logrus.Debugf("Extraction look good for file %s at location %s", tgzFile, wp)
	}

	return err
}

/* Read from store */
func ReadFromStore(ctx context.Context, fname string, fversion *semver.Version, fext string, fstore string, wp string) error {
	loc, err := url.Parse(fstore)
	if err != nil {
		return fmt.Errorf("unable to parse store location %s : %s", fstore, err)
	}

	switch loc.Scheme {
	case "s3+http", "s3+https":
		err := ReadFromS3Store(ctx, fname, fversion, fext, fstore, wp)
		if err != nil {
			return err
		}
	default:
		err := ReadFromLocalStore(ctx, fname, fversion, fext, fstore, wp)
		if err != nil {
			return err
		}
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
	u1, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	wp := "/tmp/" + u1.String() + "/"
	err = os.Mkdir(wp, 0755)
	if err != nil {
		removeWorkplace(wp)
		return "", err
	}

	return wp, nil
}

/* Read contents to be chunked from remote or S3 server and store them on locally*/
func ReadRemoteContents(ctx context.Context, fname string, fversion *semver.Version, fext string, fstore string, wp string) (string, bool, error) {
	isDir := false

	/* Read file from store */
	err := ReadFromStore(ctx, fname, fversion, fext, fstore, wp)
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

/* Create chunks for given tar file. */
func CreateArchivedChunk(ctx context.Context, storeCfg *pkg.StoreConfig, chunkCfg *pkg.ChunkConfig, content string, store string, wp string) (*casync.Index, error) {
	var fs casync.FilesystemReader
	var err error

	var tarOpt tarOptions

	opt := storeOptions(storeCfg, chunkCfg)

	logrus.Debugf("Starting chunking process for %s from store %s, opt %+v tarOpt %+v.", content, store, opt, tarOpt)
	if store == "" {
		return nil, errors.New("requires store location from where contents needs to be copied")
	}

	/* What to expect from content */
	switch chunkCfg.InFormat {
	case "disk":
		local := casync.NewLocalFS(content, tarOpt.LocalFSOptions)
		fs = local
	case "tar":
		logrus.Debugf("Reading contents of %s", content)
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
		return nil, fmt.Errorf("invalid input format '%s'", chunkCfg.InFormat)
	}

	r, w := io.Pipe()

	/*Open the target store */
	s, err := WritableStore(store, *opt)
	if err != nil {
		logrus.Errorf("Error while opening writable store %s", err.Error())
		return nil, err
	}
	defer s.Close()

	if s == nil {
		logrus.Errorf("Error Writable store %s not found", err.Error())
		return nil, fmt.Errorf("store '%s' not found", store)
	}

	/* Get chunker */
	c, err := casync.NewChunker(r, chunkCfg.MinChunkSize, chunkCfg.AvgChunkSize, chunkCfg.MaxChunkSize)
	if err != nil {
		logrus.Errorf("Error while getting new chunker %s", err.Error())
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
		logrus.Errorf("Error while chunking %s", err.Error())
		return nil, err
	}

	index.Index.FeatureFlags |= casync.TarFeatureFlags

	/* Any issues with tar */
	if tarErr != nil {
		logrus.Errorf("Error realted to tar %s", tarErr.Error())
		return nil, tarErr
	}

	return &index, nil
}

/* Create chunks for the blobs */
func CreateChunkForBlob(ctx context.Context, storeCfg *pkg.StoreConfig, chunkCfg *pkg.ChunkConfig, content string, store string, wp string) (*casync.Index, error) {

	// Open the target store if one was given
	var s casync.WriteStore
	opt := storeOptions(storeCfg, chunkCfg)

	s, err := WritableStore(store, *opt)
	if err != nil {
		return nil, err
	}
	defer s.Close()

	if s == nil {
		logrus.Errorf("Err:: No store avalibale.")
		return nil, fmt.Errorf("not able to find store")
	}

	/* Create a index file. */
	logrus.Debugf("Creating index for a  file %s.", content)
	//pbi := NewProgressBar("index ")
	index, stats, err := casync.IndexFromFile(ctx, content, chunkCfg.N, chunkCfg.MinChunkSize, chunkCfg.AvgChunkSize, chunkCfg.MaxChunkSize, nil)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("Stats:: Chunk Accepted: %d Chunk Produced: %d", stats.ChunksAccepted, stats.ChunksProduced)

	/* Create chunks for the file */
	logrus.Debugf("Creating chunks for a  file %s and storing to %s.", content, store)
	if s != nil {
		err = casync.ChopFile(ctx, content, index.Chunks, s, chunkCfg.N, nil)
		if err != nil {
			return nil, err
		}
	}
	return &index, nil
}

/* Handler for creating chunks */
func CreateChunks(ctx context.Context, storeCfg *pkg.StoreConfig, chunkCfg *pkg.ChunkConfig, fname string, fversion *semver.Version, fstore string) (*casync.Index, error) {

	var (
		index     *casync.Index
		err       error
		indexFile string
	)

	/* Prepare workplace */
	wp, err := prepareWorkplace()
	if err != nil {
		logrus.Errorf("Failed to prepare workplace %s", err.Error())
		return nil, err
	}

	/* Only first store is considered */
	storePath := chunkCfg.Stores[0]

	/* Read contents */
	content, isFS, err := ReadRemoteContents(ctx, fname, fversion, chunkCfg.Extension, fstore, wp)
	if err != nil {
		logrus.Errorf("Failed to read contents for chunking %s", err.Error())
		return nil, err
	}
	logrus.Debugf("Workplace %s, Contents %s FS: %t Store %s", wp, content, isFS, storePath)

	/* Start chunking process */
	if isFS {

		logrus.Debugf("Creating archived chunks for FS.")
		index, err = CreateArchivedChunk(ctx, storeCfg, chunkCfg, content, storePath, wp)
		indexFile = wp + "index.caidx"

	} else {

		logrus.Debugf("Creating chunks.")
		index, err = CreateChunkForBlob(ctx, storeCfg, chunkCfg, content, storePath, wp)
		indexFile = wp + "index.caibx"

	}

	if err != nil {
		logrus.Errorf("Error while creating chunks for %s from %s: %s", fname, storePath, err.Error())
		return nil, err
	}

	/* Store index file */
	err = storeIndex(indexFile, *index)
	if err != nil {
		logrus.Errorf("failed to write index file.")
		return nil, fmt.Errorf("failed to write index file")
	}

	/* Remove work place */
	//removeWorkplace(wp)

	return index, err
}
