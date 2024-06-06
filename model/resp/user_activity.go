package resp

import (
	"teachers-awards/dao"
	"time"
)

type ActivityListToUser struct {
	dao.Activity
	Status    int  `json:"status"`     //1:待开始，2：进行中，3：已结束
	IsDeclare bool `json:"is_declare"` //是否已申请
}

type GetActivityListToUserResp struct {
	Total int64                `json:"total"`
	List  []ActivityListToUser `json:"list"`
}

type GetActivityInfoToUserResp struct {
	dao.Activity
}

type UserActivity struct {
	ActivityId   int    `json:"activity_id"`   //活动id
	UserId       string `json:"user_id"`       //用户id
	UserName     string `json:"user_name"`     //用户姓名
	UserSex      int    `json:"user_sex"`      //1：男，2：女
	Birthday     string `json:"birthday"`      //出生日期
	IdentityCard string `json:"identity_card"` //身份证号
	Phone        string `json:"phone"`         //手机号
	SubjectCode  string `json:"subject_code"`  //科目code
	SchoolId     string `json:"school_id"`     //学校id
	SchoolName   string `json:"school_name"`   //学校名称
	DeclareType  int    `json:"declare_type"`  //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
}

type UserActivityIndicator struct {
	dao.UserActivityIndicator
	OneIndicatorId int `json:"one_indicator_id"` //一级指标id
}

type GetUserActivityDeclareDetailResp struct {
	UserActivity
	List []UserActivityIndicator `json:"list"`
}

type DeclareStatus struct {
	UserActivityIndicatorId int64      `json:"user_activity_indicator_id"` //申报id
	TwoIndicatorId          int        `json:"two_indicator_id"`           //二级指标id
	TwoIndicatorName        string     `json:"two_indicator_name"`         //二级指标名称
	Score                   float64    `json:"score"`                      //分值
	Status                  int        `json:"status"`                     //审核状态:1.已提交 ,2.已通过,3.学校未通过,4.专家未通过,5.教育局未通过,6.活动已结束，未审批
	DeclareTime             *time.Time `json:"declare_time"`               //申报时间
}

type OneIndicatorDeclareStatusList struct {
	OneIndicatorId    int             `json:"one_indicator_id"`   //一级指标id
	OneIndicatorName  string          `json:"one_indicator_name"` //一级指标名称
	Content           string          `json:"content"`            //评分标准说明
	DeclareStatusList []DeclareStatus `json:"declare_status_list"`
}

type GetUserActivityDeclareStatusListResp struct {
	ActivityStatus int                             `json:"activity_status"` // 活动状态，1：待开始，2：进行中，3已结束
	DeclareNum     int                             `json:"declare_num"`     //总申报项数
	PassNum        int                             `json:"pass_num"`        //已通过申报项数
	NoPassNum      int                             `json:"no_pass_num"`     //未通过申报项数
	WaitReview     int                             `json:"wait_review"`     //待审核申报项数
	List           []OneIndicatorDeclareStatusList `json:"list"`
}

type GetUserActivityIndicatorResp struct {
	UserId   string `json:"user_id"`   //用户id
	UserName string `json:"user_name"` //用户姓名
	dao.UserActivityIndicator
	ActivityIndicatorInfo
	List []dao.JudgesVerify
}

type GetUserActivityDeclareResultResp struct {
	DeclareNum int     `json:"declare_num"`  //总申报项数
	FinalScore float64 `json:"final_score" ` //最终得分（各项通过的审核）
	RankPrize  int     `json:"rank_prize"`   //0：无，1：一等奖，2：二等奖，3：三等奖
	Prize      int     `json:"prize"`        //奖金
	List       []ActivityIndicatorInfo
}
type JudgesVerifyPass struct {
	JudgesId   string `json:"judges_id"`   //评委id
	JudgesName string `json:"judges_name"` //评委姓名
	JudgesType int    `json:"judges_type"` //评委类型，1：学校，2：专家，3：教育局
	IsPass     int    `json:"is_pass"`     //0：未通过，1：通过
}

type UserDeclareRecord struct {
	dao.TwoIndicatorInfo
	UserActivityId          int64              `json:"user_activity_id"`             //用户申报活动的id
	UserActivityIndicatorId int64              `json:"user_activity_indicator_id"`   //用户活动申报二级指标id
	Status                  int                `json:"status"`                       //审核状态:1.待审核 ,2.已通过,3.学校未通过,5.未通过（教育局）,
	CreateTime              *time.Time         `json:"create_time"`                  //创建时间
	JudgesVerifyList        []JudgesVerifyPass `json:"judges_verify_list,omitempty"` //审核列表
}

type GetUserDeclareRecordListByYearResp struct {
	Total int64               `json:"total"`
	List  []UserDeclareRecord `json:"list"`
}

type GetUserHistoryDeclareResultListResp struct {
	Total int64                          `json:"total"`
	List  []dao.UserHistoryDeclareResult `json:"list"`
}

type GetUserDeclareResultAppResp struct {
	dao.UserHistoryDeclareResult
	TeacherScore      float64 `json:"teacher_score"`        //老师自评
	SchoolScore       float64 `json:"school_score"`         //学校评分
	ExpertPassScore   float64 `json:"expert_pass_score"`    //专家评分（通过）
	ExpertNoPassScore float64 `json:"expert_no_pass_score"` //专家评分（未通过）
	EdbScore          float64 `json:"edb_score"`            //教育局评分
}

type UserDeclaresToEdb struct {
	dao.TwoIndicatorInfo
	ReviewProcess           int                `json:"review_process"`               //当前审核进程状态，1：学校，2：专家，3：教育局，4：结束
	UserActivityIndicatorId int64              `json:"user_activity_indicator_id"`   //用户活动申报二级指标id
	CertificateUrl          string             `json:"certificate_url"`              //证书url
	JudgesVerifyList        []JudgesVerifyPass `json:"judges_verify_list,omitempty"` //审核列表
}

type GetUserDeclaresToEdbResp struct {
	FinalScore       float64                    `json:"final_score"` //最终得分（各项通过的审核）
	List             []UserDeclaresToEdb        `json:"list"`
	OneIndicatorList []dao.ActivityOneIndicator `json:"one_indicator_list"` //一级指标列表
}
