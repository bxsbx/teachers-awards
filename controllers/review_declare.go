package controllers

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/services"
)

type ReviewDeclareController struct {
}

// @Summary	获取待办事项列表
// @Tags		ReviewDeclare
// @Produce	json
// @Param		review_process	query		int	true	"当前审核进程，1：学校，2：专家，3：教育局，4：结束"
// @Param		page			query		int	true	"页数"
// @Param		limit			query		int	true	"每页大小"
// @Response	200				{object}	resp.Response{data=resp.GetWaitReviewListResp}
// @Router		/v1/activity/review/wait/list [GET]
func (u *ReviewDeclareController) GetWaitReviewList(c *gin.Context) {
	var params req.GetWaitReviewListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	data, err := reviewDeclareService.GetWaitReviewList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	审核列表
// @Tags		ReviewDeclare
// @Produce	json
// @Param		judges_id			query		string	false	"评委id"
// @Param		user_name			query		string	false	"用户姓名"
// @Param		user_sex			query		int		false	"1：男，2：女"
// @Param		subject_code		query		string	false	"科目code"
// @Param		school_id			query		string	false	"学校id"
// @Param		school_name			query		string	false	"学校名称"
// @Param		year				query		int		false	"年份"
// @Param		one_indicator_name	query		string	false	"一级指标名称"
// @Param		two_indicator_name	query		string	false	"二级指标名称"
// @Param		review_status		query		int		false	"审核状态 0:全部 1:待审核，2：未通过，3：已通过，4：已审核"
// @Param		review_process		query		int		true	"当前审核进程，1：学校，2：专家，3：教育局"
// @Param		only_count			query		bool	false	"是否仅获取数量"
// @Param		page				query		int		true	"页数"
// @Param		limit				query		int		true	"每页大小"
// @Response	200					{object}	resp.Response{data=resp.GetReviewListResp}
// @Router		/v1/activity/review/list [GET]
func (u *ReviewDeclareController) GetReviewList(c *gin.Context) {
	var params req.GetReviewListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	data, err := reviewDeclareService.GetReviewList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	提交审核
// @Tags		ReviewDeclare
// @Produce	json
// @Param		user_activity_indicator_id	formData	int64	true	"用户活动申报二级指标id"
// @Param		is_pass						formData	int		true	"0：未通过，1：通过"
// @Param		opinion						formData	string	true	"审核意见"
// @Response	200							{object}	resp.Response
// @Router		/v1/activity/review/commit [POST]
func (u *ReviewDeclareController) CommitReview(c *gin.Context) {
	var params req.CommitReviewReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	err := reviewDeclareService.CommitReview(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	获取历史活动评审结果列表
// @Tags		ReviewDeclare
// @Produce	json
// @Param		user_name		query		string	false	"用户姓名"
// @Param		user_sex		query		int		false	"1：男，2：女"
// @Param		subject_code	query		string	false	"科目code"
// @Param		school_id		query		string	false	"学校id"
// @Param		school_name		query		string	false	"学校名称"
// @Param		year			query		int		false	"年份"
// @Param		declare_type	query		int		false	"申报类型"
// @Param		identity_card	query		string	false	"身份证号"
// @Param		rank_prize		query		int		false	"0：无，1：一等奖，2：二等奖，3：三等奖"
// @Param		rank			query		int		false	"排名"
// @Param		final_score		query		number	false	"最终得分（各项通过的审核）"
// @Param		page			query		int		true	"页数"
// @Param		limit			query		int		true	"每页大小"
// @Response	200				{object}	resp.Response{data=resp.GetHistoryActivityListResp}
// @Router		/v1/activity/review/history [GET]
func (u *ReviewDeclareController) GetHistoryActivityList(c *gin.Context) {
	var params req.GetHistoryActivityListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	data, err := reviewDeclareService.GetHistoryActivityList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	获取奖项设置列表（教育局）
// @Tags		ReviewDeclare
// @Produce	json
// @Param		activity_id		query		int		true	"活动id"
// @Param		user_name		query		string	false	"用户姓名"
// @Param		user_sex		query		int		false	"1：男，2：女"
// @Param		subject_code	query		string	false	"科目code"
// @Param		school_id		query		string	false	"学校id"
// @Param		school_name		query		string	false	"学校名称"
// @Param		declare_type	query		int		false	"申报类型"
// @Param		identity_card	query		string	false	"身份证号"
// @Param		rank_prize		query		int		false	"名次，0：无，1：一等奖，2：二等奖，3：三等奖"
// @Param		final_score		query		number	false	"最终得分（各项通过的审核）"
// @Param		page			query		int		true	"页数"
// @Param		limit			query		int		true	"每页大小"
// @Response	200				{object}	resp.Response{data=resp.GetAwardsSetListResp}
// @Router		/v1/activity/review/awards/list [GET]
func (u *ReviewDeclareController) GetAwardsSetList(c *gin.Context) {
	var params req.GetAwardsSetListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	data, err := reviewDeclareService.GetAwardsSetList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	名次评定（设置奖项）
// @Tags		ReviewDeclare
// @Produce	json
// @Param		activity_id	formData	int		true	"活动id"
// @Param		user_id		formData	string	true	"用户id"
// @Param		rank_prize	formData	int		true	"名次"
// @Param		prize		formData	int		true	"奖金"
// @Response	200			{object}	resp.Response
// @Router		/v1/activity/review/set/awards [POST]
func (u *ReviewDeclareController) SetAwards(c *gin.Context) {
	var params req.SetAwardsReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	err := reviewDeclareService.SetAwards(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	提交活动评审结果
// @Tags		ReviewDeclare
// @Produce	json
// @Param		activity_id		formData	int	true	"活动id"
// @Param		declare_type	formData	int	true	"申报类型"
// @Response	200				{object}	resp.Response
// @Router		/v1/activity/review/commit/result [POST]
func (u *ReviewDeclareController) CommitActivityResult(c *gin.Context) {
	var params req.CommitActivityResultReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	err := reviewDeclareService.CommitActivityResult(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	更新二级指标id
// @Tags		ReviewDeclare
// @Produce	json
// @Param		data	body		req.UpdateTwoIndicatorIdReq	true	"body请求体"
// @Response	200		{object}	resp.Response
// @Router		/v1/activity/review/update/two/id [PUT]
func (u *ReviewDeclareController) UpdateTwoIndicatorId(c *gin.Context) {
	var params req.UpdateTwoIndicatorIdReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	util.TimeToLocation(params.AwardDate)
	util.TimeToLocation(params.CertificateStartDate)
	util.TimeToLocation(params.CertificateEndDate)
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	err := reviewDeclareService.UpdateTwoIndicatorId(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	教育局给用户添加申报
// @Tags		ReviewDeclare
// @Produce	json
// @Param		data	body		req.EdbDeclareToUserReq	true	"body请求体"
// @Response	200		{object}	resp.Response
// @Router		/v1/activity/review/edb/declare/user [POST]
func (u *ReviewDeclareController) EdbDeclareToUser(c *gin.Context) {
	var params req.EdbDeclareToUserReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	util.TimeToLocation(params.AwardDate)
	util.TimeToLocation(params.CertificateStartDate)
	util.TimeToLocation(params.CertificateEndDate)
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	err := reviewDeclareService.EdbDeclareToUser(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	批量审核通过
// @Tags		ReviewDeclare
// @Produce	json
// @Param		data	body		req.BatchPassReq	true	"body请求体"
// @Response	200		{object}	resp.Response
// @Router		/v1/activity/review/batch/pass [POST]
func (u *ReviewDeclareController) BatchPass(c *gin.Context) {
	var params req.BatchPassReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	err := reviewDeclareService.BatchPass(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	审核列表（教育局）
// @Tags		ReviewDeclare
// @Produce	json
// @Param		user_name		query		string	false	"用户姓名"
// @Param		user_sex		query		int		false	"1：男，2：女"
// @Param		subject_code	query		string	false	"科目code"
// @Param		school_id		query		string	false	"学校id"
// @Param		school_name		query		string	false	"学校名称"
// @Param		status			query		int		false	"审核状态 0:全部 1:待审核，2：已审核"
// @Param		only_count		query		bool	false	"是否仅获取数量"
// @Param		page			query		int		true	"页数"
// @Param		limit			query		int		true	"每页大小"
// @Response	200				{object}	resp.Response{data=resp.GetEdbReviewListResp}
// @Router		/v1/activity/edb/review/list [GET]
func (u *ReviewDeclareController) GetEdbReviewList(c *gin.Context) {
	var params req.GetEdbReviewListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	reviewDeclareService := services.NewReviewDeclareService(global.GetContext(c))
	data, err := reviewDeclareService.GetEdbReviewList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}
