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
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage interface {
	Get(ctx context.Context, bucket, objectName string) ([]byte, error)
}

type MinioStorage struct {
	client *minio.Client
}

func NewMinioStorage(endpoint, accessKey, secretKey, region string, secure bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	return &MinioStorage{client: client}, nil
}

func (m *MinioStorage) Get(ctx context.Context, bucket, objectName string) ([]byte, error) {
	obj, err := m.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	defer obj.Close()

	return io.ReadAll(obj)
}
