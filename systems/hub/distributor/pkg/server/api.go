/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

type ChunkRequest struct {
	Name    string `path:"name" validate:"required"`
	Version string `path:"version" validate:"required"`
	Store   string `json:"store" validate:"required"`
}

type ChunkResponse struct {
}

type CApp struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
