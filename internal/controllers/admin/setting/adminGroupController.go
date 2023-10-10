/*
 * @Description:用户组管理 進db管理群組權限內容
 */

package setting

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"eGame-demo-back-office-api/internal/controllers/admin"
	"eGame-demo-back-office-api/internal/menu"
	"eGame-demo-back-office-api/internal/models"
	services "eGame-demo-back-office-api/internal/services/admin"
	"eGame-demo-back-office-api/pkg/casbinauth"
	"eGame-demo-back-office-api/pkg/paginater"
	"eGame-demo-back-office-api/pkg/redisx"

	"github.com/gin-gonic/gin"
)

type adminGroupController struct {
	admin.BaseController
}

func NewAdminGroupController() adminGroupController {
	return adminGroupController{}
}

func (con adminGroupController) Routes(rg *gin.RouterGroup) {
	//**/admin/setting/admingroup/index **
	rg.GET("/index", con.index) //列表
	rg.GET("/add", con.addIndex)
	rg.POST("/save", con.save)
	rg.GET("/edit", con.edit)
	rg.GET("/del", con.del)
	//***** 進DB的接口 ******
	rg.POST("/dbindex", con.dbindex)       //單一管理員權限列表
	rg.POST("/dbsave", con.dbsave)         //修改權限
	rg.POST("/onlyindex", con.onlydbindex) //暴力列全部權限 ok

}

/*
*
角色列表
*/
func (con adminGroupController) index(c *gin.Context) {

	var groups []string

	key := c.Query("keyword")
	if key != "" {
		for _, v := range casbinauth.GetGroups() {
			if strings.Contains(v, key) {
				groups = append(groups, key)
			}
		}
	} else {
		groups = casbinauth.GetGroups()
	}

	con.Html(c, http.StatusOK, "setting/group.html", gin.H{
		"adminGroups": groups,
		"keyword":     key,
	})
}

/*
*
添加角色
*/
func (con adminGroupController) addIndex(c *gin.Context) {
	con.Html(c, http.StatusOK, "setting/group_form.html", gin.H{
		"menuList": menu.GetMenu(),
		"id":       "",
	})
}

/*
*
保存角色
*/
func (con adminGroupController) save(c *gin.Context) {

	var req models.AdminGroupSaveReq //request內容
	err := con.FormBind(c, &req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	tmp := make([]string, 0, len(req.Privs))
	for _, v := range req.Privs {
		if strings.Contains(v, "|") {
			tmp = append(tmp, strings.Split(v, "|")...)
		} else {
			tmp = append(tmp, v)
		}
	}
	req.Privs = tmp

	err = services.NewAdminGroupService().SaveGroup(req)
	if err != nil {
		con.Error(c, "操作失败")
		return
	}

	con.Success(c, "/admin/setting/admingroup/index", "操作成功")
}

/*
*
编辑
*/
func (con adminGroupController) edit(c *gin.Context) {
	id := c.Query("id")
	con.Html(c, http.StatusOK, "setting/group_form.html", gin.H{
		"menuList": menu.GetMenu(),
		"id":       id,
	})
}

/*
*
删除
*/
func (con adminGroupController) del(c *gin.Context) {

	id := c.Query("id")
	dbOk, dbErr := services.NewAdminGroupService().DelGroup(id)
	if dbErr != nil || !dbOk {
		con.Error(c, "删除失败")
	} else {
		con.Success(c, "", "删除成功")
	}
}

// *******************進DB操作權限**********************
// 列出指定admin權限內容
// 定義 adminGroupController 類型的 index 方法，用於處理 HTTP GET 請求。
func (con adminGroupController) dbindex(c *gin.Context) {

	var (
		err            error                     // 用于存储可能出现的错误。
		req            models.SuperAdminIndexReq // 请求结构，用于从请求中解析参数。
		adminGropPrivs []models.SuperAdmin       // 群組權限列表，用于存储查询结果。
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

	// 调用 services.NewAdminGroupService() 创建 adminGroupService 的实例，然后调用其 Dao 属性上的 GetAdminGroup 方法。

	adminGropDb := services.NewAdminGroupService().Dao.GetAdminGroup(ctx.(context.Context), req.GroupName, req.Username)

	// 使用 paginater.PageOperation 方法进行分页处理，将分页结果存储到 adminUserData 和 adminUserList 变量中。
	admingroupData, err := paginater.PageOperation(c, adminGropDb, 1, &adminGropPrivs)
	if err != nil {
		// 如果分页处理出现错误，将错误信息返回给客户端并退出函数。
		con.ErrorHtml(c, err)
	}

	// 使用 con.Html 方法返回 HTML 响应，将 adminUserData、c.Query("created_at") 和 c.Query("nickname")
	//作为模板变量传递给模板文件 "setting/adminuser.html"。
	con.Html(c, http.StatusOK, "setting/adminuser.html", gin.H{
		"adminUserPrvis": admingroupData,
		"username":       c.Query("username"),
		"groupname":      c.Query("gropuname"),
	})
}

/*
*從db撈所有群組跟權限********************************看看我這裡************************************************
 */
type SuperAdmin struct {
	GroupName       string          `json:"group_name"`
	PermissionsJSON map[string]bool `json:"permissions"`
}

func (con adminGroupController) onlydbindex(c *gin.Context) {
	var (
		err error
	)

	// 1. 从Authorization标头中获取令牌
	if c.Request.Method == "POST" {
		token := c.GetHeader("token")
		if token == "" {
			// 处理缺少令牌的情况
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少令牌"})
			return
		}

		var userData redisx.UserData
		// 2. 使用令牌从Redis检索用户数据
		userData, err = redisx.GetUserDataFromRedis(token)
		if err != nil {
			// 处理从Redis检索用户数据时出现错误的情况
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法檢索用戶數據"})
			return
		}

		// 检查userData中的groupname是否为"superadmin"
		if userData.Groupname == "superadmin" {
			// 查询数据库以获取组名和权限内容
			adminGroupDb, err := services.NewAdminGroupService().GetGroupIndex()
			if err != nil {
				con.Error(c, "非最高管理superadmin")
				return
			}

			// 构建所需的格式
			groupPermissions := make(map[string]interface{})
			for _, group := range adminGroupDb {
				groupName, ok := group["group_name"].(string)
				if !ok {
					// 处理无效的组名
					continue
				}
				permissionsStr, ok := group["permissions_json"].(string)
				if !ok {
					// 处理无效的权限内容
					continue
				}

				// 解析权限JSON字符串为map
				var permissions map[string]interface{}
				err := json.Unmarshal([]byte(permissionsStr), &permissions)
				if err != nil {
					// 处理解析错误
					continue
				}

				groupPermissions[groupName] = permissions
			}

			// 返回结果给客户端
			c.JSON(http.StatusOK, gin.H{
				"groupname":   userData.Groupname,
				"permissions": groupPermissions,
			})
		}
	}
}

// ****************需要大改的地方 修改指定admin權限******************************
// 定義 adminGroupController 類型的 save 方法，用於處理 HTTP POST 請求。
func (con adminGroupController) dbsave(c *gin.Context) {
	// 声明错误变量和请求结构。
	var (
		err error                    // 用于存储可能出现的错误。
		req models.AdminGroupSaveReq // 请求结构，用于从请求中解析参数。
	)

	// 使用 con.Bind 方法解析请求参数并将其存储到 req 变量中。
	if err := c.Bind(&req); err != nil {
		// 如果解析参数时出现错误，将错误信息返回给客户端并退出函数。
		con.Error(c, err.Error())
		return
	}

	// 调用 services.NewAdminGroupService() 创建 adminGroupService 的实例，然后调用其 SaveGroup 方法来保存或更新群組權限信息。
	err = services.NewAdminGroupService().SaveDbGroup(req)
	if err != nil {
		// 如果保存或更新操作出现错误，将错误信息返回给客户端并退出函数。
		con.Error(c, err.Error())
		return
	}

	// 使用 con.Success 方法返回成功的响应，包括操作成功的消息。
	con.Success(c, "", "操作成功")
}

/*
request body ex
{
    "username": "admin",
    "permissions": {
        "/admin/setting/adminuser/index:get": true,
        "/admin/setting/adminuser/add:get": true,
        "/admin/setting/adminuser/edit:get": true
    }
}
*/
