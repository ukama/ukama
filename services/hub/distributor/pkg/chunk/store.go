package chunk

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Masterminds/semver/v3"
	casync "github.com/folbricht/desync"
	minio "github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/hub/distributor/pkg"
	mc "github.com/ukama/ukama/services/hub/hub/pkg"
)

/* Local store path for artifact
storepath/filename-version.tgz */
func localStoreFilePath(fname string, fversion *semver.Version, fstore string) string {
	n := fname + "-" + fversion.String() + ".tar.gz"
	cpath := filepath.Join(fstore, n)
	return cpath
}

/* Read artifact from local store */
func GetArtifactFromLocalStore(ctx context.Context, fname string, fversion *semver.Version, fstore string, dest string) error {

	cpath := localStoreFilePath(fname, fversion, fstore)

	logrus.Debugf("Copying a image file %s to %s", cpath, dest)
	sourceFile, err := os.Open(cpath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	/* Create a copy of file */
	osFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer osFile.Close()

	bytes, err := io.Copy(osFile, sourceFile)
	if err != nil {
		return err
	}
	logrus.Debugf("Copied a file %s with %d bytes.", dest, bytes)

	return nil
}

/* Preparing a store access to read artifacts from */
func GetArtifactFromS3(ctx context.Context, fname string, fversion *semver.Version, fstore string, dest string) error {

	/* Get store config */
	artCfg, err := pkg.GetLocalStoreCredentialsFor(fstore)
	if err != nil {
		return err
	}

	if artCfg == nil {
		return fmt.Errorf("no config for artifact store found")
	}

	/* get minio client for store */
	store := mc.NewMinioWrapper(&artCfg.MinioConfig)
	if store == nil {
		return fmt.Errorf("could not connect to artifact store")
	}

	logrus.Debugf("Reading file %s from s3 store to %s.", fname, dest)
	r, err := store.GetFile(ctx, fname, fversion, mc.TarGzExtension)
	if err != nil {
		logrus.Errorf("failed to read a file %s from store %s : %s", fname, fstore, err.Error())
		return err
	}

	defer r.Close()

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("opening source archive: %v", err)
	}
	defer file.Close()

	bytes, err := io.Copy(file, r)
	if err != nil {
		logrus.Errorf("Failed to get artifact data from the artifact store.")
		return err
	}

	logrus.Debugf("Copied %s %d bytes from the store %s to %s.", fname, bytes, fstore, dest)

	return err
}

/* Multiple store with chache */
func MultiStore(cmdOpt casync.StoreOptions, storeLocations ...string) (casync.Store, error) {
	// Combine all stores into one router
	store, err := multiStoreWithRouter(cmdOpt, storeLocations...)
	if err != nil {
		return nil, err
	}
	return store, nil
}

/* multistore with router */
func multiStoreWithRouter(cmdOpt casync.StoreOptions, storeLocations ...string) (casync.Store, error) {
	var stores []casync.Store
	for _, location := range storeLocations {
		s, err := storeGroup(location, cmdOpt)
		if err != nil {
			return nil, err
		}
		stores = append(stores, s)
	}

	return casync.NewStoreRouter(stores...), nil
}

/* store group */
func storeGroup(location string, cmdOpt casync.StoreOptions) (casync.Store, error) {
	if !strings.ContainsAny(location, "|") {
		return storeFromLocation(location, cmdOpt)
	}
	var stores []casync.Store
	members := strings.Split(location, "|")
	for _, m := range members {
		s, err := storeFromLocation(m, cmdOpt)
		if err != nil {
			return nil, err
		}
		stores = append(stores, s)
	}
	return casync.NewFailoverGroup(stores...), nil
}

/* Writable store */
func WritableStore(location string, opt casync.StoreOptions) (casync.WriteStore, error) {
	s, err := storeFromLocation(location, opt)
	if err != nil {
		return nil, err
	}
	castore, ok := s.(casync.WriteStore)
	if !ok {
		return nil, fmt.Errorf("store '%s' does not support writing", location)
	}
	return castore, nil
}

/* Parse a single store URL or path and return an initialized instance of it */
func storeFromLocation(location string, opt casync.StoreOptions) (casync.Store, error) {
	loc, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("unable to parse store location %s : %s", location, err)
	}

	var s casync.Store
	switch loc.Scheme {
	case "s3+http", "s3+https":

		/* Credentials */
		s3Creds, region, err := pkg.GetS3CredentialsFor(location)
		if err != nil {
			logrus.Errorf("Failed to get credintilas for S3 store: %s", err.Error())
			return nil, err
		}

		/* Lookup */
		lookup := minio.BucketLookupAuto
		ls := loc.Query().Get("lookup")
		switch ls {
		case "dns":
			lookup = minio.BucketLookupDNS
		case "path":
			lookup = minio.BucketLookupPath
		case "", "auto":
		default:
			return nil, fmt.Errorf("unknown S3 bucket lookup type: %q", s)
		}

		/* S3 Store */
		s, err = casync.NewS3Store(loc, s3Creds, *region, opt, lookup)
		if err != nil {
			return nil, err
		}
	default:

		/*local store */
		local, err := casync.NewLocalStore(location, opt)
		if err != nil {
			return nil, err
		}
		s = local
		if runtime.GOOS == "windows" {
			s = casync.NewWriteDedupQueue(local)
		}
	}
	return s, nil
}
