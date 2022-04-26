package pkg

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"

	"io"
	"log"
	"regexp"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

// app name regex. Follows OCI image naming standards
var NameRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_.-]*$")

const BucketNamePrefix = "artifact-hub-"
const TarGzExtension = ".tar.gz"
const ChunkIndexExtension = ".caidx"
const cappsRoot = "capps/"

type InvalidInputError struct {
	Message string
}

func (e InvalidInputError) Error() string {
	return e.Message
}

type AritfactInfo struct {
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	SizeBytes int64     `json:"size_bytes"`
	Chunked   bool      `json:"chunked"`
}

type CappInfo struct {
	Name string `json:"name"`
}

type Storage interface {
	PutFile(ctx context.Context, artifactName string, version *semver.Version, ext string, content io.Reader) (string, error)
	GetFile(ctx context.Context, artifactName string, version *semver.Version, ext string) (reader io.ReadCloser, err error)
	ListVersions(ctx context.Context, artifactName string) (*[]AritfactInfo, error)
	ListApps(ctx context.Context) (*[]CappInfo, error)
}

type MinioWrapper struct {
	minioClient *minio.Client
	bucketName  string
	region      string
}

// host in host:port format
func NewMinioWrapper(options *MinioConfig) *MinioWrapper {
	m := &MinioWrapper{
		minioClient: getClient(options.Endpoint, options.AccessKey, options.SecretKey),
		bucketName:  BucketNamePrefix + options.BucketSuffix,
		region:      options.Region,
	}

	if !options.SkipBucketCreation {
		logrus.Infof("Creating bucket %s", m.bucketName)
		err := m.createBucketIfMissing()
		if err != nil {
			logrus.Fatalf("Failed to create bucket %s: %v", m.bucketName, err)
		}
	} else {
		logrus.Infof("Skipping bucket creation")
	}

	return m
}

func formatCappPath(artifactName string) string {
	return cappsRoot + artifactName
}

func formatCappFilename(artifactName string, version *semver.Version, ext string) string {
	return formatCappPath(artifactName) + "/" + version.String() + ext
}

// PutFile stores the file in storage. Based on input params we build the file path: /<artifactName>/<version>.<ext>
// artifactName - name of the artifact without extension
// version - artifact version
// ext - extension, use consts declared in this package to stay consistent
// content - content of file
// returns remote location of the file or error
func (m *MinioWrapper) PutFile(ctx context.Context, artifactName string, version *semver.Version, ext string, content io.Reader) (string, error) {
	if !NameRegex.MatchString(artifactName) {
		return "", InvalidInputError{Message: "artifact name should not contain dot"}
	}

	n, err := m.minioClient.PutObject(ctx, m.bucketName, formatCappFilename(artifactName, version, ext), content, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	logrus.Infof("Successfully uploaded %s of size %v\n", artifactName, n.Size)
	if IsDebugMode {
		logrus.Infof("File info: %+v", n)
	}
	return n.Location, nil
}

func (m *MinioWrapper) GetFile(ctx context.Context, artifactName string, version *semver.Version, ext string) (reader io.ReadCloser, err error) {
	fPath := formatCappFilename(artifactName, version, ext)
	logrus.Infof("Downloading %s from bucket %s", fPath, m.bucketName)
	o, err := m.minioClient.GetObject(ctx, m.bucketName, fPath, minio.GetObjectOptions{})

	if err != nil {
		return nil, err
	}

	return o, nil
}

func (m *MinioWrapper) ListVersions(ctx context.Context, artifactName string) (*[]AritfactInfo, error) {
	path := formatCappPath(artifactName) + "/"
	logrus.Infof("Listing objects in %s", path)
	objectCh := m.minioClient.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{
		Prefix:       path,
		Recursive:    false,
		WithMetadata: false,
	})

	ls := map[string]AritfactInfo{}
	chunked := map[string]bool{}

	for object := range objectCh {
		if object.Err != nil {
			logrus.Errorf("Failed to list objects: %v", object.Err)
			return nil, object.Err
		}
		logrus.Infof("Listing object %s", object.Key)

		if strings.HasSuffix(object.Key, ChunkIndexExtension) {
			chunked[strings.TrimSuffix(object.Key, ChunkIndexExtension)] = true
		}

		if strings.HasSuffix(object.Key, TarGzExtension) {
			version := strings.TrimSuffix(strings.TrimPrefix(object.Key, formatCappPath(artifactName)+"/"), TarGzExtension)
			_, err := semver.NewVersion(version)
			if err != nil {
				logrus.Errorf("Failed to parse version %s: %v", version, err)
				version = "INVALID_VERSION_FORMAT"
			}

			ls[strings.TrimSuffix(object.Key, TarGzExtension)] = AritfactInfo{
				Version:   version,
				CreatedAt: object.LastModified,
				SizeBytes: object.Size,
			}
		}
	}

	var result []AritfactInfo
	for k, v := range ls {
		v.Chunked = chunked[k]
		result = append(result, v)
	}

	return &result, nil
}

func (m *MinioWrapper) ListApps(ctx context.Context) (*[]CappInfo, error) {

	logrus.Infof("Listing all objects")
	objectCh := m.minioClient.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{
		Prefix:       cappsRoot,
		Recursive:    false,
		WithMetadata: false,
	})

	ls := []CappInfo{}

	for object := range objectCh {
		if object.Err != nil {
			logrus.Errorf("Failed to list objects: %v", object.Err)
			return nil, object.Err
		}

		ls = append(ls, CappInfo{Name: strings.TrimSuffix(strings.TrimPrefix(object.Key, cappsRoot), "/")})
	}

	return &ls, nil
}

// host in host:port format
func getClient(host string, accessKeyId string, accessSecret string) *minio.Client {

	// Initialize minio client object.
	c, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, accessSecret, ""),
		Secure: false,
	})
	if err != nil {
		log.Printf("Minio client initialization failed")
		log.Fatalln(err)
	}

	return c
}

func (m *MinioWrapper) createBucketIfMissing() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !strings.HasPrefix(m.bucketName, BucketNamePrefix) {
		return fmt.Errorf("bucket name should start with artifact_hub")
	}

	exists, err := m.minioClient.BucketExists(ctx, m.bucketName)
	if err == nil && exists {
		logrus.Infof("Bucket %s already exists", m.bucketName)
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to check if bucket exists")
	}

	logrus.Infof("Bucket %s does not exist, creating it", m.bucketName)
	objLocking := true
	if IsDebugMode {
		objLocking = false
	}

	err = m.minioClient.MakeBucket(ctx, m.bucketName, minio.MakeBucketOptions{
		ObjectLocking: objLocking,
	})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NotImplemented" {
			logrus.Errorf("Bucket creation is not supported by the server. Try enabling debug mode if minio is running in FS mode")
			return fmt.Errorf("bucket creation not supported by the server")
		} else {
			return err
		}
	}
	return nil
}
