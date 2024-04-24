package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func ActivityRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.ActivityController{}

	//创建或者更新活动
	router.POST("/v1/activity/save", api.CreateOrUpdateActivity)
	//获取活动详情
	router.GET("/v1/activity/get/detail", api.GetActivityDetail)
	//删除活动下的一级指标
	router.DELETE("/v1/activity/one/indicator/delete", api.DeleteActivityOneIndicator)
	//删除活动下的二级指标
	router.DELETE("/v1/activity/two/indicator/delete", api.DeleteActivityTwoIndicator)
	//删除活动
	router.DELETE("/v1/activity/delete", api.DeleteActivity)
	//获取活动列表
	router.GET("/v1/activity/list", api.GetActivityList)
	//获取活动二级指标列表
	router.GET("/v1/activity/two/indicator/list", api.GetActivityTwoIndicatorList)
	//获取年度列表
	router.GET("/v1/activity/year/list", api.GetActivityYearList)
	//获取最新活动
	router.GET("/v1/activity/latest", api.GetLatestActivity)
	//提前结束活动
	router.PUT("/v1/activity/early/end", api.EndActivity)
	// router general tag
}
