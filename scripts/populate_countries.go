package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vonmutinda/organono/app/logger"
)

type Country struct {
	Code         string `json:"code"`
	Currency     string `json:"currency"`
	Name         string `json:"name"`
	DiallingCode string `json:"dialling_code"`
}

func main() {

	var envFilePath string

	flag.StringVar(&envFilePath, "e", "", "Path to .env file")
	flag.Parse()

	if envFilePath != "" {
		err := godotenv.Load(envFilePath)
		if err != nil {
			logger.Fatalf("Failed to load env file err = %v", err)
		}
	}

	populateCountries()
}

func populateCountries() {

	var countries []*Country

	f, err := os.Open("./scripts/countries.json")
	checkErr(err)

	defer f.Close()

	err = json.NewDecoder(f).Decode(&countries)
	checkErr(err)

	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})

	dB, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	checkErr(err)

	defer dB.Close()

	err = dB.Ping()
	checkErr(err)

	tx, err := dB.Begin()
	checkErr(err)

	query := "INSERT INTO countries (code, currency, name, dialling_code) VALUES ($1, $2, $3, $4)"

	for _, country := range countries {
		fmt.Printf("Inserting %v\n", country.Name)

		_, err := tx.Exec(query, country.Code, country.Currency, country.Name, country.DiallingCode)
		if err != nil {
			err = tx.Rollback()
			checkErr(err)
			return
		}
	}

	err = tx.Commit()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
