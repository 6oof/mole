package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zulubit/mole/pkg/actions"
)

func init() {
	RootCmd.AddCommand(transformProjectCmd)

	transformProjectCmd.AddCommand(transformProjectShCmd)
	transformProjectCmd.AddCommand(transformProjectComposeCmd)
}

var transformProjectCmd = &cobra.Command{
	Use:   "templates",
	Short: "Transform project mole templates",
	Long: `Generates project-specific configuration files by transforming 
"mole-compose.yaml" and "mole.sh" templates into ready-to-use files.
Variables from the project's secrets are injected during the transformation.`,
}

var transformProjectComposeCmd = &cobra.Command{
	Use:   "compose [project name/id]",
	Short: "Transforms mole-compose.yaml into mole-compose-ready.yaml",
	Long: `Processes the mole-compose.yaml template for a specified project 
and generates a ready-to-use mole-compose-ready.yaml file. This transformation 
injects project-specific variables derived from the project's secrets.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectNOI := strings.Join(args, "")

		if err := actions.TransformCompose(projectNOI); err != nil {
			return fmt.Errorf("Error transforming mole-compose.yaml: %s\n", err.Error())
		}

		fmt.Println("Transformation completed successfully.")
		return nil
	},
}

var transformProjectShCmd = &cobra.Command{
	Use:   "deploy [project name/id]",
	Short: "Transforms mole.sh into a ready-to-execute script",
	Long: `Processes the mole.sh template for a specified project and generates 
a fully configured shell script. This transformation replaces placeholders 
with project-specific variables derived from the project's secrets.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectNOI := strings.Join(args, "")

		if err := actions.TransformDeploy(projectNOI); err != nil {
			return fmt.Errorf("Error transforming mole.sh: %s\n", err.Error())
		}

		fmt.Println("Transformation completed successfully.")
		return nil
	},
}
