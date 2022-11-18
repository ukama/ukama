package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	dbRate := ParseToModel(rawRates, "2023-10-10", "inter_mno_data")
	assert.Equal(t, rawRates[0].Country, dbRate[0].Country)
	assert.Equal(t, "2023-10-10", dbRate[0].Effective_at)
	assert.Equal(t, "inter_mno_data", dbRate[0].Sim_type)
}

func TestRateService_IsFutureDate(t *testing.T) {
	//Success case
	assert.True(t, IsFutureDate("2023-10-10T00:00:00Z"))
	//Error case
	assert.False(t, IsFutureDate("2023-10-10"))
	assert.False(t, IsFutureDate("2020-10-10T00:00:00Z"))
}

func TestRateService_FetchData(t *testing.T) {
	mockFileUrl := "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv"
	//Success case
	rawRates, err := FetchData(mockFileUrl)
	assert.NoError(t, err)
	assert.Equal(t, "The lunar maria", rawRates[0].Country)

	//Error case
	rateError1, err := FetchData("/fail" + mockFileUrl)
	assert.Error(t, err)
	assert.Nil(t, rateError1)

	rateError2, err := FetchData("https://raw.githubusercontent.com/ukama/ukama/baserate-test/systems/data-plan/docs/template/failed_template.csv")
	assert.Error(t, err)
	assert.Nil(t, rateError2)
}
