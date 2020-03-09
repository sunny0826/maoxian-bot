package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sunny0826/maoxian-bot/lab"
)

type AddHookCommand struct {
	baseCommand
	hookUrl  string
	access   string
	repo     string
	site     string
	baseUrl  string
	username string
	confirm  bool
}

func (ah *AddHookCommand) Init() {
	ah.command = &cobra.Command{
		Use:     "addhook",
		Short:   "add webhook to github or gitlab",
		Long:    "add webhook to github or gitlab",
		Aliases: []string{"a"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ah.runAdd(cmd, args)
		},
		Example: addhookExample(),
	}
	ah.command.Flags().StringVarP(&ah.hookUrl, "hook", "w", "http://test-bot.keking.cn", "url of webhook server")
	ah.command.Flags().StringVarP(&ah.access, "access", "a", "NQUQC8KqzQWzYNmsTGtK", "admin secret token of github or gitlab")
	ah.command.Flags().StringVarP(&ah.repo, "repo", "r", "kk-devops-cicd", "repo of github or gitlab")
	ah.command.Flags().StringVarP(&ah.baseUrl, "baseurl", "b", "https://git.keking.cn", "url of github or gitlab")
	ah.command.Flags().StringVarP(&ah.username, "username", "u", "kk-devops", "lab name in github or gitlab")
	ah.command.Flags().StringVarP(&ah.site, "site", "s", "gitlab", "github or gitlab")
	ah.command.Flags().BoolVarP(&ah.confirm, "confirm", "c", false, "whether to configure")
}

func (ah *AddHookCommand) runAdd(command *cobra.Command, args []string) error {
	switch ah.site {
	case "gitlab":
		gitlabBot := lab.GitlabBot{
			HookUrl:  ah.hookUrl,
			Repo:     ah.repo,
			Username: ah.username,
		}
		client := lab.GitlabClient(ah.access, ah.baseUrl)
		if ah.confirm {
			lab.CheckGitlab(client, gitlabBot)
		} else {
			lab.InitGitlab(client, gitlabBot)
		}
	case "github":
		//todo github
	default:
		//todo github
	}
	return nil
}

func addhookExample() string {
	return `
mxbot addhook --hook=http://test-lab.keking.cn --access=XXXXX --repo=kk-devops --baseurl=https://git.keking.cn --username=kk-devops --site=site --confirm=false
`
}
