package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zulubit/mole/pkg/actions"
)

func init() {
	RootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolVar(&deployDown, "down", false, "Try to run docker compose down on mole-compose-ready.yaml or fail.")
}

var deployCmd = &cobra.Command{
	Use:   "deploy [project name/id]",
	Short: "Deploy triggers project deployment",
	Long: `Deploy triggers project deployment.
This will transform your mole.sh and run it.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !deployDown {
			succ, err := actions.RunDeployment(strings.Join(args, ""))
			if err != nil {
				return err
			}

			fmt.Println(succ)
		} else if deployDown {
			succ, err := actions.RundDeplyDown(strings.Join(args, ""))
			if err != nil {
				return err
			}
			fmt.Println(succ)
		}

		return nil
	},
}
