package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var score, total int

func main() {
	filePath := flag.String("file", "problems.csv", "Problem file path")
	timeDelay := flag.Int("timer", 30, "Quiz timer (in seconds)")
	flag.Parse()

	if *timeDelay < 1 {
		fmt.Println("Oopsie! Timer value is lower than 1!")
		os.Exit(0)
	}

	timerCh := make(chan struct{})
	quizCh := make(chan struct{})

	go quiz(*filePath, quizCh)
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

func quiz(filePath string, done chan struct{}) {
	f, err := openCSV(filePath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(f)
	scanner := bufio.NewScanner(os.Stdin)
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
	done <- struct{}{}
}
