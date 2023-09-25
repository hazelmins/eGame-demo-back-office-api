/*
 * @Description:
 * @Author: gphper
 * @Date: 2021-06-01 20:15:04
 */
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

// @title GinAdmin Api
// @version 1.0
// @description GinAdmin 示例项目

// @contact.name gphper
// @contact.url https://eGame-demo-back-office-api

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:20011
// @basepath /api
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
