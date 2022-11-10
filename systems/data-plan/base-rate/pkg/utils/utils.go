package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
)

func DeleteFile(fileName string) error {
	e := os.Remove(fileName)
	if e != nil {
		return e
	}
	return nil
}

func FetchData(url string, destinationFileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, e := os.Create(destinationFileName)
	if e != nil {
		return e
	}

	defer f.Close()

	f.ReadFrom(resp.Body)

	return nil
}

func trimHeader(columnName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		strings.ReplaceAll(strings.ToLower(columnName),
			" ", "_"), "-", "_"), "2g", "x2g"), "3g", "x3g"), "5g", "x5g")
}

func CreateData(rows [][]string, effective_at string, sim_type string) []map[string]interface{} {
	headerStr := make([]string, 0, len(rows[0]))

	for _, value := range rows[0] {
		headerStr = append(headerStr, trimHeader(value))
	}

	var query []map[string]interface{}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		_tempRate := make(map[string]interface{})
		for j, value := range row {
			if value == "" {
				_tempRate[headerStr[j]] = ""
				continue
			}
			_tempRate[headerStr[j]] = strings.ReplaceAll(strings.ReplaceAll(value, "'", ""), ",", "")
		}
		_tempRate["effective_at"] = effective_at
		_tempRate["sim_type"] = sim_type
		query = append(query, _tempRate)
	}
	return query
}

func ParseToModel(slice []map[string]interface{}) []db.Rate {
	fmt.Println(slice)
	var rates []db.Rate
	for _, value := range slice {
		rates = append(rates, db.Rate{
			Country:      value["country"].(string),
			Network:      value["network"].(string),
			Vpmn:         value["vpmn"].(string),
			Imsi:         value["imsi"].(string),
			Sms_mo:       value["sms_mo"].(string),
			Sms_mt:       value["sms_mt"].(string),
			Data:         value["data"].(string),
			X2g:          value["x2g"].(string),
			X3g:          value["x3g"].(string),
			X5g:          value["x5g"].(string),
			Lte:          value["lte"].(string),
			Lte_m:        value["lte_m"].(string),
			Apn:          value["apn"].(string),
			Effective_at: value["effective_at"].(string),
			End_at:       "",
			Sim_type:     value["sim_type"].(string),
		})
	}
	return rates
}

func IsFutureDate(date string) bool {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return false
	}
	today := time.Now().UnixNano() / int64(time.Millisecond)
	return today < t.UnixNano()/int64(time.Millisecond)
}
