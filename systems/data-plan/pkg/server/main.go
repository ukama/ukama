package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234-password"
	dbname   = "ukama_rates"
  )

func insertDataInDB(query string){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_,error := db.Query(query)
	if error != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to DB!")
}

func BulkInsert(rows [][]string, effective_at string) string  {
	headerStr := ""
    valueStrings := make([]string, 0, len(rows))
    for i, row := range rows {
		if i == 0 {
			headerStr = "(" + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
				strings.ReplaceAll(strings.ToLower(strings.Join(row[:], ",")), " ", "_"),
				 "-", "_"), "2g", "_2g"), "3g", "_3g") + ",effective_at,end_at" + ")"
			continue
		}
		values := row
		str := ""
		for j, value := range values {
			if j == len(values) - 1 {
				str = str + "'" + strings.ReplaceAll(strings.ReplaceAll(value, "'", ""), ",", " ") + "'"
			} else {
				str = str + "'" + strings.ReplaceAll(strings.ReplaceAll(value, "'", ""), ",", "") + "', "
			}
		}
		str = str + ", '" + effective_at + "', ''"
		valueStrings = append(valueStrings, "(" + strings.ReplaceAll(str, "''", "NULL") + ")")
    }
	stmt := fmt.Sprintf("INSERT INTO rates %s VALUES %s", headerStr, strings.Join(valueStrings, ","))
	return stmt
}

func main(){
	f, err := os.Open("rates.csv")
	if err != nil {
        panic(err)
    }

	defer f.Close()

	csvReader := csv.NewReader(f)
    data, err := csvReader.ReadAll()

    if err != nil {
        log.Fatal(err)
    }

	query := BulkInsert(data, "2022-11-01 00:00:00")
	insertDataInDB(query)
}