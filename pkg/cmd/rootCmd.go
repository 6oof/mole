package cmd

import (
	"github.com/zulubit/mole/pkg/helpers"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}

var RootCmd = &cobra.Command{
	Use:   "mole",
	Short: "Micro-PaaS minimal in size and complexity but not in its power",
	Long:  helpers.MoleAsciiArt() + "\nMole is a lightweight micro-PaaS solution optimized for managing services via systemd.\nFind more information at https://github.com/zulubit/mole.",
}
