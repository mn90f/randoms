//The goal of this program is to get a bunch of ISO currency codes stored in a databse
//And query them against a yahoo database
//To get all the exchnage rates for the curruncies we found from the DB VS a main currency code.

package main

import (
	sql "database/sql"
	json "encoding/json"
	_ "github.com/go-sql-driver/mysql"
	//ioutil "io/ioutil"
	bytes "bytes"
	fmt "fmt"
	log "log"
	http "net/http"
	//time "time"
)

type YahooCurrencyRateResultStruct struct {
	Id   interface{}
	Name interface{}
	Rate interface{}
	Date interface{}
	Time interface{}
	Ask  interface{}
	Bid  interface{}
}
type YahooCurrencyQueryResultsStruct struct {
	Rate []YahooCurrencyRateResultStruct
}
type YahooCurrencyQueryStruct struct {
	Count   interface{}
	Created interface{}
	Lang    interface{}
	Results YahooCurrencyQueryResultsStruct
}

type YahooCurrencyStruct struct {
	Query YahooCurrencyQueryStruct
}

func main() {
	//Db connection and sql connection
	dbUsername := "root"
	dbPassword := "passsssssssssssssssssssssssssssssswordddddddddddddddddddddddddddddddhaaaaaaaaaaaahaaaaaaaaaaaaa"
	dbIp := "10.0.3.33"
	dbPort := "3306"
	dbName := "UNLOCODE_TEST"

	db, dbErr := sql.Open("mysql", dbUsername+":"+dbPassword+"@tcp("+dbIp+":"+dbPort+")/"+dbName)
	defer db.Close()

	if dbErr != nil {
		log.Println(dbErr)
	}

	stmt, dbErr := db.Prepare("SELECT currency_code FROM tbl_currency")
	if dbErr != nil {
		log.Println(dbErr)
	}

	rows, dbErr := stmt.Query()
	if dbErr != nil {
		log.Println(dbErr)
	}

	//Prepare the query to submit to yahoo servers.
	convertToCurrencyCode := "USD"
	var currencyCodeAppendBuffer bytes.Buffer

	currencyCodeAppendBuffer.WriteString("https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.xchange%20where%20pair%20in%20(")

	var allCurruncies []string
	var eachCurrency string

	//Append all the curruncies, say we want USD to ALL the curruncies in the DB, it would produce something like USDLKR,USDINR... etc
	for rows.Next() {
		rows.Scan(&eachCurrency)
		allCurruncies = append(allCurruncies, eachCurrency)

		currencyCodeAppendBuffer.WriteString("\"")
		currencyCodeAppendBuffer.WriteString(eachCurrency)
		currencyCodeAppendBuffer.WriteString(convertToCurrencyCode)
		currencyCodeAppendBuffer.WriteString("\"")
		currencyCodeAppendBuffer.WriteString(",")
	}
	//cheating ;)
	currencyCodeAppendBuffer.WriteString("\"" + convertToCurrencyCode + convertToCurrencyCode + "\"")

	currencyCodeAppendBuffer.WriteString(")&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback=")

	//fmt.Println(currencyCodeAppendBuffer.String())
	fmt.Println("Sending the following request to yahoo server ...")
	fmt.Println("\n\n", currencyCodeAppendBuffer.String(), "\n\n\n")
	response, httpErr := http.Get(currencyCodeAppendBuffer.String())
	fmt.Println("Server has responded :)")

	if httpErr != nil {
		log.Println(httpErr)
	} else {
		defer response.Body.Close()

		var tmp YahooCurrencyStruct
		decorder := json.NewDecoder(response.Body)
		jsonErr := decorder.Decode(&tmp)
		if jsonErr != nil {
			log.Println(jsonErr)
		} else {

			//Basically at this point, we can use a TYPE ASSERTION if it passes we can then do an update to the database at this point.
			for _, element := range tmp.Query.Results.Rate {
				fmt.Println(element.Ask)
				fmt.Println(element.Bid)
				fmt.Println(element.Date)
				fmt.Println(element.Id)
				fmt.Println(element.Name)
				fmt.Println(element.Rate)
				fmt.Println(element.Time)
				fmt.Println("==============================\n")
			}
		}
	}
}
