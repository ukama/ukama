package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"
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

func CreateQuery(rows [][]string, effectiveAt string, simType pb.SimType) string {
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
		fmt.Println("PBSIM", simType)
		fmt.Println("SIM STR", validations.ReqPbToStr(simType))
		str = str + " '" + effectiveAt + "', '" + validations.ReqPbToStr(simType) + "'"
		valueStrings = append(valueStrings, "("+strings.ReplaceAll(str, "''", "NULL")+")")
	}
	stmt := fmt.Sprintf("INSERT INTO rates %s VALUES %s", headerStr, strings.Join(valueStrings, ","))
	return stmt
}
