package req

import "time"

type GetUserInfoReq struct {
	UserId string `form:"user_id"  binding:"required"` //用户id
}

type SaveUserInfoReq struct {
	UserId       string `form:"user_id"  binding:"required"`           //用户id
	UserName     string `form:"user_name"  binding:"required"`         //用户名称
	UserSex      int    `form:"user_sex" binding:"required,oneof=1 2"` //1：男，2：女
	Birthday     string `form:"birthday"  binding:"required"`          //出生日期
	IdentityCard string `form:"identity_card"  binding:"required"`     //身份证号
	Phone        string `form:"phone"  binding:"required"`             //手机号
	Avatar       string `form:"avatar"  binding:"required"`            //头像
	SubjectCode  string `form:"subject_code"  binding:"required"`      //科目code
	SchoolId     string `form:"school_id"  binding:"required"`         //学校id
	SchoolName   string `form:"school_name"  binding:"required"`       //学校名称
	DeclareType  int    `form:"declare_type" binding:"required"`       //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
}

type GetUsersByNameReq struct {
	UserName string `form:"user_name"  binding:"required"`  //用户名称
	Page     int    `form:"page" binding:"required,gte=1"`  // 页数
	Limit    int    `form:"limit" binding:"required,gte=1"` // 每页大小
}

type SetRoleToUserReq struct {
	Roles   []int    `json:"roles"  binding:"required"`    //角色列表，1：学校，2：专家，3：教育局，4：老师，5：管理员
	UserIds []string `json:"user_ids"  binding:"required"` //用户id列表
}

type GetUserListByWhereReq struct {
	UserName     string `form:"user_name"`                              //用户名称
	UserSex      int    `form:"user_sex"`                               //1：男，2：女
	IdentityCard string `form:"identity_card"`                          //身份证号
	Phone        string `form:"phone"`                                  //手机号
	SchoolId     string `form:"school_id"`                              //学校id
	SchoolName   string `form:"school_name"`                            //学校名称
	Role         int    `form:"role" binding:"oneof=0 1 2 3 5"`         //0：全部，1：学校，2：专家，3：教育局，5：超级管理员
	Page         int    `form:"page" binding:"required,gte=1"`          // 页数
	Limit        int    `form:"limit" binding:"required,gte=1,lte=100"` // 每页大小
}

type GetExpertAuthListByWhereReq struct {
	UserName        string     `form:"user_name"`                              //用户名称
	UserSex         int        `form:"user_sex"`                               //1：男，2：女
	IdentityCard    string     `form:"identity_card"`                          //身份证号
	Phone           string     `form:"phone"`                                  //手机号
	TwoIndicatorIds string     `form:"two_indicator_ids"`                      //二级指标ids，, 隔开
	ExportAuth      int        `form:"export_auth"`                            //专家是否授权 0：全部 1：未授权 2：已授权
	AuthDay         *time.Time `form:"auth_day" time_format:"2006-01-02"`      //授权日期，格式2006-01-02
	Page            int        `form:"page" binding:"required,gte=1"`          // 页数
	Limit           int        `form:"limit" binding:"required,gte=1,lte=100"` // 每页大小
}

type SetExpertAuthReq struct {
	UserId string `json:"user_id"  binding:"required"` //用户id
	TwoIds []int  `json:"two_ids"  binding:"required"` //二级指标ids
}

type CancelExpertAuthReq struct {
	UserId string `form:"user_id"  binding:"required"` //用户id
}
