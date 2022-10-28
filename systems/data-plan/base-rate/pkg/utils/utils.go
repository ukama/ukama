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
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		strings.ReplaceAll(strings.ToLower(columnName),
			" ", "_"), "-", "_"), "2g", "x2g"), "3g", "x3g"), "5g", "x5g")
}

func CreateQuery(rows [][]string, effectiveAt string, simType string) string {
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
		str = str + " '" + effectiveAt + "', '" + simType + "'"
		valueStrings = append(valueStrings, "("+strings.ReplaceAll(str, "''", "NULL")+")")
	}
	stmt := fmt.Sprintf("INSERT INTO rates %s VALUES %s", headerStr, strings.Join(valueStrings, ","))
	return stmt
}
