package cmd

import (
	"fmt"
	"strings"

	"github.com/zulubit/mole/pkg/actions"
	"github.com/zulubit/mole/pkg/enums"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(servicesRootCmd)

	servicesRootCmd.AddCommand(reloadServicesCmd)
	servicesRootCmd.AddCommand(unlinkProjectServicesCmd)
	servicesRootCmd.AddCommand(listServicesCmd)

	servicesRootCmd.AddCommand(linkProjectServicesCmd)
	linkProjectServicesCmd.Flags().StringVarP(&pTypeFlag, "type", "t", "", "Type of services to be linked")
	linkProjectServicesCmd.MarkFlagRequired("type")

	servicesRootCmd.AddCommand(serviceActionCmd)
	serviceActionCmd.Flags().BoolVarP(&serviceStartFlag, "start", "", false, "Start the specified service")
	serviceActionCmd.Flags().BoolVarP(&serviceStopFlag, "stop", "", false, "Stop the specified service")
	serviceActionCmd.Flags().BoolVarP(&serviceEnableFlag, "enable", "", false, "Enable the specified service")
	serviceActionCmd.Flags().BoolVarP(&serviceDisableFlag, "disable", "", false, "Disable the specified service")

	servicesRootCmd.AddCommand(restartServicesCmd)
	restartServicesCmd.Flags().BoolVarP(&hardRerloadFlag, "full", "f", false, "Stop and start the service instead of reloading")
}

var servicesRootCmd = &cobra.Command{
	Use:   "services",
	Short: "Manage systemd services",
	Long: `Manage systemd services with commands to link, unlink, start, 
stop, enable, disable, reload, and restart services. This command provides 
tools to interact with services managed by the mole application.`,
}

var listServicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List managed services",
	Long: `Lists all services managed by the mole application. 
This includes only those services marked as "mole" services, which are 
automatically tagged when a project is managed by mole.`,
	Run: func(cmd *cobra.Command, args []string) {
		actions.ListServices()
	},
}

var reloadServicesCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload service unit files",
	Long: `Reloads the systemd manager configuration, registering any changes made 
to service unit files. This ensures that the system recognizes any updates 
to service definitions.`,
	Run: func(cmd *cobra.Command, args []string) {
		actions.ReloadServicesDaemon()
	},
}

var serviceActionCmd = &cobra.Command{
	Use:   "action [service name]",
	Short: "Control service state",
	Long: `Performs an action (start, stop, enable, disable) on the specified service. 
At least one action flag must be set to indicate the desired operation.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !serviceDisableFlag && !serviceEnableFlag && !serviceStartFlag && !serviceStopFlag {
			fmt.Println("Error: At least one action flag must be set")
			return
		}

		serviceName := strings.Join(args, "")

		if serviceStartFlag {
			if err := actions.StartService(serviceName); err != nil {
				fmt.Println(err.Error())
			}
		}

		if serviceEnableFlag {
			if err := actions.EnableService(serviceName); err != nil {
				fmt.Println(err.Error())
			}
		}

		if serviceStopFlag {
			if err := actions.StopService(serviceName); err != nil {
				fmt.Println(err.Error())
			}
		}

		if serviceDisableFlag {
			if err := actions.DisableService(serviceName); err != nil {
				fmt.Println(err.Error())
			}
		}
	},
}

var restartServicesCmd = &cobra.Command{
	Use:   "restart [service name]",
	Short: "Restart a service",
	Long: `Restarts a specified service. 
If the --full flag is set, the service will be stopped and then started, 
otherwise it will be reloaded instead.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := strings.Join(args, "")
		if hardRerloadFlag {
			if err := actions.RestartService(serviceName); err != nil {
				fmt.Println(err.Error())
			}
		} else {
			if err := actions.ReloadService(serviceName); err != nil {
				fmt.Println(err.Error())
			}
		}
	},
}

var unlinkProjectServicesCmd = &cobra.Command{
	Use:   "unlink [project name/id]",
	Short: "Unlink services from a project",
	Long: `Unlinks and removes symbolic links for services associated with the 
specified project. This operation affects services located in both 
~/.config/containers/systemd and ~/.config/systemd/user.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := actions.UnlinkServices(strings.Join(args, " ")); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var linkProjectServicesCmd = &cobra.Command{
	Use:   "link [project name/id]",
	Short: "Link services to a project",
	Long: `Links the specified project's services located in "mole_services" into the appropriate directory. 
Depending on the service type defined in the .env file, services will be linked to 
~/.config/containers/systemd (for Podman) or ~/.config/systemd/user (for systemd).`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pt, err := enums.IsProjectType(pTypeFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if err = actions.LinkServices(strings.Join(args, ""), pt); err != nil {
			fmt.Println(err.Error())
		}
	},
}
