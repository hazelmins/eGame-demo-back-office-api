/*
 * @Description:用户服务
 * @Author: gphper
 * @Date: 2021-07-18 13:59:07
 */
package admin

import (
	"encoding/json"
	"errors"
	"sync"

	"eGame-demo-back-office-api/internal/dao"
	"eGame-demo-back-office-api/internal/models"
	"eGame-demo-back-office-api/pkg/casbinauth"
	gstrings "eGame-demo-back-office-api/pkg/utils/strings"
)

// 定義 adminUserService 結構，包含一個指向 AdminUserDao 的指針
type adminUserService struct {
	Dao *dao.AdminUserDao
}

// 創建一個全局的 instanceAdminUserService 變數，用於存儲 adminUserService 的唯一實例
var (
	instanceAdminUserService *adminUserService
	onceAdminUserService     sync.Once
)

// NewAdminUserService 函數用於創建或獲取 adminUserService 的唯一實例
func NewAdminUserService() *adminUserService {
	// 使用 sync.Once 保證該函數只被執行一次，從而實現單例模式
	onceAdminUserService.Do(func() {
		instanceAdminUserService = &adminUserService{
			Dao: dao.NewAdminUserDao(), // 創建 AdminUserDao 實例並設置給 adminUserService 的 Dao 屬性
		}
	})
	return instanceAdminUserService // 返回唯一的 adminUserService 實例
}

// 添加或保存管理员信息
func (ser *adminUserService) SaveAdminUser(req models.AdminUserSaveReq) (err error) {
	// 定義變數
	var (
		adminUser models.AdminUsers
		ok        bool
	)

	// 將 req.GroupName 轉換為 JSON 字符串
	groupnameStr, _ := json.Marshal(req.GroupName)

	// 創建一個二維字符串切片 rules，用於設置管理員的組（群組）規則
	var rules = make([][]string, 0)
	for _, v := range req.GroupName {
		rules = append(rules, []string{req.Username, v})
	}

	// 開始事務，以確保操作的原子性
	tx := ser.Dao.DB.Begin()

	if req.Uid > 0 {
		// 更新現有管理員信息
		var groupOldName []string
		adminUser, err = ser.Dao.GetAdminUser(map[string]interface{}{"uid": req.Uid})
		if err != nil {
			return
		}

		// 解析原始的 GroupName 字符串
		json.Unmarshal([]byte(adminUser.GroupName), &groupOldName)

		// 設置要更新的字段
		fields := map[string]interface{}{
			"group_name": string(groupnameStr),
			"username":   req.Username,
			"nickname":   req.Nickname,
			"phone":      req.Phone,
		}

		// 如果 req.Password 不為空，則更新密碼
		if req.Password != "" {
			salt := gstrings.RandString(6)
			fields["salt"] = salt
			fields["password"] = gstrings.Encryption(req.Password, salt)
		}

		// 更新管理員信息
		err = ser.Dao.UpdateColumns(map[string]interface{}{"uid": req.Uid}, fields, tx)
		if err != nil {
			tx.Rollback()
			return
		}

		// 更新管理員的組（群組）規則
		_, err = casbinauth.UpdateGroups(req.Username, groupOldName, req.GroupName, tx)
		if err != nil {
			tx.Rollback()
			return
		}

	} else {
		// 創建新的管理員信息
		salt := gstrings.RandString(6)
		passwordSalt := gstrings.Encryption(req.Password, salt)
		adminUser := models.AdminUsers{
			GroupName: string(groupnameStr),
			Nickname:  req.Nickname,
			Username:  req.Username,
			Password:  passwordSalt,
			Phone:     req.Phone,
			Salt:      salt,
		}
		err = tx.Save(&adminUser).Error
		if err != nil {
			tx.Rollback()
			return
		}

		// 將管理員的組（群組）規則添加到 casbin 中
		ok, err = casbinauth.AddGroups("g", rules, tx)
		if err != nil || !ok {
			tx.Rollback()
			return
		}
	}

	// 提交事務
	tx.Commit()
	return
}

//获取单个管理员用户信息
func (ser *adminUserService) GetAdminUser(conditions map[string]interface{}) (adminUser models.AdminUsers, err error) {
	adminUser, err = ser.Dao.GetAdminUser(conditions)
	return
}

//删除管理员
func (ser *adminUserService) DelAdminUser(id string) (err error) {
	return ser.Dao.Del(map[string]interface{}{"uid": id})
}

//修改密码
func (ser *adminUserService) EditPass(req models.AdminUserEditPassReq) (err error) {

	var adminUser models.AdminUsers

	if req.NewPassword != req.SubPassword {
		err = errors.New("请再次确认新密码是否正确")
		return
	}

	adminUser, err = ser.GetAdminUser(map[string]interface{}{"uid": req.Uid})
	if err != nil {
		return
	}

	oldPass := gstrings.Encryption(req.OldPassword, adminUser.Salt)
	if oldPass != adminUser.Password {
		err = errors.New("原密码错误")
		return
	}

	newPass := gstrings.Encryption(req.NewPassword, adminUser.Salt)
	err = ser.Dao.UpdateColumn(adminUser.Uid, "password", newPass)

	return
}

//根究用户保存自定义皮肤
func (ser *adminUserService) EditSkin(req models.AdminUserSkinReq) (err error) {

	var skinMap = map[string]string{
		"data-logobg":    "logo",
		"data-sidebarbg": "side",
		"data-headerbg":  "header",
	}

	err = ser.Dao.UpdateColumn(uint(req.Uid), skinMap[req.Type], req.Color)

	return
}
