/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type GetNodeRequest struct {
	NodeId string `example:"{{NodeId}}" path:"nodeId" validate:"required"`
}

type ReflectorPingRequest struct {}

type ReflectorGetRequest struct {
	NodeId string `path:"nodeId" validate:"required"`
}

type ReflectorDownloadRequest struct {
	NodeId       string `json:"nodeId" body:"nodeId" validate:"required"`
	Bytes        int64  `json:"bytes" body:"bytes" validate:"required,gt=0"`
	ChunkBytes   int32  `json:"chunkBytes,omitempty" body:"chunkBytes"`
	ChunkDelayMs int32  `json:"chunkDelayMs,omitempty" body:"chunkDelayMs"`
}

type ReflectorUploadRequest struct {
	NodeId  string `json:"nodeId" body:"nodeId" validate:"required"`
	Payload []byte `json:"payload" body:"payload" validate:"required"`
}