package actions

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/6oof/mole/pkg/consts"
	"github.com/6oof/mole/pkg/data"
	"github.com/6oof/mole/pkg/enums"
	"github.com/6oof/mole/pkg/helpers"
)

func ListServices() error {
	filterString := "mole"

	conn, err := helpers.ContactDbus()
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to DBus: %v", err)
		return err
	}

	units, err := conn.ListUnitsContext(context.Background())
	if err != nil {
		fmt.Printf("Failed to list units: %v", err)
		return err
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

func EnableService(serviceName string) error {
	// cmd := exec.Command("systemctl", "--user", "enable", serviceName)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	//
	// if err != nil {
	// 	return fmt.Errorf("failed to enable service %s: %v", serviceName, err)
	// }
	// return nil

	conn, err := helpers.ContactDbus()
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to DBus: %v", err)
		return err
	}

	_, _, err = conn.EnableUnitFilesContext(context.Background(), []string{serviceName}, false, true)
	if err != nil {
		fmt.Printf("Failed to enable service %s: %v", serviceName, err)
		return err
	}

	return nil
}

func DisableService(serviceName string) error {
	// cmd := exec.Command("systemctl", "--user", "disable", serviceName)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	//
	// if err != nil {
	// 	return fmt.Errorf("failed to disable service %s: %v", serviceName, err)
	// }

	conn, err := helpers.ContactDbus()
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to DBus: %v", err)
		return err
	}

	_, err = conn.DisableUnitFilesContext(context.Background(), []string{serviceName}, false)
	if err != nil {
		fmt.Printf("Failed to enable service %s: %v", serviceName, err)
		return err
	}

	return nil
}

func StartService(serviceName string) error {
	// cmd := exec.Command("systemctl", "--user", "start", serviceName)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	//
	// if err != nil {
	// 	return fmt.Errorf("failed to start service %s: %v", serviceName, err)
	// }

	conn, err := helpers.ContactDbus()
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to DBus: %v", err)
		return err
	}

	_, err = conn.StartUnitContext(context.Background(), serviceName, "replace", nil)
	if err != nil {
		fmt.Printf("Failed to send a start signal to %s: %v", serviceName, err)
		return err
	}

	return nil
}

func StopService(serviceName string) error {
	// cmd := exec.Command("systemctl", "--user", "stop", serviceName)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	//
	// if err != nil {
	// 	return fmt.Errorf("failed to stop service %s: %v", serviceName, err)
	// }

	conn, err := helpers.ContactDbus()
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to DBus: %v", err)
		return err
	}

	_, err = conn.StopUnitContext(context.Background(), serviceName, "replace", nil)
	if err != nil {
		fmt.Printf("Failed to send a stop signal to %s: %v", serviceName, err)
		return err
	}

	return nil
}

func ReloadServicesDaemon() error {
	// cmd := exec.Command("systemctl", "--user", "daemon-reload")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	//
	// if err != nil {
	// 	return fmt.Errorf("failed to reload services: %v", err)
	// }

	conn, err := helpers.ContactDbus()
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to DBus: %v", err)
		return err
	}

	err = conn.ReloadContext(context.Background())
	if err != nil {
		fmt.Printf("Failed to reload systemd daemon: %v", err)
		return err
	}

	return nil
}

func ReloadService(serviceName string) error {
	// if err := ReloadServices(); err != nil {
	// 	return fmt.Errorf("failed to reload services: %v", err)
	// }
	//
	// cmd := exec.Command("systemctl", "--user", "restart", serviceName)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	//
	// if err != nil {
	// 	return fmt.Errorf("failed to restart service %s: %v", serviceName, err)
	// }

	conn, err := helpers.ContactDbus()
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to DBus: %v", err)
		return err
	}

	_, err = conn.ReloadUnitContext(context.Background(), serviceName, "replace", nil)
	if err != nil {
		fmt.Printf("Failed to send a reload signal to %s: %v", serviceName, err)
		return err
	}

	return nil
}

func RestartService(serviceName string) error {
	conn, err := helpers.ContactDbus()
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to DBus: %v", err)
		return err
	}

	_, err = conn.RestartUnitContext(context.Background(), serviceName, "replace", nil)
	if err != nil {
		fmt.Printf("Failed to send a restart signal to %s: %v", serviceName, err)
		return err
	}

	return nil
}

func LinkServices(projectNOI string, sType enums.ProjectType) error {
	p, err := data.FindProject(projectNOI)
	if err != nil {
		return err
	}

	sourceDir := path.Join(consts.BasePath, "projects", p.Name, "mole", "services")

	destDir := ""

	if sType == enums.Systemd {
		destDir = path.Join(consts.BasePath, ".config", "systemd", "user")
	} else if sType == enums.Podman {
		destDir = path.Join(consts.BasePath, ".config", "containers", "systemd")
	}

	if destDir == "" {
		return fmt.Errorf("invalid service type %s or linking not necessary for the project of type %s", sType.String(), sType.String())
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

		dropInDir := filepath.Join(destDir, "mole-"+p.Name+"-"+info.Name()+".d")
		err = os.MkdirAll(dropInDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create drop-in directory %s: %v", dropInDir, err)
		}

		dropInFile := filepath.Join(dropInDir, "override.conf")
		dropInContent := fmt.Sprintf("[Service]\nEnvironment=\"MOLE_PROJECT_NAME=%s\"\n\n[Container]\nEnvironment=\"MOLE_PROJECT_NAME=%s\"\n\n[Build]\nEnvironment=\"MOLE_PROJECT_NAME=%s\"\n", p.Name, p.Name, p.Name)
		err = os.WriteFile(dropInFile, []byte(dropInContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to create drop-in file %s: %v", dropInFile, err)
		}
		fmt.Printf("Created drop-in file %s\n", dropInFile)

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

			if strings.Contains(info.Name(), "mole-"+p.Name) {
				err := os.RemoveAll(path)
				if err != nil {
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
