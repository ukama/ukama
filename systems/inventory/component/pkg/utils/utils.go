package utils

type Component struct {
	Category      string `json:"category"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	UserId        string `json:"ownerId"`
	ImagesURL     string `json:"imagesURL" yaml:"imagesURL"`
	DatasheetURL  string `json:"datasheetURL" yaml:"datasheetURL"`
	InventoryID   string `json:"inventoryID" yaml:"inventoryID"`
	PartNumber    string `json:"partNumber" yaml:"partNumber"`
	Manufacturer  string `json:"manufacturer"`
	Managed       string `json:"managed"`
	Warranty      uint32 `json:"warranty"`
	Specification string `json:"specification" yaml:"specification"`
}
