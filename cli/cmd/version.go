/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sunny0826/maoxian-bot/version"
)

type VersionCommand struct {
	baseCommand
}

func (vc *VersionCommand) Init() {
	vc.command = &cobra.Command{
		Use:     "version",
		Short:   "Print version info",
		Long:    "Print version info",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Version: %s\n", version.Ver)
			cmd.Printf("GitCommit: %s\n", version.GitCommit)
			cmd.Printf("BuildDate: %s\n", version.BuildDate)
			cmd.Printf("GoOs: %s\n", version.GoOs)
			cmd.Printf("GoArch: %s\n", version.GoArch)
		},
	}
}
