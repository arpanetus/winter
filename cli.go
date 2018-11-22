package main

import (
	"github.com/steplems/winter/core"
	"os"
	"strings"
)

var (
	log = core.NewLogger("cli")
)

const (
	command_help 	= "help"
	command_init 	= "init"
	command_build 	= "build"
	command_run 	= "run"
)

func main() {
	command := ""
	args := []string{}

	if len(os.Args) > 1 {
		command = os.Args[1]
		if len(os.Args) > 2 {
			args = os.Args[2:]
		}
	}

	switch command {
	case command_help:
		helpCommand(args)
	case command_init:
		initCommand(args)
	case command_build:
		buildCommand(args)
	case command_run:
		runCommand(args)
	default:
		helpCommand(args)
		if len(strings.Trim(command, " ")) > 0 {
			log.Err("Unknown command '" + command + "'\n")
		}
	}
}
