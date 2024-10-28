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

// RunDeployment executes the deployment process for a given project.
func RunDeployment(projectNOI string, restart bool) error {
	ds, err := prepareDeployment(projectNOI)
	if err != nil {
		return err
	}

	switch p := ds.projectType; p {
	case enums.Static:
		if err := deployStatic(ds); err != nil {
			return err
		}

	case enums.Podman, enums.Systemd:
		if err := deploySystemdOrPodman(ds, restart); err != nil {
			return err
		}

	case enums.Script:
		fmt.Println("script")
	}

	return nil
}

// deploySystemdOrPodman deploys the project using either Systemd or Podman.
func deploySystemdOrPodman(project projectDeployment, restart bool) error {
	if err := gitPullProject(project); err != nil {
		return err
	}

	if err := UnlinkServices(project.projectName); err != nil {
		return err
	}

	if err := LinkServices(project.projectName, project.projectType); err != nil {
		return err
	}

	if err := ReloadServicesDaemon(); err != nil {
		return err
	}

	for _, s := range project.services {
		sn := helpers.ServiceNameModifier(s, project.projectName)

		if err := EnableService(sn); err != nil {
			return err
		}
		if err := StartService(sn); err != nil {
			return err
		}

		if restart {
			if err := RestartService(sn); err != nil {
				return err
			}
		} else {
			if err := ReloadService(sn); err != nil {
				return err
			}
		}
	}

	return nil
}

// deployStatic deploys a static project.
func deployStatic(project projectDeployment) error {
	return gitPullProject(project) // Return the error directly for simplicity.
}

// gitPullProject pulls the latest changes from the Git repository.
func gitPullProject(project projectDeployment) error {
	var errOut, stOut bytes.Buffer

	cmd := exec.Command("git", "pull")
	cmd.Dir = path.Join(consts.BasePath, "projects", project.projectName)
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

	projectEnv := path.Join(consts.BasePath, "projects", p.Name, ".env")
	env, err := godotenv.Read(projectEnv)
	if err != nil {
		return projectDeployment{}, err
	}

	dp := projectDeployment{
		envVars:     env,
		projectName: p.Name,
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

	return dp, nil
}

// hasDefinedServices checks for defined services in the deployment project.
func hasDefinedServices(pro projectDeployment) (projectDeployment, error) {
	es := pro.envVars["MOLE_SERVICES"]
	services := strings.Split(es, ",")
	if len(services) == 0 {
		return pro, errors.New("no services specified in .env file. Please add at least one service to deploy")
	}

	pro.services = services
	return pro, nil
}
