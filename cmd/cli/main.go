package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/zulubit/mole/pkg/cmd"
	"github.com/zulubit/mole/pkg/consts"
	"github.com/spf13/cobra/doc"
)

func main() {
	env, ep := os.LookupEnv("MOLE_ENV_PROD")
	if ep && env == "1" {
		consts.BasePath = "/home/mole"
	}

	doc, dg := os.LookupEnv("MOLE_DOC")
	if dg && doc == "1" {
		generateDocs()
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateDocs() {
	err := doc.GenMarkdownTree(cmd.RootCmd, path.Join("docs", "cli"))
	if err != nil {
		log.Fatal(err)
	}

}
