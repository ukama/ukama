/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	ukama "github.com/ukama/ukama/systems/common/ukama"
)

func TestRateService_ParseToModel(t *testing.T) {
	rawRates := []RawRates{{
		Country: "Tycho crater",
		Network: "Multi Tel",
		Apn:     "Manual entry required",
		Imsi:    "1",
		Vpmn:    "TTC",
		X2g:     "2G",
		X3g:     "3G",
		Lte:     "LTE",
		Data:    "$0.4",
		Sms_mo:  "$0.1",
		Sms_mt:  "$0.1",
	}}
	dbRate, err := ParseToModel(rawRates, "2023-04-10T20:05:29Z", "2024-04-10T20:05:29Z", "ukama_data")
	assert.NoError(t, err)
	assert.Equal(t, rawRates[0].Country, dbRate[0].Country)
	assert.Equal(t, "2023-04-10T20:05:29Z", dbRate[0].EffectiveAt.Format(time.RFC3339Nano))
	assert.Equal(t, ukama.ParseSimType("ukama_data"), dbRate[0].SimType)
}

func TestParseToModel_ErrorCases(t *testing.T) {
	effectiveAt := "2023-04-10T20:05:29Z"
	endAt := "2024-04-10T20:05:29Z"
	simType := "ukama_data"

	tests := []struct {
		name      string
		rawRates  []RawRates
		effAt     string
		endAt     string
		wantError string
	}{
		{
			name: "IMSI parsing error",
			rawRates: []RawRates{{
				Country: "Country",
				Network: "Net",
				Imsi:    "notanumber",
				Sms_mo:  "$0.1",
				Sms_mt:  "$0.1",
				Data:    "$0.1",
			}},
			effAt:     effectiveAt,
			endAt:     endAt,
			wantError: "failed parsing imsi value",
		},
		{
			name: "SMS MO parsing error",
			rawRates: []RawRates{{
				Country: "Country",
				Network: "Net",
				Imsi:    "1",
				Sms_mo:  "notanumber",
				Sms_mt:  "$0.1",
				Data:    "$0.1",
			}},
			effAt:     effectiveAt,
			endAt:     endAt,
			wantError: "failed parsing SMS MO rate",
		},
		{
			name: "SMS MT parsing error",
			rawRates: []RawRates{{
				Country: "Country",
				Network: "Net",
				Imsi:    "1",
				Sms_mo:  "$0.1",
				Sms_mt:  "notanumber",
				Data:    "$0.1",
			}},
			effAt:     effectiveAt,
			endAt:     endAt,
			wantError: "failed parsing SMS MT rate",
		},
		{
			name: "Data parsing error",
			rawRates: []RawRates{{
				Country: "Country",
				Network: "Net",
				Imsi:    "1",
				Sms_mo:  "$0.1",
				Sms_mt:  "$0.1",
				Data:    "notanumber",
			}},
			effAt:     effectiveAt,
			endAt:     endAt,
			wantError: "failed parsing Data rate",
		},
		{
			name: "Invalid effective_at date format",
			rawRates: []RawRates{{
				Country: "Country",
				Network: "Net",
				Imsi:    "1",
				Sms_mo:  "$0.1",
				Sms_mt:  "$0.1",
				Data:    "$0.1",
			}},
			effAt:     "notadate",
			endAt:     endAt,
			wantError: "invalid time format for effective at",
		},
		{
			name: "Invalid endAt date format",
			rawRates: []RawRates{{
				Country: "Country",
				Network: "Net",
				Imsi:    "1",
				Sms_mo:  "$0.1",
				Sms_mt:  "$0.1",
				Data:    "$0.1",
			}},
			effAt:     effectiveAt,
			endAt:     "notadate",
			wantError: "invalid time format for end at",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseToModel(tt.rawRates, tt.effAt, tt.endAt, simType)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantError)
		})
	}
}

// Fetch data success case
func TestRateService_FetchData_Success(t *testing.T) {
	mockFileUrl := "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv"

	rawRates, err := FetchData(mockFileUrl)
	assert.NoError(t, err)
	assert.Equal(t, "The lunar maria", rawRates[0].Country)
}

// Fetch data error case invalid file url
func TestRateService_FetchData_error1(t *testing.T) {
	mockFileUrl := "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv"

	rateError1, err := FetchData("/fail" + mockFileUrl)
	assert.Error(t, err)
	assert.Nil(t, rateError1)
}

// Fetch data error case file with invalid data
func TestRateService_FetchData_error2(t *testing.T) {
	failMockFileUrl := "https://raw.githubusercontent.com/ukama/ukama/baserate-test/systems/data-plan/docs/template/failed_template.csv"

	rateError2, err := FetchData(failMockFileUrl)
	assert.Error(t, err)
	assert.Nil(t, rateError2)
}
