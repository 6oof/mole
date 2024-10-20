package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/6oof/mole/pkg/data"
	"github.com/spf13/cobra"
)

var (
	repository  string
	description string
	branch      string
	confirm     bool
)

func init() {
	RootCmd.AddCommand(projectsRootCmd)

	projectsRootCmd.AddCommand(listProjectsCmd)
	projectsRootCmd.AddCommand(findProjectCmd)
	projectsRootCmd.AddCommand(projectEnvCmd)

	addProjectCmd.Flags().StringVarP(&repository, "repository", "r", "", "Repository URL *required")
	addProjectCmd.MarkFlagRequired("repository")
	addProjectCmd.Flags().StringVarP(&branch, "branch", "b", "", "Branch *required")
	addProjectCmd.MarkFlagRequired("branch")
	addProjectCmd.Flags().StringVarP(&description, "description", "d", "", "Description")
	projectsRootCmd.AddCommand(addProjectCmd)

	editProjectCmd.Flags().StringVarP(&description, "description", "d", "", "Change description")
	editProjectCmd.Flags().StringVarP(&branch, "branch", "b", "", "Change branch")
	projectsRootCmd.AddCommand(editProjectCmd)

	deleteProjectCmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Confirms intent of delition *required")
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
		fmt.Println(data.FindProject(strings.Join(args, " ")))
	},
}

// TODO: This sould imidiately pull the project down
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
			Description:   description,
			RepositoryUrl: repository,
			Branch:        branch,
		}

		id, err := data.AddProject(np)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("New project added with ID: " + id)
		}
	},
}

var deleteProjectCmd = &cobra.Command{
	Use:   "delete [project id]",
	Short: "Delete a project by ID",
	Long:  `Delete is for finding a project by id and deleting it.`,
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
	Use:   "edit [project id]",
	Short: "Edit a project by ID",
	Long: `Edit is for finding a project by id and editing its properties.
	You won't be able to change it's repository or name.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := data.EditProject(strings.Join(args, ""), description, branch)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Project with id " + args[0] + " was updated")
		}
	},
}

// TODO: find the right project folder to do this in
var projectEnvCmd = &cobra.Command{
	Use:   "env [name/id]",
	Short: "Edit project .env",
	Long:  `Env opens the project's .env file in vi.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("nano", ".env")
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		// Run the command and handle any error
		err := c.Run()
		if err != nil {
			fmt.Println("Error running vi:", err.Error())
		}
	},
}
