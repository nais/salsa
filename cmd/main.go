package main

import (
	"github.com/nais/salsa/pkg/commands"
)

var (
	// Is set during build
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	commands.Execute(version, commit, date, builtBy)
}
