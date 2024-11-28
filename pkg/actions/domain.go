package actions

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/zulubit/mole/pkg/consts"
	"github.com/zulubit/mole/pkg/helpers"
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
    root * /home/mole/projects/{{.ProjectName}}/{{.Location}}
    file_server

    encode gzip zstd

    @htmlFiles {
        file {
            try_files {path}.html
        }
    }

    @blockedFiles {
        path *.env
    }
    respond @blockedFiles 403

    rewrite @htmlFiles {path}.html
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
}
`

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

	domainDirPath := path.Join(consts.GetBasePath(), "domains")

	if err := os.MkdirAll(domainDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", domainDirPath, err)
	}

	domainFilePath := path.Join(consts.GetBasePath(), "domains", "empty.caddy")
	if err := os.WriteFile(domainFilePath, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to write static domain configuration file %s: %w", domainFilePath, err)
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

// ReloadCaddy reads the main Caddyfile and partials, consolidates them, and sends to the API.
func ReloadCaddy() error {
	mainFilePath := path.Join(consts.GetBasePath(), "caddy", "main.caddy")
	domainsDir := path.Join(consts.GetBasePath(), "domains")
	apiURL := "http://localhost:2019"

	var caddyfileBuilder strings.Builder

	// Read the main Caddyfile
	mainCaddyContent, err := os.ReadFile(mainFilePath)
	if err != nil {
		return fmt.Errorf("failed to read main Caddyfile %s: %w", mainFilePath, err)
	}
	caddyfileBuilder.Write(mainCaddyContent)
	caddyfileBuilder.WriteString("\n\n")

	// Read all partial Caddyfiles in the domains directory
	err = filepath.Walk(domainsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".caddy") {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return fmt.Errorf("failed to read Caddyfile fragment %s: %w", path, readErr)
			}
			caddyfileBuilder.Write(content)
			caddyfileBuilder.WriteString("\n\n")
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to gather Caddyfiles: %w", err)
	}

	// Send the consolidated Caddyfile to the API
	resp, err := http.Post(fmt.Sprintf("%s/load", apiURL), "text/caddyfile", bytes.NewBufferString(caddyfileBuilder.String()))
	if err != nil {
		return fmt.Errorf("failed to send consolidated Caddyfile to Caddy API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Caddy API returned status: %s\nDetails: %s", resp.Status, string(body))
	}

	return nil
}
