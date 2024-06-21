/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

package server

import (
	"time"
)

type ArtifactUploadRequest struct {
	ArtifactName string `path:"name" validate:"required"`
	ArtifactType string `path:"type" validate:"required"`
	Version      string `path:"version" validate:"required"`
}

type ArtifactRequest struct {
	Name         string `path:"name" validate:"required"`
	FileName     string `path:"filename" validate:"required"`
	ArtifactType string `path:"type" validate:"required"`
}

type ArtifactVersionListRequest struct {
	Name         string `path:"name" validate:"required"`
	ArtifactType string `path:"type" validate:"required"`
}

type ArtifactVersionListResponse struct {
	Name         string         `json:"name"`
	ArtifactType string         `path:"type" validate:"required"`
	Versions     *[]VersionInfo `json:"artifacts"`
}

type VersionInfo struct {
	Version string       `json:"version"`
	Formats []FormatInfo `json:"formats"`
}

type FormatInfo struct {
	Type      string            `json:"type"`
	Url       string            `json:"url"`
	CreatedAt time.Time         `json:"created_at"`
	SizeBytes int64             `json:"size_bytes,omitempty"`
	ExtraInfo map[string]string `json:"extra_info,omitempty"`
}

type ArtifactLocationRequest struct {
	Name         string `path:"name" validate:"required"`
	ArtifactType string `path:"type" validate:"required"`
	Version      string `path:"version" validate:"required"`
}

type ArtifactLocationResponse struct {
	Endpoint string `json:"endpoint"`
}

type ArtifactListRequest struct {
	ArtifactType string `path:"type" validate:"required"`
}

type ArtifactListResponse struct {
	Names []string `json:"names" validate:"required"`
}
