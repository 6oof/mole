package execs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/6oof/mole/pkg/consts"
	"github.com/6oof/mole/pkg/data"
)

func ListServices(pager bool) {
	cmd := exec.Command("sh", "-c", "systemctl --user list-units --type=service --all --no-legend --plain --no-pager | grep mole")

	if pager {
		cmd = exec.Command("sh", "-c", "systemctl --user list-units --type=service --all --no-legend --plain | grep mole | less")
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

func EnableService(serviceName string) error {
	cmd := exec.Command("systemctl", "--user", "enable", serviceName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to enable service %s: %v", serviceName, err)
	}
	return nil
}

func DisableService(serviceName string) error {
	cmd := exec.Command("systemctl", "--user", "disable", serviceName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to disable service %s: %v", serviceName, err)
	}
	return nil
}

func StartService(serviceName string) error {
	cmd := exec.Command("systemctl", "--user", "start", serviceName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to start service %s: %v", serviceName, err)
	}
	return nil
}

func StopService(serviceName string) error {
	cmd := exec.Command("systemctl", "--user", "stop", serviceName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to stop service %s: %v", serviceName, err)
	}
	return nil
}

func ReloadServices() error {
	cmd := exec.Command("systemctl", "--user", "daemon-reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to reload services: %v", err)
	}
	return nil
}

func RestartService(serviceName string) error {
	if err := ReloadServices(); err != nil {
		return fmt.Errorf("failed to reload services: %v", err)
	}

	cmd := exec.Command("systemctl", "--user", "restart", serviceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to restart service %s: %v", serviceName, err)
	}
	return nil
}

func RestartServiceHard(serviceName string) error {
	err := StopService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to restart service %s: %v", serviceName, err)
	}

	err = ReloadServices()
	if err != nil {
		return fmt.Errorf("failed to reload services after stopping %s: %v", serviceName, err)
	}

	err = StartService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to start service %s: %v", serviceName, err)
	}

	return nil
}

func LinkServices(projectNOI, sType string) error {
	p, err := data.FindProject(projectNOI)
	if err != nil {
		return err
	}

	sourceDir := path.Join(consts.BasePath, "projects", p.Name, "mole", "services")

	destDir := ""

	if sType == "systemd" {
		destDir = path.Join(consts.BasePath, ".config", "systemd", "user")
	} else if sType == "podman" {
		destDir = path.Join(consts.BasePath, ".config", "containers", "systemd")
	}

	if destDir == "" {
		return fmt.Errorf("invalid service type %s", sType)
	}

	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination directory %s: %v", destDir, err)
	}

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk through directory %s: %v", sourceDir, err)
		}

		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		destPath := filepath.Join(destDir, "mole-"+p.Name+"-"+info.Name())

		if _, err := os.Lstat(destPath); err == nil {
			fmt.Printf("Link already exists for %s, skipping.\n", destPath)
			return nil
		}

		err = os.Symlink(path, destPath)
		if err != nil {
			return fmt.Errorf("failed to create link %s -> %s: %v", destPath, path, err)
		}

		fmt.Printf("Linked %s to %s\n", path, destPath)
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path %s: %v", sourceDir, err)
	}

	return nil
}

func UnlinkServices(projectNOI string) error {
	p, err := data.FindProject(projectNOI)
	if err != nil {
		return err
	}

	sourceDirs := []string{
		path.Join(consts.BasePath, ".config", "systemd", "user"),
		path.Join(consts.BasePath, ".config", "containers", "systemd"),
	}

	for _, destDir := range sourceDirs {
		err = filepath.Walk(destDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to walk through directory %s: %v", destDir, err)
			}

			if info.IsDir() {
				return nil
			}

			if strings.Contains(info.Name(), "mole-"+p.Name) {
				err := os.Remove(path)
				if err != nil {
					return fmt.Errorf("failed to remove service %s: %v", path, err)
				}
				fmt.Printf("Removed service link %s\n", path)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking the path %s: %v", destDir, err)
		}
	}

	return nil
}
