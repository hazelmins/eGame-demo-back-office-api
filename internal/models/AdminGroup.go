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
		"/admin/setting/adminuser/index:get": true,
		"/admin/setting/adminuser/add:get": true,
		"/admin/setting/adminuser/edit:get": true
	}}`

	superAdmin := SuperAdmin{
		GroupName:       "superadmin", // 添加用户组名称
		PermissionsJSON: permissionsJSON,
	}

	// 使用 Create 方法插入数据
	if err := db.Create(&superAdmin).Error; err != nil {
		// 处理错误，例如记录日志或返回错误信息
		log.Fatalf("Failed to insert superAdmin: %s", err.Error())
	} else {
		// 插入成功
		log.Println("SuperAdmin inserted successfully")
	}
}
func (au *SuperAdmin) GetConnName() string {
	return "default"
}
