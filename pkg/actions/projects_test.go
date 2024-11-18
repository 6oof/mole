package actions

import (
	"os"
	"path"
	"testing"

	"github.com/6oof/mole/pkg/consts"
	"github.com/6oof/mole/pkg/enums"
	"github.com/stretchr/testify/assert"
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

func TestCheckEnvGitignore(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	np := Project{
		Name: "test-project",
	}

	addProject(np)

	os.MkdirAll(path.Join(tmp, "projects", np.Name), 0755)

	fp, _ := FindProject(np.Name)
	err := checkEnvGitignore(fp)
	assert.ErrorContains(t, err, "gitignore is missing from this project.", "error is thrown because .gitignore is missing")

	os.WriteFile(path.Join(tmp, "projects", np.Name, ".gitignore"), []byte("asdf"), 0755)
	err = checkEnvGitignore(fp)
	assert.ErrorContains(t, err, "does not include an entry for '.env'", "error is thrown because .gitignore is missing")

	os.WriteFile(path.Join(tmp, "projects", np.Name, ".gitignore"), []byte(".env"), 0755)
	err = checkEnvGitignore(fp)
	assert.Nil(t, err, "nothing is returned")
}

func TestCreateProjectBaseEnv(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	np := Project{
		Name: "test-project",
	}

	addProject(np)

	os.MkdirAll(path.Join(tmp, "projects", np.Name), 0755)

	fp, _ := FindProject(np.Name)
	err := createProjectBaseEnv(fp, enums.Podman)
	assert.Nil(t, err, "base env created")

	f, err := os.ReadFile(path.Join(tmp, "projects", np.Name, ".env"))
	assert.Nil(t, err, "env can be read")
	assert.Contains(t, string(f), "MOLE_PROJECT_TYPE=podman", "env contains the correct type")

}

func TestCreateProjectBaseEnvWithMerge(t *testing.T) {
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

	// Create the .env.mole file with a test entry
	envMoleContent := "MOLE_TEST_KEY=test_value"
	err := os.WriteFile(path.Join(projectPath, ".env.mole"), []byte(envMoleContent), 0644)
	assert.Nil(t, err, "example env file created")

	fp, _ := FindProject(np.Name)
	err = createProjectBaseEnv(fp, enums.Podman)
	assert.Nil(t, err, "base env created")

	f, err := os.ReadFile(path.Join(projectPath, ".env"))
	assert.Nil(t, err, "env can be read")

	// Check that the .env contains both the generated and merged values
	assert.Contains(t, string(f), "MOLE_PROJECT_TYPE=podman", "env contains the correct type")
	assert.Contains(t, string(f), "MOLE_TEST_KEY=test_value", "env contains the merged value from .env.mole")
}
