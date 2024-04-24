package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func ExportRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.ExportController{}

	//导出审核列表
	router.GET("/v1/export/review/list", api.ExportReviewList)
	//导出历史活动评审结果列表
	router.GET("/v1/export/history/activity/list", api.ExportHistoryActivityList)
	//导出用户的申报记录
	router.GET("/v1/export/user/declare/record/list", api.ExportUserDeclareRecordListByYear)
	//导出审核列表（教育局）
	router.GET("/v1/export/edb/review/list", api.ExportEdbReviewList)
	//导出奖项设置列表（教育局）
	router.GET("/v1/export/awards/set/list", api.ExportAwardsSetList)
	// router general tag
}
