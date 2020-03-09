package cmd

func CmdInit() *baseCommand {
	cli := NewCli()
	baseCmd := &baseCommand{
		command: cli.rootCmd,
	}
	// add version command
	baseCmd.AddCommand(&VersionCommand{})
	// add add-hook command
	addHookCommand := &AddHookCommand{}
	baseCmd.AddCommand(addHookCommand)
	// add server command
	serverCommand := &ServerCommand{}
	baseCmd.AddCommand(serverCommand)

	return baseCmd
}
