package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zulubit/mole/pkg/actions"
)

func init() {
	RootCmd.AddCommand(transformProjectCmd)
}

// TODO: Add links to main readme
var transformProjectCmd = &cobra.Command{
	Use:   "templates [project name/id]",
	Short: "Transform project mole templates",
	Long: `Generates project-specific configuration files by transforming 
"mole-compose.yaml" and "mole-deploy.sh" templates into ready-to-use files. 
Environment variables from the project's .env file are injected during the transformation.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectNOI := strings.Join(args, "")

		if err := actions.TransformCompose(projectNOI); err != nil {
			fmt.Printf("Error transforming mole-compose.yaml: %s\n", err.Error())
			return
		}

		if err := actions.TransformDeploy(projectNOI); err != nil {
			fmt.Printf("Error transforming mole-deploy.sh: %s\n", err.Error())
			return
		}

		fmt.Println("Transformation completed successfully.")
	},
}
