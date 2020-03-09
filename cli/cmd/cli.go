/*
Copyright © 2019 Guo Xudong

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
	"fmt"
	"github.com/mbndr/figlet4go"
	"github.com/spf13/cobra"
	"os"
)

type Cli struct {
	rootCmd *cobra.Command
}

//NewCli returns the cli instance used to register and execute command
func NewCli() *Cli {
	ascii := figlet4go.NewAsciiRender()
	logo, _ := ascii.Render("Mao Xian Bot")
	cli := &Cli{
		rootCmd: &cobra.Command{
			Use:   "maoxian",
			Short: "Bot of hub.",
			Long:  fmt.Sprintf("Bot of hub to Github & Gitlab\n%sFind more information at: https://github.com/sunny0826/maoxian-lab", logo),
		},
	}
	cli.rootCmd.SetOutput(os.Stdout)
	cli.setFlags()
	return cli
}

func (cli *Cli) setFlags() {
	flags := cli.rootCmd.PersistentFlags()
	flags.BoolVarP(&Debug, "debug", "d", false, "Set client to DEBUG mode")
}

//Run command
func (cli *Cli) Run() error {
	return cli.rootCmd.Execute()
}
