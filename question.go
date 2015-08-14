package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/tinzenite/shared"
)

type question struct {
	question string
	answers  []answer
}

type answer struct {
	option int
	valid  []string
}

func (q *question) ask() int {
	// prepare console reader
	reader := bufio.NewReader(os.Stdin)
	log.Println(q.question)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.Trim(input, "\n")
		for _, a := range q.answers {
			if shared.Contains(a.valid, input) {
				return a.option
			}
		}
		log.Println("Invalid reply!\n", q.question)
	}
}

func (q *question) createAnswer(option int, valid ...string) {
	q.answers = append(q.answers, answer{option: option, valid: valid})
}
