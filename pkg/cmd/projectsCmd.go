package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/6oof/mole/pkg/actions"
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
	addProjectCmd.Flags().StringVarP(&pTypeFlag, "type", "t", "", "Type *required")
	addProjectCmd.MarkFlagRequired("type")
	addProjectCmd.Flags().StringVarP(&descriptionFlag, "description", "d", "", "Description")
	projectsRootCmd.AddCommand(addProjectCmd)

	editProjectCmd.Flags().StringVarP(&descriptionFlag, "description", "d", "", "Change description")
	editProjectCmd.Flags().StringVarP(&branchFlag, "branch", "b", "", "Change branch")
	projectsRootCmd.AddCommand(editProjectCmd)

	deleteProjectCmd.Flags().BoolVarP(&confirmFlag, "confirm", "y", false, "Confirms intent of deletion *required")
	deleteProjectCmd.MarkFlagRequired("confirm")
	projectsRootCmd.AddCommand(deleteProjectCmd)
}

var projectsRootCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long: `Manage projects within the application. 
This command provides subcommands for creating, listing, 
finding, editing, and deleting projects.`,
}

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long: `Lists all projects in the system. 
This provides an overview of available projects for further actions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(actions.ListProjects())
	},
}

var findProjectCmd = &cobra.Command{
	Use:   "find [project name/id]",
	Short: "Find a project by name or ID",
	Long: `Searches for a project using its name or ID. 
This command is case insensitive and returns the project details 
to assist with further management commands.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		p, err := actions.FindProject(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(p.Stringify())
	},
}

var addProjectCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new project",
	Long: `Adds a new project to the system. 
You must provide a name, repository URL, branch, and type. 
Optionally, you can add a description.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.Join(args, " ")
		re := regexp.MustCompile(`^[a-z0-9_-]+$`)
		if !re.MatchString(projectName) {
			fmt.Println("Error: Project name can only contain lowercase letters, digits, underscores, and hyphens.")
			return
		}

		np := actions.Project{
			Name:          projectName,
			Description:   descriptionFlag,
			RepositoryURL: repositoryFlag,
			Branch:        branchFlag,
		}

		err := actions.CreateProject(np, pTypeFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("New project successfully added")
	},
}

var deleteProjectCmd = &cobra.Command{
	Use:   "delete [name/id]",
	Short: "Delete a project by name or ID",
	Long: `Deletes a project specified by its name or ID. 
This command marks the project as deleted and can also disable any 
associated services to ensure clean removal.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := actions.DisableStopAndUnlinkServices(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = actions.DeleteProject(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Project with id " + args[0] + " was marked as deleted")
	},
}

var editProjectCmd = &cobra.Command{
	Use:   "edit [name/id]",
	Short: "Edit a project by name or ID",
	Long: `Edits properties of a project identified by its name or ID. 
You can change the description or branch, but not the repository, ID, or name.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := actions.EditProject(strings.Join(args, ""), descriptionFlag, branchFlag)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Project with id " + args[0] + " was updated")
	},
}

var projectEnvCmd = &cobra.Command{
	Use:   "dotenv [name/id]",
	Short: "Edit project .env file",
	Long: `Opens the project's .env file for editing in the nano editor. 
This allows you to modify environment variables for the specified project.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := actions.FindAndEditEnv(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}
