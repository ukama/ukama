package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	uuid "github.com/ukama/ukama/systems/common/uuid"

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

func ParseToModel(slice []RawRates, effective_at, sim_type string) ([]db.Rate, error) {
	var rates []db.Rate
	for _, value := range slice {

		imsi, err := ParsedToInt(value.Imsi)
		if err != nil {
			return nil, fmt.Errorf("failed parsing imsi value." + err.Error())
		}

		smo, err := ParseToRates(value.Sms_mo, "$")
		if err != nil {
			return nil, fmt.Errorf("failed parsing SMS MO rate." + err.Error())
		}

		smt, err := ParseToRates(value.Sms_mt, "$")
		if err != nil {
			return nil, fmt.Errorf("failed parsing SMS MT rate." + err.Error())
		}

		data, err := ParseToRates(value.Data, "$")
		if err != nil {
			return nil, fmt.Errorf("failed parsing Data rate." + err.Error())
		}

		rates = append(rates, db.Rate{
			Uuid:        uuid.NewV4(),
			Country:     value.Country,
			Network:     value.Network,
			Vpmn:        value.Vpmn,
			Imsi:        imsi,
			SmsMo:       smo,
			SmsMt:       smt,
			Data:        data,
			X2g:         ParseToBoolean(value.X2g, "2G"),
			X3g:         ParseToBoolean(value.X3g, "3G"),
			X5g:         ParseToBoolean(value.X5g, "5G"),
			Lte:         ParseToBoolean(value.Lte, "LTE"),
			LteM:        ParseToBoolean(value.Lte_m, "LTE-M"),
			Apn:         value.Apn,
			EffectiveAt: effective_at,
			EndAt:       "",
			SimType:     db.ParseType(sim_type),
		})
	}
	return rates, nil
}

func ParseToRates(str string, s string) (float64, error) {
	var val float64 = 0
	var err error
	sr := strings.Split(str, s)
	for _, s := range sr {
		val, err = strconv.ParseFloat(s, 64)
	}
	if err != nil {
		log.Errorf("Failed to parse rate from %s with rate symbol %s.Error: %v", str, s, err)
		return 0, err
	}
	return val, nil
}

func ParsedToInt(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Errorf("Failed to parse int64 value from  %s.Error: %v", s, err)
		return 0, err
	}
	return val, nil
}

func ParseToBoolean(val string, s string) bool {
	return (val == s)
}
