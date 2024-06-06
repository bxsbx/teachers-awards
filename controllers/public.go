package controllers

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/services"
)

type PublicController struct {
}

//	@Summary	获取token和用户的基本信息
//	@Tags		Public
//	@Produce	json
//	@Param		user_id		query		string	true	"用户id"
//	@Param		platform	query		int		true	"登录平台，1：手机端 2：web端"
//	@Param		jxy_token	query		string	true	"教学研token，600000：中台1.0 公版，6000001：中台2.0 诸暨版"
//	@Response	200			{object}	resp.Response{data=resp.GetTokenAndUserInfoResp}
//	@Router		/v1/public/get/info/token/user [GET]
func (u *PublicController) GetTokenAndUserInfo(c *gin.Context) {
	var params req.GetTokenAndUserInfoReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	publicService := services.NewPublicService(global.GetContext(c))
	data, err := publicService.GetTokenAndUserInfo(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	通过code获取用户信息和token
//	@Tags		Public
//	@Produce	json
//	@Param		from			query		string	true	"数据来源，600000：中台1.0 公版，6000001：中台2.0 诸暨版，6000002：中台2.0 新昌"
//	@Param		code			query		string	true	"验证码"
//	@Param		redirect_url	query		string	true	"跳转链接"
//	@Response	200				{object}	resp.Response{data=resp.GetTokenAndUserInfoResp}
//	@Router		/v1/public/get/user/info/by/code [GET]
func (u *PublicController) GetTokenAndUserInfoByCode(c *gin.Context) {
	var params req.GetTokenAndUserInfoByCodeReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	publicService := services.NewPublicService(global.GetContext(c))
	data, err := publicService.GetTokenAndUserInfoByCode(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}
