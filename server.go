package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type problem struct {
	question string
	answer   int
}

type quiz struct {
	problems []problem
	score    int
}

type settings struct {
	filename string
	duration int
	shuffle  bool
}

const (
	defaultQuizFile      = "problems.csv"
	defaultQuizTimeInSec = 30
	defaultShuffle       = false
)

func main() {
	game := quiz{}
	setting := settings{}
	flag.StringVar(&setting.filename, "qf", defaultQuizFile, "quiz file name (only csv)")
	flag.IntVar(&setting.duration, "t", defaultQuizTimeInSec, "quiz duration in seconds")
	flag.BoolVar(&setting.shuffle, "s", defaultShuffle, "shuffle problems")
	flag.Parse()
	problems, err := getProblems(setting.filename)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if setting.shuffle {
		problems = shuffleProblems(problems)
	}
	game.problems = problems
	timer := startQuiz(setting.duration)
	game = playQuiz(game, timer)
	fmt.Println("Thanks for Playing!")
	fmt.Printf("You scored %d out of %d\n", game.score, len(game.problems))
}

func getProblems(quizFile string) (problems []problem, err error) {
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

func shuffleProblems(problems []problem) []problem {
	rand.Seed(int64(time.Now().Nanosecond()))
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
	return problems
}

func startQuiz(duration int) <-chan time.Time {
	fmt.Println("Let's Play a Tiny Math Quiz!")
	fmt.Printf("You have %d seconds to complete\n", duration)
	fmt.Print("Press ENTER to start!")
	fmt.Scanln()
	return time.After(time.Duration(duration) * time.Second)
}

func playQuiz(game quiz, timer <-chan time.Time) quiz {
	s := bufio.NewScanner(os.Stdin)
	ans := make(chan string)
loop:
	for _, record := range game.problems {
		fmt.Printf("Question: %s?\nYour Answer: ", record.question)
		go getAnswer(s, ans)
		select {
		case <-timer:
			fmt.Println("\nOops! Time has expired!")
			break loop
		case yourAns := <-ans:
			correctAns := record.answer
			yourAnsAsInt, err := strconv.Atoi(yourAns)
			if err != nil {
				continue
			}
			if yourAnsAsInt == correctAns {
				game.score++
			}
		}

	}
	return game
}

func getAnswer(s *bufio.Scanner, ans chan<- string) {
	s.Scan()
	yourAns := s.Text()
	ans <- yourAns
}
