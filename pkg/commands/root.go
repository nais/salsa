package commands

import (
	"os"
)

func Execute() {
	args := os.Args[1:]

	if len(args) == 0 {
		println("Please specify subcommand")
		os.Exit(1)
	}

	subCmd := args[0]

	switch subCmd {
	case "create":
		CreateYolo()
	default:
		println("Unknown subcommand")
		os.Exit(1)
	}
}
