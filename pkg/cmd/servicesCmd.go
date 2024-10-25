package cmd

import (
	"fmt"
	"strings"

	"github.com/6oof/mole/pkg/execs"
	"github.com/spf13/cobra"
)

// TODO: I like it when the execs are in a separate file and just a single function is in the body of the cmd.
// TODO: We should probably split execs and "actions" they perform
// TODO: add instructions and validation for types
// TODO: Errors and prints should be handled at the exec level so we can double up later for json
var (
	pType                                                    string
	serviceStart, serviceStop, serviceEnable, serviceDisable bool
	hardRerload                                              bool
)

func init() {
	RootCmd.AddCommand(servicesRootCmd)

	servicesRootCmd.AddCommand(reloadServicesCmd)
	servicesRootCmd.AddCommand(unlinkProjectServicesCmd)

	servicesRootCmd.AddCommand(listServicesCmd)

	servicesRootCmd.AddCommand(linkProjectServicesCmd)
	linkProjectServicesCmd.Flags().StringVarP(&pType, "type", "t", "", "Type of services to be linked")
	linkProjectServicesCmd.MarkFlagRequired("type")

	servicesRootCmd.AddCommand(serviceActionCmd)
	serviceActionCmd.Flags().BoolVarP(&serviceStart, "start", "", false, "Executes service start command")
	serviceActionCmd.Flags().BoolVarP(&serviceStop, "stop", "", false, "Executes service stop command")
	serviceActionCmd.Flags().BoolVarP(&serviceEnable, "enable", "", false, "Executes service enable command")
	serviceActionCmd.Flags().BoolVarP(&serviceDisable, "disable", "", false, "Executes service disable command")

	servicesRootCmd.AddCommand(restartServicesCmd)
	restartServicesCmd.Flags().BoolVarP(&hardRerload, "full", "f", false, "Stop and start the service insted of reloading")
}

var servicesRootCmd = &cobra.Command{
	Use:   "services",
	Short: "Interact with systemd services",
	Long:  `Interact with systemd services...`,
}

var listServicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List services",
	Long: `List services is for listing all services.
	It only lists the services that are marked as "mole" services.
	This marking happesn automatically when a project is managed by mole.`,
	Run: func(cmd *cobra.Command, args []string) {
		execs.ListServices()
	},
}

var reloadServicesCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload service unit files",
	Long:  `Reload service unit files. This will register any unit file changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		execs.ReloadServicesDaemon()
	},
}

// TODO: redo this to just be a string name of the action and handle it in another file
var serviceActionCmd = &cobra.Command{
	Use:   "action [service name]",
	Short: "Action is used to start/stop/enable/disable services",
	Long:  `Action is for starting, stopping, enabling, disabling services.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !serviceDisable && !serviceEnable && !serviceStart && !serviceStop {
			fmt.Println("at least one action flag must be set")
		}

		if serviceStart {
			err := execs.StartService(strings.Join(args, ""))
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		if serviceEnable {
			err := execs.EnableService(strings.Join(args, ""))
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		if serviceStop {
			err := execs.StopService(strings.Join(args, ""))
			if err != nil {
				fmt.Println(err.Error())
			}

		}

		if serviceDisable {
			err := execs.DisableService(strings.Join(args, ""))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	},
}

// TODO: handle hard in another file
var restartServicesCmd = &cobra.Command{
	Use:   "restart [service name]",
	Short: "Restart service",
	Long:  `Restart reloads systemd daemon and restarts a service`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if hardRerload {
			err := execs.RestartService(strings.Join(args, ""))
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			err := execs.ReloadService(strings.Join(args, ""))
			if err != nil {
				fmt.Println(err.Error())
			}
		}

	},
}

var unlinkProjectServicesCmd = &cobra.Command{
	Use:   "unlink [project name / id]",
	Short: "Unlink destroys symbolic links",
	Long: `Unlink is for unlinking the services from 
	~/.config/containers/systemd and ~/.config/systemd/user`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := execs.UnlinkServices(strings.Join(args, ""))
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

var linkProjectServicesCmd = &cobra.Command{
	Use:   "link [project name / id]",
	Short: "Link project services",
	Long: `Link is for linking the services from
	/mole/services in the project ~/.config/containers/systemd if type is "podman"
	or ~/.config/systemd/user if type is "systemd"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := execs.LinkServices(strings.Join(args, ""), pType)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}
