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

	// Setup project secrets
	err := createProjectSecretsJson(np)
	assert.Nil(t, err, "Failed to setup project secrets")

	projectDir := path.Join(tmp, "projects", projectName)

	// Ensure the project directory exists
	err = os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("failed to create project directory: %v", err)
	}

	// Simulate missing .env file and check for failure
	_, err = RunDeployment(projectName)
	assert.NotNil(t, err, "deployment should fail due to missing .env file")

	// Create a placeholder `mole-ready.sh` file
	deployScriptPath := path.Join(projectDir, "mole.sh")
	err = os.WriteFile(deployScriptPath, []byte("#!/bin/bash\necho 'Deployment script running'"), 0755)
	if err != nil {
		t.Fatalf("failed to write mole-ready.sh: %v", err)
	}

	// Test successful deployment
	output, err := RunDeployment(projectName)
	assert.Nil(t, err, "deployment should be prepared and executed without error")
	assert.Contains(t, output, "Deployment script running", "output should indicate that the deployment script ran")
}
