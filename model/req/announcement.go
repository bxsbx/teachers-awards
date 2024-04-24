package req

import "time"

type GetAnnouncementByIdReq struct {
	UserId         string `form:"user_id" binding:"required"`         // 用户id
	AnnouncementId int    `form:"announcement_id" binding:"required"` // 公告id
}

type SaveAnnouncementReq struct {
	AnnouncementId int    `form:"announcement_id"`              // 公告id
	UserId         string `form:"user_id" binding:"required"`   // 用户id
	UserName       string `form:"user_name" binding:"required"` // 用户名称
	Title          string `form:"title" binding:"required"`     // 标题
	Content        string `form:"content" binding:"required"`   // 内容
	Annex          string `form:"annex"`                        // 附件
}

type DeleteAnnouncementByIdReq struct {
	AnnouncementId int `form:"announcement_id" binding:"required"` // 公告id
}

type GetAnnouncementListReq struct {
	UserId         string     `form:"user_id" binding:"required"`                    // 用户id
	InputContent   string     `form:"input_content"`                                 // 输入的内容
	PublishTime    *time.Time `form:"publish_time" time_format:"2006-01-02"`         // 发布时间，格式2006-01-02
	CreatTimeOrder string     `form:"create_time_order" binding:"oneof=desc asc ''"` // 按时间升序或者降序
	Page           int        `form:"page" binding:"required,gte=1"`                 // 页数
	Limit          int        `form:"limit" binding:"required,gte=1,lte=50"`         // 每页大小
}

type AllAnnouncementReadReq struct {
	UserId string `form:"user_id" binding:"required"` // 用户id
}
