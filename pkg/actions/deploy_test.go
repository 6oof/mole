package actions

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zulubit/mole/pkg/consts"
)

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
	_, err = RunDeployment(projectName)
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
	deployScriptPath := path.Join(projectDir, "mole-deploy.sh")
	err = os.WriteFile(deployScriptPath, []byte("#!/bin/bash\necho 'Deployment script running'"), 0755)
	if err != nil {
		t.Fatalf("failed to write mole-deploy-ready.sh: %v", err)
	}

	// Create a placeholder `mole-compose.yaml` file
	composeFilePath := path.Join(projectDir, "mole-compose.yaml")
	err = os.WriteFile(composeFilePath, []byte(""), 0600)
	if err != nil {
		t.Fatalf("failed to write mole-compose.yaml: %v", err)
	}

	// Test successful deployment
	output, err := RunDeployment(projectName)
	assert.Nil(t, err, "deployment should be prepared and executed without error")
	assert.Contains(t, output, "Deployment script running", "output should indicate that the deployment script ran")
}
