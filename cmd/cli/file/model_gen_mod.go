package file

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"

	"eGame-demo-back-office-api/configs"

	"github.com/spf13/cobra"
)

var modelStr string = `package models

import (
	"gorm.io/gorm"
)

type {{.Name}} struct {
	BaseModle
}

// 设置数据表名称
func ({{.Short}} *{{.Name}}) TableName() string {
	return "{{.Table}}"
}

// 填充数据
func ({{.Short}} *{{.Name}}) FillData(db *gorm.DB) {
	
}

// 设置数据库链接名称
func ({{.Short}} *{{.Name}}) GetConnName() string {
	return "default"
}
`

func writeModel(fileName string, firstName string) error {
	// 定義一個叫做 parms 的結構，用於存儲生成模型（model）代碼所需的參數。
	parms := struct {
		Name  string // 模型名稱，通常是文件名
		Table string // 資料庫表格名稱
		Short string // 模型簡短描述
	}{
		Name:  fileName,  // 模型名稱，由函數的 fileName 參數指定
		Table: modelName, // 資料庫表格名稱，通常是 modelName
		Short: firstName, // 模型的簡短描述，由函數的 firstName 參數指定
	}

	// 設定新模型文件的路徑，這是在指定的配置根路徑下的 "internal/models" 目錄中的文件
	newPath := configs.RootPath + "internal" + string(filepath.Separator) + "models" + string(filepath.Separator) + fileName + ".go"

	// 檢查是否已經存在相同名稱的模型文件，如果存在則返回錯誤
	_, err := os.Lstat(newPath)
	if err == nil {
		return errors.New("file already exist")
	}

	// 創建新的模型文件
	file, err := os.Create(newPath)
	if err != nil {
		cobra.CompError(err.Error())
		return err
	}
	defer file.Close()

	// 使用模板將參數填充到模型文件中，生成模型代碼
	tem, _ := template.New("models_file").Parse(modelStr)
	tem.ExecuteTemplate(file, "models_file", parms)

	return nil
}
