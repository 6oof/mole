package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zulubit/mole/pkg/actions"
	"github.com/zulubit/mole/pkg/helpers"
)

func init() {
	RootCmd.AddCommand(keysRootCmd)

	keysRootCmd.AddCommand(readDeployKeyCmd)
	keysRootCmd.AddCommand(getActionsKeyCmd)

	keysRootCmd.AddCommand(addAuthorizedKeyCmd)
	addAuthorizedKeyCmd.Flags().StringVarP(&keyName, "name", "n", "", "name the key for future reference *required")
	addAuthorizedKeyCmd.MarkFlagRequired("name")
}

var keysRootCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage SSH keys for secure server access",
	Long: `The "keys" command group provides options for managing SSH keys, 
including deploying and authorizing keys for secure server access.

To remove an authorized key, navigate to /home/mole/.ssh and delete the 
relevant entry in the authorized_keys file.

To regenerate the deploy key, navigate to /home/mole/.ssh and delete 
the id_rsa and id_rsa.pub files. This will allow the system to create 
a new deploy key as needed.`,
}

var readDeployKeyCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Retrieve or create the deploy key for SSH access",
	Long: `The "deploy" command displays the current deploy key (id_rsa.pub),
used for accessing private repositories securely. 

If no deploy key is found, a new one will be generated automatically 
and saved to the standard SSH key path. This deploy key enables 
secure, automated interactions with external repositories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := actions.FindOrCreateDeployKey()
		if err != nil {
			return err
		}

		fmt.Println(key)
		return nil
	},
}

var addAuthorizedKeyCmd = &cobra.Command{
	Use:   "authorize [public RSA key]",
	Short: "Add a new public key to the authorized_keys file",
	Long: `The "authorize" command validates and appends a given public RSA key 
to the authorized_keys file. This allows the specified key to be used 
for SSH access to the server.

Ensure the key provided is correctly formatted, as it will be validated 
before being added to prevent errors. Only unique keys will be appended 
to avoid duplicates in the authorized_keys file.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !helpers.ValidateProjectName(keyName) {
			return errors.New("Error: Key name can only contain lowercase letters, digits, underscores, and hyphens. It should start and end with a letter or a number")
		}

		err := actions.AddAuthorizedKeys(strings.Join(args, " "), keyName)
		if err != nil {
			return err
		}
		fmt.Println("Key added")
		return nil
	},
}

var getActionsKeyCmd = &cobra.Command{
	Use:   "actions",
	Short: "Retrieve or create the SSH key for actions and add it to authorized_keys",
	Long: `The "actions" command generates or retrieves the private SSH key used for 
server-to-server communication or other automated tasks. 

If no key is found, a new private key (actions_rsa) and its corresponding 
public key (actions_rsa.pub) are created and stored in the standard SSH 
directory. The public key is automatically added to the authorized_keys 
file, allowing the associated private key to be used for secure access.

This command is particularly useful for enabling secure access for CI/CD 
pipelines, automated scripts, or other server-to-server operations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := actions.FindOrCreateActionsKey()
		if err != nil {
			return err
		}

		fmt.Println(key)
		return nil
	},
}
