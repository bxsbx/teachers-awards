package resp

import "teachers-awards/dao"

type OneIndicator struct {
	dao.OneIndicator
	Count int                `json:"count"`          //二级指标个数
	List  []dao.TwoIndicator `json:"list,omitempty"` //二级指标列表
}

type GetOneIndicatorListResp struct {
	Total int64          `json:"total"`
	List  []OneIndicator `json:"list"`
}

type TwoIndicator struct {
	dao.TwoIndicator
	OneIndicatorName string `json:"one_indicator_name"` //一级指标名称
}

type GetTwoIndicatorListResp struct {
	Total int64          `json:"total"`
	List  []TwoIndicator `json:"list"`
}

type CreateOrUpdateTwoIndicatorResp struct {
	dao.TwoIndicator
}
