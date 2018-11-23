package main

import (
	"fmt"
	"runtime"
)

const (
	cli_config_file = "winter.go"
	cli_config_file_so = "winter.so"
	cli_config_dir = ".winter"
	cli_app_config = "App"
)

var (
	help = map[string]string{
		"global": "\n" +
			"Welcome to the Winter CLI!\n" +
			"\n" +
			"Usage:  winter <command> [arguments]\n" +
			"\n" +
			"Commands:\n" +
			"        help        Print this message\n" +
			"        init        Create new project from templates\n" +
			"        build       Build winter project\n" +
			"        run         Run winter project\n" +
			"\n" +
			"Use \"winter help <command>\" to see more information about that command\n",
		command_init: getDoc(
			"winter init <template>",
			"Init creates new project with given template name.\n" +
				"        If the default template is not found, then it finds it in the github repository",
			"        -d --dir    Path to the project directory\n" +
				 "        --dep       Create project with dep"),
		command_run: getDoc(
			"winter run [options]",
			"Run compiles plugins at the working dir into tmp dir\n" +
				"        and runs app with given config in required winter.go file",
			"        -d --dir    Path to the project directory"),
		command_build: getDoc(
			"winter build [options]",
			"Build compiles app with given config in required winter.go file\n" +
				"        at the working dir plugins into one executable",
			"        -d --dir       Path to the project directory\n" +
				 "        -o --output    Build output path\n" +
				 "        -a --args      Arguments for the 'go build' command"),
	}
)

func helpCommand(args []string) {
	if len(args) > 0 {
		cmd := args[0]
		doc, ok := help[cmd]
		if ok {
			fmt.Println(doc)
		} else {
			fmt.Println()
			log.Err("Unknown command '" + cmd + "'\n")
		}
		return
	}
	fmt.Println(help["global"])
}

func initCommand(args []string) {
	if len(args) == 0 {
		helpCommand([]string{command_init})
		log.Err("Template name's required (mvc - recommended template for beginners)\n")
		return
	}
}

func buildCommand(args []string) {
	if runtime.GOOS == "windows" {
		log.Err("Not suitable for Windows, sorry :(")
		return
	}

	winterBuild()
}

func runCommand(args []string) {
	if runtime.GOOS == "windows" {
		log.Err("Not suitable for Windows, sorry :(")
		return
	}
	log.Info("Work in progress")
}
