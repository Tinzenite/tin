package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
getInput poses a request to the user and returns his entry.
*/
func getString(request string) string {
	fmt.Println(request)
	// read input
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.Trim(input, "\n")
	return input
}

func getInt(request string) int {
	for {
		fmt.Println(request)
		// read input
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.Trim(input, "\n")
		value, err := strconv.ParseInt(input, 10, 0)
		if err == nil {
			return int(value)
		}
	}
}

/*
cmd is the enum for which operation the program should execute. Satisfies the
Value interface so that it can be used in flag.
*/
type cmd int

const (
	cmdNone cmd = iota
	cmdCreate
	cmdLoad
	cmdBootstrap
)

func (c cmd) String() string {
	switch c {
	case cmdNone:
		return "none"
	case cmdCreate:
		return "create"
	case cmdLoad:
		return "load"
	case cmdBootstrap:
		return "bootstrap"
	default:
		return "unknown"
	}
}

/*
cmdParse parses a string to cmd. If illegal or can not be matched will simply
return cmdNone.
*/
func cmdParse(value string) cmd {
	switch value {
	case "create":
		return cmdCreate
	case "load":
		return cmdLoad
	case "bootstrap":
		return cmdBootstrap
	default:
		return cmdNone
	}
}
