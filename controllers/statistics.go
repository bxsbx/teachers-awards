package controllers

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/services"
)

type StatisticsController struct {
}

//	@Summary	获取首页下总的统计
//	@Tags		Statistics
//	@Produce	json
//	@Param		year		query		int		true	"年份"
//	@Param		school_id	query		string	false	"学校id"
//	@Response	200			{object}	resp.Response{data=resp.GetSimpleSumStatsResp}
//	@Router		/v1/statistics/simple/sum/stats [GET]
func (u *StatisticsController) GetSimpleSumStats(c *gin.Context) {
	var params req.GetSimpleSumStatsReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	statisticsService := services.NewStatisticsService(global.GetContext(c))
	data, err := statisticsService.GetSimpleSumStats(params.Year, params.SchoolId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取获奖比例
//	@Tags		Statistics
//	@Produce	json
//	@Param		year		query		int		false	"年份"
//	@Param		school_id	query		string	false	"学校id"
//	@Response	200			{object}	resp.Response{data=resp.GetAwardRateResp}
//	@Router		/v1/statistics/rate/award [GET]
func (u *StatisticsController) GetAwardRate(c *gin.Context) {
	var params req.GetAwardRateReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	statisticsService := services.NewStatisticsService(global.GetContext(c))
	data, err := statisticsService.GetAwardRate(params.Year, params.SchoolId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取申报比例
//	@Tags		Statistics
//	@Produce	json
//	@Param		year		query		int		false	"年份"
//	@Param		school_id	query		string	false	"学校id"
//	@Response	200			{object}	resp.Response{data=resp.GetDeclareRateResp}
//	@Router		/v1/statistics/rate/declare [GET]
func (u *StatisticsController) GetDeclareRate(c *gin.Context) {
	var params req.GetDeclareRateReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	statisticsService := services.NewStatisticsService(global.GetContext(c))
	data, err := statisticsService.GetDeclareRate(params.Year, params.SchoolId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	各年获奖情况
//	@Tags		Statistics
//	@Produce	json
//	@Param		school_id	query		string	false	"学校id"
//	@Response	200			{object}	resp.Response{data=resp.GetEveryYearAwardNumResp}
//	@Router		/v1/statistics/award/every/year [GET]
func (u *StatisticsController) GetEveryYearAwardNum(c *gin.Context) {
	var params req.GetEveryYearAwardNumReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	statisticsService := services.NewStatisticsService(global.GetContext(c))
	data, err := statisticsService.GetEveryYearAwardNum(params.SchoolId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	各校获奖情况
//	@Tags		Statistics
//	@Produce	json
//	@Param		year	query		int	false	"年份"
//	@Response	200		{object}	resp.Response{data=resp.GetEverySchoolAwardNumResp}
//	@Router		/v1/statistics/award/every/school [GET]
func (u *StatisticsController) GetEverySchoolAwardNum(c *gin.Context) {
	var params req.GetEverySchoolAwardNumReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	statisticsService := services.NewStatisticsService(global.GetContext(c))
	data, err := statisticsService.GetEverySchoolAwardNum(params.Year)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	各类教师获奖情况
//	@Tags		Statistics
//	@Produce	json
//	@Param		year	query		int	false	"年份"
//	@Response	200		{object}	resp.Response{data=resp.GetEveryTeacherTypeAwardNumResp}
//	@Router		/v1/statistics/award/teacher/every/type [GET]
func (u *StatisticsController) GetEveryTeacherTypeAwardNum(c *gin.Context) {
	var params req.GetEveryTeacherTypeAwardNumReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	statisticsService := services.NewStatisticsService(global.GetContext(c))
	data, err := statisticsService.GetEveryTeacherTypeAwardNum(params.Year)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	总表（数据统计）
//	@Tags		Statistics
//	@Produce	json
//	@Param		school_id	query		string	false	"学校id"
//	@Response	200			{object}	resp.Response{data=resp.GetYearDeclareAwardRankResp}
//	@Router		/v1/statistics/data/detail/year [GET]
func (u *StatisticsController) GetYearDeclareAwardRank(c *gin.Context) {
	var params req.GetYearDeclareAwardRankReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	statisticsService := services.NewStatisticsService(global.GetContext(c))
	data, err := statisticsService.GetYearDeclareAwardRank(params.SchoolId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	各校获奖情况（具体）
//	@Tags		Statistics
//	@Produce	json
//	@Param		school_name	query		string	false	"学校名称"
//	@Param		year		query		int		true	"年份"
//	@Param		page		query		int		true	"页数"
//	@Param		limit		query		int		true	"每页大小"
//	@Response	200			{object}	resp.Response{data=resp.GetSchoolDeclareAwardRankResp}
//	@Router		/v1/statistics/data/detail/school [GET]
func (u *StatisticsController) GetSchoolDeclareAwardRank(c *gin.Context) {
	var params req.GetSchoolDeclareAwardRankReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	statisticsService := services.NewStatisticsService(global.GetContext(c))
	data, err := statisticsService.GetSchoolDeclareAwardRank(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	各类教师获奖情况（具体）
//	@Tags		Statistics
//	@Produce	json
//	@Param		year	query		int	true	"年份"
//	@Response	200		{object}	resp.Response{data=resp.GetTeacherTypeDeclareAwardRankResp}
//	@Router		/v1/statistics/data/detail/teacher/type [GET]
func (u *StatisticsController) GetTeacherTypeDeclareAwardRank(c *gin.Context) {
	var params req.GetTeacherTypeDeclareAwardRankReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	statisticsService := services.NewStatisticsService(global.GetContext(c))
	data, err := statisticsService.GetTeacherTypeDeclareAwardRank(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取活动统计结果
//	@Tags		Statistics
//	@Produce	json
//	@Param		activity_id	query		int	true	"活动id"
//	@Response	200			{object}	resp.Response{data=resp.GetResultGroupByDeclareTypeResp}
//	@Router		/v1/statistics/activity/result [GET]
func (u *StatisticsController) GetResultGroupByDeclareType(c *gin.Context) {
	var params req.GetResultGroupByDeclareTypeReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	statisticsService := services.NewStatisticsService(global.GetContext(c))
	data, err := statisticsService.GetResultGroupByDeclareType(params.ActivityId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}
