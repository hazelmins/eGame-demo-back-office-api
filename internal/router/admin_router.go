/*
 * @Description:
 * @Author: gphper
 * @Date: 2021-07-13 19:45:15
 */
package router

import (
	"eGame-demo-back-office-api/internal/controllers/admin"
	"eGame-demo-back-office-api/internal/controllers/admin/article"
	"eGame-demo-back-office-api/internal/controllers/admin/demo"
	"eGame-demo-back-office-api/internal/controllers/admin/setting"
	"eGame-demo-back-office-api/internal/controllers/admin/upload"
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

		adminSystem := adminSettingRouter.Group("/system")
		{
			ar.addRouter(setting.NewAdminSystemController(), adminSystem)
		}

	}

	//Demo演示文件上传
	adminDemoRouter := ar.root.Group("/demo")
	adminDemoRouter.Use(middleware.AdminUserAuth(), middleware.AdminUserPrivs())
	{
		ar.addRouter(demo.NewUploadController(), adminDemoRouter)
	}

	//Article文章管理
	adminArticleRouter := ar.root.Group("/article")
	adminArticleRouter.Use(middleware.AdminUserAuth(), middleware.AdminUserPrivs())
	{
		ar.addRouter(article.NewArticleController(), adminArticleRouter)
	}

	//文件上传
	adminUploadRouter := ar.root.Group("/upload")
	adminUploadRouter.Use(middleware.AdminUserAuth())
	{
		ar.addRouter(upload.NewUploadController(), adminUploadRouter)
	}

}
