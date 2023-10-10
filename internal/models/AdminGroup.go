/*
 * @Description:用户组相关model
 */

package models

import (
	"eGame-demo-back-office-api/pkg/mysqlx"
	"log"

	"gorm.io/gorm"
)

type AdminGroupSaveReq struct {
	Privs     []string `form:"privs[]" label:"权限" json:"privs" binding:"required"`
	Username  string   `form:"username" label:"管理者名稱" json:"username" binding:"required"`
	GroupName string   `form:"groupname" label:"用户组名" json:"groupname" binding:"required"`
	GroupId   uint     `form:"groupid"`
}

// speradmin DB表格建立
type SuperAdmin struct {
	mysqlx.BaseModle
	Uid             uint            `gorm:"primary_key;auto_increment"`
	GroupName       string          `gorm:"column:group_name"` // 使用標籤指定列名
	Permissions     map[string]bool `gorm:"-"`
	PermissionsJSON string          `gorm:"type:json"`
	CreatedAt       int64           `gorm:"type:bigint"`
	UpdatedAt       int64           `gorm:"type:bigint"`
}

// superadmin表名稱
func (au *SuperAdmin) TableName() string {
	return "super_admin"
}

type SuperAdminIndexReq struct {
	GroupName string `form:"groupname"`
	Username  string `form:"username"`
}

// 新增或更新superadmin權限請求內容
type SuperAdminSaveReq struct {
	Username    string `form:"username" label:"用户名" binding:"required"`
	Uid         uint   `form:"uid"`
	Permissions map[string]bool
}

// 建立種子admin權限內容
func (au *SuperAdmin) FillData(db *gorm.DB) {
	permissionsJSON := `{"permissions": {
		"/admin/setting/adminuser/index:post": true,
		"/ctrl/system/account/rolelist:post": true
	}}`

	groupName := "nobody" // 设置要插入的 GroupName

	// 检查数据库中是否已存在相同 GroupName 的记录
	var existingRecord SuperAdmin
	err := db.Where("group_name = ?", groupName).First(&existingRecord).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// 处理查询错误
		log.Fatalf("Failed to query database: %s", err.Error())
		return
	}

	// 如果已存在相同 GroupName 的记录，则不执行插入
	if err != gorm.ErrRecordNotFound {
		log.Println("SuperAdmin with the same GroupName already exists. Skipping insertion.")
		return
	}

	superAdmin := SuperAdmin{
		GroupName:       groupName,
		PermissionsJSON: permissionsJSON,
	}

	// 使用 Create 方法插入数据
	if err := db.Create(&superAdmin).Error; err != nil {
		// 处理其他插入错误，例如记录日志或返回错误信息
		log.Fatalf("Failed to insert superAdmin: %s", err.Error())
	} else {
		// 插入成功
		log.Println("SuperAdmin inserted successfully")
	}
}

func (au *SuperAdmin) GetConnName() string {
	return "default"
}
