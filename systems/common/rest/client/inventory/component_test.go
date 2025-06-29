/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package inventory

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/uuid"
)

func TestComponentClient_Get(t *testing.T) {
	// Mock data
	componentID, err := uuid.FromString("6006846a-c371-48c6-989c-9e029e144bb5")
	if err != nil {
		return
	}
	mockResponse := Component{
		ComponentInfo: &ComponentInfo{
			Id:            componentID,
			Inventory:     "200",
			UserId:        "bc082789-bef7-4baf-9cd1-d479fdb3184b",
			Category:      "BACKHAUL",
			Type:          "backhaul",
			Description:   "A 100uF capacitor",
			DatasheetURL:  "http://example.com/datasheet2",
			ImagesURL:     "http://example.com/image2",
			PartNumber:    "C-100uF",
			Manufacturer:  "Capacitors Ltd.",
			Managed:       "no",
			Warranty:      uint32(12),
			Specification: "100uF, 16V",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(mockResponse)
		if err != nil {
			t.Fatalf("Failed to encode mock response: %v", err)
		}
	}))
	defer ts.Close()

	client := NewComponentClient(ts.URL)

	result, err := client.Get(componentID.String())

	assert.NoError(t, err)

	assert.NotNil(t, result)
	assert.Equal(t, mockResponse.ComponentInfo.Id, result.Id)
	assert.Equal(t, mockResponse.ComponentInfo.Inventory, result.Inventory)
	assert.Equal(t, mockResponse.ComponentInfo.UserId, result.UserId)
	assert.Equal(t, mockResponse.ComponentInfo.Category, result.Category)
	assert.Equal(t, mockResponse.ComponentInfo.Type, result.Type)
	assert.Equal(t, mockResponse.ComponentInfo.Description, result.Description)
	assert.Equal(t, mockResponse.ComponentInfo.DatasheetURL, result.DatasheetURL)
	assert.Equal(t, mockResponse.ComponentInfo.ImagesURL, result.ImagesURL)
	assert.Equal(t, mockResponse.ComponentInfo.PartNumber, result.PartNumber)
	assert.Equal(t, mockResponse.ComponentInfo.Manufacturer, result.Manufacturer)
	assert.Equal(t, mockResponse.ComponentInfo.Managed, result.Managed)
	assert.Equal(t, mockResponse.ComponentInfo.Warranty, result.Warranty)
	assert.Equal(t, mockResponse.ComponentInfo.Specification, result.Specification)
}
