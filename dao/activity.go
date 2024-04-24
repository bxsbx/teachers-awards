package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"time"
)

const (
	ACTIVITY_TABLE = "activity"
)

type Activity struct {
	ActivityId   int            `json:"activity_id" gorm:"column:activity_id;primary_key"` //主键id
	ActivityName string         `json:"activity_name" gorm:"column:activity_name"`         //活动名称
	Year         int            `json:"year" gorm:"column:year"`                           //年份
	Description  string         `json:"description,omitempty" gorm:"column:description"`   //申报须知
	Url          string         `json:"url" gorm:"column:url"`                             //活动文件url
	ReviewNum    int            `json:"review_num" gorm:"column:review_num"`               //当前已审核人数
	StartTime    *time.Time     `json:"start_time" gorm:"column:start_time"`               //开始时间
	EndTime      *time.Time     `json:"end_time" gorm:"column:end_time"`                   //结束时间
	CreateTime   *time.Time     `json:"create_time,omitempty" gorm:"column:create_time"`   //创建时间
	UpdateTime   *time.Time     `json:"-" gorm:"column:update_time"`                       //更新时间
	DeleteTime   gorm.DeletedAt `json:"-" gorm:"column:delete_time"`                       //删除时间

}

func (Activity) TableName() string {
	return ACTIVITY_TABLE
}

type ActivityDao struct {
	BaseMysql
}

func NewActivityDao(appCtx context.Context) *ActivityDao {
	return &ActivityDao{
		BaseMysql{
			db:     global.GormDB.Model(&Activity{}),
			appCtx: appCtx,
		},
	}
}

func NewActivityDaoWithDB(db *gorm.DB, appCtx context.Context) *ActivityDao {
	return &ActivityDao{
		BaseMysql{
			db:     db.Model(&Activity{}),
			appCtx: appCtx,
		},
	}
}

func (t *ActivityDao) GetYearsByGroup() (years []int, err error) {
	err = t.db.WithContext(t.appCtx).Group("year").Order("year desc").Pluck("year", &years).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *ActivityDao) GetLastYear(year int) (lastYear int, err error) {
	var item Activity
	err = t.db.WithContext(t.appCtx).Where("year < ?", year).Order("year desc").First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		} else {
			err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		}
	}
	lastYear = item.Year
	if lastYear == 0 {
		lastYear = year - 1
	}
	return
}

func (t *ActivityDao) GetOngoingActivity() (activityIds []int, err error) {
	nowTime := util.NowTime()
	err = t.db.WithContext(t.appCtx).
		Where("start_time <= ? and end_time > ?", &nowTime, &nowTime).
		Pluck("activity_id", &activityIds).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *ActivityDao) GetActivityIdsByYear(year int) (activityIds []int, err error) {
	err = t.db.WithContext(t.appCtx).
		Where(Activity{Year: year}).
		Pluck("activity_id", &activityIds).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}
