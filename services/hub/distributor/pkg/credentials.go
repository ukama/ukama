package pkg

import (
	"fmt"

	"github.com/minio/minio-go/v6/pkg/credentials"
	mc "github.com/ukama/ukama/services/hub/hub/pkg"
)

type StoreCredentialsOptions struct {
	mc.MinioConfig
}

var store *StoreCredentialsOptions

// StaticCredentialsProvider implements credentials.Provider from github.com/minio/minio-go/pkg/credentials
type S3CredentialsProvider struct {
	creds credentials.Value
}

/* IsExpired returns true when the credentials are expired*/
func (cp *S3CredentialsProvider) IsExpired() bool {
	return false
}

/* Retrieve returns credentials */
func (cp *S3CredentialsProvider) Retrieve() (credentials.Value, error) {
	return cp.creds, nil
}

func InitStoreCredentialsOptions(c *MinioConfig) {
	store = &StoreCredentialsOptions{
		mc.MinioConfig{
			TimeoutSecond:      c.TimeoutSecond,
			Endpoint:           c.Endpoint,
			AccessKey:          c.AccessKey,
			SecretKey:          c.SecretKey,
			BucketSuffix:       c.BucketSuffix,
			Region:             c.Region,
			SkipBucketCreation: c.SkipBucketCreation,
		},
	}
}

/* NewS3Credentials initializes a new set of S3 credentials */
func NewS3Credentials(accessKey, secretKey string) *credentials.Credentials {
	p := &S3CredentialsProvider{
		credentials.Value{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
		},
	}
	return credentials.New(p)
}

/* Get S3 credentilas */
func GetS3CredentialsFor(fstore string) (*credentials.Credentials, *string, error) {

	/* Get store config */
	if store == nil {
		return nil, nil, fmt.Errorf("no config for artifact store found")
	}
	region := &store.Region

	return NewS3Credentials(store.AccessKey, store.SecretKey), region, nil

}

/* Get credentials for local store */
func GetLocalStoreCredentialsFor(fstore string) (*StoreCredentialsOptions, error) {
	/* Get store config */
	if store == nil {
		return nil, fmt.Errorf("no config for artifact store found")
	}
	return store, nil
}
