package drone_ci

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/sirupsen/logrus"
)

const (
	buildCommitPath = "%s/api/repos/%s/builds?branch=%s&commit=%s"
)

type BuildOptions struct {
	PathWithNamespace string
	Branch    string
	Commit    string
}

func (c *Client) GetAddr() string {
	return c.addr
}

// Drone API: http://{droneUrl}/api/repos/{namespace}/{name}/builds?branch={branch}&commit={commit}
func (c *Client) BuildCommit(options BuildOptions) (*drone.Build, error) {
	out := new(drone.Build)
	logrus.Info("options:", options)
	uri := fmt.Sprintf(buildCommitPath, c.addr, options.PathWithNamespace, options.Branch, options.Commit)
	logrus.Info("url:", uri)
	err := c.post(uri, nil, out)
	return out, err
}
