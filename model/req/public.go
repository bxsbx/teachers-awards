package req

type GetTokenAndUserInfoReq struct {
	UserId   string `form:"user_id" binding:"required"`         //用户id
	JXYToken string `form:"jxy_token" binding:"required,min=7"` //教学研token，600000：中台1.0 公版，6000001：中台2.0 诸暨版
}
