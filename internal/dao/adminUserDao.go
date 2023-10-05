//資料庫操作

package dao

import (
	"context"
	"sync"
	"time"

	"eGame-demo-back-office-api/internal/models"
	"eGame-demo-back-office-api/pkg/mysqlx"

	"gorm.io/gorm"
)

type AdminUserDao struct {
	DB *gorm.DB
}

type SuperAdminDao struct {
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

// 關聯了GroupName 為了在login返回權限
func (dao *AdminUserDao) GetAdminUser(conditions map[string]interface{}) (adminUser models.AdminUsers, err error) {
	err = dao.DB.Where(conditions).First(&adminUser).Error
	return
}

func (dao *AdminUserDao) GetAdminUsers(ctx context.Context, nickname string, created_time int64) *gorm.DB {
	// 创建一个初始的数据库查询对象
	db := dao.DB.WithContext(ctx).Table("admin_users")

	// 根据提供的条件动态构建查询
	if nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+nickname+"%")
	}

	if created_time != 0 {
		// 将时间戳转换为时间对象
		startTime := time.Unix(created_time, 0)
		endTime := startTime.Add(24 * time.Hour).Add(-time.Second)

		// 添加创建时间的过滤条件
		db = db.Where("admin_users.created_at BETWEEN ? AND ?", startTime, endTime)
	}

	// 返回数据库查询对象，以便在其他地方进行进一步操作
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
