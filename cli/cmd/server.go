package cmd

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"github.com/mbndr/figlet4go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	drone_ci "github.com/sunny0826/maoxian-bot/drone-ci"
	"github.com/sunny0826/maoxian-bot/lab"
	"io/ioutil"
	"net/http"
)

type ServerCommand struct {
	baseCommand
	port         string
	access       string
	baseUrl      string
	droneUrl     string
	droneToken   string
	webhookToken string
}

func (s *ServerCommand) Init() {
	s.command = &cobra.Command{
		Use:     "server",
		Short:   "lab runner",
		Long:    "lab runner",
		Aliases: []string{"s"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.runServer(cmd, args)
		},
		Example: serverExample(),
	}
	s.command.Flags().StringVarP(&s.port, "port", "p", "9000", "server listen port")
	s.command.Flags().StringVarP(&s.access, "access", "a", "Xy2pyb4_KG_2YN3sxPfx", "lab access token of gitlab or github")
	s.command.Flags().StringVarP(&s.baseUrl, "baseurl", "b", "https://git.keking.cn", "url of github or gitlab")
	s.command.Flags().StringVar(&s.droneUrl, "droneurl", "http://drone.keking.cn", "url of drone server")
	s.command.Flags().StringVar(&s.droneToken, "dronetoken", "DZYzrsfUTyKHrenElWxVrJG1Xun48D6t", "drone token of bot")
	s.command.Flags().StringVarP(&s.webhookToken, "token", "t", "ba628c268ea4007ed57660b28853ffabc19b2038c16e0616e7d63849aee627d3", "wehbook token of gitlab or github")
}

func (s *ServerCommand) runServer(command *cobra.Command, args []string) error {
	ascii := figlet4go.NewAsciiRender()
	logo, _ := ascii.Render("Mao Xian Bot")
	command.Printf("MaoXian Bot server start.\n%s", logo)
	http.HandleFunc("/", s.handler)
	logrus.Infof("listen port:%s", s.port)
	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		logrus.Error(err)
	}
	return nil
}

func (s *ServerCommand) handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Fatal(err)
	}
	if _, ok := r.Header["X-Gitlab-Event"]; ok {
		logrus.Info("GitLab Event")
		if _, ok := r.Header["X-Gitlab-Token"]; ok {
			if r.Header["X-Gitlab-Token"][0] == s.webhookToken && r.Header["X-Gitlab-Event"][0] == "Note Hook" {
				event := lab.IssueCommentEvent{}
				err := json.Unmarshal(body, &event)
				if err != nil {
					logrus.Fatal(err)
				}
				gitlabClient := lab.GitlabClient(s.access, s.baseUrl)
				droneClient := drone_ci.DroneClient(s.droneUrl, s.droneToken)
				err = lab.Process(gitlabClient, droneClient, event)
				if err != nil {
					logrus.Fatal(err)
				}
			} else {
				logrus.Warning("illegal token,please check you webhook secretToken")
			}
		} else {
			logrus.Info("miss webhookToken")
		}
	} else if _, ok := r.Header["X-GitHub-Event"]; ok {
		logrus.Info("Github Event")
		event := github.IssueCommentEvent{}
		err := json.Unmarshal(body, &event)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		logrus.Info("Up")
		return
	}
}

func serverExample() string {
	return `
`
}
