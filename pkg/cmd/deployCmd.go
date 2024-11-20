package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zulubit/mole/pkg/actions"
)

func init() {
	RootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy [project name/id]",
	Short: "Deploy triggers project deployment",
	Long: `Deploy triggers project deployment.
This will resolve your mole templates and run the deploy script.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		succ, err := actions.RunDeployment(strings.Join(args, ""))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(succ)
	},
}
