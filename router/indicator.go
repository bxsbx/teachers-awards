package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func IndicatorRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.IndicatorController{}

	//创建或者更新一级指标
	router.POST("/v1/indicator/one/save", api.CreateOrUpdateOneIndicator)
	//获取一级指标列表
	router.GET("/v1/indicator/one/list", api.GetOneIndicatorList)
	//删除一级指标
	router.DELETE("/v1/indicator/one/delete/ids", api.DeleteOneIndicatorByIds)
	//创建或者更新二级指标
	router.POST("/v1/indicator/two/save", api.CreateOrUpdateTwoIndicator)
	//删除二级指标
	router.DELETE("/v1/indicator/two/delete/ids", api.DeleteTwoIndicatorByIds)
	//获取二级指标列表
	router.GET("/v1/indicator/two/list", api.GetTwoIndicatorList)
	// router general tag
}
