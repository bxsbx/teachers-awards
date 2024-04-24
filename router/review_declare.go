package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func ReviewDeclareRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.ReviewDeclareController{}

	//获取待办事项列表
	router.GET("/v1/activity/review/wait/list", api.GetWaitReviewList)
	//审核列表
	router.GET("/v1/activity/review/list", api.GetReviewList)
	//提交审核
	router.POST("/v1/activity/review/commit", api.CommitReview)
	//获取历史活动评审结果列表
	router.GET("/v1/activity/review/history", api.GetHistoryActivityList)
	//获取奖项设置列表（教育局）
	router.GET("/v1/activity/review/awards/list", api.GetAwardsSetList)
	//名次评定（设置奖项）
	router.POST("/v1/activity/review/set/awards", api.SetAwards)
	//提交活动评审结果
	router.POST("/v1/activity/review/commit/result", api.CommitActivityResult)
	//更新二级指标id
	router.PUT("/v1/activity/review/update/two/id", api.UpdateTwoIndicatorId)
	//教育局给用户添加申报
	router.POST("/v1/activity/review/edb/declare/user", api.EdbDeclareToUser)
	//批量审核通过
	router.POST("/v1/activity/review/batch/pass", api.BatchPass)
	//审核列表（教育局）
	router.GET("/v1/activity/edb/review/list", api.GetEdbReviewList)
	// router general tag
}
