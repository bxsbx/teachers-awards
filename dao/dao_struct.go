package dao

import (
	"time"
)

type WaitReview struct {
	ActivityId              int        `json:"activity_id" gorm:"column:activity_id"`                                           //活动id
	UserId                  string     `json:"user_id" gorm:"column:user_id"`                                                   //用户id
	UserName                string     `json:"user_name" gorm:"column:user_name"`                                               //用户姓名
	UserActivityIndicatorId int64      `json:"user_activity_indicator_id" gorm:"column:user_activity_indicator_id;primary_key"` //主键id
	TwoIndicatorId          int        `json:"two_indicator_id" gorm:"column:two_indicator_id"`                                 //二级指标id
	CreateTime              *time.Time `json:"create_time" gorm:"column:create_time"`                                           //创建时间
}

type YearDeclareAwardNum struct {
	Year       int   `json:"year" gorm:"column:year"`
	DeclareNum int64 `json:"declare_num"  gorm:"column:declare_num"` //申报数量
	AwardNum   int64 `json:"award_num"  gorm:"column:award_num"`     //获奖人数
}

type SchoolAwardNum struct {
	SchoolId   string `json:"school_id" gorm:"column:school_id"`     //学校id，-1：其他学校
	SchoolName string `json:"school_name" gorm:"column:school_name"` //学校名称
	AwardNum   int64  `json:"award_num" gorm:"column:award_num"`     //获奖人数
}

type TeacherTypeAwardNum struct {
	DeclareType int   `json:"declare_type" gorm:"column:declare_type"` //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
	AwardNum    int64 `json:"award_num" gorm:"column:award_num"`       //获奖人数
}

type DeclareAwardRankNum struct {
	DeclareNum     int64 `json:"declare_num"  gorm:"column:declare_num"`          //申报数量
	AwardNum       int64 `json:"award_num" gorm:"column:award_num"`               //获奖人数
	FirstPrizeNum  int64 `json:"first_prize_num" gorm:"column:first_prize_num"`   //一等奖
	SecondPrizeNum int64 `json:"second_prize_num" gorm:"column:second_prize_num"` //二等奖
	ThirdPrizeNum  int64 `json:"third_prize_num" gorm:"column:third_prize_num"`   //三等奖
}

type YearDeclareAwardRankNum struct {
	Year int `json:"year" gorm:"column:year"`
	DeclareAwardRankNum
}

type SchoolDeclareAwardRankNum struct {
	SchoolId   string `json:"school_id" gorm:"column:school_id"`     //学校id
	SchoolName string `json:"school_name" gorm:"column:school_name"` //学校名称
	DeclareAwardRankNum
}

type TypeDeclareAwardRankNum struct {
	DeclareType int `json:"declare_type" gorm:"column:declare_type"` //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
	DeclareAwardRankNum
}

type TwoIndicatorInfo struct {
	ActivityId       int     `json:"activity_id,omitempty" gorm:"column:activity_id"`             //活动id
	TwoIndicatorId   int     `json:"two_indicator_id" gorm:"column:two_indicator_id;primary_key"` //二级指标id
	TwoIndicatorName string  `json:"two_indicator_name" gorm:"column:two_indicator_name"`         //二级指标名称
	Score            float64 `json:"score" gorm:"column:score"`                                   //分值
	OneIndicatorId   int     `json:"one_indicator_id" gorm:"column:one_indicator_id"`             //所属一级指标
	OneIndicatorName string  `json:"one_indicator_name" gorm:"column:one_indicator_name"`         //一级指标名称
}

type ReviewItem struct {
	ActivityId              int        `json:"activity_id" gorm:"column:activity_id"`                               //活动id
	Year                    int        `json:"year" gorm:"column:year"`                                             //年份
	UserId                  string     `json:"user_id" gorm:"column:user_id"`                                       //用户id
	UserName                string     `json:"user_name" gorm:"column:user_name"`                                   //用户姓名
	UserSex                 int        `json:"user_sex" gorm:"column:user_sex"`                                     //1：男，2：女
	SubjectCode             string     `json:"subject_code" gorm:"column:subject_code"`                             //科目code
	SchoolId                string     `json:"school_id" gorm:"column:school_id"`                                   //学校id
	SchoolName              string     `json:"school_name" gorm:"column:school_name"`                               //学校名称
	DeclareType             int        `json:"declare_type" gorm:"column:declare_type"`                             //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
	UserActivityIndicatorId int64      `json:"user_activity_indicator_id" gorm:"column:user_activity_indicator_id"` //主键id
	TwoIndicatorId          int        `json:"two_indicator_id" gorm:"column:two_indicator_id"`                     //二级指标id
	Status                  int        `json:"status" gorm:"column:status"`                                         //审核状态:1.已提交 ,2.已通过,3.学校未通过,4.专家未通过,5.教育局未通过,6.活动已结束，未审批,
	ReviewStatus            int        `json:"review_status" gorm:"column:review_status"`                           //审核状态 1:待审核，2：未通过，3：通过
	CreateTime              *time.Time `json:"create_time" gorm:"column:create_time"`                               //创建时间
}

type ActivityAndUserInfo struct {
	ActivityId   int    `json:"activity_id" gorm:"column:activity_id"`     //活动id
	ActivityName string `json:"activity_name" gorm:"column:activity_name"` //活动名称
	Year         int    `json:"year" gorm:"column:year"`                   //年份
	UserActivity
}

type UserDeclareRecord struct {
	ActivityId              int        `json:"activity_id" gorm:"column:activity_id"`                               //活动id
	UserActivityId          int64      `json:"user_activity_id" gorm:"column:user_activity_id"`                     //用户申报活动的id
	UserActivityIndicatorId int64      `json:"user_activity_indicator_id" gorm:"column:user_activity_indicator_id"` //用户活动申报二级指标id
	TwoIndicatorId          int        `json:"two_indicator_id" gorm:"column:two_indicator_id"`                     //二级指标id
	Status                  int        `json:"status" gorm:"column:status"`                                         //审核状态:1.已提交 ,2.已通过,3.学校未通过,4.专家未通过,5.教育局未通过,6.活动已结束，未审批,
	ReviewProcess           int        `json:"review_process" gorm:"column:review_process"`                         //当前审核进程状态，1：学校，2：专家，3：教育局，4：结束
	CreateTime              *time.Time `json:"create_time" gorm:"column:create_time"`                               //创建时间
}

type UserHistoryDeclareResult struct {
	ActivityId   int     `json:"activity_id" gorm:"column:activity_id"`     //活动id
	ActivityName string  `json:"activity_name" gorm:"column:activity_name"` //活动名称
	FinalScore   float64 `json:"final_score" gorm:"column:final_score"`     //最终得分（各项通过的审核）
	Rank         int     `json:"rank" gorm:"column:rank"`                   //排名
	RankPrize    int     `json:"rank_prize" gorm:"column:rank_prize"`       //0：无，1：一等奖，2：二等奖，3：三等奖
	Prize        int     `json:"prize" gorm:"column:prize"`                 //奖金
}

type ResultGroupByDeclareType struct {
	DeclareType int   `json:"declare_type" gorm:"column:declare_type"` //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
	SchoolNum   int64 `json:"school_num"  gorm:"column:school_num"`    //申报学校数量
	RankNum     int64 `json:"status" gorm:"column:rank_num"`           //排名总数，判断是否提交状态，0：未提交，1：已提交
	DeclareAwardRankNum
}

type EdbReviewItem struct {
	UserActivityId int64      `json:"user_activity_id" gorm:"column:user_activity_id"` //用户活动id
	UserId         string     `json:"user_id" gorm:"column:user_id"`                   //用户id
	UserName       string     `json:"user_name" gorm:"column:user_name"`               //用户姓名
	UserSex        int        `json:"user_sex" gorm:"column:user_sex"`                 //1：男，2：女
	SubjectCode    string     `json:"subject_code" gorm:"column:subject_code"`         //科目code
	SchoolId       string     `json:"school_id" gorm:"column:school_id"`               //学校id
	SchoolName     string     `json:"school_name" gorm:"column:school_name"`           //学校名称
	DeclareType    int        `json:"declare_type" gorm:"column:declare_type"`         //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
	IdentityCard   string     `json:"identity_card" gorm:"column:identity_card"`       //身份证号
	FinalScore     float64    `json:"final_score" gorm:"column:final_score"`           //最终得分（各项通过的审核）
	Status         int        `json:"status" gorm:"column:status"`                     //审核状态:1:待审核，2：已审核
	CreateTime     *time.Time `json:"create_time" gorm:"column:create_time"`           //创建时间
}
