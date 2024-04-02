package main

import (
	"fmt"
	"os"

	"github.com/pablo-costanzo/git-go/pkg/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mygit <command> [<args>...]")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	args := os.Args[2:]

	var command cmd.Executable

	switch cmdName {
	case "init":
		command = cmd.NewInit()
	case "cat-file":
		if len(args) < 2 {
			command.Help()
			os.Exit(1)
		}
		command = cmd.NewCatFile(args)
	default:
		fmt.Printf("Unknown command: %s\n", cmdName)
		os.Exit(1)
	}

	if err := command.Exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
