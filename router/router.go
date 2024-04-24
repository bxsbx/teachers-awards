package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/global"
	"teachers-awards/middleware"
)

func Routers() *gin.Engine {
	// 设置gin模式
	gin.SetMode(global.GetEnvMode())
	var Router = gin.New()

	//仅在dev模式下才使用，其他环境下注释掉Swagger
	//swagger.Swagger(Router)

	// 跨域处理
	Router.Use(middleware.Cors())
	// 上下文处理、日志记录
	Router.Use(middleware.Common())

	//不需要授权组
	NoAuthGroup := Router.Group("")
	PublicRouter(NoAuthGroup)

	//需要授权组
	AuthGroup := Router.Group("")
	{
		//jwt授权
		AuthGroup.Use(middleware.JWTAuth())
		//路由权限
		AuthGroup.Use(middleware.RouterAuth())

		AnnouncementRouter(AuthGroup)
		NoticeRouter(AuthGroup)
		IndicatorRouter(AuthGroup)
		ActivityRouter(AuthGroup)
		UserActivityRouter(AuthGroup)
		ReviewDeclareRouter(AuthGroup)
		StatisticsRouter(AuthGroup)
		ExportRouter(AuthGroup)
		OtherRouter(AuthGroup)
		UserRouter(AuthGroup)
	}

	return Router
}
