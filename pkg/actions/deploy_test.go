package actions

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zulubit/mole/pkg/consts"
)

var expectedDeployment = projectDeployment{
	envVars: map[string]string{
		"MOLE_PROJECT_NAME": "test-project",
		"MOLE_ROOT_PATH":    "/home/mole/projects/test",
		"MOLE_SERVICES":     "app",
		"MOLE_PORT_APP":     "9000",
		"MOLE_PORT_TWO":     "9001",
		"MOLE_PORT_THREE":   "9002",
		"MOLE_APP_KEY":      "94yiGlTP1vhny0nvTuUmCw23EskugBdw",
		"MOLE_DB_NAME":      "testdbZRN7fP4b",
		"MOLE_DB_USER":      "testusermjpMrM",
		"MOLE_DB_PASS":      "nBFB8SmUGXDRbN0EDRDHDIJM",
	},
	projectName: "test-project",
}

func TestPrepareDeployment(t *testing.T) {
	consts.Testing = true

	// Create a temporary directory for testing
	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	// Set up a test project
	projectName := "test-project"
	np := Project{
		Name: projectName,
	}
	addProject(np)

	projectDir := path.Join(tmp, "projects", projectName)

	// Ensure the project directory exists
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("failed to create project directory: %v", err)
	}

	// Simulate missing .env file and check for failure
	_, err = prepareDeployment(projectName)
	assert.NotNil(t, err, "deployment should fail due to missing .env file")

	// Write a valid .env file
	envContent := `
MOLE_PROJECT_NAME=test-project
MOLE_ROOT_PATH=/home/mole/projects/test
MOLE_SERVICES=app
MOLE_PORT_APP=9000
MOLE_PORT_TWO=9001
MOLE_PORT_THREE=9002
MOLE_APP_KEY=94yiGlTP1vhny0nvTuUmCw23EskugBdw
MOLE_DB_NAME=testdbZRN7fP4b
MOLE_DB_USER=testusermjpMrM
MOLE_DB_PASS=nBFB8SmUGXDRbN0EDRDHDIJM
`
	envPath := path.Join(projectDir, ".env")
	err = os.WriteFile(envPath, []byte(envContent), 0600)
	if err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	// Create a placeholder `mole-deploy-ready.sh` file
	deployScriptPath := path.Join(projectDir, "mole-deploy-ready.sh")
	err = os.WriteFile(deployScriptPath, []byte("#!/bin/bash\necho 'Deployment script running'"), 0755)
	if err != nil {
		t.Fatalf("failed to write mole-deploy-ready.sh: %v", err)
	}

	// Test successful deployment preparation
	pd, err := prepareDeployment(projectName)
	assert.Nil(t, err, "deployment should be prepared without error")
	assert.Equal(t, expectedDeployment.projectName, pd.projectName, "project name should match")
	assert.Equal(t, expectedDeployment.envVars, pd.envVars, "environment variables should match")
}
