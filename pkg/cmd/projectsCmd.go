package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/6oof/mole/pkg/actions"
	"github.com/6oof/mole/pkg/data"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(projectsRootCmd)

	projectsRootCmd.AddCommand(listProjectsCmd)
	projectsRootCmd.AddCommand(findProjectCmd)
	projectsRootCmd.AddCommand(projectEnvCmd)

	addProjectCmd.Flags().StringVarP(&repositoryFlag, "repository", "r", "", "Repository URL *required")
	addProjectCmd.MarkFlagRequired("repository")
	addProjectCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "Branch *required")
	addProjectCmd.MarkFlagRequired("branch")
	addProjectCmd.Flags().StringVarP(&descriptionFlag, "description", "d", "", "Description")
	projectsRootCmd.AddCommand(addProjectCmd)

	editProjectCmd.Flags().StringVarP(&descriptionFlag, "description", "d", "", "Change description")
	editProjectCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "Change branch")
	projectsRootCmd.AddCommand(editProjectCmd)

	deleteProjectCmd.Flags().BoolVarP(&confirmFlag, "confirm", "y", false, "Confirms intent of delition *required")
	deleteProjectCmd.MarkFlagRequired("confirm")
	projectsRootCmd.AddCommand(deleteProjectCmd)
}

var projectsRootCmd = &cobra.Command{
	Use:   "projects",
	Short: "Interact with projects",
	Long:  `Interact with projects`,
}

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long: `List projects is for listing all projects.
	It only returns the projets not marked as deleted.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(data.ListProjects())
	},
}

var findProjectCmd = &cobra.Command{
	Use:   "find [project name / id]",
	Short: "Find a project by name",
	Long: `Find is for finding a project by name.
	Usefull bedause many commands reqire the project ID.
	The method is NOT case sensitive.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		p, err := data.FindProject(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(p.Stringify())
		}
	},
}

var addProjectCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new project",
	Long: `Adds a new project.
	Name and repository are required.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.Join(args, " ")
		re := regexp.MustCompile(`^[a-z0-9_-]+$`)
		if !re.MatchString(projectName) {
			fmt.Println("Error: Project name can only contain lowercase letters, digits, underscores, and hyphens.")
			return
		}

		np := data.Project{
			Name:          projectName,
			Description:   descriptionFlag,
			RepositoryUrl: repositoryFlag,
			Branch:        branchFlag,
		}

		err := data.CreateProject(np)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("New project successfully added")
		}
	},
}

// TODO: bring the project down before you delete it. Delete the domain partial before you delete it
var deleteProjectCmd = &cobra.Command{
	Use:   "delete [name/id]",
	Short: "Delete a project by name or ID",
	Long:  `Delete is for finding a project by id or name and deleting it.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := data.DeleteProject(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Project with id " + args[0] + " was marked as deleted")
		}
	},
}

var editProjectCmd = &cobra.Command{
	Use:   "edit [name/id]",
	Short: "Edit a project by name or ID",
	Long: `Edit is for finding a project by id or name and editing its properties.
	You won't be able to change it's repository, id, or name.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := data.EditProject(strings.Join(args, ""), descriptionFlag, branchFlag)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Project with id " + args[0] + " was updated")
		}
	},
}

var projectEnvCmd = &cobra.Command{
	Use:   "dotenv [name/id]",
	Short: "Edit project .env",
	Long:  `Env opens the project's .env file in nano.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := actions.FindAndEditEnv(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}
