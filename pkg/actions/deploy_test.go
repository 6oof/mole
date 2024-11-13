package actions

import (
	"os"
	"path"
	"testing"

	"github.com/6oof/mole/pkg/consts"
	"github.com/6oof/mole/pkg/enums"
	"github.com/stretchr/testify/assert"
)

var expectedDeployment = projectDeployment{
	projectType: enums.Podman,
	envVars: map[string]string{
		"MOLE_PROJECT_TYPE": "podman",
		"MOLE_VOLUME_PATH":  "/home/mole/volumes/test",
		"MOLE_SERVICES":     "app",
		"MOLE_APP_PORT":     "9000",
		"MOLE_TWO_PORT":     "9001",
		"MOLE_THREE_PORT":   "9002",
		"MOLE_APP_KEY":      "2NC0bCFYDVtlLW2c2AimhjmS2DsSL4vE",
	},
	services:    []string{"app"},
	projectName: "test-project",
}

func TestPrepareDeployment(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	np := Project{
		Name: "test-project",
	}

	addProject(np)

	_, err := prepareDeployment(np.Name)
	assert.NotNil(t, err, "deployment fails due to .env missing")

	envContent, err := os.ReadFile(path.Join("testing", "test-project", ".env"))
	if err != nil {
		t.Fatalf("failed to read .env file: %v", err)
	}

	err = os.MkdirAll(path.Join(tmp, "projects", np.Name), 0755)
	if err != nil {
		t.Fatalf("failed to create tmp directories: %v", err)
	}

	envPath := path.Join(tmp, "projects", np.Name, ".env")
	err = os.WriteFile(envPath, envContent, 0600)
	if err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	pd, err := prepareDeployment(np.Name)
	assert.Nil(t, err, "deployment should be prepared without error")
	assert.Equal(t, expectedDeployment, pd, "deployment struct should match expected values")
}
