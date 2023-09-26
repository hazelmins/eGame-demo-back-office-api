/*
 * @Description:用户组服务
 */
package admin

import (
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

//保存角色
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


//删除角色
func (ser *adminGroupService) DelGroup(id string) (ok bool, err error) {
	polices := casbinauth.GetPoliceByGroup(id)
	ok, err = casbinauth.DelGroups("p", polices)
	return
}
