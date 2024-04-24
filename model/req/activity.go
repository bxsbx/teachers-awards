package req

import "time"

type CreateOrUpdateActivityReq struct {
	ActivityId      int        `json:"activity_id"`                           //活动id
	ActivityName    string     `json:"activity_name" binding:"required"`      //活动名称
	Year            int        `json:"year" binding:"required"`               //年份
	StartTime       *time.Time `json:"start_time" binding:"required"`         //开始时间
	EndTime         *time.Time `json:"end_time" binding:"required"`           //结束时间
	Url             string     `json:"url" binding:"required"`                //活动文件url
	Description     string     `json:"description" binding:"required"`        //申报须知
	TwoIndicatorIds string     `json:"two_indicator_ids"  binding:"required"` //二级指标ids，","隔开
}

type GetActivityDetailReq struct {
	ActivityId int `form:"activity_id"  binding:"required"` //活动id
}

type DeleteActivityOneIndicatorReq struct {
	ActivityId     int `form:"activity_id"  binding:"required"`      //活动id
	OneIndicatorId int `form:"one_indicator_id"  binding:"required"` //一级指标id
}

type DeleteActivityTwoIndicatorReq struct {
	ActivityId     int `form:"activity_id"  binding:"required"`      //活动id
	TwoIndicatorId int `form:"two_indicator_id"  binding:"required"` //二级指标id
}

type DeleteActivityReq struct {
	ActivityId int `form:"activity_id"  binding:"required"` //活动id
}

type GetActivityListReq struct {
	ActivityName string `form:"activity_name"`                         //活动名称
	Year         int    `form:"year"`                                  //活动名称
	Page         int    `form:"page" binding:"required,gte=1"`         // 页数
	Limit        int    `form:"limit" binding:"required,gte=1,lte=50"` // 每页大小
}

type GetActivityTwoIndicatorListReq struct {
	ActivityId     int `form:"activity_id"  binding:"required"`      //活动id
	OneIndicatorId int `form:"one_indicator_id"  binding:"required"` //一级指标id
}

type EndActivityReq struct {
	ActivityId int `form:"activity_id"  binding:"required"` //活动id
}
