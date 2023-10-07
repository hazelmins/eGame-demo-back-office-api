/*
 * @Description:router路徑
 */
package router

import (
	"eGame-demo-back-office-api/internal/controllers/admin"
	"eGame-demo-back-office-api/internal/controllers/admin/setting"
	"eGame-demo-back-office-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AdminRouter struct {
	root *gin.RouterGroup
}

func NewAdminRouter() *AdminRouter {
	return &AdminRouter{}
}

func (ar AdminRouter) addRouter(con IAdminController, router *gin.RouterGroup) {
	con.Routes(router)
}

func (ar AdminRouter) AddRouters() {

	ar.addRouter(admin.NewLoginController(), ar.root)

	adminHomeRouter := ar.root.Group("/home")
	adminHomeRouter.Use(middleware.AdminUserAuth())
	{
		ar.addRouter(admin.NewHomeController(), adminHomeRouter)
	}

	adminSettingRouter := ar.root.Group("/setting")
	adminSettingRouter.Use(middleware.AdminUserAuth(), middleware.AdminUserPrivs())
	{
		adminGroup := adminSettingRouter.Group("/admingroup")
		{
			ar.addRouter(setting.NewAdminGroupController(), adminGroup)
		}

		adminUser := adminSettingRouter.Group("/adminuser")
		{
			ar.addRouter(setting.NewAdminUserController(), adminUser)
		}

	}

}
