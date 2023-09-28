package db

import (
	"eGame-demo-back-office-api/configs"
	"eGame-demo-back-office-api/internal/models"
	"eGame-demo-back-office-api/pkg/mysqlx"
	"eGame-demo-back-office-api/pkg/redisx"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var cmdSeed = &cobra.Command{
	//seed [-t table] 表示命令的基本用法是 seed，並且可以包含 -t 選項，後跟表格名稱（例如 -t users）
	Use:   "seed [-t table]",
	Short: "DB Seed",
	Run:   seedFunc,
}

//當用戶運行 seed 命令時，將調用 seedFunc 函數執行實際的初始化操作。

var tableSeed string
var confPath string

func init() {
	cmdSeed.Flags().StringVarP(&confPath, "config path", "c", "", "config path")
	cmdSeed.Flags().StringVarP(&tableSeed, "table", "t", "", "input a table name")
}

func seedFunc(cmd *cobra.Command, args []string) {
	var tableMap map[string]struct{} // 創建一個空的 map 變數 tableMap，用於存儲表格名稱
	var err error                    // 宣告一個用於錯誤處理的變數 err

	// 初始化配置（可能是讀取配置文件）並處理錯誤
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

	// 如果提供了 tableSeed，將其拆分為字符串切片並填充 tableMap
	if tableSeed != "" {
		tablesSlice := strings.Split(tableSeed, ",")
		for _, v := range tablesSlice {
			tableMap[v] = struct{}{}
		}
	}

	// 遍歷所有模型，可能是用於數據庫表格的模型
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
