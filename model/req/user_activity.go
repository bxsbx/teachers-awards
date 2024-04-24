package req

import "time"

type GetActivityListToUserReq struct {
	UserId string `form:"user_id" binding:"required"`            // 用户id
	Page   int    `form:"page" binding:"required,gte=1"`         // 页数
	Limit  int    `form:"limit" binding:"required,gte=1,lte=50"` // 每页大小
}

type GetActivityInfoToUserReq struct {
	ActivityId int `form:"activity_id"  binding:"required"` //活动id
}

type GetUserActivityDetailReq struct {
	UserId     string `form:"user_id"  binding:"required"`     //用户id
	ActivityId int    `form:"activity_id"  binding:"required"` //活动id
}

type UserActivityIndicator struct {
	TwoIndicatorId       int        `json:"two_indicator_id" binding:"required"` //二级指标id
	AwardDate            *time.Time `json:"award_date" binding:"required"`       //获奖日期
	CertificateType      int        `json:"certificate_type" binding:"required"` //1.证书；需要填写证书有效期，2.证明；不需要有效期
	CertificateUrl       string     `json:"certificate_url" binding:"required"`  //证书url
	CertificateStartDate *time.Time `json:"certificate_start_date"`              //证书有效期——开始时间
	CertificateEndDate   *time.Time `json:"certificate_end_date"`                //证书有效期——结束时间
}

type UserActivity struct {
	ActivityId   int    `json:"activity_id" binding:"required"`        //活动id
	UserId       string `json:"user_id" binding:"required"`            //用户id
	UserName     string `json:"user_name" binding:"required"`          //用户姓名
	UserSex      int    `json:"user_sex" binding:"required,oneof=1 2"` //1：男，2：女
	Birthday     string `json:"birthday" binding:"required"`           //出生日期
	IdentityCard string `json:"identity_card" binding:"required"`      //身份证号
	Phone        string `json:"phone" binding:"required"`              //手机号
	SubjectCode  string `json:"subject_code" binding:"required"`       //科目code
	SchoolId     string `json:"school_id" binding:"required"`          //学校id
	SchoolName   string `json:"school_name" binding:"required"`        //学校名称
	DeclareType  int    `json:"declare_type" binding:"required"`       //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
}

type CreateUserActivityDeclareReq struct {
	UserActivity
	List []UserActivityIndicator `json:"list"`
}

type GetUserActivityDeclareDetailReq struct {
	ActivityId int    `form:"activity_id" binding:"required"` //活动id
	UserId     string `form:"user_id" binding:"required"`     //用户id
}

type GetUserActivityDeclareStatusListReq struct {
	ActivityId int    `form:"activity_id" binding:"required"` //活动id
	UserId     string `form:"user_id" binding:"required"`     //用户id
}

type GetUserActivityIndicatorReq struct {
	UserActivityIndicatorId int64 `form:"user_activity_indicator_id" binding:"required"` //用户活动申报下单个项目id
}

type UpdateUserActivityIndicatorReq struct {
	UserActivityIndicatorId int64 `json:"user_activity_indicator_id" binding:"required"` //用户活动申报下单个项目id
	UserActivityIndicator
}

type DeleteUserActivityIndicatorReq struct {
	UserActivityIndicatorId int64 `form:"user_activity_indicator_id" binding:"required"` //用户活动申报下单个项目id
}

type GetUserActivityDeclareResultReq struct {
	ActivityId int    `form:"activity_id" binding:"required"` //活动id
	UserId     string `form:"user_id" binding:"required"`     //用户id
}

type GetUserDeclareRecordListByYearReq struct {
	Year          int    `form:"year" binding:"required"`               //年份
	DeclareUserId string `form:"declare_user_id" binding:"required"`    //申报的用户id
	Role          int    `form:"role" binding:"required"`               //角色类型，1：学校，2：专家，3：教育局，4：老师
	Page          int    `form:"page" binding:"required,gte=1"`         // 页数
	Limit         int    `form:"limit" binding:"required,gte=1,lte=50"` // 每页大小
}

type GetUserHistoryDeclareResultListReq struct {
	DeclareUserId string `form:"declare_user_id" binding:"required"`    //申报的用户id
	Page          int    `form:"page" binding:"required,gte=1"`         // 页数
	Limit         int    `form:"limit" binding:"required,gte=1,lte=50"` // 每页大小
}

type GetUserDeclareResultAppReq struct {
	ActivityId int    `form:"activity_id" binding:"required"` //活动id
	UserId     string `form:"user_id" binding:"required"`     //用户id
}

type GetUserDeclaresToEdbReq struct {
	UserActivityId int64 `form:"user_activity_id"  binding:"required"` //用户活动id
}
