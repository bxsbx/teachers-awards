package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func UserActivityRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.UserActivityController{}

	//用户获取首页活动列表
	router.GET("/v1/user/activity/activity/list", api.GetActivityListToUser)
	//获取活动基本信息
	router.GET("/v1/user/activity/activity/get", api.GetActivityInfoToUser)
	//创建用户活动申报
	router.POST("/v1/user/activity/declare/create", api.CreateUserActivityDeclare)
	//获取用户活动申报详情
	router.GET("/v1/user/activity/declare/list", api.GetUserActivityDeclareDetail)
	//获取用户活动申报状态列表
	router.GET("/v1/user/activity/declare/status/list", api.GetUserActivityDeclareStatusList)
	//获取用户单个申报的详情
	router.GET("/v1/user/activity/indicator/declare/detail", api.GetUserActivityIndicator)
	//用户活动申报单个项目修改
	router.PUT("/v1/user/activity/indicator/update", api.UpdateUserActivityIndicator)
	//撤销用户活动单个项目
	router.DELETE("/v1/user/activity/indicator/delete", api.DeleteUserActivityIndicator)
	//获取用户活动申报结果
	router.GET("/v1/user/activity/declare/result", api.GetUserActivityDeclareResult)
	//获取用户的申报记录列表
	router.GET("/v1/user/activity/declare/record", api.GetUserDeclareRecordListByYear)
	//获取用户的历史申报结果列表
	router.GET("/v1/user/activity/history/declare/result", api.GetUserHistoryDeclareResultList)
	//获取用户申报结果(app端)
	router.GET("/v1/user/activity/declare/result/app", api.GetUserDeclareResultApp)
	//获取用户的活动申报项目（教育局）
	router.GET("/v1/user/activity/declares/edb", api.GetUserDeclaresToEdb)
	// router general tag
}
