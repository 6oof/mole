package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/gofrs/flock"
	"github.com/lithammer/shortuuid/v4"
	"github.com/zulubit/mole/pkg/consts"
	"github.com/zulubit/mole/pkg/helpers"
)

// TODO:
// 2. fix all the names in mole secrets
// 3. review all cmds to change wording
// 4. fix the documentations now that nothing is necessary
// 5. change all cobra errors to fatalf

// Projects represents a collection of Project.
type Projects struct {
	Projects []Project `json:"projects"`
}

// Project represents an individual project with its details.
type Project struct {
	ProjectID     string `json:"projectId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	RepositoryURL string `json:"repositoryUrl"`
	Branch        string `json:"branch"`
}

// getMoleJSONPath returns the full path to mole.json based on consts.GetBasePath().
func getMoleJSONPath() string {
	return path.Join(consts.GetBasePath(), "mole.json")
}

// getMoleLockPath returns the full path to mole.lock based on consts.GetBasePath().
func getMoleLockPath() string {
	return path.Join(consts.GetBasePath(), "mole.lock")
}

// readProjectsFromFile reads the project data from the mole.json file.
// It uses a file lock to ensure that no other process is modifying the file at the same time.
func readProjectsFromFile() (Projects, error) {
	fileLock := flock.New(getMoleLockPath())
	locked, err := fileLock.TryLock()
	defer fileLock.Unlock()

	if err != nil {
		return Projects{}, err
	}

	if locked {
		f, err := os.ReadFile(getMoleJSONPath())
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				err = os.WriteFile(getMoleJSONPath(), []byte{}, 0644)
				if err != nil {
					return Projects{}, errors.New("mole.json was missing. An attempt to create it failed")
				}
			}
			return Projects{}, err
		}

		var p Projects
		if err := json.Unmarshal(f, &p); err != nil {
			return Projects{}, fmt.Errorf("failed to unmarshal projects: %w", err)
		}

		return p, nil
	} else {
		return Projects{}, errors.New("someone else is trying to work with the project store, please try again")
	}
}

// ListProjects returns a string representation of all projects in mole.json.
func ListProjects() string {
	p, err := readProjectsFromFile()
	if err != nil {
		return err.Error()
	}

	return p.Stringify()
}

// FindProject searches for a project by name or ID and returns it.
func FindProject(searchTerm string) (Project, error) {
	p, err := readProjectsFromFile()
	if err != nil {
		return Project{}, err
	}

	var foundProject Project
	for _, pro := range p.Projects {
		if strings.EqualFold(pro.Name, searchTerm) || pro.ProjectID == searchTerm {
			foundProject = pro
			break
		}
	}

	if foundProject == (Project{}) {
		return foundProject, errors.New("sorry, no project was found!\nYou can use the \"mole projects list\" command to see all projects")
	}

	return foundProject, nil
}

// saveProjectsToFile saves the current state of the Projects to mole.json.
func (p Projects) saveProjectsToFile() error {
	f, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal projects data: %w", err)
	}

	if err := os.WriteFile(getMoleJSONPath(), f, 0644); err != nil {
		return fmt.Errorf("failed to write projects to file: %w", err)
	}

	return nil
}

// addProject adds a new project to the list of projects and saves it to the file.
func addProject(newProject Project) error {

	p, err := readProjectsFromFile()
	if err != nil {
		return err
	}

	fileLock := flock.New(getMoleLockPath())
	locked, err := fileLock.TryLock()
	defer fileLock.Unlock()

	if err != nil {
		return err
	}

	if locked {
		newProject.ProjectID = shortuuid.New() // Assign a new unique project ID.
		p.Projects = append(p.Projects, newProject)

		if err := p.saveProjectsToFile(); err != nil {
			return err
		}

		return nil
	} else {
		return errors.New("someone else is trying to work with the project store, please try again")
	}
}

// cloneProject clones a project from a given repository URL into the local file system.
func cloneProject(project Project) error {
	if !consts.Testing {
		_, err := exec.LookPath("git")
		if err != nil {
			return errors.New("git is not installed or not available in PATH")
		}

		clonePath := path.Join(consts.GetBasePath(), "projects", project.Name)

		var stErr bytes.Buffer
		c := exec.Command("git", "clone", "--depth", "1", "-b", project.Branch, project.RepositoryURL, clonePath)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = &stErr

		if err := c.Run(); err != nil {
			return errors.New(stErr.String())
		}
		return nil
	} else {
		return nil
	}
}

// containsEnvEntry checks if the .gitignore content contains an entry for .env.
func containsEnvEntry(content string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == ".env" {
			return true
		}
	}
	return false
}

// projectSecrets holds the data for generating the environment configuration.
type projectSecrets struct {
	EnvPath    string
	RootPath   string
	LogPath    string
	PName      string
	AppKey     string
	PortApp    int
	PortTwo    int
	PortThree  int
	DbName     string
	DbUser     string
	DbPassword string
}

func createProjectSecretsJson(project Project) error {
	mp, err := FindAndReserveMolePorts()
	if err != nil {
		return err
	}

	key := helpers.GenerateRandomKey(32)
	dbName := project.Name + "db" + helpers.GenerateRandomKey(8)
	dbUser := project.Name + "user" + helpers.GenerateRandomKey(6)
	dbPass := helpers.GenerateRandomKey(24)

	be := projectSecrets{
		EnvPath:    "/home/mole/projects/" + project.Name + "/.env",
		RootPath:   "/home/mole/projects/" + project.Name,
		LogPath:    "/home/mole/logs/" + project.Name,
		PName:      project.Name,
		PortApp:    mp[0],
		PortTwo:    mp[1],
		PortThree:  mp[2],
		AppKey:     key,
		DbName:     dbName,
		DbUser:     dbUser,
		DbPassword: dbPass,
	}

	jbe, err := json.Marshal(be)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(consts.GetBasePath(), "secrets"), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(consts.GetBasePath(), "secrets", project.Name+".json"), jbe, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadProjectSecrets(projectNOI string) (*projectSecrets, error) {
	p, err := FindProject(projectNOI)
	if err != nil {
		return nil, err
	}

	s, err := os.ReadFile(path.Join(consts.GetBasePath(), "secrets", p.Name+".json"))
	if err != nil {
		return nil, err
	}

	var sec projectSecrets

	err = json.Unmarshal(s, &sec)
	if err != nil {
		return nil, err
	}

	return &sec, nil
}

// createProjectBaseEnv generates the base environment file for the project if .env.mole is present.
func createProjectBaseEnv(project Project) error {

	epath := path.Join(consts.BasePath, "projects", project.Name, ".env.mole")

	env, err := os.ReadFile(epath)
	if err != nil {
		return nil
	}

	if env != nil {
		var ft bytes.Buffer

		ft.Write(env)

		efp := path.Join(consts.GetBasePath(), "projects", project.Name, ".env")

		if err := os.WriteFile(efp, ft.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write environment file: %w", err)
		}
	}

	return nil
}

// mole.sh is mandatory for the project to be deployable by mole
func ensureMoleSh(newProject Project) error {

	shPath := path.Join(consts.BasePath, "projects", newProject.Name, "mole.sh")

	_, err := os.ReadFile(shPath)
	if err != nil {
		return fmt.Errorf("Project does not conatain or is unable to read mole.sh: %v", err)
	}

	return nil
}

// CreateProject creates a new project by cloning a repository and setting it up.
func CreateProject(newProject Project) error {

	clonePath := path.Join(consts.GetBasePath(), "projects", newProject.Name)

	if err := cloneProject(newProject); err != nil {
		return err
	}

	// ensure project hase the mole.sh
	if err := ensureMoleSh(newProject); err != nil {
		return err
	}

	if err := createProjectBaseEnv(newProject); err != nil {
		os.RemoveAll(clonePath) // Clean up on error
		return err
	}

	if err := createProjectLogDirectory(newProject); err != nil {
		os.RemoveAll(clonePath) // Clean up on error
		return err
	}

	if err := createProjectBaseEnv(newProject); err != nil {
		os.RemoveAll(clonePath) // Clean up on error
		return err
	}

	if err := createProjectSecretsJson(newProject); err != nil {
		os.RemoveAll(clonePath) // Clean up on error
		return err
	}

	if err := addProject(newProject); err != nil {
		os.RemoveAll(clonePath) // Clean up on error
		return err
	}

	return nil
}

func createProjectLogDirectory(p Project) error {
	return os.MkdirAll(path.Join(consts.GetBasePath(), "logs", p.Name), 0755)
}

// EditProject updates the details of an existing project by its name or ID.
func EditProject(proNOI, desc, branch string) error {
	p, err := readProjectsFromFile()
	if err != nil {
		return err
	}

	found := false

	for i, pro := range p.Projects {
		if proNOI == pro.ProjectID || proNOI == pro.Name {
			found = true
			if desc != "" {
				p.Projects[i].Description = desc
			}
			if branch != "" {
				p.Projects[i].Branch = branch
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("project with ID %s not found", proNOI)
	}

	return p.saveProjectsToFile()
}

// DeleteProject removes a project from the list by its ID.
func DeleteProject(proId string) error {
	p, err := readProjectsFromFile()
	if err != nil {
		return err
	}

	found := false

	for i, pro := range p.Projects {
		if proId == pro.ProjectID {
			found = true
			p.Projects = append(p.Projects[:i], p.Projects[i+1:]...) // Remove the project from the slice.
			break
		}
	}

	if !found {
		return fmt.Errorf("project with ID %s not found", proId)
	}

	return p.saveProjectsToFile()
}

// Stringify returns a string representation of the Project.
func (ps Project) Stringify() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(" |ID     : " + ps.ProjectID + "\n")
	b.WriteString(" |Name   : " + ps.Name + "\n")
	b.WriteString(" |Desc.  : " + ps.Description + "\n")
	b.WriteString(" |Git    : " + ps.RepositoryURL + "\n")
	b.WriteString(" |Branch : " + ps.Branch + "\n")
	return b.String()
}

// Stringify returns a string representation of all projects.
func (p Projects) Stringify() string {
	var b strings.Builder
	for i, pro := range p.Projects {
		b.WriteString("\n")
		b.WriteString(strconv.Itoa(i) + ":\n")
		b.WriteString(pro.Stringify())
	}
	return b.String()
}

// FindAndEditEnv opens the .env file for editing using nano.
func FindAndEditEnv(pName string) error {
	p, err := FindProject(pName)
	if err != nil {
		return err
	}

	c := exec.Command("nano", ".env")
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = path.Join(consts.GetBasePath(), "projects", p.Name)

	// Run the command and handle any error
	if err := c.Run(); err != nil {
		return fmt.Errorf("error running nano: %w", err)
	}

	return nil
}
