package cmd

import (
	"fmt"

	"github.com/zulubit/mole/pkg/consts"
	"github.com/zulubit/mole/pkg/helpers"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mole",
	Long:  `All software has versions. This is mole's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(helpers.MoleAsciiArt() + "\n" + consts.Version)
	},
}
