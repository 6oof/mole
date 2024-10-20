package main

import (
	"fmt"
	"os"

	"github.com/6oof/mole/pkg/cmd"
	"github.com/6oof/mole/pkg/consts"
)

func main() {
	env, ep := os.LookupEnv("MOLE_ENV_PROD")
	if ep && env == "1" {
		consts.BasePath = ""
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
