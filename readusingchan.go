package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type Industry struct {
	CompanyName                string `csv:"Organization Name"`
	LinkedIn                   string `csv:"LinkedIn"`
	Website                    string `csv:"Website"`
	TotalFundingAmount         string `csv:"Total Funding Amount"`
	TotalFundingAmountCurrency string `csv:"Total Funding Amount Currency"`
	HeadquartersLocation       string `csv:"Headquarters Location"`
}

// 5.57ms -> 600 records (read)
func main() {
	now := time.Now()
	readChannel := make(chan Industry, 1)

	readFilePath := "process.csv"

	// Open the CSV readFile
	readFile, err := os.OpenFile(readFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer readFile.Close()

	count := 0
	readFromCSV(readFile, readChannel)

	// Print the records
	for r := range readChannel {
		fmt.Println("========================================")
		fmt.Println(r)
		fmt.Println("========================================")
		fmt.Println()

		count++
	}

	fmt.Println(time.Since(now), count)
}

func readFromCSV(file *os.File, c chan Industry) {
	gocsv.SetCSVReader(func(r io.Reader) gocsv.CSVReader {
		reader := csv.NewReader(r)
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1
		return reader
	})

	// Read the CSV file into a slice of Record structs
	go func() {
		err := gocsv.UnmarshalToChan(file, c)
		if err != nil {
			panic(err)
		}
	}()
}
