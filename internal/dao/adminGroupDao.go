package dao

import (
	"context"
	"sync"

	"eGame-demo-back-office-api/internal/models"
	"eGame-demo-back-office-api/pkg/mysqlx"

	"gorm.io/gorm"
)

// 定義一個名為 AdminGroupDao 的結構體，該結構體包含一個指向 gorm.DB 的指針。
type AdminGroupDao struct {
	DB *gorm.DB
}

// 定義了兩個變數：instanceAdminGroup 和 onceAdminGroup。
var (
	instanceAdminGroup *AdminGroupDao
	onceAdminGroup     sync.Once
)

// 定義一個名為 NewAdminGroupDao 的函數，該函數返回一個指向 AdminGroupDao 結構體的指針。
func NewAdminGroupDao() *AdminGroupDao {
	// 使用 sync.Once 確保以下代碼只會執行一次。
	onceAdminGroup.Do(func() {
		// 初始化 instanceAdminGroup，將其設置為一個指向 AdminGroupDao 結構體的指針，並設置 DB 字段。
		instanceAdminGroup = &AdminGroupDao{DB: mysqlx.GetDB(&mysqlx.BaseModle{ConnName: "default"})}
	})
	// 返回 instanceAdminGroup 的指針。
	return instanceAdminGroup
}

/*
*進DB依照群組跟玩家撈取權限返回
 */
func (dao *AdminGroupDao) GetAdminGroup(ctx context.Context, groupName string, username string) *gorm.DB {
	// 創建一個初始的數據庫查詢對象
	db := dao.DB.WithContext(ctx).Table("super_admin")

	// 根據提供的條件動態構建查詢
	if groupName != "" {
		db = db.Where("group_name = ?", groupName)
	}

	if username != "" {
		db = db.Where("username = ?", username)
	}

	// 返回數據庫查詢對象，以便在其他地方進行進一步操作
	return db
}

/*
*單純進db撈群組資料返回沒有任何條件設定
 */
func (dao *AdminGroupDao) GetGroupIndex(ctx context.Context) *gorm.DB {
	// 创建一个初始的数据库查询对象
	db := dao.DB.WithContext(ctx).Table("super_admin")

	// 执行查询
	var results []map[string]interface{}
	db = db.Select("group_name, permissions_json").Find(&results)

	// 返回查询结果
	return db
}

func (dao *AdminGroupDao) GetPermissionsByGroupName(groupName string) (permission models.SuperAdmin, err error) {
	err = dao.DB.Where("group_name = ?", groupName).First(&permission).Error
	return
}
