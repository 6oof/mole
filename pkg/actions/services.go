package actions

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/6oof/mole/pkg/consts"
	"github.com/6oof/mole/pkg/enums"
	"github.com/6oof/mole/pkg/helpers"
)

// ListServices lists all services containing the filter string.
func ListServices() error {
	filterString := "mole"

	conn, err := helpers.ContactDbus()
	if err != nil {
		return fmt.Errorf("failed to connect to DBus: %v", err)
	}
	defer conn.Close() // Ensure the connection is closed after the function execution.

	units, err := conn.ListUnitsContext(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list units: %v", err)
	}

	var output strings.Builder
	for _, unit := range units {
		if strings.Contains(unit.Name, filterString) {
			output.WriteString(fmt.Sprintf("%s - %s\n  LoadState: %s, ActiveState: %s\n",
				unit.Name, unit.Description, unit.LoadState, unit.ActiveState))
		}
	}

	if output.Len() > 0 {
		fmt.Print(output.String())
	} else {
		fmt.Println("No services matching " + filterString + " found.")
	}

	return nil
}

// EnableService enables a specified service.
func EnableService(serviceName string) error {
	conn, err := helpers.ContactDbus()
	if err != nil {
		return fmt.Errorf("failed to connect to DBus: %v", err)
	}
	defer conn.Close()

	_, _, err = conn.EnableUnitFilesContext(context.Background(), []string{serviceName}, false, true)
	if err != nil {
		return fmt.Errorf("failed to enable service %s: %v", serviceName, err)
	}

	return nil
}

// DisableService disables a specified service.
func DisableService(serviceName string) error {
	conn, err := helpers.ContactDbus()
	if err != nil {
		return fmt.Errorf("failed to connect to DBus: %v", err)
	}
	defer conn.Close()

	_, err = conn.DisableUnitFilesContext(context.Background(), []string{serviceName}, false)
	if err != nil {
		return fmt.Errorf("failed to disable service %s: %v", serviceName, err)
	}

	return nil
}

// StartService starts a specified service.
func StartService(serviceName string) error {
	conn, err := helpers.ContactDbus()
	if err != nil {
		return fmt.Errorf("failed to connect to DBus: %v", err)
	}
	defer conn.Close()

	_, err = conn.StartUnitContext(context.Background(), serviceName, "replace", nil)
	if err != nil {
		return fmt.Errorf("failed to send a start signal to %s: %v", serviceName, err)
	}

	return nil
}

// StopService stops a specified service.
func StopService(serviceName string) error {
	conn, err := helpers.ContactDbus()
	if err != nil {
		return fmt.Errorf("failed to connect to DBus: %v", err)
	}
	defer conn.Close()

	_, err = conn.StopUnitContext(context.Background(), serviceName, "replace", nil)
	if err != nil {
		return fmt.Errorf("failed to send a stop signal to %s: %v", serviceName, err)
	}

	return nil
}

// ReloadServicesDaemon reloads the systemd daemon.
func ReloadServicesDaemon() error {
	conn, err := helpers.ContactDbus()
	if err != nil {
		return fmt.Errorf("failed to connect to DBus: %v", err)
	}
	defer conn.Close()

	if err := conn.ReloadContext(context.Background()); err != nil {
		return fmt.Errorf("failed to reload systemd daemon: %v", err)
	}

	return nil
}

// ReloadService reloads a specified service.
func ReloadService(serviceName string) error {
	conn, err := helpers.ContactDbus()
	if err != nil {
		return fmt.Errorf("failed to connect to DBus: %v", err)
	}
	defer conn.Close()

	_, err = conn.ReloadUnitContext(context.Background(), serviceName, "replace", nil)
	if err != nil {
		return fmt.Errorf("failed to send a reload signal to %s: %v", serviceName, err)
	}

	return nil
}

// RestartService restarts a specified service.
func RestartService(serviceName string) error {
	conn, err := helpers.ContactDbus()
	if err != nil {
		return fmt.Errorf("failed to connect to DBus: %v", err)
	}
	defer conn.Close()

	_, err = conn.RestartUnitContext(context.Background(), serviceName, "replace", nil)
	if err != nil {
		return fmt.Errorf("failed to send a restart signal to %s: %v", serviceName, err)
	}

	return nil
}

// ensureDirsExist creates directories if they do not exist.
func ensureDirsExist(dirs []string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}
	return nil
}

// DisableStopAndUnlinkServices disables and stops services related to a project.
func DisableStopAndUnlinkServices(projectNOI string) error {
	p, err := FindProject(projectNOI)
	if err != nil {
		return err
	}

	searchString := "mole-" + p.Name

	conn, err := helpers.ContactDbus()
	if err != nil {
		return fmt.Errorf("failed to connect to DBus: %v", err)
	}
	defer conn.Close()

	units, err := conn.ListUnitsContext(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list units: %v", err)
	}

	for _, unit := range units {
		if strings.Contains(unit.Name, searchString) {
			if err := DisableService(unit.Name); err != nil {
				return err
			}
			if err := StopService(unit.Name); err != nil {
				return err
			}
		}
	}

	return UnlinkServices(projectNOI)
}

// LinkServices creates symbolic links for the services of a project.
func LinkServices(projectNOI string, sType enums.ProjectType) error {
	p, err := FindProject(projectNOI)
	if err != nil {
		return err
	}

	sourceDir := path.Join(consts.BasePath, "projects", p.Name, "mole_services")
	var destDir string

	switch sType {
	case enums.Systemd:
		destDir = path.Join(consts.BasePath, ".config", "systemd", "user")
	case enums.Podman:
		destDir = path.Join(consts.BasePath, ".config", "containers", "systemd")
	default:
		return fmt.Errorf("invalid service type %s or linking not necessary for the project of type %s", sType.String(), sType.String())
	}

	if err := ensureDirsExist([]string{destDir}); err != nil {
		return err
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

		if err := os.Symlink(path, destPath); err != nil {
			return fmt.Errorf("failed to create link %s -> %s: %v", destPath, path, err)
		}
		fmt.Printf("Linked %s to %s\n", path, destPath)

		dropInDir := filepath.Join(destDir, "mole-"+p.Name+"-"+info.Name()+".d")
		if err := ensureDirsExist([]string{dropInDir}); err != nil {
			return err
		}

		dropInFile := filepath.Join(dropInDir, "override.conf")
		dropInContent := fmt.Sprintf("[Service]\nEnvironment=\"MOLE_PROJECT_NAME=%s\"\n\n[Container]\nEnvironment=\"MOLE_PROJECT_NAME=%s\"\n\n[Build]\nEnvironment=\"MOLE_PROJECT_NAME=%s\"\n", p.Name, p.Name, p.Name)
		if err := os.WriteFile(dropInFile, []byte(dropInContent), 0644); err != nil {
			return fmt.Errorf("failed to create drop-in file %s: %v", dropInFile, err)
		}
		fmt.Printf("Created drop-in file %s\n", dropInFile)

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path %s: %v, ensure mole_services is present in your repository", sourceDir, err)
	}

	return nil
}

// UnlinkServices removes symbolic links for the services of a project.
func UnlinkServices(projectNOI string) error {
	p, err := FindProject(projectNOI)
	if err != nil {
		return err
	}

	sourceDirs := []string{
		path.Join(consts.BasePath, ".config", "systemd", "user"),
		path.Join(consts.BasePath, ".config", "containers", "systemd"),
	}

	if err := ensureDirsExist(sourceDirs); err != nil {
		return err
	}

	for _, destDir := range sourceDirs {
		err = filepath.Walk(destDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to walk through directory %s: %v", destDir, err)
			}

			if strings.Contains(info.Name(), "mole-"+p.Name) {
				if err := os.RemoveAll(path); err != nil {
					return fmt.Errorf("failed to remove service %s: %v", path, err)
				}
				fmt.Printf("Removed service link %s\n", path)

				if info.IsDir() {
					return filepath.SkipDir
				}
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking the path %s: %v", destDir, err)
		}
	}

	return nil
}
