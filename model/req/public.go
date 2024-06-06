package req

type GetTokenAndUserInfoReq struct {
	UserId   string `form:"user_id" binding:"required"`         //用户id
	Platform int    `form:"platform" binding:"required"`        //登录平台，1：手机端 2：web端
	JXYToken string `form:"jxy_token" binding:"required,min=7"` //教学研token，600000：中台1.0 公版，6000001：中台2.0 诸暨版，6000002：中台2.0 新昌
}

type GetTokenAndUserInfoByCodeReq struct {
	From        string `form:"from" binding:"required"`         //数据来源，600000：中台1.0 公版，6000001：中台2.0 诸暨版，6000002：中台2.0 新昌
	Code        string `form:"code" binding:"required"`         //验证码
	RedirectUrl string `form:"redirect_url" binding:"required"` //跳转链接
}
