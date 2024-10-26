package enums

import "fmt"

type ProjectType int

const (
	Podman ProjectType = iota
	Static
	Systemd
)

var typeName = map[ProjectType]string{
	Podman:  "podman",
	Static:  "static",
	Systemd: "systemd",
}

func (p ProjectType) String() string {
	return typeName[p]
}

func IsProjectType(t string) (ProjectType, error) {
	for pt, name := range typeName {
		if name == t {
			return pt, nil
		}
	}

	return -1, fmt.Errorf("invalid project type: %s", t)
}
