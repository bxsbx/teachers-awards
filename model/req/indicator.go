package req

type CreateOrUpdateOneIndicatorReq struct {
	OneIndicatorId   int    `form:"one_indicator_id"`                      //一级指标id
	OneIndicatorName string `form:"one_indicator_name" binding:"required"` //一级指标名称
	Content          string `form:"content" binding:"required"`            //评分标准说明
}

type GetOneIndicatorListReq struct {
	InputName    string `form:"input_name"`                   //输入的一级指标名称
	InputContent string `form:"input_content"`                //输入的评分标准说明
	WithTwo      bool   `form:"with_two"`                     //是否获取二级指标
	Page         int    `form:"page" binding:"gte=0"`         // 页数
	Limit        int    `form:"limit" binding:"gte=0,lte=50"` // 每页大小
}

type DeleteOneIndicatorByIdsReq struct {
	OneIndicatorIds string `form:"one_indicator_ids" binding:"required"` //一级指标ids
}

type CreateOrUpdateTwoIndicatorReq struct {
	TwoIndicatorId   int     `form:"two_indicator_id"`                      //二级指标id
	TwoIndicatorName string  `form:"two_indicator_name" binding:"required"` //二级指标名称
	Score            float64 `form:"score" binding:"required"`              //分值
	OneIndicatorId   int     `form:"one_indicator_id" binding:"required"`   //所属一级指标
}

type GetTwoIndicatorListReq struct {
	OneIndicatorId int     `form:"one_indicator_id"`             //所属一级指标
	InputName      string  `form:"input_name"`                   //输入的二级指标名称
	InputScore     float64 `form:"input_score"`                  //输入的分值
	Page           int     `form:"page" binding:"gte=0"`         // 页数
	Limit          int     `form:"limit" binding:"gte=0,lte=50"` // 每页大小
}

type DeleteTwoIndicatorByIdsReq struct {
	TwoIndicatorIds string `form:"two_indicator_ids" binding:"required"` //二级指标ids
}
