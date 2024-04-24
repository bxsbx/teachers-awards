package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/global"
	"time"
)

const (
	ANNOUNCEMENT_TABLE = "announcement"
)

type Announcement struct {
	AnnouncementId int        `json:"announcement_id" gorm:"column:announcement_id;primary_key"` //主键id
	Title          string     `json:"title" gorm:"column:title"`                                 //标题
	Content        string     `json:"content" gorm:"column:content"`                             //内容
	Annex          string     `json:"annex" gorm:"column:annex"`                                 //附件链接
	UserId         string     `json:"user_id" gorm:"column:user_id"`                             //用户id
	UserName       string     `json:"user_name" gorm:"column:user_name"`                         //用户名称
	CreateTime     *time.Time `json:"create_time" gorm:"column:create_time"`                     //创建时间

}

func (Announcement) TableName() string {
	return ANNOUNCEMENT_TABLE
}

type AnnouncementDao struct {
	BaseMysql
}

func NewAnnouncementDao(appCtx context.Context) *AnnouncementDao {
	return &AnnouncementDao{
		BaseMysql{
			db:     global.GormDB.Model(&Announcement{}),
			appCtx: appCtx,
		},
	}
}

func NewAnnouncementDaoWithDB(db *gorm.DB, appCtx context.Context) *AnnouncementDao {
	return &AnnouncementDao{
		BaseMysql{
			db:     db.Model(&Announcement{}),
			appCtx: appCtx,
		},
	}
}
