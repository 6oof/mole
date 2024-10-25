package enums

import "fmt"

type ProjectType int

const (
	podman ProjectType = iota
	static
	systemd
)

var typeName = map[ProjectType]string{
	podman:  "podman",
	static:  "static",
	systemd: "systemd",
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
