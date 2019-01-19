package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
)

type quizRecord struct {
	question string
	answer   int
}

var quizFile string

func init() {
	const defaultQuizFile = "problems.csv"
	flag.StringVar(&quizFile, "qf", defaultQuizFile, "quiz file name (only csv)")
}

func main() {
	flag.Parse()
	quiz, err := buildQuiz(quizFile)
	_ = quiz
	_ = err
}

func buildQuiz(quizFile string) (quiz []quizRecord, err error) {
	f, err := os.Open(quizFile)
	defer f.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v", err)
		}
		ansAsInt, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatalf("%v", err)
		}
		quiz = append(quiz, quizRecord{record[0], ansAsInt})
	}
	return
}
