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

type questionBankRecord struct {
	question string
	answer   int
}

type quiz struct {
	totalQuestions   int
	correctQuestions int
}

var quizFile string

func init() {
	const defaultQuizFile = "problems.csv"
	flag.StringVar(&quizFile, "qf", defaultQuizFile, "quiz file name (only csv)")
}

func main() {
	flag.Parse()
	questionBank, err := buildQuiz(quizFile)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println("Let's Play a Tiny Math Quiz!")
	game := playQuiz(questionBank)
	fmt.Println("Thanks for Playing!")
	fmt.Printf("Your Score: %d\nTotal Questions: %d\n", game.correctQuestions, game.totalQuestions)
}

func buildQuiz(quizFile string) (questionBank []questionBankRecord, err error) {
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
		questionBank = append(questionBank, questionBankRecord{record[0], ansAsInt})
	}
	return
}

func playQuiz(questionBank []questionBankRecord) (game quiz) {
	s := bufio.NewScanner(os.Stdin)
	for _, record := range questionBank {
		fmt.Printf("Question: %s?\nYour Answer: ", record.question)
		s.Scan()
		yourAns := s.Text()
		correctAns := record.answer
		yourAnsAsInt, err := strconv.Atoi(yourAns)
		if err != nil {
			continue
		}
		if yourAnsAsInt == correctAns {
			game.correctQuestions++
		}
		game.totalQuestions++
	}
	return
}
