package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func DeleteFile(fileName string) {
	e := os.Remove(fileName)
	Check(e)
}

func FetchData(url string, destinationFileName string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	f, e := os.Create(destinationFileName)
	if e != nil {
		panic(e)
	}

	defer f.Close()

	f.ReadFrom(resp.Body)
}

func trimHeader(columnName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(columnName),
			" ", "_"), "-", "_"), "2g", "X2g"), "3g", "X3g"), "5g", "X5g")
}

func CreateQuery(rows [][]string, effective_at string, sim_type string) string {
	headerStr := ""
	valueStrings := make([]string, 0, len(rows))

	for _, value := range rows[0] {
		headerStr = headerStr + trimHeader(value) + ","
	}
	headerStr = "(" + headerStr + "effective_at, sim_type)"

	for i, row := range rows {
		if i == 0 {
			continue
		}

		str := ""
		for _, value := range row {
			str = str + "'" + strings.ReplaceAll(strings.ReplaceAll(value, "'", ""), ",", "") + "', "
		}
		str = str + " '" + effective_at + "', '" + sim_type + "'"
		valueStrings = append(valueStrings, "("+strings.ReplaceAll(str, "''", "NULL")+")")
	}
	stmt := fmt.Sprintf("INSERT INTO rates %s VALUES %s", headerStr, strings.Join(valueStrings, ","))
	return stmt
}
