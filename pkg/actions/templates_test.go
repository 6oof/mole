package actions

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zulubit/mole/pkg/consts"
)

func TestTransformTemplates(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	projectName := "test-project"
	np := Project{
		Name: projectName,
	}
	addProject(np)

	// Create necessary directories and files for the test project
	projectDir := path.Join(tmp, "projects", projectName)
	os.MkdirAll(projectDir, 0755)

	sourceCompose := path.Join(projectDir, "mole-compose.yaml")
	sourceDeploy := path.Join(projectDir, "mole.sh")

	err := os.WriteFile(sourceCompose, []byte("service_name: {{.PName}}"), 0644)
	assert.Nil(t, err, "Failed to create mole-compose.yaml")

	err = os.WriteFile(sourceDeploy, []byte("#!/bin/bash\necho {{.PName}}"), 0755)
	assert.Nil(t, err, "Failed to create mole.sh")

	// Setup project secrets
	err = createProjectSecretsJson(np)
	assert.Nil(t, err, "Failed to setup project secrets")

	// Test TransformCompose
	err = TransformCompose(projectName)
	assert.Nil(t, err, "TransformCompose should complete without error")

	destCompose := path.Join(projectDir, "mole-compose-ready.yaml")
	assert.FileExists(t, destCompose, "mole-compose-ready.yaml should be generated")

	content, err := os.ReadFile(destCompose)
	assert.Nil(t, err, "Should be able to read generated mole-compose-ready.yaml")
	assert.Contains(t, string(content), "service_name: test-project", "Transformed mole-compose.yaml should contain injected variables")

	// Test TransformDeploy
	err = TransformDeploy(projectName)
	assert.Nil(t, err, "TransformDeploy should complete without error")

	destDeploy := path.Join(projectDir, "mole-ready.sh")
	assert.FileExists(t, destDeploy, "mole-ready.sh should be generated")

	content, err = os.ReadFile(destDeploy)
	assert.Nil(t, err, "Should be able to read generated mole-ready.sh")
	assert.Contains(t, string(content), "echo test-project", "Transformed mole.sh should contain injected variables")
}
