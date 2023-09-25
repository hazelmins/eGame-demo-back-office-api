/*
 * @Description:
 * @Author: gphper
 * @Date: 2021-06-01 20:15:04
 */
package router

import (
	"time"

	"eGame-demo-back-office-api/internal/controllers"
	"eGame-demo-back-office-api/internal/middleware"
	"eGame-demo-back-office-api/pkg/loggers/facade"
	"eGame-demo-back-office-api/pkg/loggers/medium"
	"eGame-demo-back-office-api/web"

	"github.com/gin-gonic/gin"
)

func Init() (*Router, error) {

	router := NewRouter(gin.Default())

	//设置404错误处理
	router.SetRouteError(controllers.NewHandleController().Handle)

	//设置全局中间件
	router.SetGlobalMiddleware(middleware.Trace(), medium.GinLog(facade.NewLogger("admin"), time.RFC3339, false), medium.RecoveryWithLog(facade.NewLogger("admin"), true))

	// 设置模板解析函数
	render, err := web.LoadTemplates()
	if err != nil {

		return nil, err
	}
	router.SetHtmlRenderer(render)

	return &router, nil
}
