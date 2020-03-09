package lab

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

type GitlabBot struct {
	HookUrl  string
	Repo     string
	Username string
}

type InitBot interface {
	InitGitlab(git *GitlabBot)
	CheckGitlab(git *GitlabBot)
}

func InitGitlab(client *gitlab.Client, git GitlabBot) {
	userId, projectId := getId(client, git.Username, git.Repo)
	err := checkAndAddMember(client, projectId, userId)
	if err != nil {
		logrus.Fatal(err)
	}
	h := generateHmac()
	err = addWebhook(client, projectId, git.HookUrl, h)
	if err != nil {
		logrus.Fatal(err)
	}
}

func CheckGitlab(client *gitlab.Client, git GitlabBot) {
	logrus.Info("Check User & Project")
	userId, projectId := getId(client, git.Username, git.Repo)
	logrus.Info("Check project member")
	err := checkMember(client, projectId, userId, git.Repo, git.Username)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Check project webhook")
	err = checkWebhook(client, projectId, git.HookUrl)
	if err != nil {
		logrus.Fatal(err)
	}

}

func getId(client *gitlab.Client, username string, repo string) (int, int) {
	userOpt := &gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{},
		Username:    gitlab.String(username),
	}
	names, _, err := client.Users.ListUsers(userOpt)
	if err != nil {
		logrus.Fatal(err)
	}
	if len(names) == 0 {
		logrus.Fatalf("can not find user:「%s」,please create lab user.", username)
	} else if len(names) > 1 {
		logrus.Fatalf("please lab username,find more than one result,username:%s", username)
	}
	userId := names[0].ID
	var projectId int
	projectOpt := &gitlab.SearchOptions{}
	projects, _, err := client.Search.Projects(repo, projectOpt)
	if err != nil {
		logrus.Fatal(err)
	}
	for _, project := range projects {
		if project.Name == repo {
			projectId = project.ID
			break
		}
	}
	if projectId == 0 {
		logrus.Fatalf("can not find project:「%s」", repo)
	}
	return userId, projectId
}

func checkAndAddMember(client *gitlab.Client, projectId int, userId int) error {
	memberOpt := &gitlab.ListProjectMembersOptions{}
	members, _, err := client.ProjectMembers.ListAllProjectMembers(projectId, memberOpt)
	if err != nil {
		return err
	}
	var isExit bool
	for _, member := range members {
		if member.ID == userId {
			isExit = true
			break
		}
	}
	if !isExit {
		add_opt := &gitlab.AddProjectMemberOptions{
			UserID:      gitlab.Int(userId),
			AccessLevel: gitlab.AccessLevel(20),
		}
		member, _, err := client.ProjectMembers.AddProjectMember(projectId, add_opt)
		if err != nil {
			return err
		}
		logrus.Infof("add member successful! userID:%v userName:%s level:%s", member.ID, member.Name, member.AccessLevel)
	} else {
		logrus.Info("member already exists")
	}
	return nil
}

func checkMember(client *gitlab.Client, projectId int, userId int, project string, user string) error {
	memberOpt := &gitlab.ListProjectMembersOptions{}
	members, _, err := client.ProjectMembers.ListAllProjectMembers(projectId, memberOpt)
	if err != nil {
		return err
	}
	var isExit bool
	for _, member := range members {
		if member.ID == userId {
			isExit = true
			break
		}
	}
	if !isExit {
		logrus.Warningf("project:「%s」do not hava member: 「%s」.After removing the --confirm flag, member:「%s」will be automatically added.", project, user, user)
	} else {
		logrus.Info("member already exists")
	}
	return nil
}

func addWebhook(client *gitlab.Client, projectId int, hookUrl string, hmacToken string) error {
	hookOpt := &gitlab.ListProjectHooksOptions{}
	hooks, _, err := client.Projects.ListProjectHooks(projectId, hookOpt)
	if err != nil {
		return err
	}
	var hookIsExit bool
	for _, hook := range hooks {
		if hook.URL == hookUrl {
			logrus.Warningf("webhook:%s already exists", hookUrl)
			if !hook.NoteEvents {
				logrus.Info("IssuesEvents is closed")
				editOpt := &gitlab.EditProjectHookOptions{
					URL:                   gitlab.String(hookUrl),
					Token:                 gitlab.String(hmacToken),
					NoteEvents:            gitlab.Bool(true),
					PushEvents:            gitlab.Bool(false),
					EnableSSLVerification: gitlab.Bool(false),
				}
				_, _, err := client.Projects.EditProjectHook(hook.ProjectID, hook.ID, editOpt)
				if err != nil {
					return err
				}
				logrus.Info("Open IssuesEvents")
			} else {
				logrus.Info("IssuesEvents already open")
			}
			hookIsExit = true
		}
	}
	if !hookIsExit {
		addOpt := &gitlab.AddProjectHookOptions{
			URL:                   gitlab.String(hookUrl),
			Token:                 gitlab.String(hmacToken),
			NoteEvents:            gitlab.Bool(true),
			PushEvents:            gitlab.Bool(false),
			EnableSSLVerification: gitlab.Bool(false),
		}
		_, _, err := client.Projects.AddProjectHook(projectId, addOpt)
		if err != nil {
			return err
		}
		logrus.Infof("add webhook:%s successful", hookUrl)
		logrus.Printf("token:%s", hmacToken)
	}
	return nil
}

func checkWebhook(client *gitlab.Client, projectId int, hookUrl string) error {
	hookOpt := &gitlab.ListProjectHooksOptions{}
	hooks, _, err := client.Projects.ListProjectHooks(projectId, hookOpt)
	if err != nil {
		return err
	}
	var hookIsExit bool
	for _, hook := range hooks {
		if hook.URL == hookUrl {
			logrus.Warningf("webhook:%s already exists", hookUrl)
			hookIsExit = true
		}
	}
	if !hookIsExit {
		logrus.Warningf("webhook:%s do not exists.After removing the --confirm flag, webhook:「%s」will be automatically added.", hookUrl, hookUrl)
	}
	return nil
}

func generateHmac() string {
	secret := "maoxian"
	data := "lab"
	logrus.Infof("Secret: %s Data: %s", secret, data)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
