package main

import (
	"fmt"
	"os"
	"time"

	"eGame-demo-back-office-api/cmd/cli/db"
	"eGame-demo-back-office-api/cmd/cli/file"
	"eGame-demo-back-office-api/cmd/cli/run"
	"eGame-demo-back-office-api/cmd/cli/version"
	"github.com/spf13/cobra"
)

var (
	release bool = true
)

func main() {

	// 设置时区
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("set location fail: %s", err.Error())
		os.Exit(1)
	}
	time.Local = local

	var rootCmd = &cobra.Command{Use: "ginadmin"}
	rootCmd.AddCommand(run.CmdRun, db.CmdDb, file.CmdFile, version.CmdVersion)
	rootCmd.Execute()

}
