package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	shuffle := flag.Int("shuffle", 0, "shuffle the problem list(if not 0)")
	timeLimit := flag.Int("limit", 3, "each time limit for the quiz in seconds")
	timeTotal := flag.Int("limitTotal", 15, "total time limit for the quiz in seconds")

	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	problems := parseLines(lines)
	if *shuffle != 0 {
		problems = shuffleProblem(problems)
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	timerTotal := time.NewTimer(time.Duration(*timeTotal) * time.Second)
	correct := 0
	answerCh := make(chan string)

problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timerTotal.C:
			fmt.Println("\n Total time up")
			break problemloop
		case <-timer.C:
			var tmp string
			timer.Reset(time.Duration(*timeLimit) * time.Second)
			fmt.Println(" -->time up")
			fmt.Scanln("%s\n", &tmp) // to clear the buffer
			break
		case answer := <-answerCh:
			if strings.TrimSpace(answer) == p.a {
				correct++
			} else {
				fmt.Println("[error] input : [" + answer + "]")
			}
			timer.Reset(time.Duration(*timeLimit) * time.Second)

		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func shuffleProblem(source []problem) []problem {
	out := make([]problem, len(source))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	perm := r.Perm(len(source))
	for i, randIndex := range perm {
		out[i] = source[randIndex]
	}
	return out
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
