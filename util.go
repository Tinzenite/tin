package main

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
