package controllers

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/services"
)

type UserActivityController struct {
}

//	@Summary	用户获取首页活动列表
//	@Tags		UserActivity
//	@Produce	json
//	@Param		user_id	query		string	true	"用户id"
//	@Param		page	query		int		true	"页数"
//	@Param		limit	query		int		true	"每页大小"
//	@Response	200		{object}	resp.Response{data=resp.GetActivityListToUserResp}
//	@Router		/v1/user/activity/activity/list [GET]
func (u *UserActivityController) GetActivityListToUser(c *gin.Context) {
	var params req.GetActivityListToUserReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	data, err := userActivityService.GetActivityListToUser(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取活动基本信息
//	@Tags		UserActivity
//	@Produce	json
//	@Param		activity_id	query		int	true	"活动id"
//	@Response	200			{object}	resp.Response{data=resp.GetActivityInfoToUserResp}
//	@Router		/v1/user/activity/activity/get [GET]
func (u *UserActivityController) GetActivityInfoToUser(c *gin.Context) {
	var params req.GetActivityInfoToUserReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	data, err := userActivityService.GetActivityInfoToUser(params.ActivityId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	创建用户活动申报
//	@Tags		UserActivity
//	@Produce	json
//	@Param		data	body		req.CreateUserActivityDeclareReq	true	"body请求体"
//	@Response	200		{object}	resp.Response
//	@Router		/v1/user/activity/declare/create [POST]
func (u *UserActivityController) CreateUserActivityDeclare(c *gin.Context) {
	var params req.CreateUserActivityDeclareReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	err := userActivityService.CreateUserActivityDeclare(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

//	@Summary	获取用户活动申报详情
//	@Tags		UserActivity
//	@Produce	json
//	@Param		activity_id	query		int		true	"活动id"
//	@Param		user_id		query		string	true	"用户id"
//	@Response	200			{object}	resp.Response{data=resp.GetUserActivityDeclareDetailResp}
//	@Router		/v1/user/activity/declare/list [GET]
func (u *UserActivityController) GetUserActivityDeclareDetail(c *gin.Context) {
	var params req.GetUserActivityDeclareDetailReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	data, err := userActivityService.GetUserActivityDeclareDetail(params.ActivityId, params.UserId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取用户活动申报状态列表
//	@Tags		UserActivity
//	@Produce	json
//	@Param		activity_id	query		int		true	"活动id"
//	@Param		user_id		query		string	true	"用户id"
//	@Response	200			{object}	resp.Response{data=resp.GetUserActivityDeclareStatusListResp}
//	@Router		/v1/user/activity/declare/status/list [GET]
func (u *UserActivityController) GetUserActivityDeclareStatusList(c *gin.Context) {
	var params req.GetUserActivityDeclareStatusListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	data, err := userActivityService.GetUserActivityDeclareStatusList(params.ActivityId, params.UserId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取用户单个申报的详情
//	@Tags		UserActivity
//	@Produce	json
//	@Param		user_activity_indicator_id	query		int64	true	"用户活动申报下单个项目id"
//	@Response	200							{object}	resp.Response{data=resp.GetUserActivityIndicatorResp}
//	@Router		/v1/user/activity/indicator/declare/detail [GET]
func (u *UserActivityController) GetUserActivityIndicator(c *gin.Context) {
	var params req.GetUserActivityIndicatorReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	data, err := userActivityService.GetUserActivityIndicator(params.UserActivityIndicatorId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	用户活动申报单个项目修改
//	@Tags		UserActivity
//	@Produce	json
//	@Param		data	body		req.UpdateUserActivityIndicatorReq	true	"body请求体"
//	@Response	200		{object}	resp.Response
//	@Router		/v1/user/activity/indicator/update [PUT]
func (u *UserActivityController) UpdateUserActivityIndicator(c *gin.Context) {
	var params req.UpdateUserActivityIndicatorReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	err := userActivityService.UpdateUserActivityIndicator(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

//	@Summary	撤销用户活动单个项目
//	@Tags		UserActivity
//	@Produce	json
//	@Param		user_activity_indicator_id	query		int64	true	"用户活动申报下单个项目id"
//	@Response	200							{object}	resp.Response
//	@Router		/v1/user/activity/indicator/delete [DELETE]
func (u *UserActivityController) DeleteUserActivityIndicator(c *gin.Context) {
	var params req.DeleteUserActivityIndicatorReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	err := userActivityService.DeleteUserActivityIndicator(params.UserActivityIndicatorId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

//	@Summary	获取用户活动申报结果
//	@Tags		UserActivity
//	@Produce	json
//	@Param		activity_id	query		int		true	"活动id"
//	@Param		user_id		query		string	true	"用户id"
//	@Response	200			{object}	resp.Response{data=resp.GetUserActivityDeclareResultResp}
//	@Router		/v1/user/activity/declare/result [GET]
func (u *UserActivityController) GetUserActivityDeclareResult(c *gin.Context) {
	var params req.GetUserActivityDeclareResultReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	data, err := userActivityService.GetUserActivityDeclareResult(params.ActivityId, params.UserId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取用户的申报记录列表
//	@Tags		UserActivity
//	@Produce	json
//	@Param		year			query		int		true	"年份"
//	@Param		declare_user_id	query		string	true	"申报的用户id"
//	@Param		role			query		int		true	"角色类型，1：学校，2：专家，3：教育局，4：老师"
//	@Param		page			query		int		true	"页数"
//	@Param		limit			query		int		true	"每页大小"
//	@Response	200				{object}	resp.Response{data=resp.GetUserDeclareRecordListByYearResp}
//	@Router		/v1/user/activity/declare/record [GET]
func (u *UserActivityController) GetUserDeclareRecordListByYear(c *gin.Context) {
	var params req.GetUserDeclareRecordListByYearReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	data, err := userActivityService.GetUserDeclareRecordListByYear(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取用户的历史申报结果列表
//	@Tags		UserActivity
//	@Produce	json
//	@Param		declare_user_id	query		string	true	"申报的用户id"
//	@Param		page			query		int		true	"页数"
//	@Param		limit			query		int		true	"每页大小"
//	@Response	200				{object}	resp.Response{data=resp.GetUserHistoryDeclareResultListResp}
//	@Router		/v1/user/activity/history/declare/result [GET]
func (u *UserActivityController) GetUserHistoryDeclareResultList(c *gin.Context) {
	var params req.GetUserHistoryDeclareResultListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	data, err := userActivityService.GetUserHistoryDeclareResultList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取用户申报结果(app端)
//	@Tags		UserActivity
//	@Produce	json
//	@Param		activity_id	query		int		true	"活动id"
//	@Param		user_id		query		string	true	"用户id"
//	@Response	200			{object}	resp.Response{data=resp.GetUserDeclareResultAppResp}
//	@Router		/v1/user/activity/declare/result/app [GET]
func (u *UserActivityController) GetUserDeclareResultApp(c *gin.Context) {
	var params req.GetUserDeclareResultAppReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	data, err := userActivityService.GetUserDeclareResultApp(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	获取用户的活动申报项目（教育局）
//	@Tags		UserActivity
//	@Produce	json
//	@Param		user_activity_id	query		int64	true	"用户活动id"
//	@Response	200					{object}	resp.Response{data=resp.GetUserDeclaresToEdbResp}
//	@Router		/v1/user/activity/declares/edb [GET]
func (u *UserActivityController) GetUserDeclaresToEdb(c *gin.Context) {
	var params req.GetUserDeclaresToEdbReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userActivityService := services.NewUserActivityService(global.GetContext(c))
	data, err := userActivityService.GetUserDeclaresToEdb(params.UserActivityId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}
