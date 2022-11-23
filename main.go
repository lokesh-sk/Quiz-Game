package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

func parseCSV(lines [][]string) []problem {
	var problems []problem
	for _, line := range lines {
		record := problem{
			question: line[0],
			answer:   line[1],
		}
		problems = append(problems, record)
	}
	return problems
}

func askQuiz(ch chan int, quiz problem) {
	var answer string
	fmt.Scanln(&answer)
	if answer == quiz.answer {
		ch <- 1
		return
	}
	ch <- 0
}

func main() {

	// csvFileName the name of the file containing the problem.
	// The default filename is problem.csv
	csvFileName := flag.String("csv", "problem.csv", "The csv filename containing problem and answer")

	// timeLimit for the quiz the default timelimit is 30 seconds
	timeLimit := flag.Int("timeLimit", 30, "The time limit for the quiz")

	flag.Parse()

	file, err := os.Open(*csvFileName)

	if err != nil {
		fmt.Println("Unable to open the file")
		os.Exit(1)
	}

	r := csv.NewReader(file)

	lines, err := r.ReadAll()

	if err != nil {
		fmt.Println("Unable to read the file")
		os.Exit(1)
	}

	problems := parseCSV(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	answerChan := make(chan int)

	correctAnswer := 0

	for i, quiz := range problems {
		fmt.Println(i, "Question -->", quiz.question)
		go askQuiz(answerChan, quiz)

		select {
		case <-timer.C:
			fmt.Println("The quiz time over")
			fmt.Println("Correct answer = ", correctAnswer)
			return

		case answer := <-answerChan:
			correctAnswer += answer
		}

	}
	fmt.Println("Correct answer = ", correctAnswer)
}
