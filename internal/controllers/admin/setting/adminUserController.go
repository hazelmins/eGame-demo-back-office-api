/*
 * @Description:用户管理 支持多权限管理
 */

package setting

import (
	"eGame-demo-back-office-api/internal/controllers/admin"
	"eGame-demo-back-office-api/internal/dao"
	"eGame-demo-back-office-api/internal/models"
	services "eGame-demo-back-office-api/internal/services/admin"
	"eGame-demo-back-office-api/pkg/casbinauth"
	"eGame-demo-back-office-api/pkg/paginater"
	"encoding/json"
	"net/http"

	"context"

	"github.com/gin-gonic/gin"
)

type adminUserController struct {
	admin.BaseController
}

func NewAdminUserController() adminUserController {
	return adminUserController{}
}

func (con adminUserController) Routes(rg *gin.RouterGroup) {
	rg.GET("/index", con.index)
	rg.GET("/add", con.addIndex)
	rg.POST("/save", con.save)
	rg.GET("/edit", con.edit)
	rg.GET("/del", con.del)
}

/**
管理员列表
首先解析了 HTTP 請求中的參數，然後根據這些參數查詢數據庫，
並使用分頁操作處理結果，最後將結果渲染到 HTML 模板中並返回給客戶端。
*/
// 定義 adminUserController 類型的 index 方法，用於處理 HTTP GET 請求。
func (con adminUserController) index(c *gin.Context) {
	// 声明错误变量、请求结构和管理员用户列表变量。
	var (
		err           error                    // 用于存储可能出现的错误。
		req           models.AdminUserIndexReq // 请求结构，用于从请求中解析参数。
		adminUserList []models.AdminUsers      // 管理员用户列表，用于存储查询结果。
	)

	// 使用 con.FormBind 方法解析请求参数并将其存储到 req 变量中。
	err = con.FormBind(c, &req)
	if err != nil {
		// 如果解析参数时出现错误，将错误信息返回给客户端并退出函数。
		con.ErrorHtml(c, err)
		return
	}

	// 从 Gin 上下文中获取请求上下文。
	ctx, _ := c.Get("ctx")

	// 声明一个未使用的 dao.AdminGroupDao 变量，可能是因为后续代码中未使用，可以删除。
	var _ dao.AdminGroupDao

	// 调用 services.NewAdminUserService() 创建 adminUserService 的实例，然后调用其 Dao 属性上的 GetAdminUsers 方法。
	adminDb := services.NewAdminUserService().Dao.GetAdminUsers(ctx.(context.Context), req.Nickname, req.CreatedAt)

	// 使用 paginater.PageOperation 方法进行分页处理，将分页结果存储到 adminUserData 和 adminUserList 变量中。
	adminUserData, err := paginater.PageOperation(c, adminDb, 1, &adminUserList)
	if err != nil {
		// 如果分页处理出现错误，将错误信息返回给客户端并退出函数。
		con.ErrorHtml(c, err)
	}

	// 使用 con.Html 方法返回 HTML 响应，将 adminUserData、c.Query("created_at") 和 c.Query("nickname") 作为模板变量传递给模板文件 "setting/adminuser.html"。
	con.Html(c, http.StatusOK, "setting/adminuser.html", gin.H{
		"adminUserData": adminUserData,
		"created_at":    c.Query("created_at"),
		"nickname":      c.Query("nickname"),
	})
}

/**
添加
*/
// 定義 adminUserController 類型的 addIndex 方法，用於處理 HTTP GET 請求。
func (con adminUserController) addIndex(c *gin.Context) {
	// 使用 con.Html 方法返回 HTML 响应，HTTP 狀態碼為 http.StatusOK。
	con.Html(c, http.StatusOK, "setting/adminuser_form.html", gin.H{
		// 將 "adminGroups" 和 casbinauth.GetGroups() 的返回值作為模板變量傳遞給模板文件 "setting/adminuser_form.html"。
		"adminGroups": casbinauth.GetGroups(),
	})
}

/**
保存
*/
// 定義 adminUserController 類型的 save 方法，用於處理 HTTP POST 請求。
func (con adminUserController) save(c *gin.Context) {
	// 声明错误变量和请求结构。
	var (
		err error                   // 用于存储可能出现的错误。
		req models.AdminUserSaveReq // 请求结构，用于从请求中解析参数。
	)

	// 使用 con.FormBind 方法解析请求参数并将其存储到 req 变量中。
	err = con.FormBind(c, &req)
	if err != nil {
		// 如果解析参数时出现错误，将错误信息返回给客户端并退出函数。
		con.Error(c, err.Error())
		return
	}

	// 调用 services.NewAdminUserService() 创建 adminUserService 的实例，然后调用其 SaveAdminUser 方法来保存或更新管理员信息。
	err = services.NewAdminUserService().SaveAdminUser(req)
	if err != nil {
		// 如果保存或更新操作出现错误，将错误信息返回给客户端并退出函数。
		con.Error(c, err.Error())
		return
	}

	// 使用 con.Success 方法返回成功的响应，包括重定向到 "/admin/setting/adminuser/index" 和操作成功的消息。
	con.Success(c, "/admin/setting/adminuser/index", "操作成功")
}

/**
编辑
*/
// 定義 adminUserController 類型的 edit 方法，用於處理 HTTP GET 請求。
func (con adminUserController) edit(c *gin.Context) {
	// 從 URL 中獲取 id 參數。
	id := c.Query("id")

	// 調用 services.NewAdminUserService() 創建 adminUserService 的實例，然後使用 GetAdminUser 方法根據 id 查詢管理員信息。
	adminUser, _ := services.NewAdminUserService().GetAdminUser(map[string]interface{}{"uid": id})

	// 解析管理員的群組信息並存儲到 groupName 變數中。
	var groupName []string
	json.Unmarshal([]byte(adminUser.GroupName), &groupName)

	// 創建一個空的 groupMap 用於存儲群組信息。
	var groupMap = make(map[string]struct{})

	// 遍歷 groupName，將每個群組名稱存儲到 groupMap 中。
	for _, v := range groupName {
		groupMap[v] = struct{}{}
	}

	// 使用 con.Html 方法返回 HTML 頁面，HTTP 狀態碼為 http.StatusOK。
	con.Html(c, http.StatusOK, "setting/adminuser_form.html", gin.H{
		// 將 "adminGroups" 和 casbinauth.GetGroups() 的返回值作為模板變量傳遞給模板文件 "setting/adminuser_form.html"。
		"adminGroups": casbinauth.GetGroups(),
		// 將管理員信息、groupMap 變數作為模板變量傳遞給模板文件。
		"adminUser": adminUser,
		"groupMap":  groupMap,
	})
}

/**
删除
*/
// 定義 adminUserController 類型的 del 方法，用於處理 HTTP GET 請求。
func (con adminUserController) del(c *gin.Context) {
	// 從 URL 中獲取 id 參數。
	id := c.Query("id")

	// 調用 services.NewAdminUserService() 創建 adminUserService 的實例，然後使用 DelAdminUser 方法刪除指定 id 的管理員信息。
	err := services.NewAdminUserService().DelAdminUser(id)
	if err != nil {
		// 如果刪除操作出現錯誤，將錯誤信息返回給客戶端並退出函數。
		con.Error(c, "删除失败")
		return
	}

	// 使用 con.Success 方法返回成功的响应，不包括重定向 URL，但包括操作成功的消息。
	con.Success(c, "", "删除成功")
}
