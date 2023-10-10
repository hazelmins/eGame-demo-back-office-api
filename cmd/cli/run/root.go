/*
已經合併db migrate and seed
*/

package run

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"eGame-demo-back-office-api/configs"
	"eGame-demo-back-office-api/internal"
	"eGame-demo-back-office-api/internal/cron"
	"eGame-demo-back-office-api/internal/models"
	"eGame-demo-back-office-api/internal/router"
	"eGame-demo-back-office-api/pkg/mysqlx"
	"eGame-demo-back-office-api/pkg/redisx"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var CmdRun = &cobra.Command{
	Use:   "run",
	Short: "Run app",
	Run:   runFunction,
}

var (
	configPath string
	crontab    string
	mode       string
	tables     string // 用于指定迁移的表格参数
	tableSeed  string // 用于指定种子数据的表格参数
)

func init() {
	CmdRun.Flags().StringVarP(&configPath, "config path", "c", "", "config path")
	CmdRun.Flags().StringVarP(&mode, "mode", "m", "debug", "debug or release")
	CmdRun.Flags().StringVarP(&crontab, "cron", "t", "open", "scheduled task control open or close")
	CmdRun.Flags().StringVarP(&tables, "table", "T", "", "input a table name")
	CmdRun.Flags().StringVarP(&tableSeed, "seed", "s", "", "input a table name for seeding")
}
func runFunction(cmd *cobra.Command, args []string) {
	var err error

	showLogo()

	//判断是否编译线上版本
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}

	//定时任务
	if crontab == "open" {
		cron.Init()
	}

	err = configs.Init(configPath)
	if err != nil {
		log.Fatalf("start fail:[Config Init] %s", err.Error())
	}

	err = redisx.Init()
	if err != nil {
		log.Fatalf("start fail:[Redis Init] %s", err.Error())
	}

	err = mysqlx.Init()
	if err != nil {
		log.Fatalf("start fail:[Mysql Init] %s", err.Error())
	}

	// 在这里直接执行数据库迁移和种子操作

	migrateDatabase()
	seedDatabase()
	showPanel()

	r, err := router.Init()
	if err != nil {
		log.Fatalf("start fail:[Route Init] %s", err.Error())
	}

	app := internal.Application{}

	r.SetEngine(&app)
	app.Run()
}

func showLogo() {
	fmt.Println("   _____ _                   _           _       ")
	fmt.Println("  / ____(_)         /\\      | |         (_)      ")
	fmt.Println(" | |  __ _ _ __    /  \\   __| |_ __ ___  _ _ __  ")
	fmt.Println(" | | |_ | | '_ \\  / /\\ \\ / _` | '_ ` _ \\| | '_ \\ ")
	fmt.Println(" | |__| | | | | |/ _____\\ (_| | | | | | | | | | |")
	fmt.Println("  \\_____|_|_| |_/_/    \\_\\__,_|_| |_| |_|_|_| |_| ")
}

func showPanel() {
	fmt.Println("App running at:")
	fmt.Printf("- Http Address:   %c[%d;%d;%dm%s%c[0m \n", 0x1B, 0, 40, 34, "http://"+configs.App.Base.Host+":"+configs.App.Base.Port, 0x1B)
	fmt.Println("")
}
func migrateDatabase() {
	// 创建一个空的 map 来存储表格名称
	tableMap := make(map[string]struct{})

	// 如果提供了表格参数，将其拆分为字符串切片并填充 tableMap
	if tables != "" {
		tablesSlice := strings.Split(tables, ",")
		for _, v := range tablesSlice {
			tableMap[v] = struct{}{}
		}
	}

	// 遍历所有模型，可能是用于数据库表格的模型
	for _, v := range models.GetModels() {
		db := mysqlx.GetDB(v.(mysqlx.GaTabler))

		// 如果提供了表格参数且模型不在 tableMap 中，则跳过该模型
		if tables != "" {
			if _, ok := tableMap[v.(mysqlx.GaTabler).TableName()]; !ok {
				continue
			}
		}

		// 使用 GORM 自动迁移数据表，并处理错误
		err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(v)
		if err != nil {
			fmt.Println("migrate database fail:", err.Error())
			os.Exit(0)
		}
	}
}

func seedDatabase() {
	// 创建一个空的 map 来存储表格名称
	tableMap := make(map[string]struct{})

	// 如果提供了表格参数，将其拆分为字符串切片并填充 tableMap
	// 如果提供了 tableSeed，將其拆分為字符串切片並填充 tableMap
	if tableSeed != "" {
		tablesSlice := strings.Split(tableSeed, ",")
		for _, v := range tablesSlice {
			tableMap[v] = struct{}{}
		}
	}

	// GetModels 中列了所以的TABLE 在這裡阿鄉親
	for _, v := range models.GetModels() {

		// 如果提供了 tableSeed 且模型不在 tableMap 中，則跳過該模型
		if tableSeed != "" {
			if _, ok := tableMap[v.(mysqlx.GaTabler).TableName()]; !ok {
				continue
			}
		}

		// 將模型轉換為 mysqlx.GaTabler 介面，從數據庫中填充數據
		//這裡會將默認的admin灌入db
		tabler := v.(mysqlx.GaTabler)
		db := mysqlx.GetDB(tabler)
		tabler.FillData(db)
	}
}
