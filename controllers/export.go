package controllers

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/services"
)

type ExportController struct {
}

// @Summary	导出审核列表
// @Tags		Export
// @Produce	octet-stream
// @Param		judges_id			query	string	false	"评委id"
// @Param		user_name			query	string	false	"用户姓名"
// @Param		user_sex			query	int		false	"1：男，2：女"
// @Param		subject_code		query	string	false	"科目code"
// @Param		school_id			query	string	false	"学校id"
// @Param		school_name			query	string	false	"学校名称"
// @Param		year				query	int		false	"年份"
// @Param		one_indicator_name	query	string	false	"一级指标名称"
// @Param		two_indicator_name	query	string	false	"二级指标名称"
// @Param		review_status		query	int		false	"审核状态 0:全部 1:待审核，2：未通过，3：已通过，4：已审核""
// @Param		review_process		query	int		true	"当前审核进程，1：学校，2：专家，3：教育局"
// @Router		/v1/export/review/list [GET]
func (u *ExportController) ExportReviewList(c *gin.Context) {
	var params req.ExportReviewListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	exportService := services.NewExportService(global.GetContext(c))
	f := excelize.NewFile()
	err := exportService.ExportReviewList(&params, f)
	ExportExcelFile(c, f, "审核列表", err)
}

// @Summary	导出历史活动评审结果列表
// @Tags		Export
// @Produce	octet-stream
// @Param		user_name		query	string	false	"用户姓名"
// @Param		user_sex		query	int		false	"1：男，2：女"
// @Param		subject_code	query	string	false	"科目code"
// @Param		school_id		query	string	false	"学校id"
// @Param		school_name		query	string	false	"学校名称"
// @Param		year			query	int		false	"年份"
// @Param		declare_type	query	int		false	"申报类型"
// @Param		identity_card	query	string	false	"身份证号"
// @Param		rank_prize		query	int		false	"0：无，1：一等奖，2：二等奖，3：三等奖"
// @Param		rank			query	int		false	"排名"
// @Param		final_score		query	number	false	"最终得分（各项通过的审核）"
// @Router		/v1/export/history/activity/list [GET]
func (u *ExportController) ExportHistoryActivityList(c *gin.Context) {
	var params req.ExportHistoryActivityListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	exportService := services.NewExportService(global.GetContext(c))
	f := excelize.NewFile()
	err := exportService.ExportHistoryActivityList(&params, f)
	ExportExcelFile(c, f, "历史活动评审结果", err)
}

// @Summary	导出用户的申报记录
// @Tags		Export
// @Produce	octet-stream
// @Param		year			query	int		true	"年份"
// @Param		declare_user_id	query	string	true	"申报的用户id"
// @Param		role			query	int		true	"角色类型，1：学校，2：专家，3：教育局，4：老师"
// @Router		/v1/export/user/declare/record/list [GET]
func (u *ExportController) ExportUserDeclareRecordListByYear(c *gin.Context) {
	var params req.ExportUserDeclareRecordListByYearReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	exportService := services.NewExportService(global.GetContext(c))
	f := excelize.NewFile()
	err := exportService.ExportUserDeclareRecordListByYear(&params, f)
	ExportExcelFile(c, f, "用户的申报记录", err)
}

// @Summary	导出审核列表（教育局）
// @Tags		Export
// @Produce	json
// @Param		user_name		query		string	false	"用户姓名"
// @Param		user_sex		query		int		false	"1：男，2：女"
// @Param		subject_code	query		string	false	"科目code"
// @Param		school_id		query		string	false	"学校id"
// @Param		school_name		query		string	false	"学校名称"
// @Param		status			query		int		false	"审核状态 0:全部 1:待审核，2：已审核"
// @Response	200				{object}	resp.Response
// @Router		/v1/export/edb/review/list [GET]
func (u *ExportController) ExportEdbReviewList(c *gin.Context) {
	var params req.ExportEdbReviewListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	exportService := services.NewExportService(global.GetContext(c))
	f := excelize.NewFile()
	err := exportService.ExportEdbReviewList(&params, f)
	ExportExcelFile(c, f, "审核列表（教育局）", err)
}

// @Summary	导出奖项设置列表（教育局）
// @Tags		Export
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
// @Response	200				{object}	resp.Response
// @Router		/v1/export/awards/set/list [GET]
func (u *ExportController) ExportAwardsSetList(c *gin.Context) {
	var params req.ExportAwardsSetListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	exportService := services.NewExportService(global.GetContext(c))
	f := excelize.NewFile()
	err := exportService.ExportAwardsSetList(&params, f)
	ExportExcelFile(c, f, "获奖列表", err)
}
