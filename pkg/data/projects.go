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

	"github.com/6oof/mole/pkg/consts"
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
	Deleted       bool   `json:"deleted"`
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
		return fp, errors.New("Sorry, no project was found!\nYou can use the list command to see all projects")
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

	c := exec.Command("git", "clone")

	c = exec.Command("git", "clone", "--depth", "1", "-b", project.Branch, project.RepositoryUrl, clonePath)
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

func CreateProject(newProject Project) error {

	clonePath := path.Join(consts.BasePath, "projects", newProject.Name)

	err := cloneProject(newProject)
	if err != nil {
		return err
	}

	err = checkEnvGitignore(newProject)
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
