package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file := "problems.csv"
	f, err := openCSV(file)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	readCSV(f)
}

func openCSV(file string) (*os.File, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func readCSV(f *os.File) {
	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", rec)
	}
}
