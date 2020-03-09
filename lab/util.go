package lab

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

func GitlabClient(secret string, baseUrl string) *gitlab.Client {
	git := gitlab.NewClient(nil, secret)
	url := fmt.Sprintf("%s/api/v4", baseUrl)
	err := git.SetBaseURL(url)
	if err != nil {
		logrus.Error(err)
	}
	return git
}

// create gitlab issue comment
func createComment(body string, projectId int, issueIid int, client *gitlab.Client) error {
	note_opt := &gitlab.CreateIssueNoteOptions{
		Body: gitlab.String(body),
	}
	_, _, err := client.Notes.CreateIssueNote(projectId, issueIid, note_opt)
	if err != nil {
		return err
	}
	logrus.Infof("create issue note %s", body)
	return nil
}

// check branch is exist
func checkBranch(projectId int, branch string, client *gitlab.Client) (bool, error) {
	bran, _, err := client.Branches.GetBranch(projectId, branch)
	if err != nil {
		return false, err
	}
	if bran != nil {
		return true, nil
	}
	return false, nil
}

func gitCommitSha(projectId int, branch, module string, client *gitlab.Client) (string, error) {
	commit_opt := &gitlab.ListCommitsOptions{
		RefName:     gitlab.String(branch),
		Path:        gitlab.String(module),
		FirstParent: gitlab.Bool(false),
	}
	listCom, _, err := client.Commits.ListCommits(projectId, commit_opt)
	if err != nil {
		return "error", err
	}

	return listCom[0].ID, nil
}

func getShaOfMergeRequests(projectId int, sha string, client *gitlab.Client) (string, error) {
	mrc, _, err := client.Commits.GetMergeRequestsByCommit(projectId, sha)
	if err != nil {
		return "error", err
	}
	return mrc[0].SHA, nil
}

// list lable
func listLable(lable []string, projectId int, client *gitlab.Client) ([]*gitlab.Label, error) {
	lab_opt := &gitlab.ListLabelsOptions{}
	lab, _, err := client.Labels.ListLabels(projectId, lab_opt)
	if err != nil {
		return nil, err
	}
	return lab, nil
}

// makedown template
func makedownTpl(command *Command, event IssueCommentEvent, out *drone.Build, addr, module string) string {
	var tpl string
	tpl = "### Start Build :100:\n"
	if module != "" {
		tpl += fmt.Sprintf("**module: `%s`**\n", module)
	}
	branchUrl := fmt.Sprintf("%s/tree/%s", event.Project.Homepage, command.Args[0])
	deployUrl := fmt.Sprintf("%s/%s/%s/%v", addr, event.Project.Namespace, event.Project.Name, out.Number)
	tpl += fmt.Sprintf("- commit: %s\n", out.After)
	tpl += fmt.Sprintf("- build path: %s\n", deployUrl)
	tpl += fmt.Sprintf("- build branch: [%s](%s)\n", command.Args[0], branchUrl)
	tpl += fmt.Sprintf("- build user: @%s\n", event.User.Username)
	tpl += fmt.Sprintf(`
<details>
<summary>详细信息</summary>
<pre>
BuildNum : %v
Status: %s
CommitAuthor: %s
Command: %s
ProjectId: %v
ProjectNmae: %s
</pre>
</details>
`, out.Number, out.Status, out.AuthorName, command.Command, event.ProjectID, event.Project.Name)
	return tpl
}

// table markdown
func tableMakedown(labels []*gitlab.Label) string {
	var tpl string
	if labels != nil {
		tpl = "### Labels Table\n\n"
		tpl += "| Title | Description | Color | Text Color |\n| ---      |  ------  |---------|---------:|\n"
		for _, label := range labels {
			tpl += fmt.Sprintf("| `%s` | %s | `%s` | `%s` |\n", label.Name, label.Description, label.Color, label.TextColor)
		}
	} else {
		tpl = "### No label for this project. :sweat_smile:"
	}
	return tpl
}
