package utils

import (
	"io"
	"net/http"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
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

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawRates []RawRates
	if err := csvutil.Unmarshal(content, &rawRates); err != nil {
		return nil, err
	}

	return rawRates, nil
}

func ParseToModel(slice []RawRates, effective_at, sim_type string) []db.Rate {
	var rates []db.Rate
	for _, value := range slice {
		rates = append(rates, db.Rate{
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

func IsFutureDate(date string) bool {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return false
	}
	return time.Now().Before(t)
}
