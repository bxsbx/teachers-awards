package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/global"
)

const (
	EXPERT_AUTH_INDICATOR_TABLE = "expert_auth_indicator"
)

type ExpertAuthIndicator struct {
	UserId         string `json:"user_id" gorm:"column:user_id;primary_key"`                   //用户id
	TwoIndicatorId int    `json:"two_indicator_id" gorm:"column:two_indicator_id;primary_key"` //二级指标id

}

func (ExpertAuthIndicator) TableName() string {
	return EXPERT_AUTH_INDICATOR_TABLE
}

type ExpertAuthIndicatorDao struct {
	BaseMysql
}

func NewExpertAuthIndicatorDao(appCtx context.Context) *ExpertAuthIndicatorDao {
	return &ExpertAuthIndicatorDao{
		BaseMysql{
			db:     global.GormDB.Model(&ExpertAuthIndicator{}),
			appCtx: appCtx,
		},
	}
}

func NewExpertAuthIndicatorDaoWithDB(db *gorm.DB, appCtx context.Context) *ExpertAuthIndicatorDao {
	return &ExpertAuthIndicatorDao{
		BaseMysql{
			db:     db.Model(&ExpertAuthIndicator{}),
			appCtx: appCtx,
		},
	}
}
