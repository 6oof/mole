package data

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

func AddDomainProxy(projectNOI, domain string, port int) error {
	if !helpers.ValidateCaddyDomain(domain) {
		return fmt.Errorf("error validating domain: %s", domain)
	}

	p, err := FindProject(projectNOI)
	if err != nil {
		return err
	}

	domainTemplate := `www.{{.Domain}} {
    redir https://{{.Domain}}{uri}
}

{{.Domain}} {
    reverse_proxy 127.0.0.1:{{.Port}}
}`

	dom := domainData{
		Domain: domain,
		Port:   strconv.Itoa(port),
	}

	tmpl, err := template.New("proxy").Parse(domainTemplate)
	if err != nil {
		return err
	}

	var ft bytes.Buffer

	err = tmpl.Execute(&ft, dom)
	if err != nil {
		return err
	}

	dfp := path.Join(consts.BasePath, "domains", p.Name+".caddy")
	pdir := path.Join(consts.BasePath, "domains")

	err = os.MkdirAll(pdir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(dfp, ft.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func AddDomainStatic(projectNOI, domain, location string) error {
	p, err := FindProject(projectNOI)
	if err != nil {
		return err
	}

	staticDomainTemplate := `www.{{.Domain}} {
    redir https://{{.Domain}}{uri}
}

{{.Domain}} {
    root * /home/projects/{{.ProjectName}}/{{.Location}}
    file_server
}`

	dom := domainData{
		Domain:      domain,
		Location:    location,
		ProjectName: p.Name,
	}

	tmpl, err := template.New("static").Parse(staticDomainTemplate)
	if err != nil {
		return err
	}

	var ft bytes.Buffer

	err = tmpl.Execute(&ft, dom)
	if err != nil {
		return err
	}

	dfp := path.Join(consts.BasePath, "domains", p.Name+".caddy")
	pdir := path.Join(consts.BasePath, "domains")

	err = os.MkdirAll(pdir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(dfp, ft.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

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

	ds := domainSetup{
		Email: email,
	}

	tmpl, err := template.New("setup").Parse(domainTemplate)
	if err != nil {
		return err
	}

	var ft bytes.Buffer

	err = tmpl.Execute(&ft, ds)
	if err != nil {
		return err
	}

	dfp := path.Join(consts.BasePath, "caddy", "main.caddy")
	pdir := path.Join(consts.BasePath, "caddy")

	err = os.MkdirAll(pdir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(dfp, ft.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func DeleteProjectDomain(projectName string) error {
	pd := path.Join(consts.BasePath, "domains", projectName+".caddy")

	err := os.Remove(pd)
	if err != nil {
		return err
	}

	return nil
}
