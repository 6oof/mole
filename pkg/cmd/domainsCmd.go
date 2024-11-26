package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zulubit/mole/pkg/actions"
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
	addStaticCaddyCmd.Flags().StringVarP(&locationFlag, "location", "l", "", "Location adds to the default path: /home/mole/projects/#project#/#provided location#")
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
	RunE: func(cmd *cobra.Command, args []string) error {
		e := strings.Join(args, " ")

		err := actions.SetupDomains(e)
		if err != nil {
			return err
		}

		fmt.Println("Domain setup complete. SSL issues will be reported to: " + e)
		return nil
	},
}

var listTakenPortsCmd = &cobra.Command{
	Use:   "ports",
	Short: "List active ports in use",
	Long: `This command lists all active ports currently in use, 
	retrieving the information using the "ss" command to display essential details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := actions.PortReport()
		if err != nil {
			return err
		}

		fmt.Println(p)
		return nil
	},
}

var reloadCaddyCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload the Caddy service configuration",
	Long: `Reload collects the main Caddyfile and all partial configurations, 
	merges them, and sends them to the Caddy API.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mainFilePath := "/home/mole/caddy/main.caddy"
		domainsDir := "/home/mole/domains"
		apiURL := "http://localhost:2019"

		err := actions.ReloadCaddy(mainFilePath, domainsDir, apiURL)
		if err != nil {
			return err
		}

		fmt.Println("Caddy configuration reloaded successfully.")
		return nil
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
	Long: `This command creates a reverse proxy configuration in Caddy for the specified project.
	if an empty on 0 port flag is set, MOLE_PORT_APP env variable will be used instead.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a := strings.Join(args, " ")

		err := actions.AddDomainProxy(a, domainFlag, portFlag)
		if err != nil {
			return err
		}

		fmt.Println("Reverse proxy added for: " + a)
		return nil
	},
}

var addStaticCaddyCmd = &cobra.Command{
	Use:   "static [project name/id]",
	Short: "Add a static file server route",
	Long:  `This command adds a static file server configuration in Caddy for serving files for the specified project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a := strings.Join(args, " ")

		err := actions.AddDomainStatic(a, domainFlag, locationFlag)
		if err != nil {
			return err
		}

		fmt.Println("Static server created for: " + a)
		return nil
	},
}

var deleteCaddyCmd = &cobra.Command{
	Use:   "delete [project name]",
	Short: "Delete a domain from the Caddy configuration",
	Long:  `This command attempts to locate and remove the specified project’s domain configuration file from Caddy.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a := strings.Join(args, " ")

		err := actions.DeleteProjectDomain(a)
		if err != nil {
			return err
		}

		fmt.Println("Caddy configuration deleted for: " + a + ".caddy")
		return nil
	},
}
