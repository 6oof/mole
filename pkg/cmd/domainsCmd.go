package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/6oof/mole/pkg/data"
	"github.com/spf13/cobra"
)

// TODO: we need to validate inputs and flags to not make an unrecoverable mess
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
	Short: "Interact with caddy reverse proxy",
	Long: `Interact with caddy, read caddy logs, list all taken ports
	and more...`,
}

var setupCaddyCmd = &cobra.Command{
	Use:   "setup [email]",
	Short: "Setup support for domains",
	Long: `Setup creates the main caddy config.
	This command will overwrite the existing config.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		e := strings.Join(args, " ")

		err := data.SetupDomains(e)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Domain setup done. ssl issues reported to: " + e)
		}

	},
}

var listTakenPortsCmd = &cobra.Command{
	Use:   "ports",
	Short: "List active ports",
	Long: `List projects is for listing all taken ports.
	It uses ss under the hood and strips most of the output.`,
	Run: func(cmd *cobra.Command, args []string) {

		c := exec.Command("sh", "-c", "ss -tulnH | awk '{print $5}'")
		co, _ := c.CombinedOutput()

		fmt.Println(string(co))
	},
}

var reloadCaddyCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload caddy config",
	Long: `Reload validates the caddy config
	and realoads the caddy sevice if everything is correct.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("sh", "-c", "caddy validate --config /home/mole/caddy/main.caddy")
		co, _ := c.CombinedOutput()
		fmt.Println(string(co))

		reloadCmd := exec.Command("sh", "-c", "systemctl --user reload caddy")
		reloadOut, err := reloadCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to reload caddy service: %v\n", err)
			return
		}
		fmt.Println("Caddy service reloaded successfully.")
		fmt.Println(string(reloadOut))
	},
}

var validateCaddyCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate caddy config",
	Long:  `Validate validates the caddy config`,
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("sh", "-c", "caddy validate --config /home/mole/caddy/main.caddy")
		co, _ := c.CombinedOutput()

		fmt.Println(string(co))
	},
}

var addCaddyCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a domain to caddy config",
	Long: `Add supports adding a reverse proxy
	or static file server domains with automatic tls`,
}

var addProxyCaddyCmd = &cobra.Command{
	Use:   "proxy [project name/id]",
	Short: "Add a reverse proxy",
	Long:  `Add supports adding a reverse proxy partial config.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		a := strings.Join(args, " ")

		err := data.AddDomainProxy(a, domainFlag, portFlag)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Reverse proxy added for: " + a)
		}

	},
}

var addStaticCaddyCmd = &cobra.Command{
	Use:   "static [project name/id]",
	Short: "Add a static route",
	Long:  `Add supports adding a static fileserver partial config.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		a := strings.Join(args, " ")

		err := data.AddDomainStatic(a, domainFlag, locationFlag)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Static server created for: " + a)
		}

	},
}

var deleteCaddyCmd = &cobra.Command{
	Use:   "delete [project name]",
	Short: "Delete a caddy partial",
	Long:  `Delete trys to find a caddy partial and deletes it.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		a := strings.Join(args, " ")

		err := data.DeleteProjectDomain(a)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Caddy partial config deleted: " + a + ".caddy")
		}

	},
}
