package lab

import (
	"fmt"
	"github.com/sirupsen/logrus"
	drone_ci "github.com/sunny0826/maoxian-bot/drone-ci"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type IssueCommentEvent gitlab.IssueCommentEvent

const (
	BuildCommand = "/build"
	DogCommand   = "/doge"
	LableCommad  = "/labels"
	dogeUrl      = "https://tva2.sinaimg.cn/large/ad5fbf65ly1g1194gx01uj23402c0e82.jpg"
)

type Command struct {
	Command string   // command
	Args    []string // args
}

func Process(client *gitlab.Client, droneClient drone_ci.Client, event IssueCommentEvent) error {
	projecId := event.ProjectID
	iid := event.Issue.IID
	cmd := event.ObjectAttributes.Note
	commands := decodeCommand(cmd)
	for _, command := range commands {
		switch command.Command {
		case BuildCommand:
			branch := command.Args[0]
			var module string
			if len(command.Args) > 1 {
				module = command.Args[1]
			}
			isExist, err := checkBranch(projecId, branch, client)
			if err != nil {
				logrus.Error(err)
			}
			var body string
			if isExist {
				sha, err := gitCommitSha(event.ProjectID, branch, module, client)
				if err != nil {
					logrus.Error(err)
				}
				if module != "" {
					sha, err = getShaOfMergeRequests(event.ProjectID, sha, client)
					if err != nil {
						logrus.Error(err)
					}
				}
				buildOpt := drone_ci.BuildOptions{
					PathWithNamespace: event.Project.PathWithNamespace,
					Branch:    branch,
					Commit:    sha,
				}
				out, err := droneClient.BuildCommit(buildOpt)
				if err != nil {
					logrus.Error(err)
					body = fmt.Sprintf("### :x: Build Error: %s", err.Error())
				} else {
					addr := droneClient.GetAddr()
					body = makedownTpl(command, event, out, addr, module)
				}
			} else {
				logrus.Warningf("Can not find branch:[%s]", branch)
				body = fmt.Sprintf("### :x: Can not find branch: [%s]", command.Args[0])
			}
			err = createComment(body, projecId, iid, client)
			if err != nil {
				logrus.Error(err)
			}
		case LableCommad:
			labs, err := listLable(command.Args, projecId, client)
			if err != nil {
				logrus.Error(err)
			}
			body := tableMakedown(labs)
			err = createComment(body, projecId, iid, client)
			if err != nil {
				logrus.Error(err)
			}
		case DogCommand:
			body := fmt.Sprintf(":sunny:\n\n![doge](%s)", dogeUrl)
			err := createComment(body, projecId, iid, client)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
	fmt.Println(commands)
	return nil
}

// decode command
func decodeCommand(s string) []*Command {
	var res []*Command
	fmt.Printf("decode cmd: %s\n", s)
	rows := strings.Split(s, "\n")
	for _, r := range rows {
		command := &Command{}
		if strings.HasPrefix(r, "/") {
			cmd := strings.Split(r, " ")
			if len(cmd) == 1 {
				command.Command = cmd[0]
			} else {
				command.Command = cmd[0]
				command.Args = cmd[1:]
			}
			res = append(res, command)
		}
	}
	return res
}
