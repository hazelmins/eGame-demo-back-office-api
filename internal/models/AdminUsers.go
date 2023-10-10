/*
 * @Description:用户相关model
 */
package models

import (
	"eGame-demo-back-office-api/pkg/mysqlx"
	"eGame-demo-back-office-api/pkg/utils/strings"

	"gorm.io/gorm"
)

// Adminuser DB表格
type AdminUsers struct {
	mysqlx.BaseModle
	Uid            uint   `gorm:"primary_key;auto_increment"`
	GroupName      string `gorm:"size:20;column:groupname"` // 假設groupname欄位名稱是'groupname'
	Username       string `gorm:"size:100;comment:'用户名'"`
	Nickname       string `gorm:"size:100;comment:'姓名'"`
	Password       string `gorm:"size:200;comment:'密码'"`
	ChangePassword bool   `gorm:"comment:'是否更换密码'"`
	LastLogin      string `gorm:"size:30;comment:'最后登录ip地址'"`
	Salt           string `gorm:"size:32;comment:'密码盐'"`
	ApiToken       string `gorm:"size:32;comment:'用户登录凭证'"`
	CreatedAt      int64  `gorm:"type:bigint"`
	UpdatedAt      int64  `gorm:"type:bigint"`
}

// admin列表
type AdminUserIndexReq struct {
	Nickname  string `form:"nickname"`
	CreatedAt int64  `form:"created_at"`
}

// 新增或更新管理員req請求內容
type AdminUserSaveReq struct {
	Username       string   `form:"username" label:"用户名" binding:"required"`
	Password       string   `form:"password"`
	Nickname       string   `form:"nickname" label:"姓名" binding:"required"`
	ChangePassword bool     `form:"changepassword"`
	GroupName      []string `form:"groupname[]" label:"用户组" binding:"required"`
	Uid            uint     `form:"uid"`
}

type AdminUserEditPassReq struct {
	Uid         int    `form:"id"`
	OldPassword string `form:"old_password" label:"原始密码" binding:"required"`
	NewPassword string `form:"new_password" label:"新密码" binding:"required"`
	SubPassword string `form:"sub_password" label:"确认密码" binding:"required"`
}

type AdminUserSkinReq struct {
	Type  string `form:"type" json:"type"`
	Color string `form:"color" json:"color"`
	Uid   int
}

func (au *AdminUsers) TableName() string {
	return "admin_users"
}

// 建立種子管理員
func (au *AdminUsers) FillData(db *gorm.DB) {

	salt := strings.RandString(6)
	passwordSalt := strings.Encryption("111111", salt)
	adminUser := AdminUsers{
		GroupName:      "superadmin",
		Username:       "admin",
		Nickname:       "管理员",
		Password:       passwordSalt,
		ChangePassword: true,
		LastLogin:      "",
		Salt:           salt,
		ApiToken:       "",
	}
	db.Save(&adminUser)
}

func (au *AdminUsers) GetConnName() string {
	return "default"
}
