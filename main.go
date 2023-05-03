package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	defaultFilePath := "problems.csv"
	filePath := flag.String("file", defaultFilePath, "Problem file path")
	flag.Parse()
	f, err := openCSV(*filePath)
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
	scanner := bufio.NewScanner(os.Stdin)
	var score, total int
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		total++
		problem, solution := rec[0], rec[1]
		fmt.Printf("Problem #%d %s? ", total, problem)
		// read user answer
		scanner.Scan()
		err = scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		answer := scanner.Text()

		if err != nil {
			log.Fatal(err)
		}
		if answer == solution {
			score++
		}
	}
	fmt.Println()
	fmt.Printf("You scored: %d out of %d!\n", score, total)
	fmt.Println()
}
