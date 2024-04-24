package controllers

import (
	"github.com/gin-gonic/gin"
	"sort"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/model/resp"
	"teachers-awards/services"
)

type OtherController struct {
}

// @Summary	获取上传token
// @Tags		Other
// @Produce	json
// @Param		user_id		query		string	true	"用户id"
// @Param		file_name	query		string	true	"文件名"
// @Response	200			{object}	resp.Response{data=resp.GetUploadTokenResp}
// @Router		/v1/other/upload/token [GET]
func (u *OtherController) GetUploadToken(c *gin.Context) {
	var params req.GetUploadTokenReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	otherService := services.NewOtherService(global.GetContext(c))
	data, err := otherService.GetUploadToken(params.UserId, params.FileName)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	获取uai表的的操作记录
// @Tags		Other
// @Produce	json
// @Param		user_activity_indicator_id	query		int64	true	"用户活动申报单项目id"
// @Response	200							{object}	resp.Response{data=resp.GetOperationRecordFromUaiResp}
// @Router		/v1/other/operation/record/uai [GET]
func (u *OtherController) GetOperationRecordFromUai(c *gin.Context) {
	var params req.GetOperationRecordFromUaiReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	otherService := services.NewOtherService(global.GetContext(c))
	data, err := otherService.GetOperationRecordFromUai(params.UserActivityIndicatorId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	获取学科枚举
// @Tags		Other
// @Produce	json
// @Param		refresh	query		int	false	"是否刷新学科枚举，0：否，1：是"
// @Response	200		{object}	resp.Response{data=resp.GetSubjectEnumResp}
// @Router		/v1/other/subject/enum [GET]
func (u *OtherController) GetSubjectEnum(c *gin.Context) {
	var params req.GetSubjectEnumReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	OtherService := services.NewOtherService(global.GetContext(c))
	subjectEnumMap, err := OtherService.GetSubjectEnum(params.Refresh)
	var list []resp.Subject
	for k, v := range subjectEnumMap {
		list = append(list, resp.Subject{SubjectCode: k, SubjectName: v})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].SubjectCode < list[j].SubjectCode
	})
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, &resp.GetSubjectEnumResp{SubjectList: list})
}

// @Summary	分页获取中台学校列表
// @Tags		Other
// @Produce	json
// @Param		page	query		int	true	"页数"
// @Param		limit	query		int	true	"每页大小"
// @Response	200		{object}	resp.Response{data=resp.GetSchoolListByPageResp}
// @Router		/v1/other/school/list/page [GET]
func (u *OtherController) GetSchoolListByPage(c *gin.Context) {
	var params req.GetSchoolListByPageReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	OtherService := services.NewOtherService(global.GetContext(c))
	data, err := OtherService.GetSchoolListByPage(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	通过学校名称模糊匹配获取学校列表
// @Tags		Other
// @Produce	json
// @Param		school_name	query		string	true	"学校名称"
// @Response	200			{object}	resp.Response{data=resp.GetSchoolListByNameResp}
// @Router		/v1/other/school/list/by/name [GET]
func (u *OtherController) GetSchoolListByName(c *gin.Context) {
	var params req.GetSchoolListByNameReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	OtherService := services.NewOtherService(global.GetContext(c))
	data, err := OtherService.GetSchoolListByName(params.SchoolName)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}
