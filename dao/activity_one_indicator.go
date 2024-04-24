package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"time"
)

const (
	ACTIVITY_ONE_INDICATOR_TABLE = "activity_one_indicator"
)

type ActivityOneIndicator struct {
	ActivityId       int        `json:"activity_id" gorm:"column:activity_id;primary_key"`           //活动id
	OneIndicatorId   int        `json:"one_indicator_id" gorm:"column:one_indicator_id;primary_key"` //一级指标id
	OneIndicatorName string     `json:"one_indicator_name" gorm:"column:one_indicator_name"`         //一级指标名称
	Content          string     `json:"content" gorm:"column:content"`                               //评分标准说明
	CreateTime       *time.Time `json:"create_time" gorm:"column:create_time"`                       //创建时间
}

func (ActivityOneIndicator) TableName() string {
	return ACTIVITY_ONE_INDICATOR_TABLE
}

type ActivityOneIndicatorDao struct {
	BaseMysql
}

func NewActivityOneIndicatorDao(appCtx context.Context) *ActivityOneIndicatorDao {
	return &ActivityOneIndicatorDao{
		BaseMysql{
			db:     global.GormDB.Model(&ActivityOneIndicator{}),
			appCtx: appCtx,
		},
	}
}

func NewActivityOneIndicatorDaoWithDB(db *gorm.DB, appCtx context.Context) *ActivityOneIndicatorDao {
	return &ActivityOneIndicatorDao{
		BaseMysql{
			db:     db.Model(&ActivityOneIndicator{}),
			appCtx: appCtx,
		},
	}
}

func (t *ActivityOneIndicatorDao) GetActivityOneIndicatorMap(activityId int) (activityOneIndicatorMap map[int]ActivityOneIndicator, err error) {
	activityOneIndicatorMap = make(map[int]ActivityOneIndicator)
	var list []ActivityOneIndicator
	err = t.db.WithContext(t.appCtx).Where(ActivityOneIndicator{ActivityId: activityId}).Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	for _, v := range list {
		activityOneIndicatorMap[v.OneIndicatorId] = v
	}
	return
}
