package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type problem struct {
	question string
	answer   int
}

type quiz struct {
	problems []problem
	score    int
}

var quizFile string

func init() {
	const defaultQuizFile = "problems.csv"
	flag.StringVar(&quizFile, "qf", defaultQuizFile, "quiz file name (only csv)")
}

func main() {
	game := quiz{}
	flag.Parse()
	problems, err := buildQuiz(quizFile)
	if err != nil {
		log.Fatalf("%v", err)
	}
	game.problems = problems
	fmt.Println("Let's Play a Tiny Math Quiz!")
	game = playQuiz(game)
	fmt.Println("Thanks for Playing!")
	fmt.Printf("You scored %d out of %d\n", game.score, len(game.problems))
}

func buildQuiz(quizFile string) (problems []problem, err error) {
	f, e := os.Open(quizFile)
	defer f.Close()
	if e != nil {
		err = e
		return
	}
	r := csv.NewReader(f)
	for {
		record, e := r.Read()
		if e == io.EOF {
			break
		}
		if e != nil {
			err = e
			return
		}
		ansAsInt, e := strconv.Atoi(record[1])
		if e != nil {
			err = e
			return
		}
		problems = append(problems, problem{record[0], ansAsInt})
	}
	return
}

func playQuiz(game quiz) quiz {
	s := bufio.NewScanner(os.Stdin)
	for _, record := range game.problems {
		fmt.Printf("Question: %s?\nYour Answer: ", record.question)
		s.Scan()
		yourAns := s.Text()
		correctAns := record.answer
		yourAnsAsInt, err := strconv.Atoi(yourAns)
		if err != nil {
			continue
		}
		if yourAnsAsInt == correctAns {
			game.score++
		}
	}
	return game
}
