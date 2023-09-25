package file

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"

	"eGame-demo-back-office-api/configs"

	"github.com/spf13/cobra"
)

var daoStr string = `package dao

import (
	"sync"

	"eGame-demo-back-office-api/internal/models"

	"gorm.io/gorm"
)

type {{.ModelName}}Dao struct {
	DB *gorm.DB
}

var (
	instance{{.ModelName}} *{{.ModelName}}Dao
	once{{.ModelName}}Dao  sync.Once
)

func New{{.ModelName}}Dao() *{{.ModelName}}Dao {
	once{{.ModelName}}Dao.Do(func() {
		instance{{.ModelName}} = &{{.ModelName}}Dao{DB: models.GetDB(&models.{{.ModelName}}{})}
	})
	return instance{{.ModelName}}
}

// 新增数据
func (dao *{{.ModelName}}Dao) Create(data models.{{.ModelName}}) error {
	return dao.DB.Create(&data).Error
}

// 获取单条数据
func (dao *{{.ModelName}}Dao) GetOne(conditions map[string]interface{}) (data models.{{.ModelName}}, err error) {

	err = dao.DB.First(&data, conditions).Error
	return
}

// 更新数据
func (dao *{{.ModelName}}Dao) UpdateColumns(conditions, field map[string]interface{}, tx *gorm.DB) error {

	if tx != nil {
		return tx.Model(&models.{{.ModelName}}{}).Where(conditions).UpdateColumns(field).Error
	}

	return dao.DB.Model(&models.{{.ModelName}}{}).Where(conditions).UpdateColumns(field).Error
}

// 删除数据
func (dao *{{.ModelName}}Dao) Del(conditions map[string]interface{}) error {
	return dao.DB.Delete(&models.{{.ModelName}}{}, conditions).Error
}
`

func writeDao(modelName string, file string) error {
	// 定義一個叫做 parms 的結構，用於存儲 DAO（資料訪問對象）生成所需的參數。
	parms := struct {
		ModelName string // 資料模型的名稱，例如 "Person"
	}{
		ModelName: modelName, // 資料模型的名稱由函數的參數 modelName 傳遞進來
	}

	// 設定新 DAO 文件的路徑，這是在指定的配置根路徑下的 "internal/dao" 目錄中的文件
	newDaoPath := configs.RootPath + "internal" + string(filepath.Separator) + "dao" + string(filepath.Separator) + file + "Dao.go"

	// 檢查是否已經存在相同名稱的 DAO 文件，如果存在則返回錯誤
	_, err := os.Lstat(newDaoPath)
	if err == nil {
		return errors.New("file already exist")
	}

	// 創建新的 DAO 文件
	fileDao, err := os.Create(newDaoPath)
	if err != nil {
		cobra.CompError(err.Error())
		return err
	}
	defer fileDao.Close()

	// 使用模板將參數填充到 DAO 文件中，生成 DAO 代碼
	tem, _ := template.New("models_file").Parse(daoStr)
	tem.ExecuteTemplate(fileDao, "models_file", parms)

	return nil
}
