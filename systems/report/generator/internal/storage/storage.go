/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package storage

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	log "github.com/sirupsen/logrus"
)

const pdfContentType = "application/pdf"

type Storage interface {
	Upload(ctx context.Context, objectName, filePath string) (string, error)
}

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucket, region string, secure bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	s := &MinioStorage{
		client: client,
		bucket: bucket,
	}

	if err := s.createBucketIfMissing(); err != nil {
		return nil, err
	}

	return s, nil
}

func (m *MinioStorage) createBucketIfMissing() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := m.client.BucketExists(ctx, m.bucket)
	if err != nil {
		return err
	}

	if exists {
		log.Infof("Bucket %s already exists", m.bucket)

		return nil
	}

	log.Infof("Bucket %s does not exist, creating it", m.bucket)

	return m.client.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{})
}

func (m *MinioStorage) Upload(ctx context.Context, objectName, filePath string) (string, error) {
	log.Infof("Uploading %s to bucket %s", objectName, m.bucket)

	info, err := m.client.FPutObject(ctx, m.bucket, objectName, filePath,
		minio.PutObjectOptions{ContentType: pdfContentType})
	if err != nil {
		return "", err
	}

	log.Infof("Successfully uploaded %s of size %d", objectName, info.Size)

	return info.Location, nil
}
