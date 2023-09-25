package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"eGame-demo-back-office-api/configs"
	"eGame-demo-back-office-api/internal/models"
	"eGame-demo-back-office-api/pkg/mysqlx"
	"eGame-demo-back-office-api/pkg/redisx"

	"github.com/spf13/cobra"
)

var cmdMigrate = &cobra.Command{
	Use:   "migrate [-t table]",
	Short: "DB Migrate",
	Run:   migrateFunc,
}

var tables string
var configPath string

func init() {
	cmdMigrate.Flags().StringVarP(&configPath, "config path", "c", "", "config path")
	cmdMigrate.Flags().StringVarP(&tables, "table", "t", "", "input a table name")
}

//碼定義了一個名為 "table" 的參數選項，簡稱 "t"，並將其與變數 tables 綁定。這個選項用於指定要遷移的數據表的名稱

func migrateFunc(cmd *cobra.Command, args []string) {
	var tableMap map[string]struct{}
	var err error

	// 初始化配置並處理錯誤
	err = configs.Init(configPath)
	if err != nil {
		log.Fatalf("start fail:[Config Init] %s", err.Error())
	}

	// 初始化 Redis 並處理錯誤
	err = redisx.Init()
	if err != nil {
		log.Fatalf("start fail:[Redis Init] %s", err.Error())
	}

	// 初始化 Mysql 並處理錯誤
	err = mysqlx.Init()
	if err != nil {
		log.Fatalf("start fail:[Mysql Init] %s", err.Error())
	}

	// 創建一個空的 map 來存儲表格名稱
	tableMap = make(map[string]struct{})
	if tables != "" {
		tablesSlice := strings.Split(tables, ",")
		for _, v := range tablesSlice {
			fmt.Println(v)
			tableMap[v] = struct{}{}
		}
	}

	// 遍歷所有模型，可能是用於數據庫表格的模型
	for _, v := range models.GetModels() {
		db := mysqlx.GetDB(v.(mysqlx.GaTabler))

		// 如果提供了 tables 參數且模型不在 tableMap 中，則跳過該模型
		if tables != "" {
			if _, ok := tableMap[v.(mysqlx.GaTabler).TableName()]; !ok {
				continue
			}
		}

		// 使用 GORM 自動遷移數據表，並處理錯誤
		err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(v)
		if err != nil {
			fmt.Println("migrate database fail:", err.Error())
			os.Exit(0)
		}
	}
}
