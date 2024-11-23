package actions

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/zulubit/mole/pkg/consts"
)

// TransformCompose generates "mole-compose-ready.yaml" by transforming "mole-compose.yaml"
// using secrets from the project's secrets file.
func TransformCompose(projectNOI string) error {
	p, err := FindProject(projectNOI)
	if err != nil {
		return err
	}

	sourcePath := path.Join(consts.GetBasePath(), "projects", p.Name, "mole-compose.yaml")
	destPath := path.Join(consts.GetBasePath(), "projects", p.Name, "mole-compose-ready.yaml")

	return injectSecrets(sourcePath, destPath, p.Name)
}

// TransformDeploy generates "mole-deploy-ready.sh" by transforming "mole-deploy.sh"
// using secrets from the project's secrets file.
func TransformDeploy(projectNOI string) error {
	p, err := FindProject(projectNOI)
	if err != nil {
		return err
	}

	sourcePath := path.Join(consts.GetBasePath(), "projects", p.Name, "mole-deploy.sh")
	destPath := path.Join(consts.GetBasePath(), "projects", p.Name, "mole-deploy-ready.sh")

	return injectSecrets(sourcePath, destPath, p.Name)
}

// TODO: optimize this method to not look for project twice
// injectSecrets reads the secrets JSON file and injects its values into a template.
func injectSecrets(sourcePath, destPath, projectName string) error {
	// Read the project secrets
	secrets, err := readProjectSecrets(projectName)
	if err != nil {
		return fmt.Errorf("failed to load secrets for project %s: %v", projectName, err)
	}

	// Read the source template file
	templateContent, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %v", sourcePath, err)
	}

	// Convert projectSecrets to a map for template rendering
	secretsMap, err := secretsToMap(secrets)
	if err != nil {
		return fmt.Errorf("failed to convert secrets to map: %v", err)
	}

	// Render the template with secrets
	renderedContent, err := renderTemplate(string(templateContent), secretsMap)
	if err != nil {
		return err
	}

	// Write the rendered content to the destination file
	err = os.WriteFile(destPath, []byte(renderedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	return nil
}

// renderTemplate renders a Go template with secrets
func renderTemplate(templateText string, data map[string]string) (string, error) {
	tmpl, err := template.New("unit").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var renderedContent strings.Builder
	err = tmpl.Execute(&renderedContent, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return renderedContent.String(), nil
}

// readProjectSecrets reads and unmarshals the secrets JSON for a given project
func readProjectSecrets(projectName string) (*projectSecrets, error) {
	secretsPath := path.Join(consts.GetBasePath(), "secrets", projectName+".json")

	data, err := os.ReadFile(secretsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets file: %v", err)
	}

	var secrets projectSecrets
	err = json.Unmarshal(data, &secrets)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secrets JSON: %v", err)
	}

	return &secrets, nil
}

// secretsToMap converts projectSecrets to a map for template rendering
func secretsToMap(secrets *projectSecrets) (map[string]string, error) {
	// Use reflection or manual mapping to convert struct fields to a map
	return map[string]string{
		"EnvPath":    secrets.EnvPath,
		"RootPath":   secrets.RootPath,
		"LogPath":    secrets.LogPath,
		"PName":      secrets.PName,
		"AppKey":     secrets.AppKey,
		"PortApp":    fmt.Sprintf("%d", secrets.PortApp),
		"PortTwo":    fmt.Sprintf("%d", secrets.PortTwo),
		"PortThree":  fmt.Sprintf("%d", secrets.PortThree),
		"DbName":     secrets.DbName,
		"DbUser":     secrets.DbUser,
		"DbPassword": secrets.DbPassword,
	}, nil
}
