package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/6oof/mole/pkg/actions"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(domainsRootCmd)

	domainsRootCmd.AddCommand(listTakenPortsCmd)
	domainsRootCmd.AddCommand(reloadCaddyCmd)
	domainsRootCmd.AddCommand(setupCaddyCmd)
	domainsRootCmd.AddCommand(deleteCaddyCmd)

	addProxyCaddyCmd.Flags().StringVarP(&domainFlag, "domain", "d", "", "Domain *required")
	addProxyCaddyCmd.MarkFlagRequired("domain")
	addProxyCaddyCmd.Flags().IntVarP(&portFlag, "port", "p", 0, "Port *required")
	addProxyCaddyCmd.MarkFlagRequired("port")
	addCaddyCmd.AddCommand(addProxyCaddyCmd)

	addStaticCaddyCmd.Flags().StringVarP(&domainFlag, "domain", "d", "", "Domain *required")
	addStaticCaddyCmd.MarkFlagRequired("domain")
	addStaticCaddyCmd.Flags().StringVarP(&locationFlag, "location", "l", "", "Location *required")
	addCaddyCmd.AddCommand(addStaticCaddyCmd)

	domainsRootCmd.AddCommand(addCaddyCmd)
}

var domainsRootCmd = &cobra.Command{
	Use:   "domains",
	Short: "Manage Caddy reverse proxy configurations",
	Long: `The "domains" command group allows interaction with Caddy, 
	including managing reverse proxy domains, viewing logs, listing used ports, 
	and more to support domain configuration.`,
}

var setupCaddyCmd = &cobra.Command{
	Use:   "setup [email]",
	Short: "Initialize Caddy with domain support",
	Long: `Setup initializes the primary Caddy configuration.
	This command will overwrite any existing configuration file.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		e := strings.Join(args, " ")

		err := actions.SetupDomains(e)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Domain setup complete. SSL issues will be reported to: " + e)
	},
}

var listTakenPortsCmd = &cobra.Command{
	Use:   "ports",
	Short: "List active ports in use",
	Long: `This command lists all active ports currently in use, 
	retrieving the information using the "ss" command to display essential details.`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := actions.PortReport()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(p)
	},
}

var reloadCaddyCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload the Caddy service configuration",
	Long: `Reload re-validates the current Caddy configuration file 
	and reloads the Caddy service if validation succeeds.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("sh", "-c", "caddy validate --config /home/mole/caddy/main.caddy")
		co, _ := c.CombinedOutput()
		fmt.Println(string(co))

		reloadCmd := exec.Command("sh", "-c", "systemctl --user reload caddy")
		reloadOut, err := reloadCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to reload Caddy service: %v\n", err)
			return
		}
		fmt.Println("Caddy service reloaded successfully.")
		fmt.Println(string(reloadOut))
	},
}

var validateCaddyCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the Caddy configuration file",
	Long:  `This command verifies the integrity and correctness of the current Caddy configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("sh", "-c", "caddy validate --config /home/mole/caddy/main.caddy")
		co, _ := c.CombinedOutput()

		fmt.Println(string(co))
	},
}

var addCaddyCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new domain to the Caddy configuration",
	Long: `This command supports adding domains to the Caddy configuration, 
	including reverse proxies and static file servers with automatic TLS configuration.`,
}

var addProxyCaddyCmd = &cobra.Command{
	Use:   "proxy [project name/id]",
	Short: "Add a reverse proxy for a domain",
	Long:  `This command creates a reverse proxy configuration in Caddy for the specified project.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		a := strings.Join(args, " ")

		err := actions.AddDomainProxy(a, domainFlag, portFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Reverse proxy added for: " + a)
	},
}

var addStaticCaddyCmd = &cobra.Command{
	Use:   "static [project name/id]",
	Short: "Add a static file server route",
	Long:  `This command adds a static file server configuration in Caddy for serving files for the specified project.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		a := strings.Join(args, " ")

		err := actions.AddDomainStatic(a, domainFlag, locationFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Static server created for: " + a)
	},
}

var deleteCaddyCmd = &cobra.Command{
	Use:   "delete [project name]",
	Short: "Delete a domain from the Caddy configuration",
	Long:  `This command attempts to locate and remove the specified project’s domain configuration file from Caddy.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		a := strings.Join(args, " ")

		err := actions.DeleteProjectDomain(a)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Caddy configuration deleted for: " + a + ".caddy")
	},
}
