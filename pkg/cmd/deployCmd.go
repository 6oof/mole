package cmd

import (
	"fmt"
	"strings"

	"github.com/6oof/mole/pkg/actions"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolVar(&restartOnDeplyFlag, "restart", false, "restarts services instead or reloading on deployment")
}

var deployCmd = &cobra.Command{
	Use:   "deploy [project name/id]",
	Short: "Deploy triggers project deployment",
	Long: `Deploy triggers project deployment.

Depending on the project type it will do the following:

1. Git pull
(2) Link all the services / run the deployment script
(3) Start the services defined in .env`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := actions.RunDeployment(strings.Join(args, ""), restartOnDeplyFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Deployment succeeded.")
	},
}
