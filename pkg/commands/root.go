package commands

import (
	"fmt"
	"os"

	"github.com/nais/salsa/pkg/vcs"
)

func Execute() {
	args := os.Args[1:]

	if len(args) == 0 {
		println("Please specify subcommand")
		os.Exit(1)
	}

	subCmd := args[0]

    // TODO validate second arg
	switch subCmd {
	case "clone":
		repo := args[1]
		err := vcs.CloneRepo(repo, "tmp")
		if err != nil {
			fmt.Printf("something went wrong %v", err)
		}
	case "scan":
		path := args[1]
		GradleScan(path)
	default:
		println("Unknown subcommand")
		os.Exit(1)
	}
}
