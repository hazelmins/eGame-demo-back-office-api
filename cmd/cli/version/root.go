package version

import (
	"fmt"

	_ "eGame-demo-back-office-api/docs"
	"github.com/spf13/cobra"
)

var CmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Get App Version",
	Run:   versionFunction,
}

var (
	version string
)

func versionFunction(cmd *cobra.Command, args []string) {
	fmt.Println(version)
}
