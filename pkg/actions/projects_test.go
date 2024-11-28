package actions

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zulubit/mole/pkg/consts"
)

func TestAddReadFindProject(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	np := Project{
		Name: "test-project",
	}

	err := addProject(np)
	assert.Nil(t, err, "project was succesfully added")

	p, err := readProjectsFromFile()
	assert.Nil(t, err, "projects can be read")
	assert.Equal(t, p.Projects[0].Name, np.Name, "project just was read from file")

	fp, err := FindProject(np.Name)
	assert.Nil(t, err, "projects was found")
	assert.Equal(t, fp.Name, np.Name, "project just added was found")
}

func TestCreateProjectBaseEnv(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	np := Project{
		Name: "test-project",
	}

	err := addProject(np)
	assert.Nil(t, err, "project was successfully added")

	projectPath := path.Join(tmp, "projects", np.Name)
	err = os.MkdirAll(projectPath, 0755)
	assert.Nil(t, err, "project directory created")

	envMoleContent := "MOLE_TEST_KEY=test_value\n"
	err = os.WriteFile(path.Join(projectPath, ".env.example"), []byte(envMoleContent), 0644)
	assert.Nil(t, err, "example env.example file created")

	fp, err := FindProject(np.Name)
	assert.Nil(t, err, "project found successfully")

	err = createProjectBaseEnv(fp)
	assert.Nil(t, err, "base env created")

	envPath := path.Join(projectPath, ".env")
	envContent, err := os.ReadFile(envPath)
	assert.Nil(t, err, "env file can be read")

	assert.Contains(t, string(envContent), "MOLE_TEST_KEY=test_value", "env contains merged value from .env.example")
}

func TestCreateProjectSecretsJson(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	np := Project{
		Name: "test-project",
	}

	err := createProjectSecretsJson(np)
	assert.Nil(t, err, "project secrets JSON created successfully")

	// Verify the secrets JSON file
	secretsPath := path.Join(tmp, "secrets", np.Name+".json")
	secretsContent, err := os.ReadFile(secretsPath)
	assert.Nil(t, err, "secrets file can be read")

	// Parse the secrets JSON
	var secrets projectSecrets
	err = json.Unmarshal(secretsContent, &secrets)
	assert.Nil(t, err, "secrets JSON unmarshalled successfully")

	// Verify the contents
	assert.Equal(t, "/home/mole/projects/test-project/.env", secrets.EnvFilePath, "EnvPath is correct")
	assert.Equal(t, "/home/mole/projects/test-project", secrets.RootDirectory, "RootPath is correct")
	assert.Equal(t, "/home/mole/logs/test-project", secrets.LogDirectory, "LogPath is correct")
	assert.Equal(t, "test-project", secrets.ProjectName, "PName is correct")
	assert.NotEmpty(t, secrets.AppKey, "AppKey is generated")
	assert.NotEmpty(t, secrets.DatabaseName, "DbName is generated")
	assert.NotEmpty(t, secrets.DatabaseUser, "DbUser is generated")
	assert.NotEmpty(t, secrets.DatabasePass, "DbPassword is generated")
}

func TestEnsureMoleSh(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	np := Project{
		Name: "test-project",
	}

	addProject(np)

	projectPath := path.Join(tmp, "projects", np.Name)
	os.MkdirAll(projectPath, 0755)

	err := ensureMoleSh(np)
	assert.ErrorContains(t, err, "Project does not conatain or is unable to read mole.sh", "error is returned because mole.sh is missing")

	moleSh := "# silence"
	os.WriteFile(path.Join(projectPath, "mole.sh"), []byte(moleSh), 0644)
	err = ensureMoleSh(np)
	assert.Nil(t, err, "mole.sh check passes without error")

}

func TestDeleteProject(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	// Create and add a new project
	np := Project{
		Name: "test-project",
	}

	err := addProject(np)
	assert.Nil(t, err, "project was successfully added")

	// Create project directories and files
	projectPath := path.Join(tmp, "projects", np.Name)
	err = os.MkdirAll(projectPath, 0755)
	assert.Nil(t, err, "project directory created")

	logPath := path.Join(tmp, "logs", np.Name)
	err = os.MkdirAll(logPath, 0755)
	assert.Nil(t, err, "log directory created")

	domainFilePath := path.Join(tmp, "domains", np.Name+".caddy")
	err = os.MkdirAll(path.Dir(domainFilePath), 0755)
	assert.Nil(t, err, "domain directory created")
	err = os.WriteFile(domainFilePath, []byte("example caddy config"), 0644)
	assert.Nil(t, err, "domain file created")

	// Verify everything exists
	_, err = os.Stat(projectPath)
	assert.Nil(t, err, "project directory exists")
	_, err = os.Stat(logPath)
	assert.Nil(t, err, "log directory exists")
	_, err = os.Stat(domainFilePath)
	assert.Nil(t, err, "domain file exists")

	// Call DeleteProject
	createdProject, err := FindProject(np.Name)
	assert.Nil(t, err, "project was found")
	err = DeleteProject(createdProject.ProjectID)
	assert.Nil(t, err, "project deleted successfully")

	// Verify the project is removed
	p, err := readProjectsFromFile()
	assert.Nil(t, err, "projects can be read")
	for _, pro := range p.Projects {
		assert.NotEqual(t, pro.ProjectID, np.ProjectID, "deleted project is not in the list")
	}

	// Verify directories and files are deleted
	_, err = os.Stat(projectPath)
	assert.True(t, os.IsNotExist(err), "project directory is deleted")
	_, err = os.Stat(logPath)
	assert.True(t, os.IsNotExist(err), "log directory is deleted")
	_, err = os.Stat(domainFilePath)
	assert.True(t, os.IsNotExist(err), "domain file is deleted")
}
