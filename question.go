package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type question struct {
	question      string
	answers       []answer
	caseSensitive bool
}

type answer struct {
	option int
	valid  []string
}

/*
createQuestion creates a new question object with the given question. Per default
the answers are not case sensitive.
*/
func createQuestion(questionText string) *question {
	return &question{question: questionText, caseSensitive: false}
}

/*
createYesNo creates a question with predeclared Yes and No answer options. The
Option for No is negative, for Yes positive.
*/
func createYesNo(questionText string) *question {
	question := &question{question: "(Y/N) " + questionText, caseSensitive: false}
	question.createAnswer(-1, "n", "no")
	question.createAnswer(1, "y", "yes")
	return question
}

/*
ask the question and returns the option value of the chosen answer.
*/
func (q *question) ask() int {
	// prepare console reader
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(q.question)
	// keep asking the question until we get an answer
	for {
		// read input
		input, _ := reader.ReadString('\n')
		input = strings.Trim(input, "\n")
		// check if one of the answers has been chosen
		for _, a := range q.answers {
			for _, value := range a.valid {
				if q.caseSensitive && value == input {
					// if case sensitive it must be an exact match
					return a.option
				} else if !q.caseSensitive && strings.EqualFold(value, input) {
					// if not case sensitive use EqualFold
					return a.option
				}
			} // end answer check
		} // end all answers check: if we reach this no legal answer was given
		fmt.Println("Invalid reply!\n", q.question)
	}
}

/*
createAnswer creates an answer for the given question.
*/
func (q *question) createAnswer(option int, valid ...string) {
	q.answers = append(q.answers, answer{option: option, valid: valid})
}
