package drone_ci

import (
	"fmt"
	"github.com/drone/drone-go/drone"
)

const (
	buildCommitPath = "%s/api/repos/%s/%s/builds?branch=%s&commit=%s"
)

type BuildOptions struct {
	Namespace string
	Name      string
	Branch    string
	Commit    string
}

func (c *Client) GetAddr() string {
	return c.addr
}

// Drone API: http://{droneUrl}/api/repos/{namespace}/{name}/builds?branch={branch}&commit={commit}
func (c *Client) BuildCommit(options BuildOptions) (*drone.Build, error) {
	out := new(drone.Build)
	uri := fmt.Sprintf(buildCommitPath, c.addr, options.Namespace, options.Name, options.Branch, options.Commit)
	err := c.post(uri, nil, out)
	return out, err
}
