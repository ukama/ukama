package utils

import (
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
)

type Component struct {
	Company       string `json:"company"`
	Category      string `json:"category"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	ImagesURL     string `json:"imagesURL" yaml:"imagesURL"`
	DatasheetURL  string `json:"datasheetURL" yaml:"datasheetURL"`
	InventoryID   string `json:"inventoryID" yaml:"inventoryID"`
	PartNumber    string `json:"partNumber" yaml:"partNumber"`
	Manufacturer  string `json:"manufacturer"`
	Managed       string `json:"managed"`
	Warranty      uint32 `json:"warranty"`
	Specification string `json:"specification" yaml:"specification"`
}

func UniqueComponentIds(components []*pb.Component) []string {
	uniqueIds := make(map[string]bool)
	for _, component := range components {
		if _, exists := uniqueIds[component.Inventory]; !exists {
			uniqueIds[component.Inventory] = true
		}
	}

	ids := make([]string, 0, len(uniqueIds))
	for id := range uniqueIds {
		ids = append(ids, id)
	}

	return ids
}
