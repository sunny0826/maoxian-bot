package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// Command is cli command interface
type Command interface {
	// Init command
	Init()

	// CobraCmd
	CobraCmd() *cobra.Command
}

// baseCommand
type baseCommand struct {
	command *cobra.Command
}

func (bc *baseCommand) Init() {
}

func (bc *baseCommand) CobraCmd() *cobra.Command {
	return bc.command
}

func (bc *baseCommand) Name() string {
	return bc.command.Name()
}

type fileWriterWithoutErr struct {
	io.Writer
}

var (
	Debug bool
)

const BotLog = "maoxian.log"

//AddCommand is add child command to the parent command
func (bc *baseCommand) AddCommand(child Command) {
	child.Init()
	childCmd := child.CobraCmd()
	//childCmd.PreRun = func(cmd *cobra.Command, args []string) {
	//	InitLog()
	//}
	bc.CobraCmd().AddCommand(childCmd)
}


