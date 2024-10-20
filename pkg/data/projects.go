package data

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/6oof/mole/pkg/consts"
	"github.com/gofrs/flock"
	"github.com/lithammer/shortuuid/v4"
)

var moleJson = path.Join(consts.BasePath, "/etc/mole/mole.json")
var moleLock = path.Join(consts.BasePath, "/etc/mole/mole.lock")

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

func FindProject(searchTerm string) string {
	p, err := readProjectsFromFile()
	if err != nil {
		return err.Error()
	}

	var fp Project

	for _, pro := range p.Projects {
		if strings.ToLower(pro.Name) == strings.ToLower(searchTerm) {
			fp = pro
		} else if pro.ProjectId == searchTerm {
			fp = pro
		}
	}

	if fp.ProjectId == "" {
		return "Sorry, no project was found!\nYou can use the list command to see all projects"
	}

	return fp.Stringify()
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

func AddProject(newProject Project) (string, error) {
	p, err := readProjectsFromFile()
	if err != nil {
		return "", err
	}

	fileLock := flock.New(moleLock)
	locked, err := fileLock.TryLock()
	defer fileLock.Unlock()

	if err != nil {
		return "", err
	}

	if locked {
		proId := shortuuid.New()

		newProject.ProjectId = proId

		p.Projects = append(p.Projects, newProject)

		err = p.saveProjectsToFile()

		if err != nil {
			return "", err
		}

		return newProject.ProjectId, nil
	} else {
		return "", errors.New("Someone else is trying to work with the project store, please try again")
	}

}

func EditProject(proId, desc, branch string) error {
	newDesc := desc
	newBranch := branch

	p, err := readProjectsFromFile()
	if err != nil {
		return err
	}

	found := false

	for i, pro := range p.Projects {
		if proId == pro.ProjectId {
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
		return errors.New("Project with ID " + proId + " not found")
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
