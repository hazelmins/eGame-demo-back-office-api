/*
 * @Description:用户组服务 只有單純管理群組 DB中的群組權限
 */
package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"eGame-demo-back-office-api/internal/dao"
	"eGame-demo-back-office-api/internal/models"
	"eGame-demo-back-office-api/pkg/casbinauth"
)

/*
定義了一個名為 adminGroupService 的結構體，它包含一個指向 dao.AdminGroupDao 的指針 Dao
*/
type adminGroupService struct {
	Dao *dao.AdminGroupDao
}

/*
定義了兩個變數：
instanceAdminGroupService：用於存儲單例 adminGroupService 的指針。
onceAdminGroupService：用於確保 instanceAdminGroupService 只被創建一次的同步變數。
*/
var (
	instanceAdminGroupService *adminGroupService
	onceAdminGroupService     sync.Once
)

/*
定義了一個名為 NewAdminGroupService 的函式，用於創建並返回 adminGroupService 的單例。
使用 sync.Once 確保這個函式只會被執行一次，並且在第一次調用時初始化 instanceAdminGroupService。
在初始化時，創建一個 adminGroupService 對象，並將其 Dao 成員設置為 dao.NewAdminGroupDao()，然後返回該對象。
*/
func NewAdminGroupService() *adminGroupService {
	onceAdminGroupService.Do(func() {
		instanceAdminGroupService = &adminGroupService{
			Dao: dao.NewAdminGroupDao(),
		}
	})
	return instanceAdminGroupService
}

type Permissions struct {
	Permissions map[string]bool `json:"permissions"`
}

// *********************撈取群組權限per groupName********************************************************
func (ser *adminGroupService) GetAdminGroup(groupName string) (map[string]bool, uint, error) {
	adminGroup, err := ser.Dao.GetPermissionsByGroupName(groupName)
	if err != nil {
		fmt.Printf("GORM 查詢錯誤：%v\n", err)
		return nil, 0, err
	}

	// 解析权限数据
	var permissions Permissions
	err = json.Unmarshal([]byte(adminGroup.PermissionsJSON), &permissions)
	if err != nil {
		fmt.Printf("解析 JSON 錯誤：%v\n", err)
		return nil, 0, err
	}

	// 返回解析后的权限数据和 Uid
	return permissions.Permissions, adminGroup.Uid, nil
}

// *******************************撈取全部群組以及權限***********************************
func (ser *adminGroupService) GetGroupIndex() ([]map[string]interface{}, error) {
	// 获取查询结果
	db := ser.Dao.GetGroupDBIndex(context.TODO()) // 使用context.TODO代替nil

	var results []map[string]interface{}
	if err := db.Select("group_name, permissions_json").Find(&results).Error; err != nil {
		fmt.Printf("GORM 查詢錯誤：%v\n", err)
		return nil, err
	}

	return results, nil
}

// ************************************ 保存角色
func (ser *adminGroupService) SaveGroup(req models.AdminGroupSaveReq) error {
	// 從 casbinauth 模組中獲取特定群組的舊角色信息。
	oldGroup := casbinauth.GetPoliceByGroup(req.GroupName)

	// 計算舊角色信息的數量。
	oldLen := len(oldGroup)

	// 創建一個字串切片，用於存儲角色信息。
	oldSlice := make([]string, oldLen)

	// 檢查是否有舊的角色信息（oldLen 大於 0）。
	if oldLen > 0 {
		// 迭代處理舊角色信息。
		for oldk, oldv := range oldGroup {
			// 將處理後的角色信息轉換為字串，並存儲在 oldSlice 中。
			oldSlice[oldk] = oldv[1] + ":" + oldv[2]
		}
	}

	// 開始一個數據庫事務，以便後續的數據庫操作可以作為一個事務進行。
	tx := ser.Dao.DB.Begin()

	// 調用 casbinauth.UpdatePolices 函式來更新角色的權限。
	// 如果更新失敗，則回滾事務並返回錯誤。
	_, err := casbinauth.UpdatePolices(req.GroupName, oldSlice, req.Privs, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事務，確保更改生效。
	tx.Commit()

	// 如果一切正常，函式返回 nil 表示成功。
	return nil
}

// 删除角色
func (ser *adminGroupService) DelGroup(id string) (ok bool, err error) {
	polices := casbinauth.GetPoliceByGroup(id)
	ok, err = casbinauth.DelGroups("p", polices)
	return
}

//**************************用不到分隔線***************************

func (ser *adminGroupService) SaveDbGroup(req models.AdminGroupSaveReq) error {
	tx := ser.Dao.DB.Begin()

	// 將需要的參數傳遞給 casbinauth.UpdateUserPolices 函數
	// 注意：這裡傳遞了三個 []string 參數
	_, err := casbinauth.UpdateUserPolices(req.GroupName, req.Username, req.Privs, []string{}, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 在這裡處理將 req.Privs 映射到 Permissions 字段的邏輯
	var superAdmin models.SuperAdmin

	superAdmin.Permissions = make(map[string]bool)
	for _, priv := range req.Privs {
		superAdmin.Permissions[priv] = true
	}

	// 使用 GORM 保存或更新 SuperAdmin 記錄，將 Permissions 映射到 PermissionsJSON 字段
	permissionsJSON, err := json.Marshal(superAdmin.Permissions)
	if err != nil {
		tx.Rollback()
		return err
	}
	superAdmin.PermissionsJSON = string(permissionsJSON)

	if err := tx.Save(&superAdmin).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
