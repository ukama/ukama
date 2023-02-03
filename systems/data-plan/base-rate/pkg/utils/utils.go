package utils

import (
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/jszwec/csvutil"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"
)

type RawRates struct {
	Country string `csv:"Country"`
	Network string `csv:"Network"`
	Vpmn    string `csv:"VPMN"`
	Imsi    string `csv:"IMSI"`
	Sms_mo  string `csv:"SMS MO"`
	Sms_mt  string `csv:"SMS MT"`
	Data    string `csv:"Data"`
	X2g     string `csv:"2G"`
	X3g     string `csv:"3G"`
	X5g     string `csv:"5G"`
	Lte     string `csv:"LTE"`
	Lte_m   string `csv:"LTE-M"`
	Apn     string `csv:"APN"`
}

func FetchData(url string) ([]RawRates, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, _ := io.ReadAll(resp.Body)

	var r []RawRates
	errorStr := "invalid CSV file data"
	err = csvutil.Unmarshal(content, &r)
	if err != nil {
		return nil, err
	}

	if len(r) == 0 || validations.IsEmpty(r[0].Country) ||
		validations.IsEmpty(r[0].Network) ||
		validations.IsEmpty(r[0].Data) {
		return nil, errors.New(errorStr)
	}

	return r, nil
}

func ParseToModel(slice []RawRates, effective_at, sim_type string) []db.Rate {
	var rates []db.Rate
	for _, value := range slice {
		rates = append(rates, db.Rate{
			Uuid:         uuid.New(),
			Country:      value.Country,
			Network:      value.Network,
			Vpmn:         value.Vpmn,
			Imsi:         value.Imsi,
			Sms_mo:       value.Sms_mo,
			Sms_mt:       value.Sms_mt,
			Data:         value.Data,
			X2g:          value.X2g,
			X3g:          value.X3g,
			X5g:          value.X5g,
			Lte:          value.Lte,
			Lte_m:        value.Lte_m,
			Apn:          value.Apn,
			Effective_at: effective_at,
			End_at:       "",
			Sim_type:     sim_type,
		})
	}
	return rates
}
