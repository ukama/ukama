/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

package server

type ArtifactUploadRequest struct {
	ArtifactName string `path:"name" validate:"required"`
	Version      string `path:"version" validate:"required"`
	ArtifactType string `path:"type" validate:"eq=app|eq=cert|eq=config"`
}

type ArtifactRequest struct {
	Name         string `path:"name" validate:"required"`
	FileName     string `path:"filename" validate:"required"`
	ArtifactType string `path:"type" validate:"eq=app|eq=cert|eq=config"`
}

type ArtifactVersionListRequest struct {
	Name         string `path:"name" validate:"required"`
	ArtifactType string `path:"type" validate:"eq=app|eq=cert|eq=config"`
}

type ArtifactListRequest struct {
	ArtifactType string `path:"type" validate:"required"`
}
