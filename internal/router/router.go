/**
 * @FilePath: \ginadmin\internal\router\router.go
 * @Description:
 */
package router

import (
	"net/http"

	"eGame-demo-back-office-api/internal"

	"github.com/gin-gonic/gin"
)

type Router struct {
	r *gin.Engine
}

func NewRouter(r *gin.Engine) Router {
	return Router{
		r: r,
	}
}

func (route Router) SetGlobalMiddleware(middlewares ...gin.HandlerFunc) {
	route.r.Use(middlewares...)
}

func (route Router) SetEngine(app *internal.Application) {
	app.Route = route.r
}

func (route Router) SetAdminRoute(ar *AdminRouter, middlewares ...gin.HandlerFunc) {
	ar.root = route.r.Group("/admin")
	if len(middlewares) > 0 {
		ar.root.Use(middlewares...)
	}
	ar.AddRouters()
}

func (route Router) SetApiRoute(ar *ApiRouter, middlewares ...gin.HandlerFunc) {
	ar.root = route.r.Group("/api")
	if len(middlewares) > 0 {
		ar.root.Use(middlewares...)
	}
	ar.AddRouters()
}

func (route Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route.r.ServeHTTP(w, req)
}

func (route Router) SetRouteError(handler gin.HandlerFunc) {
	route.r.NoMethod(handler)
	route.r.NoRoute(handler)
}
