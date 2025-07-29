/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package utils

type Accounting struct {
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
