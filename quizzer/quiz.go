package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type que struct {
	question string
	answer   string
}

//	func quizzer(quiz_paper []que, score *int) {
//		var ans_buff string
//		for _, value := range quiz_paper {
//			fmt.Println(value.question)
//			_, err := fmt.Scanln(&ans_buff)
//			if err != nil {
//				return
//			}
//			if ans_buff == value.answer {
//				*score += 1
//			}
//		}
//	}
func main() {
	csvFile := flag.String("csv", "problems.csv", "Input file name")
	timer := flag.Int("time", 30, "Duration of Quiz")
	shuffle := flag.Bool("shuffle", false, "Do you want to shuffle")
	flag.Parse()
	quiz_paper := readCSV(csvFile)
	if *shuffle == true {
		for i := range quiz_paper {
			j := rand.Intn(len(quiz_paper))                             // Generate a random index
			quiz_paper[i], quiz_paper[j] = quiz_paper[j], quiz_paper[i] // Swap elements
		}
	}
	score := 0
	total_score := len(quiz_paper)
	var begin string

	fmt.Println("Do you want to Start")
	fmt.Scanln(&begin)
	if begin == "Yes" || begin == "yes" {
		timeman := time.NewTimer(time.Second * time.Duration(*timer))
		for _, value := range quiz_paper {
			fmt.Println(value.question)
			answerChan := make(chan string)

			go func() {
				var answer_buff string
				_, err := fmt.Scanln(&answer_buff)
				if err != nil {
					return
				}
				answerChan <- answer_buff
			}()
			select {
			case <-timeman.C:
				fmt.Println("You have scored", score, "Out of", total_score)
				return
			case answer := <-answerChan:
				if value.answer == answer {
					score++
				}
			}

		}

	}
	fmt.Println("You scored", score, "Out of--", total_score)

}

func readCSV(csvFile *string) []que {
	csv_data, err := os.ReadFile(*csvFile)
	if err != nil {
		panic(err)
	}
	quiz_paper := []que{}
	r := csv.NewReader(strings.NewReader(string(csv_data))) //NewReader implements a readble interface
	for {
		record, e := r.Read()
		if e == io.EOF {
			break
		}
		if e != nil {
			panic(e)
		}
		temp_struct := que{record[0], record[1]}
		quiz_paper = append(quiz_paper, temp_struct)
	}
	return quiz_paper
}

//TODO -- Implement channels instead of this unstable solution
