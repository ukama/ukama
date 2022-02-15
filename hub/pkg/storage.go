package pkg

import (
	"context"
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
const cappsRoot = "capps/"

type InvalidInputError struct {
	Message string
}

func (e InvalidInputError) Error() string {
	return e.Message
}

type AritfactInfo struct {
	Url       string
	Version   string
	CreatedAt time.Time
	SizeBytes int64
}

type CappInfo struct {
	Name string `json:"name"`
}

type Storage interface {
	PutFile(ctx context.Context, artifactName string, version *semver.Version, content io.Reader) error
	GetFile(ctx context.Context, artifactName string, version *semver.Version) (reader io.ReadCloser, err error)
	List(ctx context.Context, artifactName string) (*[]AritfactInfo, error)
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

	err := m.createBucketIfMissing()
	if err != nil {
		logrus.Fatalf("Failed to create bucket %s: %v", m.bucketName, err)
	}

	return m
}

func formatCappPath(artifactName string) string {
	return cappsRoot + artifactName
}

func formatCappFilename(artifactName string, version *semver.Version) string {
	return formatCappPath(artifactName) + "/" + version.String() + TarGzExtension
}

func (m *MinioWrapper) PutFile(ctx context.Context, artifactName string, version *semver.Version, content io.Reader) error {
	if !NameRegex.MatchString(artifactName) {
		return InvalidInputError{Message: "artifact name should not contain dot"}
	}

	n, err := m.minioClient.PutObject(ctx, m.bucketName, formatCappFilename(artifactName, version), content, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	logrus.Infof("Successfully uploaded %s of size %v\n", artifactName, n.Size)
	return nil
}

func (m *MinioWrapper) GetFile(ctx context.Context, artifactName string, version *semver.Version) (reader io.ReadCloser, err error) {
	o, err := m.minioClient.GetObject(ctx, m.bucketName, formatCappFilename(artifactName, version), minio.GetObjectOptions{})

	if err != nil {
		return nil, err
	}

	return o, nil
}

func (m *MinioWrapper) List(ctx context.Context, artifactName string) (*[]AritfactInfo, error) {
	path := formatCappPath(artifactName) + "/"
	logrus.Infof("Listing objects in %s", path)
	objectCh := m.minioClient.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{
		Prefix:       path,
		Recursive:    false,
		WithMetadata: false,
	})

	ls := []AritfactInfo{}

	for object := range objectCh {
		if object.Err != nil {
			logrus.Errorf("Failed to list objects: %v", object.Err)
			return nil, object.Err
		}
		if !strings.HasSuffix(object.Key, "") {
			continue
		}
		version := strings.TrimSuffix(strings.TrimPrefix(object.Key, formatCappPath(artifactName)+"/"), TarGzExtension)
		_, err := semver.NewVersion(version)
		if err != nil {
			logrus.Errorf("Failed to parse version %s: %v", version, err)
			version = "INVALID_VERSION_FORMAT"
		}

		ls = append(ls, AritfactInfo{
			Url:       formatCappPath(artifactName),
			Version:   version,
			CreatedAt: object.LastModified,
			SizeBytes: object.Size,
		})
	}

	return &ls, nil
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
		logrus.Fatalf("Bucket name should start with artifact_hub")
	}

	exists, err := m.minioClient.BucketExists(ctx, m.bucketName)
	if err == nil && exists {
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
			logrus.Warnf("Bucket creation is not supported by the server. Try enabling debug mode if minio is running in FS mode")
			return nil
		} else {
			return err
		}
	}
	return nil
}
