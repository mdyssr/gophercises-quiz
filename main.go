package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var score, total int

func main() {
	filePath := flag.String("file", "problems.csv", "Problem file path")
	timeDelay := flag.Int("timer", 30, "Quiz timer (in seconds)")
	shuffle := flag.Bool("shuffle", false, "Shuffle the problems")
	flag.Parse()

	if *timeDelay < 1 {
		fmt.Println("Oopsie! Timer value is lower than 1!")
		os.Exit(0)
	}

	timerCh := make(chan struct{})
	quizCh := make(chan struct{})

	go quiz(*filePath, *shuffle, quizCh)
	go timer(*timeDelay, timerCh)

	select {
	case <-timerCh:
		fmt.Printf("\n\nTime out!\n")
		quizOver()
	case <-quizCh:
		fmt.Printf("\nQuiz completed :)\n")
		quizOver()
	}

}

func quizOver() {
	fmt.Printf("You scored: %d out of %d!\n\n", score, total)
}

func openCSV(file string) (*os.File, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func timer(delay int, ch chan struct{}) {
	time.Sleep(time.Duration(delay) * time.Second)
	ch <- struct{}{}
}

func quiz(filePath string, shuffle bool, done chan struct{}) {
	f, err := openCSV(filePath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(f)
	scanner := bufio.NewScanner(os.Stdin)
	records, err := csvReader.ReadAll()

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}

	if err != nil {
		log.Fatal(err)
	}
	total = len(records)

	for _, rec := range records {
		problem, solution := rec[0], rec[1]
		fmt.Printf("Problem #%d %s? ", total, problem)
		// read user answer
		scanner.Scan()
		err = scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		answer := scanner.Text()
		cleanAnswer := strings.ToLower(strings.TrimSpace(answer))
		if err != nil {
			log.Fatal(err)
		}
		if cleanAnswer == solution {
			score++
		}
	}
	done <- struct{}{}
}
