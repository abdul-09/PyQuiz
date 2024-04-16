package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	filename, timeLimit, shuffle := parseFlags()
	problems, err := readProblems(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if shuffle {
		shuffleProblems(problems)
	}

	score := conductQuiz(problems, timeLimit)
	fmt.Printf("You scored %d out of %d.\n", score, len(problems))
}

func parseFlags() (string, int, bool) {
	filename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "shuffle the order of the quiz questions")
	flag.Parse()
	return *filename, *timeLimit, *shuffle
}

func readProblems(filename string) ([]problem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open the CSV file: %v", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.LazyQuotes = true // Allow multi-line quoted fields
	lines, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse the CSV file")
	}

	var problems []problem
	for _, line := range lines {
		if len(line) != 2 {
			return nil, fmt.Errorf("invalid format in CSV file: each line should have two fields")
		}
		question := strings.TrimSpace(line[0])
		answer := strings.TrimSpace(line[1])
		problems = append(problems, problem{q: question, a: answer})
	}

	return problems, nil
}

func shuffleProblems(problems []problem) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}

func conductQuiz(problems []problem, timeLimit int) int {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	correct := 0

	input := bufio.NewReader(os.Stdin)
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)

		answerChan := make(chan string)
		go func() {
			answer, _ := input.ReadString('\n')
			answerChan <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			return correct
		case answer := <-answerChan:
			if strings.TrimSpace(answer) == p.a {
				correct++
			}
		}
	}
	return correct
}
