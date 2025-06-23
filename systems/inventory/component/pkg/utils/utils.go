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
