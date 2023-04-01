package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
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
	dbRate, err := ParseToModel(rawRates, "2023-10-10", "ukama_data")
	assert.NoError(t, err)
	assert.Equal(t, rawRates[0].Country, dbRate[0].Country)
	assert.Equal(t, "2023-10-10", dbRate[0].EffectiveAt)
	assert.Equal(t, db.ParseType("ukama_data"), dbRate[0].SimType)
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
