package main

import (
	"fmt"
	"os"

	"github.com/zwing99/change-log-magician/internal/cli"
)

var version = "devel"

func main() {
	command := cli.NewRootCommand()
	command.Version = version
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
