package cmd

import (
	"github.com/6oof/mole/pkg/helpers"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}

var RootCmd = &cobra.Command{
	Use:   "mole",
	Short: "micro-PaaS focused around systemd",
	Long:  helpers.MoleAsciiArt() + "\nMicro-PaaS solution minimal in size and footprint.\nFocused around systemd.\nfind more info at github.com/6oof/mole",
}
