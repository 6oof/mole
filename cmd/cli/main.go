package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/spf13/cobra/doc"
	"github.com/zulubit/mole/pkg/cmd"
	"github.com/zulubit/mole/pkg/consts"
)

func main() {
	// Restrict execution to non-root users in production
	if consts.Prod == 1 {
		if os.Geteuid() == 0 {
			fmt.Println("Error: This CLI cannot be run as root to avoid permission issues.")
			os.Exit(1)
		}

		currentUser, err := user.Current()
		if err != nil {
			log.Fatalf("Failed to get current user: %v", err)
		}

		if currentUser.Username != "mole" {
			fmt.Printf("Error: This CLI can only be executed by the user 'mole', not '%s'.\n", currentUser.Username)
			os.Exit(1)
		}
	}

	// Check for MOLE_DOC environment variable and generate documentation if set
	doc, dg := os.LookupEnv("MOLE_DOC")
	if dg && doc == "1" {
		generateDocs()
	}

	// Execute the root command
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Function to generate CLI documentation
func generateDocs() {
	err := doc.GenMarkdownTree(cmd.RootCmd, path.Join(consts.GetBasePath(), "docs", "cli"))
	if err != nil {
		log.Fatal(err)
	}
}
