package utils

import (
	pb "github.com/ukama/ukama/systems/inventory/account/pb/gen"
)

type Account struct {
	EffectiveDate        string               `json:"effective_date"`
	ConnectivityProvider ConnectivityProvider `json:"connectivityProvider"`
	Nodes                Nodes                `json:"nodes"`
	Ukama                []Item               `json:"ukama"`
	Backhaul             []Item               `json:"backhaul"`
}

type ConnectivityProvider struct {
	Company string `json:"company"`
	Poc     string `json:"poc"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type Nodes struct {
	Inventory string `json:"inventory"`
	OnOrder   string `json:"onOrder"`
}

type Item struct {
	Item          string `json:"item"`
	Description   string `json:"description"`
	Inventory     string `json:"inventory"`
	OpexFee       string `json:"opex_fee"`
	Vat           string `json:"vat"`
	EffectiveDate string `json:"effective_date"`
}

func UniqueInventoryIds(accounts []*pb.Account) []string {
	uniqueIds := make(map[string]bool)
	for _, account := range accounts {
		if _, exists := uniqueIds[account.Inventory]; !exists {
			uniqueIds[account.Inventory] = true
		}
	}

	ids := make([]string, 0, len(uniqueIds))
	for id := range uniqueIds {
		ids = append(ids, id)
	}

	return ids
}
