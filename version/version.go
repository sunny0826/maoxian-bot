package version

import "runtime"

var (
	Ver       = "unknown"
	GoOs      = runtime.GOOS
	GoArch    = runtime.GOARCH
	GitCommit = "$Format:%H$"          // sha1 from git, output of $(git rev-parse HEAD)
	BuildDate = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)
