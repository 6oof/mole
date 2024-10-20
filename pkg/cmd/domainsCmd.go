package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(domainsRootCmd)

	domainsRootCmd.AddCommand(listTakenPortsCmd)
	domainsRootCmd.AddCommand(reloadCaddyCmd)
}

var domainsRootCmd = &cobra.Command{
	Use:   "domains",
	Short: "Interact with caddy reverse proxy",
	Long: `Interact with caddy, read caddy logs, list all taken ports
	and more...`,
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

// TODO: figure out where to put the main caddy config and validate that to get the total validation picture
var reloadCaddyCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload caddy config",
	Long: `Reload validates the caddy config
	and realoads the caddy sevice if everything is correct.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("sh", "-c", "caddy validate --config example_data/Caddyfile")
		co, _ := c.CombinedOutput()

		fmt.Println(string(co))
	},
}

// TODO: figure out where to put the main caddy config and validate that to get the total validation picture
var validateCaddyCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate caddy config",
	Long:  `Validate validates the caddy config`,
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("sh", "-c", "caddy validate --config example_data/Caddyfile")
		co, _ := c.CombinedOutput()

		fmt.Println(string(co))
	},
}
