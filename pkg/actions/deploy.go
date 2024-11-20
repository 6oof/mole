package actions

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/zulubit/mole/pkg/consts"
)

type projectDeployment struct {
	envVars     map[string]string
	projectName string
}

// TODO: check if transform is necessary for both

// RunDeployment executes the deployment process for a given project.
func RunDeployment(projectNOI string) (string, error) {
	ds, err := prepareDeployment(projectNOI)
	if err != nil {
		return "", err
	}

	err = gitPullProject(ds)
	if err != nil {
		return "", err
	}

	err = TransformCompose(projectNOI)
	if err != nil {
		return "", err
	}
	err = TransformDeploy(projectNOI)
	if err != nil {
		return "", err
	}

	succ, err := runDeploymentScript(ds.projectName)

	return succ, nil
}

// runDeploymentScript executes the mole-deploy-ready.sh script for the given project.
// It captures and returns the entire output (both stdout and stderr).
// The output is also written to a log file in the deploy_logs directory.
func runDeploymentScript(projectNOI string) (string, error) {
	p, err := FindProject(projectNOI)
	if err != nil {
		return "", fmt.Errorf("failed to find project: %w", err)
	}

	scriptPath := path.Join(consts.GetBasePath(), "projects", p.Name, "mole-deploy-ready.sh")
	logsDir := path.Join(consts.GetBasePath(), "deploy_logs")
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create logs directory: %w", err)
	}

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		logFile := path.Join(logsDir, fmt.Sprintf("%s-%s-failure.log", timestamp, p.Name))
		writeLog(logFile, "Deployment script not found")
		return "", fmt.Errorf("deployment script not found at %s", scriptPath)
	}

	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Dir = path.Join(consts.GetBasePath(), "projects", p.Name) // Set the working directory to the project folder

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	fmt.Println("Deploying, this might take a while...")

	err = cmd.Run()
	logFile := path.Join(logsDir, fmt.Sprintf("%s-%s-%s.log", timestamp, p.Name, status(err)))

	writeLog(logFile, output.String())

	if err != nil {
		return output.String(), fmt.Errorf("deployment script failed: %w", err)
	}

	return output.String(), nil
}

// status returns "success" or "failure" based on the error value.
func status(err error) string {
	if err != nil {
		return "failure"
	}
	return "success"
}

// writeLog writes the given content to the specified log file.
func writeLog(filePath, content string) {
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		fmt.Printf("Failed to write log file: %s, Error: %v\n", filePath, err)
	}
}

// gitPullProject pulls the latest changes from the Git repository.
func gitPullProject(project projectDeployment) error {
	var errOut, stOut bytes.Buffer

	cmd := exec.Command("git", "pull")
	cmd.Dir = path.Join(consts.GetBasePath(), "projects", project.projectName)
	cmd.Stderr = &errOut
	cmd.Stdout = &stOut

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git pull failed: %w, %s", err, errOut.String())
	}

	return nil
}

// prepareDeployment initializes the project deployment configuration.
func prepareDeployment(projectNOI string) (projectDeployment, error) {
	p, err := FindProject(projectNOI)
	if err != nil {
		return projectDeployment{}, err
	}

	projectEnv := path.Join(consts.GetBasePath(), "projects", p.Name, ".env")
	env, err := godotenv.Read(projectEnv)
	if err != nil {
		return projectDeployment{}, err
	}

	_, err = os.ReadFile(path.Join(consts.BasePath, "projects", p.Name, "mole-deploy-ready.sh"))
	if err != nil {
		return projectDeployment{}, err
	}

	dp := projectDeployment{
		envVars:     env,
		projectName: p.Name,
	}
	return dp, nil
}
