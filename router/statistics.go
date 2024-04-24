package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func StatisticsRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.StatisticsController{}

	//获取首页下总的统计
	router.GET("/v1/statistics/simple/sum/stats", api.GetSimpleSumStats)
	//获取获奖比例
	router.GET("/v1/statistics/rate/award", api.GetAwardRate)
	//获取申报比例
	router.GET("/v1/statistics/rate/declare", api.GetDeclareRate)
	//各年获奖情况
	router.GET("/v1/statistics/award/every/year", api.GetEveryYearAwardNum)
	//各校获奖情况
	router.GET("/v1/statistics/award/every/school", api.GetEverySchoolAwardNum)
	//各类教师获奖情况
	router.GET("/v1/statistics/award/teacher/every/type", api.GetEveryTeacherTypeAwardNum)
	//总表（数据统计）
	router.GET("/v1/statistics/data/detail/year", api.GetYearDeclareAwardRank)
	//各校获奖情况（具体）
	router.GET("/v1/statistics/data/detail/school", api.GetSchoolDeclareAwardRank)
	//各类教师获奖情况（具体）
	router.GET("/v1/statistics/data/detail/teacher/type", api.GetTeacherTypeDeclareAwardRank)
	//获取活动统计结果
	router.GET("/v1/statistics/activity/result", api.GetResultGroupByDeclareType)
	// router general tag
}
