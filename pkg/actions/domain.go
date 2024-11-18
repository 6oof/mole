package actions

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"text/template"

	"github.com/6oof/mole/pkg/consts"
	"github.com/6oof/mole/pkg/helpers"
	"github.com/joho/godotenv"
)

type domainData struct {
	Domain      string
	Port        string
	Location    string
	ProjectName string
}

type domainSetup struct {
	Email string
}

// AddDomainProxy generates and adds a reverse proxy configuration for the specified domain and port.
func AddDomainProxy(projectNOI, domain string, port int) error {
	if !helpers.ValidateCaddyDomain(domain) {
		return fmt.Errorf("invalid domain format: %s", domain)
	}

	project, err := FindProject(projectNOI)
	if err != nil {
		return fmt.Errorf("failed to find project %s: %w", projectNOI, err)
	}

	var templatePort int
	if port == 0 {
		templatePort = readDefaultProtFromEnv(projectNOI)
	} else {
		templatePort = port
	}

	domainTemplate := `www.{{.Domain}} {
    redir https://{{.Domain}}{uri}
}

{{.Domain}} {
    reverse_proxy 127.0.0.1:{{.Port}}
}`

	domainData := domainData{
		Domain: domain,
		Port:   strconv.Itoa(templatePort),
	}

	templateInstance, err := template.New("proxy").Parse(domainTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse proxy template: %w", err)
	}

	var configBuffer bytes.Buffer
	if err := templateInstance.Execute(&configBuffer, domainData); err != nil {
		return fmt.Errorf("failed to execute template for domain %s: %w", domain, err)
	}

	domainFilePath := path.Join(consts.GetBasePath(), "domains", project.Name+".caddy")
	domainDirPath := path.Join(consts.GetBasePath(), "domains")

	if err := os.MkdirAll(domainDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", domainDirPath, err)
	}

	if err := os.WriteFile(domainFilePath, configBuffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write domain configuration file %s: %w", domainFilePath, err)
	}

	return nil
}

func readDefaultProtFromEnv(projectNOI string) int {
	p, err := FindProject(projectNOI)
	if err != nil {
		return 0
	}

	projectEnv := path.Join(consts.GetBasePath(), "projects", p.Name, ".env")
	env, err := godotenv.Read(projectEnv)
	if err != nil {
		return 0
	}

	port := env["MOLE_PORT_APP"]
	i, err := strconv.Atoi(port)
	if err != nil {
		return 0
	}
	return i
}

// AddDomainStatic creates a static file server configuration for a specified domain and location.
func AddDomainStatic(projectNOI, domain, location string) error {
	if !helpers.ValidateCaddyDomain(domain) {
		return fmt.Errorf("invalid domain format: %s", domain)
	}

	project, err := FindProject(projectNOI)
	if err != nil {
		return fmt.Errorf("failed to find project %s: %w", projectNOI, err)
	}

	staticDomainTemplate := `www.{{.Domain}} {
    redir https://{{.Domain}}{uri}
}

{{.Domain}} {
    root * /home/projects/{{.ProjectName}}/{{.Location}}
    file_server
}`

	domainData := domainData{
		Domain:      domain,
		Location:    location,
		ProjectName: project.Name,
	}

	templateInstance, err := template.New("static").Parse(staticDomainTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse static template: %w", err)
	}

	var configBuffer bytes.Buffer
	if err := templateInstance.Execute(&configBuffer, domainData); err != nil {
		return fmt.Errorf("failed to execute template for static domain %s: %w", domain, err)
	}

	domainFilePath := path.Join(consts.GetBasePath(), "domains", project.Name+".caddy")
	domainDirPath := path.Join(consts.GetBasePath(), "domains")

	if err := os.MkdirAll(domainDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", domainDirPath, err)
	}

	if err := os.WriteFile(domainFilePath, configBuffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write static domain configuration file %s: %w", domainFilePath, err)
	}

	return nil
}

// SetupDomains initializes the main Caddy configuration, enabling domain support with TLS.
func SetupDomains(email string) error {
	if !helpers.ValidateEmail(email) {
		return errors.New("invalid email provided")
	}

	domainTemplate := `{
    email {{.Email}}
    servers {
        protocol {
            experimental_http3
        }
    }
}

tls {
    on_demand
}

header {
    Accept-Encoding gzip, br
    Content-Type * gzip
    Content-Type * brotli
}

import /home/mole/domains/*.caddy`

	domainSetupData := domainSetup{Email: email}

	templateInstance, err := template.New("setup").Parse(domainTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse setup template: %w", err)
	}

	var configBuffer bytes.Buffer
	if err := templateInstance.Execute(&configBuffer, domainSetupData); err != nil {
		return fmt.Errorf("failed to execute setup template for email %s: %w", email, err)
	}

	caddyFilePath := path.Join(consts.GetBasePath(), "caddy", "main.caddy")
	caddyDirPath := path.Join(consts.GetBasePath(), "caddy")

	if err := os.MkdirAll(caddyDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", caddyDirPath, err)
	}

	if err := os.WriteFile(caddyFilePath, configBuffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write caddy configuration file %s: %w", caddyFilePath, err)
	}

	return nil
}

// DeleteProjectDomain removes the Caddy configuration file for the specified project domain.
func DeleteProjectDomain(projectName string) error {
	domainFilePath := path.Join(consts.GetBasePath(), "domains", projectName+".caddy")

	if err := os.Remove(domainFilePath); err != nil {
		return fmt.Errorf("failed to delete project domain configuration %s: %w", domainFilePath, err)
	}

	return nil
}
