package resp

import (
	"teachers-awards/dao"
)

type ActivityOneIndicator struct {
	dao.ActivityOneIndicator
	ActivityTwoIndicators []dao.ActivityTwoIndicator `json:"activity_two_indicators"` //二级指标列表
}

type GetActivityDetailResp struct {
	dao.Activity
	OneIndicators []ActivityOneIndicator `json:"activity_one_indicators"` //一级指标列表
}

type Activity struct {
	dao.Activity
	Status    int `json:"status"  gorm:"column:status"` //1:待开始，2：进行中，3：已结束
	AttendNum int `json:"attend_num"`                   //参加人数
}

type GetActivityListResp struct {
	Total int64      `json:"total"`
	List  []Activity `json:"list"`
}

type GetActivityTwoIndicatorListResp struct {
	List []dao.ActivityTwoIndicator `json:"list"`
}

type ActivityIndicatorInfo struct {
	OneIndicatorId   int     `json:"one_indicator_id"`   //一级指标id
	OneIndicatorName string  `json:"one_indicator_name"` //一级指标名称
	Content          string  `json:"content"`            //评分标准说明
	TwoIndicatorId   int     `json:"two_indicator_id"`   //二级指标id
	TwoIndicatorName string  `json:"two_indicator_name"` //二级指标名称
	Score            float64 `json:"score"`              //分值
}

type GetActivityYearListResp struct {
	Years []int `json:"years"`
}

type GetLatestActivityResp struct {
	HasDeclare bool `json:"has_declare"` //是否有申报过
	dao.Activity
}
