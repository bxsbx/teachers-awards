package controllers

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/services"
)

type ActivityController struct {
}

//	@Summary	创建或者更新活动
//	@Tags		Activity
//	@Produce	json
//
//	@Param		data	body		req.CreateOrUpdateActivityReq	true	"请求body"
//
//	@Response	200		{object}	resp.Response
//	@Router		/v1/activity/save [POST]
func (u *ActivityController) CreateOrUpdateActivity(c *gin.Context) {
	var params req.CreateOrUpdateActivityReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	util.TimeToLocation(params.StartTime)
	util.TimeToLocation(params.EndTime)
	if params.StartTime.After(*params.EndTime) {
		OutputError(c, errorz.CodeMsg(errorz.RESP_PARAM_ERR, "开始时间不能大于结束时间"))
		return
	}
	intArray, err := util.StrToIntArray(params.TwoIndicatorIds, ",")
	if err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	activityService := services.NewActivityService(global.GetContext(c))
	err = activityService.CreateOrUpdateActivity(&params, intArray)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

//	@Summary	获取活动详情
//	@Tags		Activity
//	@Produce	json
//	@Param		activity_id	query		int	true	"活动id"
//	@Response	200			{object}	resp.Response{data=resp.GetActivityDetailResp}
//	@Router		/v1/activity/get/detail [GET]
func (u *ActivityController) GetActivityDetail(c *gin.Context) {
	var params req.GetActivityDetailReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	activityService := services.NewActivityService(global.GetContext(c))
	data, err := activityService.GetActivityDetail(params.ActivityId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	删除活动下的一级指标
//	@Tags		Activity
//	@Produce	json
//	@Param		activity_id			query		int	true	"活动id"
//	@Param		one_indicator_id	query		int	true	"一级指标id"
//	@Response	200					{object}	resp.Response
//	@Router		/v1/activity/one/indicator/delete [DELETE]
func (u *ActivityController) DeleteActivityOneIndicator(c *gin.Context) {
	var params req.DeleteActivityOneIndicatorReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	activityService := services.NewActivityService(global.GetContext(c))
	err := activityService.DeleteActivityOneIndicator(params.ActivityId, params.OneIndicatorId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

//	@Summary	删除活动下的二级指标
//	@Tags		Activity
//	@Produce	json
//	@Param		activity_id			query		int	true	"活动id"
//	@Param		two_indicator_id	query		int	true	"二级指标id"
//	@Response	200					{object}	resp.Response
//	@Router		/v1/activity/two/indicator/delete [DELETE]
func (u *ActivityController) DeleteActivityTwoIndicator(c *gin.Context) {
	var params req.DeleteActivityTwoIndicatorReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	activityService := services.NewActivityService(global.GetContext(c))
	err := activityService.DeleteActivityTwoIndicator(params.ActivityId, params.TwoIndicatorId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

//	@Summary	删除活动
//	@Tags		Activity
//	@Produce	json
//	@Param		activity_id	query		int	true	"活动id"
//	@Response	200			{object}	resp.Response
//	@Router		/v1/activity/delete [DELETE]
func (u *ActivityController) DeleteActivity(c *gin.Context) {
	var params req.DeleteActivityReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	activityService := services.NewActivityService(global.GetContext(c))
	err := activityService.DeleteActivity(params.ActivityId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

//	@Summary	获取活动列表
//	@Tags		Activity
//	@Produce	json
//	@Param		activity_name	query		string	false	"活动名称"
//	@Param		year			query		int		false	"活动名称"
//	@Param		page			query		int		true	"页数"
//	@Param		limit			query		int		true	"每页大小"
//	@Response	200				{object}	resp.Response{data=resp.GetActivityListResp}
//	@Router		/v1/activity/list [GET]
func (u *ActivityController) GetActivityList(c *gin.Context) {
	var params req.GetActivityListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	activityService := services.NewActivityService(global.GetContext(c))
	data, err := activityService.GetActivityList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取活动二级指标列表
//	@Tags		Activity
//	@Produce	json
//	@Param		activity_id			query		int	true	"活动id"
//	@Param		one_indicator_id	query		int	true	"一级指标id"
//	@Response	200					{object}	resp.Response{data=resp.GetActivityTwoIndicatorListResp}
//	@Router		/v1/activity/two/indicator/list [GET]
func (u *ActivityController) GetActivityTwoIndicatorList(c *gin.Context) {
	var params req.GetActivityTwoIndicatorListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	activityService := services.NewActivityService(global.GetContext(c))
	data, err := activityService.GetActivityTwoIndicatorList(params.ActivityId, params.OneIndicatorId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取年度列表
//	@Tags		Activity
//	@Produce	json
//	@Response	200	{object}	resp.Response{data=resp.GetActivityYearListResp}
//	@Router		/v1/activity/year/list [GET]
func (u *ActivityController) GetActivityYearList(c *gin.Context) {
	activityService := services.NewActivityService(global.GetContext(c))
	data, err := activityService.GetActivityYearList()
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取最新活动
//	@Tags		Activity
//	@Produce	json
//	@Param		user_id	query		int	true	"用户id"
//	@Response	200		{object}	resp.Response{data=resp.GetLatestActivityResp}
//	@Router		/v1/activity/latest [GET]
func (u *ActivityController) GetLatestActivity(c *gin.Context) {
	var params req.GetLatestActivityReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	activityService := services.NewActivityService(global.GetContext(c))
	data, err := activityService.GetLatestActivity(params.UserId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	提前结束活动
//	@Tags		Activity
//	@Produce	json
//	@Param		activity_id	formData	int	true	"活动id"
//	@Response	200			{object}	resp.Response
//	@Router		/v1/activity/early/end [PUT]
func (u *ActivityController) EndActivity(c *gin.Context) {
	var params req.EndActivityReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	activityService := services.NewActivityService(global.GetContext(c))
	err := activityService.EndActivity(params.ActivityId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}
