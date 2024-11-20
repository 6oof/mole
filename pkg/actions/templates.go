package actions

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/zulubit/mole/pkg/consts"
)

// transform mole-compose.ymal to mole-compose-ready.yaml
func TransformCompose(projectNOI string) error {
	p, err := FindProject(projectNOI)
	if err != nil {
		return err
	}

	sourcePath := path.Join(consts.GetBasePath(), "projects", p.Name, "mole-compose.yaml")
	destPath := path.Join(consts.GetBasePath(), "projects", p.Name, "mole-compose-ready.yaml")
	envFilePath := path.Join(consts.GetBasePath(), "projects", p.Name, ".env")

	return injectEnv(sourcePath, destPath, envFilePath)

}

// TransformDeploy generates "mole-deploy-ready.sh" by transforming "mole-deploy.sh"
// using environment variables from the .env file.
func TransformDeploy(projectNOI string) error {
	p, err := FindProject(projectNOI)
	if err != nil {
		return err
	}

	sourcePath := path.Join(consts.GetBasePath(), "projects", p.Name, "mole-deploy.sh")
	destPath := path.Join(consts.GetBasePath(), "projects", p.Name, "mole-deploy-ready.sh")
	envFilePath := path.Join(consts.GetBasePath(), "projects", p.Name, ".env")

	return injectEnv(sourcePath, destPath, envFilePath)
}

// LinkServices creates symbolic links for the services of a project.
func injectEnv(sourcePath, destPath, envFilePath string) error {
	env, err := godotenv.Read(envFilePath)
	if err != nil {
		return fmt.Errorf("failed to load environment variables from %s: %v", envFilePath, err)
	}
	templateContent, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %v", sourcePath, err)
	}

	renderedContent, err := renderTemplate(string(templateContent), env)
	if err != nil {
		return err
	}

	err = os.WriteFile(destPath, []byte(renderedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	return nil
}

// RenderTemplate renders a Go template with environment variables
func renderTemplate(templateText string, env map[string]string) (string, error) {
	tmpl, err := template.New("unit").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var renderedContent strings.Builder
	err = tmpl.Execute(&renderedContent, env)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return renderedContent.String(), nil
}
