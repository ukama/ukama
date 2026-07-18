/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/errors"

	"github.com/Masterminds/semver/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	log "github.com/sirupsen/logrus"
)

// app name regex. Follows OCI image naming standards
var NameRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_.-]*$")

const BucketNamePrefix = "hub-"
const TarGzExtension = ".tar.gz"
const ChunkIndexExtension = ".caibx"
const appsRoot = "apps/"

// Immutability: the content digest stored as object user-metadata and the
// canonical response header MinIO returns it under.
const ContentDigestMetaKey = "content-sha256"
const contentDigestHeader = "X-Amz-Meta-Content-Sha256"

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

// AppLatest is one app plus its newest version, for the list?latest=true path.
type AppLatest struct {
	Name   string       `json:"name"`
	Latest AritfactInfo `json:"latest"`
}

type Storage interface {
	// PutFile now takes an optional user-metadata map (pass nil when none).
	PutFile(ctx context.Context, artifactName string, artifactType string, version *semver.Version, ext string, content io.Reader, meta map[string]string) (string, error)
	GetFile(ctx context.Context, artifactName string, artifactType string, version *semver.Version, ext string) (reader io.ReadCloser, err error)
	// StatFile reports existence and the stored content digest ("" if the object has none).
	StatFile(ctx context.Context, artifactName string, artifactType string, version *semver.Version, ext string) (exists bool, digest string, err error)
	ListVersions(ctx context.Context, artifactName string, artifactType string) (*[]AritfactInfo, error)
	ListApps(ctx context.Context, artifactType string) ([]string, error)
	ListLatestPerApp(ctx context.Context, artifactType string) ([]AppLatest, error)
	GetEndpoint() string
	// StoreBaseURL is the s3 base the distributor uses to locate an artifact's source bucket.
	StoreBaseURL(artifactType string) string
	ValidateArtifactType(artifactType string) error
}

type MinioWrapper struct {
	minioClient   *minio.Client
	bucketSuffix  string
	region        string
	typeToNameMap map[string]string
}

// host in host:port format
func NewMinioWrapper(options *MinioConfig) *MinioWrapper {
	m := &MinioWrapper{
		minioClient:   getClient(options.Endpoint, options.AccessKey, options.SecretKey),
		bucketSuffix:  options.BucketSuffix,
		region:        options.Region,
		typeToNameMap: options.ArtifactTypeBucketMap,
	}

	if !options.SkipBucketCreation {
		for _, bucket := range options.ArtifactTypeBucketMap {
			bucketName := BucketNamePrefix + bucket + m.bucketSuffix
			log.Infof("Creating bucket %s", bucketName)

			err := m.createBucketIfMissing(bucketName)
			if err != nil {
				log.Fatalf("Failed to create bucket %s: %v", bucketName, err)
			}
		}
	} else {
		log.Infof("Skipping bucket creation")
	}

	return m
}

func formatAppPath(artifactName string) string {
	return appsRoot + artifactName
}

func formatAppFilename(artifactName string, version *semver.Version, ext string) string {
	return formatAppPath(artifactName) + "/" + version.String() + ext
}

func (m *MinioWrapper) GetBucketName(artifactType string) string {
	return BucketNamePrefix + m.typeToNameMap[artifactType] + m.bucketSuffix
}

func (m *MinioWrapper) ValidateArtifactType(artifactType string) error {
	if _, ok := m.typeToNameMap[artifactType]; !ok {
		return fmt.Errorf("%s type artifact not supported", artifactType)
	}
	return nil
}

// PutFile stores the file in storage. Based on input params we build the file path: /<artifactName>/<version>.<ext>
// artifactName - name of the artifact without extension
// version - artifact version
// ext - extension, use consts declared in this package to stay consistent
// content - content of file
// returns remote location of the file or error
func (m *MinioWrapper) PutFile(ctx context.Context, artifactName string, artifactType string, version *semver.Version, ext string, content io.Reader, meta map[string]string) (string, error) {
	log.Infof("Uploading %s-%s to storage", artifactName, version.String())

	bucket := m.GetBucketName(artifactType)
	log.Infof("Uploading %s-%s to storage to bucket %s", artifactName, version.String(), bucket)
	if !NameRegex.MatchString(artifactName) {
		return "", InvalidInputError{Message: "artifact name should not contain dot"}
	}

	opts := minio.PutObjectOptions{}
	if len(meta) > 0 {
		opts.UserMetadata = meta
	}

	n, err := m.minioClient.PutObject(ctx, bucket, formatAppFilename(artifactName, version, ext), content, -1, opts)
	if err != nil {
		return "", err
	}

	log.Infof("Successfully uploaded %s of size %v\n", artifactName, n.Size)
	if IsDebugMode {
		log.Infof("File info: %+v", n)
	}

	return n.Location, nil
}

// StatFile reports whether the object exists and its stored content digest.
func (m *MinioWrapper) StatFile(ctx context.Context, artifactName string, artifactType string, version *semver.Version, ext string) (bool, string, error) {
	fPath := formatAppFilename(artifactName, version, ext)

	info, err := m.minioClient.StatObject(ctx, m.GetBucketName(artifactType), fPath, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, "", nil
		}
		return false, "", err
	}

	digest := info.Metadata.Get(contentDigestHeader)
	if digest == "" && info.UserMetadata != nil {
		digest = info.UserMetadata[ContentDigestMetaKey]
	}
	return true, digest, nil
}

func (m *MinioWrapper) GetFile(ctx context.Context, artifactName string, artifactType string, version *semver.Version, ext string) (reader io.ReadCloser, err error) {
	fPath := formatAppFilename(artifactName, version, ext)

	log.Infof("Downloading %s from bucket %s", fPath, m.GetBucketName(artifactType))
	o, err := m.minioClient.GetObject(ctx, m.GetBucketName(artifactType), fPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (m *MinioWrapper) ListVersions(ctx context.Context, artifactName string, artifactType string) (*[]AritfactInfo, error) {
	path := formatAppPath(artifactName) + "/"

	log.Infof("Listing objects in %s", path)
	objectCh := m.minioClient.ListObjects(ctx, m.GetBucketName(artifactType), minio.ListObjectsOptions{
		Prefix:       path,
		Recursive:    false,
		WithMetadata: false,
	})

	ls := map[string]AritfactInfo{}
	chunked := map[string]bool{}

	for object := range objectCh {
		if object.Err != nil {
			log.Errorf("Failed to list objects: %v", object.Err)

			return nil, object.Err
		}

		log.Infof("Listing object %s", object.Key)
		if strings.HasSuffix(object.Key, ChunkIndexExtension) {
			chunked[strings.TrimSuffix(object.Key, ChunkIndexExtension)] = true
		}

		if strings.HasSuffix(object.Key, TarGzExtension) {
			version := strings.TrimSuffix(strings.TrimPrefix(object.Key,
				formatAppPath(artifactName)+"/"), TarGzExtension)

			_, err := semver.NewVersion(version)
			if err != nil {
				log.Errorf("Failed to parse version %s: %v", version, err)
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

func (m *MinioWrapper) ListApps(ctx context.Context, artifactType string) ([]string, error) {
	log.Infof("Listing all objects")

	objectCh := m.minioClient.ListObjects(ctx, m.GetBucketName(artifactType), minio.ListObjectsOptions{
		Prefix:       appsRoot,
		Recursive:    false,
		WithMetadata: false,
	})

	ls := []string{}

	for object := range objectCh {
		if object.Err != nil {
			log.Errorf("Failed to list objects: %v", object.Err)
			return nil, object.Err
		}

		ls = append(ls, strings.TrimSuffix(strings.TrimPrefix(object.Key, appsRoot), "/"))
	}

	return ls, nil
}

func (m *MinioWrapper) GetEndpoint() string {
	return m.minioClient.EndpointURL().String() + "/" + appsRoot
}

// StoreBaseURL is the s3 base for an artifact type's source bucket. Only the scheme
// drives the distributor's read path (it resolves creds from its own config), but we
// return an accurate URL for clarity.
func (m *MinioWrapper) StoreBaseURL(artifactType string) string {
	return m.minioClient.EndpointURL().String() + "/" + m.GetBucketName(artifactType) + "/"
}

// ListLatestPerApp returns each app with its highest valid-semver version.
func (m *MinioWrapper) ListLatestPerApp(ctx context.Context, artifactType string) ([]AppLatest, error) {
	apps, err := m.ListApps(ctx, artifactType)
	if err != nil {
		return nil, err
	}

	result := make([]AppLatest, 0, len(apps))
	for _, name := range apps {
		if name == "" {
			continue
		}
		versions, err := m.ListVersions(ctx, name, artifactType)
		if err != nil {
			log.Errorf("Failed to list versions for %s: %v", name, err)
			continue
		}
		latest := pickLatest(*versions)
		if latest == nil {
			continue // no valid semver version for this app
		}
		result = append(result, AppLatest{Name: name, Latest: *latest})
	}
	return result, nil
}

// pickLatest returns the highest valid-semver version, skipping INVALID_VERSION_FORMAT entries.
func pickLatest(vs []AritfactInfo) *AritfactInfo {
	var best *AritfactInfo
	var bestV *semver.Version
	for i := range vs {
		v, err := semver.NewVersion(vs[i].Version)
		if err != nil {
			continue
		}
		if bestV == nil || v.GreaterThan(bestV) {
			bestV = v
			best = &vs[i]
		}
	}
	return best
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

func (m *MinioWrapper) createBucketIfMissing(bucketName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !strings.HasPrefix(bucketName, BucketNamePrefix) {
		return fmt.Errorf("bucket name should start with artifact_hub")
	}

	exists, err := m.minioClient.BucketExists(ctx, bucketName)
	if err == nil && exists {
		log.Infof("Bucket %s already exists", bucketName)

		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to check if bucket exists")
	}

	log.Infof("Bucket %s does not exist, creating it", bucketName)
	objLocking := !IsDebugMode

	err = m.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		ObjectLocking: objLocking,
	})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NotImplemented" {
			log.Errorf("Bucket creation is not supported by the server. Try enabling debug mode if minio is running in FS mode")

			return fmt.Errorf("bucket creation not supported by the server")
		} else {
			return err
		}
	}

	return nil
}
