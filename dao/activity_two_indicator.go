package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"time"
)

const (
	ACTIVITY_TWO_INDICATOR_TABLE = "activity_two_indicator"
)

type ActivityTwoIndicator struct {
	ActivityId       int        `json:"activity_id" gorm:"column:activity_id;primary_key"`           //活动id
	TwoIndicatorId   int        `json:"two_indicator_id" gorm:"column:two_indicator_id;primary_key"` //二级指标id
	TwoIndicatorName string     `json:"two_indicator_name" gorm:"column:two_indicator_name"`         //二级指标名称
	Score            float64    `json:"score" gorm:"column:score"`                                   //分值
	OneIndicatorId   int        `json:"one_indicator_id" gorm:"column:one_indicator_id"`             //所属一级指标
	CreateTime       *time.Time `json:"create_time" gorm:"column:create_time"`                       //创建时间
}

func (ActivityTwoIndicator) TableName() string {
	return ACTIVITY_TWO_INDICATOR_TABLE
}

type ActivityTwoIndicatorDao struct {
	BaseMysql
}

func NewActivityTwoIndicatorDao(appCtx context.Context) *ActivityTwoIndicatorDao {
	return &ActivityTwoIndicatorDao{
		BaseMysql{
			db:     global.GormDB.Model(&ActivityTwoIndicator{}),
			appCtx: appCtx,
		},
	}
}

func NewActivityTwoIndicatorDaoWithDB(db *gorm.DB, appCtx context.Context) *ActivityTwoIndicatorDao {
	return &ActivityTwoIndicatorDao{
		BaseMysql{
			db:     db.Model(&ActivityTwoIndicator{}),
			appCtx: appCtx,
		},
	}
}

func (t *ActivityTwoIndicatorDao) GetTwoToOneMap(activityId int) (twoToOneMap map[int]int, err error) {
	twoToOneMap = make(map[int]int)
	rows, err := t.db.WithContext(t.appCtx).Select("one_indicator_id,two_indicator_id").
		Where(ActivityTwoIndicator{ActivityId: activityId}).Rows()
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var oneId, twoId int
		rows.Scan(&oneId, &twoId)
		twoToOneMap[twoId] = oneId
	}
	return
}

func (t *ActivityTwoIndicatorDao) GetActivityTwoIndicatorMap(activityId int, twoIds []int) (activityTwoIndicatorMap map[int]ActivityTwoIndicator, err error) {
	activityTwoIndicatorMap = make(map[int]ActivityTwoIndicator)
	var list []ActivityTwoIndicator
	err = t.db.WithContext(t.appCtx).Where("activity_id = ?", activityId).
		Where("two_indicator_id in (?)", twoIds).Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	for _, v := range list {
		activityTwoIndicatorMap[v.TwoIndicatorId] = v
	}
	return
}

// activityId-twoId
func (t *ActivityTwoIndicatorDao) GetTwoIndicatorInfoMap(activityIds []int, twoIds []int) (activityTwoIndicatorsMap map[int]map[int]TwoIndicatorInfo, err error) {
	activityTwoIndicatorsMap = make(map[int]map[int]TwoIndicatorInfo)
	var list []TwoIndicatorInfo
	err = t.db.WithContext(t.appCtx).Table("activity_two_indicator ati").
		Joins("join activity_one_indicator aoi on aoi.activity_id = ati.activity_id and aoi.one_indicator_id = ati.one_indicator_id").
		Select("ati.*,aoi.one_indicator_name").
		Where("ati.activity_id in (?)", activityIds).
		Where("ati.two_indicator_id in (?)", twoIds).Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	for _, v := range list {
		tempMap, ok := activityTwoIndicatorsMap[v.ActivityId]
		if !ok {
			tempMap = make(map[int]TwoIndicatorInfo)
		}
		tempMap[v.TwoIndicatorId] = v
		activityTwoIndicatorsMap[v.ActivityId] = tempMap
	}
	return
}
