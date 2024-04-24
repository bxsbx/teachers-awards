package req

import "time"

type GetWaitReviewListReq struct {
	ReviewProcess int `form:"review_process" binding:"required"`     // 当前审核进程，1：学校，2：专家，3：教育局，4：结束
	Page          int `form:"page" binding:"required,gte=1"`         // 页数
	Limit         int `form:"limit" binding:"required,gte=1,lte=50"` // 每页大小
}

type GetReviewListReq struct {
	JudgesId         string `form:"judges_id"`                                     //评委id
	UserName         string `form:"user_name"`                                     //用户姓名
	UserSex          int    `form:"user_sex"`                                      //1：男，2：女
	SubjectCode      string `form:"subject_code"`                                  //科目code
	SchoolId         string `form:"school_id"`                                     //学校id
	SchoolName       string `form:"school_name"`                                   //学校名称
	Year             int    `form:"year"`                                          //年份
	OneIndicatorName string `form:"one_indicator_name"`                            //一级指标名称
	TwoIndicatorName string `form:"two_indicator_name"`                            //二级指标名称
	ReviewStatus     int    `form:"review_status" binding:"oneof=0 1 2 3 4"`       //审核状态 0:全部 1:待审核，2：未通过，3：已通过，4：已审核"
	ReviewProcess    int    `form:"review_process" binding:"required,oneof=1 2 3"` // 当前审核进程，1：学校，2：专家，3：教育局
	OnlyCount        bool   `form:"only_count"`                                    //是否仅获取数量
	Page             int    `form:"page" binding:"required,gte=1"`                 // 页数
	Limit            int    `form:"limit" binding:"required,gte=1,lte=50"`         // 每页大小
}

type CommitReviewReq struct {
	UserActivityIndicatorId int64  `form:"user_activity_indicator_id" binding:"required"` //用户活动申报单项目id
	IsPass                  int    `form:"is_pass" binding:"oneof=0 1"`                   //0：未通过，1：通过
	Opinion                 string `form:"opinion" binding:"required"`                    //审核意见
}

type GetHistoryActivityListReq struct {
	UserName     string `form:"user_name"`                             //用户姓名
	UserSex      int    `form:"user_sex"`                              //1：男，2：女
	SubjectCode  string `form:"subject_code"`                          //科目code
	SchoolId     string `form:"school_id"`                             //学校id
	SchoolName   string `form:"school_name"`                           //学校名称
	Year         int    `form:"year"`                                  //年份
	DeclareType  int    `form:"declare_type"`                          //申报类型
	IdentityCard string `form:"identity_card"`                         //身份证号
	RankPrize    *int   `form:"rank_prize"`                            //0：无，1：一等奖，2：二等奖，3：三等奖
	Rank         int    `form:"rank"`                                  //排名
	FinalScore   int    `form:"final_score"`                           //最终得分（各项通过的审核）
	Page         int    `form:"page" binding:"required,gte=1"`         // 页数
	Limit        int    `form:"limit" binding:"required,gte=1,lte=50"` // 每页大小
}

type GetAwardsSetListReq struct {
	ActivityId   int    `form:"activity_id"  binding:"required"`       //活动id
	UserName     string `form:"user_name"`                             //用户姓名
	UserSex      int    `form:"user_sex"`                              //1：男，2：女
	SubjectCode  string `form:"subject_code"`                          //科目code
	SchoolId     string `form:"school_id"`                             //学校id
	SchoolName   string `form:"school_name"`                           //学校名称
	DeclareType  int    `form:"declare_type"`                          //申报类型
	IdentityCard string `form:"identity_card"`                         //身份证号
	RankPrize    *int   `form:"rank_prize"`                            //名次，0：无，1：一等奖，2：二等奖，3：三等奖
	FinalScore   int    `form:"final_score"`                           //最终得分（各项通过的审核）
	Page         int    `form:"page" binding:"required,gte=1"`         // 页数
	Limit        int    `form:"limit" binding:"required,gte=1,lte=50"` // 每页大小
}

type SetAwardsReq struct {
	ActivityId int    `form:"activity_id"  binding:"required"` //活动id
	UserId     string `form:"user_id" binding:"required"`      //用户id
	RankPrize  *int   `form:"rank_prize"  binding:"required"`  //名次
	Prize      int    `form:"prize" binding:"required"`        //奖金
}

type CommitActivityResultReq struct {
	ActivityId  int `form:"activity_id"  binding:"required"` //活动id
	DeclareType int `form:"declare_type" binding:"required"` //申报类型
}

type UpdateTwoIndicatorIdReq struct {
	UserActivityIndicatorId int64      `json:"user_activity_indicator_id"  binding:"required"` //用户活动申报单项目id
	TwoIndicatorId          int        `json:"two_indicator_id" binding:"required"`            //二级指标id
	AwardDate               *time.Time `json:"award_date"`                                     //获奖日期
	CertificateType         int        `json:"certificate_type"`                               //1.证书；需要填写证书有效期，2.证明；不需要有效期
	CertificateUrl          string     `json:"certificate_url"`                                //证书url
	CertificateStartDate    *time.Time `json:"certificate_start_date"`                         //证书有效期——开始时间
	CertificateEndDate      *time.Time `json:"certificate_end_date"`                           //证书有效期——结束时间
}

type EdbDeclareToUserReq struct {
	UserActivityId       int64      `json:"user_activity_id"  binding:"required"` //用户活动id
	TwoIndicatorId       int        `json:"two_indicator_id" binding:"required"`  //二级指标id
	AwardDate            *time.Time `json:"award_date"`                           //获奖日期
	CertificateType      int        `json:"certificate_type"`                     //1.证书；需要填写证书有效期，2.证明；不需要有效期
	CertificateUrl       string     `json:"certificate_url"`                      //证书url
	CertificateStartDate *time.Time `json:"certificate_start_date"`               //证书有效期——开始时间
	CertificateEndDate   *time.Time `json:"certificate_end_date"`                 //证书有效期——结束时间
}

type GetEdbReviewListReq struct {
	UserName    string `form:"user_name"`                             //用户姓名
	UserSex     int    `form:"user_sex"`                              //1：男，2：女
	SubjectCode string `form:"subject_code"`                          //科目code
	SchoolId    string `form:"school_id"`                             //学校id
	SchoolName  string `form:"school_name"`                           //学校名称
	Status      int    `form:"status" binding:"oneof=0 1 2"`          //审核状态 0:全部 1:待审核，2：已审核
	OnlyCount   bool   `form:"only_count"`                            //是否仅获取数量
	Page        int    `form:"page" binding:"required,gte=1"`         // 页数
	Limit       int    `form:"limit" binding:"required,gte=1,lte=50"` // 每页大小
}

type BatchPassReq struct {
	UserActivityIndicatorIds []int64 `json:"user_activity_indicator_ids"  binding:"required"` //用户活动申报单项目id
}
