package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zulubit/mole/pkg/actions"
)

func init() {
	RootCmd.AddCommand(templatesRootCmd)

	templatesRootCmd.AddCommand(transformProjectCmd)
}

var templatesRootCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage project templates",
	Long: `Manage project templates by providing tools to transform files such as 
"mole-compose.yaml" and "mole-deploy.sh" into ready-to-use configurations. These transformations 
inject environment variables from the project's .env file to create deployment-ready files.`,
}

var transformProjectCmd = &cobra.Command{
	Use:   "transform [project name/id]",
	Short: "Transform project service templates",
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
