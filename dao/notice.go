package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/global"
	"time"
)

const (
	NOTICE_TABLE = "notice"
)

type Notice struct {
	NoticeId   int64      `json:"notice_id" gorm:"column:notice_id;primary_key"` //主键
	NoticeType int        `json:"notice_type" gorm:"column:notice_type"`         //通知类型，1：系统通知
	Content    string     `json:"content" gorm:"column:content"`                 //通知内容
	UserId     string     `json:"user_id" gorm:"column:user_id"`                 //接收者id
	IsRead     int        `json:"is_read" gorm:"column:is_read"`                 //0：未读，1：已读
	CreateTime *time.Time `json:"create_time" gorm:"column:create_time"`         //创建时间

}

func (Notice) TableName() string {
	return NOTICE_TABLE
}

type NoticeDao struct {
	BaseMysql
}

func NewNoticeDao(appCtx context.Context) *NoticeDao {
	return &NoticeDao{
		BaseMysql{
			db:     global.GormDB.Model(&Notice{}),
			appCtx: appCtx,
		},
	}
}

func NewNoticeDaoWithDB(db *gorm.DB, appCtx context.Context) *NoticeDao {
	return &NoticeDao{
		BaseMysql{
			db:     db.Model(&Notice{}),
			appCtx: appCtx,
		},
	}
}
