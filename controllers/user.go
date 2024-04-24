package controllers

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/services"
)

type UserController struct {
}

// @Summary	获取用户信息
// @Tags		User
// @Produce	json
// @Param		user_id	query		string	true	"用户id"
// @Response	200		{object}	resp.Response{data=resp.GetUserInfoResp}
// @Router		/v1/user/get/info [GET]
func (u *UserController) GetUserInfo(c *gin.Context) {
	var params req.GetUserInfoReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userService := services.NewUserService(global.GetContext(c))
	data, err := userService.GetUserInfo(params.UserId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	保存当前用户信息
// @Tags		User
// @Produce	json
// @Param		user_id			formData	string	true	"用户id"
// @Param		user_name		formData	string	true	"用户名称"
// @Param		user_sex		formData	int		true	"1：男，2：女"
// @Param		birthday		formData	string	true	"出生日期"
// @Param		identity_card	formData	string	true	"身份证号"
// @Param		phone			formData	string	true	"手机号"
// @Param		avatar			formData	string	true	"头像"
// @Param		subject_code	formData	string	true	"科目code"
// @Param		school_id		formData	string	true	"学校id"
// @Param		school_name		formData	string	true	"学校名称"
// @Param		declare_type	formData	int		true	"1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研"
// @Response	200				{object}	resp.Response
// @Router		/v1/user/save/info [POST]
func (u *UserController) SaveUserInfo(c *gin.Context) {
	var params req.SaveUserInfoReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userService := services.NewUserService(global.GetContext(c))
	err := userService.SaveUserInfo(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	通过用户名称模糊匹配获取用户列表
// @Tags		User
// @Produce	json
// @Param		user_name	query		string	true	"用户名称"
// @Param		page		query		int		true	"页数"
// @Param		limit		query		int		true	"每页大小"
// @Response	200			{object}	resp.Response{data=resp.GetUsersByNameResp}
// @Router		/v1/user/list/by/name [GET]
func (u *UserController) GetUsersByName(c *gin.Context) {
	var params req.GetUsersByNameReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userService := services.NewUserService(global.GetContext(c))
	data, err := userService.GetUsersByName(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	设置用户角色
// @Tags		User
// @Produce	json
// @Param		data	body		req.SetRoleToUserReq	true	"body请求体"
// @Response	200		{object}	resp.Response
// @Router		/v1/user/set/role [POST]
func (u *UserController) SetRoleToUser(c *gin.Context) {
	var params req.SetRoleToUserReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userService := services.NewUserService(global.GetContext(c))
	err := userService.SetRoleToUser(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	获取已设置角色的用户列表
// @Tags		User
// @Produce	json
// @Param		user_name		query		string	false	"用户名称"
// @Param		user_sex		query		int		false	"1：男，2：女"
// @Param		identity_card	query		string	false	"身份证号"
// @Param		phone			query		string	false	"手机号"
// @Param		school_id		query		string	false	"学校id"
// @Param		school_name		query		string	false	"学校名称"
// @Param		role			query		int		false	"0：全部，1：学校，2：专家，3：教育局，5：超级管理员"
// @Param		page			query		int		true	"页数"
// @Param		limit			query		int		true	"每页大小"
// @Response	200				{object}	resp.Response{data=resp.GetUserListByWhereResp}
// @Router		/v1/user/get/role/list [GET]
func (u *UserController) GetUserListByWhere(c *gin.Context) {
	var params req.GetUserListByWhereReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userService := services.NewUserService(global.GetContext(c))
	data, err := userService.GetUserListByWhere(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	获取专家角色列表
// @Tags		User
// @Produce	json
// @Param		user_name			query		string	false	"用户名称"
// @Param		user_sex			query		int		false	"1：男，2：女"
// @Param		identity_card		query		string	false	"身份证号"
// @Param		phone				query		string	false	"手机号"
// @Param		two_indicator_ids	query		string	false	"二级指标ids，, 隔开"
// @Param		export_auth			query		int		false	"专家是否授权 0：全部 1：未授权 2：已授权"
// @Param		auth_day			query		string	false	"授权日期，格式2006-01-02"
// @Param		page				query		int		true	"页数"
// @Param		limit				query		int		true	"每页大小"
// @Response	200					{object}	resp.Response{data=resp.GetExpertAuthListByWhereResp}
// @Router		/v1/user/get/expert/auth/list [GET]
func (u *UserController) GetExpertAuthListByWhere(c *gin.Context) {
	var params req.GetExpertAuthListByWhereReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	twoIds, err := util.StrToIntArray(params.TwoIndicatorIds, ",")
	if err != nil {
		OutputError(c, err)
		return
	}
	userService := services.NewUserService(global.GetContext(c))
	data, err := userService.GetExpertAuthListByWhere(&params, twoIds)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	给专家授权指标项
// @Tags		User
// @Produce	json
// @Param		data	body		req.SetExpertAuthReq	true	"body请求体"
// @Response	200		{object}	resp.Response
// @Router		/v1/user/set/expert/auth [POST]
func (u *UserController) SetExpertAuth(c *gin.Context) {
	var params req.SetExpertAuthReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userService := services.NewUserService(global.GetContext(c))
	err := userService.SetExpertAuth(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	取消专家授权
// @Tags		User
// @Produce	json
// @Param		user_id	query		string	true	"用户id"
// @Response	200		{object}	resp.Response
// @Router		/v1/user/cancel/expert/auth [DELETE]
func (u *UserController) CancelExpertAuth(c *gin.Context) {
	var params req.CancelExpertAuthReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	userService := services.NewUserService(global.GetContext(c))
	err := userService.CancelExpertAuth(params.UserId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}
