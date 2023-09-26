package dao

import (
	"context"
	"strings"
	"sync"

	"eGame-demo-back-office-api/internal/models"
	"eGame-demo-back-office-api/pkg/mysqlx"

	"gorm.io/gorm"
)

type AdminUserDao struct {
	DB *gorm.DB
}

var (
	instanceAdminUser *AdminUserDao
	onceAdminUserDao  sync.Once
)

func NewAdminUserDao() *AdminUserDao {
	onceAdminUserDao.Do(func() {
		instanceAdminUser = &AdminUserDao{DB: mysqlx.GetDB(&models.AdminUsers{})}
	})
	return instanceAdminUser
}

func (dao *AdminUserDao) GetAdminUser(conditions map[string]interface{}) (adminUser models.AdminUsers, err error) {
	err = dao.DB.Where(conditions).First(&adminUser).Error
	return
}

func (dao *AdminUserDao) GetAdminUsers(ctx context.Context, nickname string, created_time string) *gorm.DB {
	// 創建一個初始的數據庫查詢對象
	db := dao.DB.WithContext(ctx).Table("admin_users")

	// 根據提供的條件動態構建查詢
	if nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+nickname+"%")
	}

	if created_time != "" {
		// 解析創建時間範圍
		period := strings.Split(created_time, " ~ ")
		start := period[0] + " 00:00:00"
		end := period[1] + " 23:59:59"

		// 添加創建時間的過濾條件
		db = db.Where("admin_users.created_at BETWEEN ? AND ?", start, end)
	}

	// 返回數據庫查詢對象，以便在其他地方進行進一步操作
	return db
}

func (dao *AdminUserDao) UpdateColumn(uid uint, key, value string) error {
	return dao.DB.Model(&models.AdminUsers{}).Where("uid = ?", uid).UpdateColumn(key, value).Error
}

func (dao *AdminUserDao) UpdateColumns(conditions, field map[string]interface{}, tx *gorm.DB) error {

	if tx != nil {
		return tx.Model(&models.AdminUsers{}).Where(conditions).UpdateColumns(field).Error
	}

	return dao.DB.Model(&models.AdminUsers{}).Where(conditions).UpdateColumns(field).Error
}

func (dao *AdminUserDao) Del(conditions map[string]interface{}) error {
	return dao.DB.Delete(&models.AdminUsers{}, conditions).Error
}
