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
	err = os.WriteFile(path.Join(projectPath, ".env.mole"), []byte(envMoleContent), 0644)
	assert.Nil(t, err, "example env.mole file created")

	fp, err := FindProject(np.Name)
	assert.Nil(t, err, "project found successfully")

	err = createProjectBaseEnv(fp)
	assert.Nil(t, err, "base env created")

	envPath := path.Join(projectPath, ".env")
	envContent, err := os.ReadFile(envPath)
	assert.Nil(t, err, "env file can be read")

	assert.Contains(t, string(envContent), "MOLE_TEST_KEY=test_value", "env contains merged value from .env.mole")
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
	assert.Equal(t, "/home/mole/projects/test-project/.env", secrets.EnvPath, "EnvPath is correct")
	assert.Equal(t, "/home/mole/projects/test-project", secrets.RootPath, "RootPath is correct")
	assert.Equal(t, "/home/mole/logs/test-project", secrets.LogPath, "LogPath is correct")
	assert.Equal(t, "test-project", secrets.PName, "PName is correct")
	assert.NotEmpty(t, secrets.AppKey, "AppKey is generated")
	assert.NotEmpty(t, secrets.DbName, "DbName is generated")
	assert.NotEmpty(t, secrets.DbUser, "DbUser is generated")
	assert.NotEmpty(t, secrets.DbPassword, "DbPassword is generated")
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
