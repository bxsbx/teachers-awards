package resp

import (
	"teachers-awards/dao"
	"time"
)

type WaitReview struct {
	UserActivityIndicatorId int64      `json:"user_activity_indicator_id"` //用户活动申报单个项目的id
	UserId                  string     `json:"user_id"`                    //用户id
	UserName                string     `json:"user_name"`                  //用户姓名
	OneIndicatorId          int        `json:"one_indicator_id"`           //一级指标id
	OneIndicatorName        string     `json:"one_indicator_name"`         //一级指标名称
	TwoIndicatorId          int        `json:"two_indicator_id"`           //二级指标id
	TwoIndicatorName        string     `json:"two_indicator_name"`         //二级指标名称
	DeclareTime             *time.Time `json:"declare_time"`               //申报时间
}

type GetWaitReviewListResp struct {
	Total int64        `json:"total"`
	List  []WaitReview `json:"list"`
}

type ReviewItem struct {
	dao.ReviewItem
	dao.TwoIndicatorInfo
}

type GetReviewListResp struct {
	Total int64        `json:"total"`
	List  []ReviewItem `json:"list"`
}

type GetHistoryActivityListResp struct {
	Total int64                     `json:"total"`
	List  []dao.ActivityAndUserInfo `json:"list"`
}

type GetAwardsSetListResp struct {
	ActivityId     int                `json:"activity_id"`      //活动id
	ActivityName   string             `json:"activity_name"`    //活动名称
	AwardNum       int64              `json:"award_num"`        //获奖人数
	FirstPrizeNum  int64              `json:"first_prize_num"`  //一等奖
	SecondPrizeNum int64              `json:"second_prize_num"` //二等奖
	ThirdPrizeNum  int64              `json:"third_prize_num"`  //三等奖
	Total          int64              `json:"total"`
	List           []dao.UserActivity `json:"list"`
}

type GetEdbReviewListResp struct {
	Total int64               `json:"total"`
	List  []dao.EdbReviewItem `json:"list"`
}
