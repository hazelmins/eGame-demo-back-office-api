/*
 * @Description:
 * @Author: gphper
 * @Date: 2021-11-14 13:39:32
 */
package dao

import (
	"sync"
	"time"

	"eGame-demo-back-office-api/internal/models"
	"eGame-demo-back-office-api/pkg/mysqlx"

	"gorm.io/gorm"
)

type ArticleDao struct {
	DB *gorm.DB
}

var (
	instanceArticle *ArticleDao
	onceArticleDao  sync.Once
)

func NewArticleDao() *ArticleDao {
	onceArticleDao.Do(func() {
		instanceArticle = &ArticleDao{DB: mysqlx.GetDB(&models.Article{})}
	})
	return instanceArticle
}

func (dao *ArticleDao) GetArticle(conditions map[string]interface{}) (article models.Article, err error) {

	err = dao.DB.First(&article, conditions).Error
	return
}

func (dao *ArticleDao) GetArticles(title string, createdAt int64) (db *gorm.DB) {
	db = dao.DB.Table("article")

	if title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}

	if createdAt != 0 {
		// 将时间戳转换为时间对象
		startTime := time.Unix(createdAt, 0)
		endTime := startTime.Add(24 * time.Hour).Add(-time.Second)

		// 添加创建时间的过滤条件
		db = db.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	return
}

func (dao *ArticleDao) UpdateColumns(conditions, field map[string]interface{}, tx *gorm.DB) error {

	if tx != nil {
		return tx.Model(&models.Article{}).Where(conditions).UpdateColumns(field).Error
	}

	return dao.DB.Model(&models.Article{}).Where(conditions).UpdateColumns(field).Error
}

func (dao *ArticleDao) Del(conditions map[string]interface{}) error {
	return dao.DB.Delete(&models.Article{}, conditions).Error
}
