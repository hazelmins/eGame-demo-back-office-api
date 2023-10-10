/*
 * @Description: 初始化router 沒啥好說的
 */
package router

import (
	"time"

	"eGame-demo-back-office-api/internal/controllers"
	"eGame-demo-back-office-api/internal/middleware"
	"eGame-demo-back-office-api/pkg/loggers/facade"
	"eGame-demo-back-office-api/pkg/loggers/medium"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Init() (*Router, error) {

	router := NewRouter(gin.Default())

	//设置404错误处理
	router.SetRouteError(controllers.NewHandleController().Handle)

	//设置全局中间件
	router.SetGlobalMiddleware(middleware.Trace(), medium.GinLog(facade.NewLogger("ctrl"), time.RFC3339, false), medium.RecoveryWithLog(facade.NewLogger("admin"), true))

	// 设置后台全局中间件
	store := cookie.NewStore([]byte("1GdFRMs4fcWBvLXT"))
	router.SetAdminRoute(NewAdminRouter(), gzip.Gzip(gzip.DefaultCompression), sessions.Sessions("mysession", store))
	router.SetApiRoute(NewApiRouter())
	return &router, nil
}
