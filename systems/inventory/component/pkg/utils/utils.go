/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package utils

type Component struct {
	Category      string `json:"category" yaml:"category"`
	Type          string `json:"type" yaml:"type"`
	Description   string `json:"description" yaml:"description"`
	UserId        string `json:"ownerId" yaml:"ownerId"`
	ImagesURL     string `json:"imagesURL" yaml:"imagesURL"`
	DatasheetURL  string `json:"datasheetURL" yaml:"datasheetURL"`
	InventoryID   string `json:"inventoryID" yaml:"inventoryID"`
	PartNumber    string `json:"partNumber" yaml:"partNumber"`
	Manufacturer  string `json:"manufacturer" yaml:"manufacturer"`
	Managed       string `json:"managed" yaml:"managed"`
	Warranty      uint32 `json:"warranty" yaml:"warranty"`
	Specification string `json:"specification" yaml:"specification"`
}
