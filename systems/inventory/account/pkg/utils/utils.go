package utils

type Account struct {
	Item          string  `json:"item"`
	Description   string  `json:"description"`
	Company       string  `json:"company"`
	Inventory     string  `json:"inventory"`
	OpexFee       float64 `json:"opex_fee" yaml:"opex_fee"`
	Vat           float64 `json:"vat" yaml:"vat"`
	EffectiveDate string  `json:"effective_date" yaml:"effective_date"`
}
