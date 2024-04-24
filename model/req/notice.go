package req

type GetUserNoticeListReq struct {
	UserId         string `form:"user_id" binding:"required"`                    // 用户id
	CreatTimeOrder string `form:"create_time_order" binding:"oneof=desc asc ''"` // 按时间升序或者降序
	Page           int    `form:"page" binding:"required,gte=1"`                 // 页数
	Limit          int    `form:"limit" binding:"required,gte=1,lte=50"`         // 每页大小
}

type UpdateUserAllNoticeToReadReq struct {
	UserId string `form:"user_id" binding:"required"` // 用户id
}

type DeleteUserAllNoticeReq struct {
	UserId string `form:"user_id" binding:"required"` // 用户id
}

type DeleteUserNoticeByIdsReq struct {
	NoticeIds string `form:"notice_ids" binding:"required"` // 通知ids，英文","隔开
	UserId    string `form:"user_id" binding:"required"`    // 用户id
}
