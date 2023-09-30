/*
 * @Description:用户组相关model
 */

package models

import "eGame-demo-back-office-api/pkg/mysqlx"

type AdminGroupSaveReq struct {
	Privs     []string `form:"privs[]" label:"权限" json:"privs" binding:"required"`
	Username  string   `form:"username" label:"管理者名稱" json:"username" binding:"required"`
	GroupName string   `form:"groupname" label:"用户组名" json:"groupname" binding:"required"`
	GroupId   uint     `form:"groupid"`
}

//speradmin DB表格建立
type SuperAdmin struct {
	// 嵌入 mysqlx.BaseModle 以包含通用的數據庫字段，如 ID、創建時間和更新時間。
	mysqlx.BaseModle

	// 使用 GORM 標籤定義 UID 字段作為主鍵，並自動增加。
	Uid uint `gorm:"primary_key;auto_increment"`

	// 定義用戶名欄位。
	Username string `gorm:"size:100;comment:'用户名'"`

	// 定義權限字段，使用映射來存儲權限路徑和布爾值。
	Permissions map[string]bool `gorm:"-"`

	// 您可以在數據庫中定義一個 JSON 或其他適合的數據庫欄位來存儲權限信息。
	// 以下示例使用 JSON 字段來存儲權限信息。
	PermissionsJSON string `gorm:"type:json"`
}

//superadmin表名稱
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
