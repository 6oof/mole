package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/6oof/mole/pkg/consts"
	"github.com/6oof/mole/pkg/enums"
	"github.com/6oof/mole/pkg/helpers"
	"github.com/gofrs/flock"
	"github.com/lithammer/shortuuid/v4"
)

var moleJson = path.Join(consts.BasePath, "mole.json")
var moleLock = path.Join(consts.BasePath, "mole.lock")

type Projects struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	ProjectId     string `json:"projectId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	RepositoryUrl string `json:"repositoryUrl"`
	Branch        string `json:"branch"`
}

func readProjectsFromFile() (Projects, error) {
	fileLock := flock.New(moleLock)
	locked, err := fileLock.TryLock()
	defer fileLock.Unlock()

	if err != nil {
		return Projects{}, err
	}

	if locked {
		f, err := os.ReadFile(moleJson)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				err = os.WriteFile(moleJson, []byte{}, 0644)
				if err != nil {
					return Projects{}, errors.New("mole.json was missing. An attempt to create it failed")
				}
			}

			return Projects{}, err
		}

		var p Projects

		err = json.Unmarshal(f, &p)
		if err != nil {
			return Projects{}, err
		}

		return p, nil
	} else {
		return Projects{}, errors.New("Someone else is trying to work with the project store, please try again")
	}
}

func ListProjects() string {
	p, err := readProjectsFromFile()
	if err != nil {
		return err.Error()
	}

	return p.Stringify()
}

func FindProject(searchTerm string) (Project, error) {
	p, err := readProjectsFromFile()
	if err != nil {
		return Project{}, err
	}

	var fp Project

	for _, pro := range p.Projects {
		if strings.ToLower(pro.Name) == strings.ToLower(searchTerm) || pro.ProjectId == searchTerm {
			fp = pro
		}
	}

	if fp == (Project{}) {
		return fp, errors.New("Sorry, no project was found!\nYou can use the \"mole projects list\" command to see all projects")
	}

	return fp, nil
}

func (p Projects) saveProjectsToFile() error {
	f, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(moleJson, f, 0644)
	if err != nil {
		return err
	}

	return nil
}

func addProject(newProject Project) error {
	p, err := readProjectsFromFile()
	if err != nil {
		return err
	}

	fileLock := flock.New(moleLock)
	locked, err := fileLock.TryLock()
	defer fileLock.Unlock()

	if err != nil {
		return err
	}

	if locked {
		proId := shortuuid.New()

		newProject.ProjectId = proId

		p.Projects = append(p.Projects, newProject)

		err = p.saveProjectsToFile()

		if err != nil {
			return err
		}

		return nil
	} else {
		return errors.New("Someone else is trying to work with the project store, please try again")
	}

}

func cloneProject(project Project) error {
	clonePath := path.Join(consts.BasePath, "projects", project.Name)

	var stErr bytes.Buffer

	c := exec.Command("git", "clone", "--depth", "1", "-b", project.Branch, project.RepositoryUrl, clonePath)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = &stErr

	err := c.Run()
	if err != nil {
		return errors.New(stErr.String())
	}
	return nil

}

func checkEnvGitignore(project Project) error {
	clonePath := path.Join(consts.BasePath, "projects", project.Name, ".gitignore")

	gi, err := os.ReadFile(clonePath)
	if err != nil {
		return errors.New(`Gitignore is missing from this project. This is mandatory for mole to work properly.`)
	}

	if !containsEnvEntry(string(gi)) {
		return errors.New(".gitignore does not include an entry for '.env'. This is mandatory for mole to work properly")
	}

	return nil
}

func containsEnvEntry(content string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == ".env" {
			return true
		}
	}
	return false
}

func ensureProjectVolume(project Project) error {
	volumePath := path.Join(consts.BasePath, "volumes", project.Name)

	err := os.MkdirAll(volumePath, 0755)
	if err != nil {
		return err
	}

	return nil
}

type baseEnvData struct {
	PType      string
	EnvPath    string
	VolumePath string
	Services   string
	PName      string
	AppKey     string
	PortApp    int
	PortTwo    int
	PortThree  int
}

func createProjectBaseEnv(project Project, pType enums.ProjectType) error {
	domainTemplate := `# Auto-generated environment configuration for {{.PName}}.
# DO NOT DELETE OR MODIFY THIS SECTION.
# This configuration is necessary for the project to work properly.
# Static path to this file on mole managed servers is:
# {{.EnvPath}}

# Available types: static, podman, systemd
MOLE_PROJECT_TYPE={{.PType}}

# Volume path to be used in podman quadlets
MOLE_VOLUME_PATH={{.VolumePath}}

# Comma separated list of services to start ("service-1,service-2").
MOLE_SERVICES={{.Services}}

# Three reserved ports for this deployment.
MOLE_APP_PORT={{.PortApp}}
MOLE_TWO_PORT={{.PortTwo}}
MOLE_THREE_PORT={{.PortThree}}

# Random string to be used as a key when necessary
MOLE_APP_KEY={{.AppKey}}

# User-defined environment variables can be added below.
# Add your own variables here:`

	mp, err := FindAndReserveMolePorts()
	if err != nil {
		return err
	}

	key := helpers.GenerateAppKey()

	be := baseEnvData{
		PType:      pType.String(),
		EnvPath:    "/home/mole/projects/" + project.Name + "/.env",
		VolumePath: "/home/mole/volumes/" + project.Name,
		Services:   "app",
		PName:      project.Name,
		PortApp:    mp[0],
		PortTwo:    mp[1],
		PortThree:  mp[2],
		AppKey:     key,
	}

	tmpl, err := template.New("env").Parse(domainTemplate)
	if err != nil {
		return err
	}

	var ft bytes.Buffer

	err = tmpl.Execute(&ft, be)
	if err != nil {
		return err
	}

	efp := path.Join(consts.BasePath, "projects", project.Name, ".env")

	err = os.WriteFile(efp, ft.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

// TODO: make sure to check if there are any services defined in the correct folder for the types that need them
// TODO: make sure to check for the deploy script for projects that need it
func CreateProject(newProject Project, projectType string) error {

	pt, err := enums.IsProjectType(projectType)
	if err != nil {
		return err
	}

	clonePath := path.Join(consts.BasePath, "projects", newProject.Name)

	err = cloneProject(newProject)
	if err != nil {
		return err
	}

	err = checkEnvGitignore(newProject)
	if err != nil {
		os.RemoveAll(clonePath)
		return err
	}

	err = ensureProjectVolume(newProject)
	if err != nil {
		os.RemoveAll(clonePath)
		return err
	}

	err = createProjectBaseEnv(newProject, pt)
	if err != nil {
		os.RemoveAll(clonePath)
		return err
	}

	err = addProject(newProject)
	if err != nil {
		os.RemoveAll(clonePath)
	}

	return err

}

// NOI stands for name or Id
func EditProject(proNOI, desc, branch string) error {
	newDesc := desc
	newBranch := branch

	p, err := readProjectsFromFile()
	if err != nil {
		return err
	}

	found := false

	for i, pro := range p.Projects {
		if proNOI == pro.ProjectId || proNOI == pro.Name {
			found = true

			if newDesc != "" {
				p.Projects[i].Description = newDesc
			}
			if newBranch != "" {
				p.Projects[i].Branch = newBranch
			}

		}
	}

	if !found {
		return errors.New("Project with ID " + proNOI + " not found")
	}

	err = p.saveProjectsToFile()

	return err

}

func DeleteProject(proId string) error {
	p, err := readProjectsFromFile()
	if err != nil {
		return err
	}

	found := false

	for i, pro := range p.Projects {
		if proId == pro.ProjectId {
			found = true
			p.Projects = append(p.Projects[:i], p.Projects[i+1:]...)
		}
	}

	if !found {
		return errors.New("Project with ID " + proId + " not found")
	}

	err = p.saveProjectsToFile()

	return err
}

func (ps Project) Stringify() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(" |ID     : " + ps.ProjectId + "\n")
	b.WriteString(" |Name   : " + ps.Name + "\n")
	b.WriteString(" |Desc.  : " + ps.Description + "\n")
	b.WriteString(" |Git    : " + ps.RepositoryUrl + "\n")
	b.WriteString(" |Branch : " + ps.Branch + "\n")

	return b.String()
}

func (p Projects) Stringify() string {
	var b strings.Builder

	for i, pro := range p.Projects {
		b.WriteString("\n")
		b.WriteString(strconv.Itoa(i) + ":\n")
		b.WriteString(pro.Stringify())
	}

	return b.String()
}
