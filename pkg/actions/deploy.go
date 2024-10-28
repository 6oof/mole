package actions

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/6oof/mole/pkg/consts"
	"github.com/6oof/mole/pkg/enums"
	"github.com/6oof/mole/pkg/helpers"
	"github.com/joho/godotenv"
)

type projectDeployment struct {
	projectType enums.ProjectType
	envVars     map[string]string
	services    []string
	projectName string
}

func RunDeployment(projectNOI string, restart bool) error {
	ds, err := prepareDeployment(projectNOI)
	if err != nil {
		return err
	}

	switch p := ds.projectType; p {
	case enums.Static:
		if err := deplyStatic(ds); err != nil {
			return err
		}

	case enums.Podman:
		if err := deploySystemdOrPodman(ds, restart); err != nil {
			return err
		}

	case enums.Systemd:
		if err := deploySystemdOrPodman(ds, restart); err != nil {
			return err
		}

	case enums.Script:
		fmt.Println("script")
	}

	return nil
}

func deploySystemdOrPodman(project projectDeployment, restart bool) error {
	err := gitPullProject(project)
	if err != nil {
		return err
	}

	err = UnlinkServices(project.projectName)
	if err != nil {
		return err
	}

	err = LinkServices(project.projectName, project.projectType)
	if err != nil {
		return err
	}

	err = ReloadServicesDaemon()
	if err != nil {
		return err
	}

	for _, s := range project.services {
		sn := helpers.ServiceNameModifier(s, project.projectName)
		err := EnableService(sn)
		if err != nil {
			return err
		}
		err = StartService(sn)
		if err != nil {
			return err
		}

		if !restart {
			err = ReloadService(sn)
			if err != nil {
				return err
			}
		} else {
			err = RestartService(sn)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func deplyStatic(project projectDeployment) error {
	err := gitPullProject(project)
	if err != nil {
		return err
	}

	return nil
}

// TODO: probably shouldn't fire this when testing. A flag should be included to run the app in testing mode
func gitPullProject(project projectDeployment) error {
	var errOut, stOut bytes.Buffer

	cmd := exec.Command("git", "pull")
	cmd.Dir = (path.Join(consts.BasePath, "projects", project.projectName))
	cmd.Stderr = &errOut
	cmd.Stdout = &stOut

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Git pull failed: %w, %s", err, errOut.String())
	}

	return nil
}

func prepareDeployment(projectNOI string) (projectDeployment, error) {
	p, err := FindProject(projectNOI)
	if err != nil {
		return projectDeployment{}, err
	}

	projectEnv := path.Join(consts.BasePath, "projects", p.Name, ".env")

	env, err := godotenv.Read(projectEnv)
	if err != nil {
		return projectDeployment{}, err
	}

	dp := projectDeployment{
		envVars: env,
	}

	pts := env["MOLE_PROJECT_TYPE"]

	pt, err := enums.IsProjectType(pts)
	if err != nil {
		return projectDeployment{}, err
	}

	dp.projectType = pt

	if dp.projectType != enums.Static {
		dp, err = hasDefinedServices(dp)
		if err != nil {
			return projectDeployment{}, err
		}
	}

	dp.projectName = p.Name

	return dp, nil
}

func hasDefinedServices(pro projectDeployment) (projectDeployment, error) {

	es := pro.envVars["MOLE_SERVICES"]
	services := strings.Split(es, ",")
	if len(services) < 0 {
		return pro, errors.New("No services specified in .env file. Please add at least one service to deploy")
	}

	pro.services = services

	return pro, nil

}
