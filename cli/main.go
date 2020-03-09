package main

import (
	"fmt"
	"github.com/sunny0826/maoxian-bot/cli/cmd"
	"os"
)

func main() {
	baseCommand := cmd.CmdInit()
	if err := baseCommand.CobraCmd().Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
