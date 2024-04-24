package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"time"
)

const (
	TWO_INDICATOR_TABLE = "two_indicator"
)

type TwoIndicator struct {
	TwoIndicatorId   int            `json:"two_indicator_id" gorm:"column:two_indicator_id;primary_key"` //主键id
	TwoIndicatorName string         `json:"two_indicator_name" gorm:"column:two_indicator_name"`         //二级指标名称
	Score            int            `json:"score" gorm:"column:score"`                                   //分值
	OneIndicatorId   int            `json:"one_indicator_id" gorm:"column:one_indicator_id"`             //所属一级指标
	CreateTime       *time.Time     `json:"create_time" gorm:"column:create_time"`                       //创建时间
	UpdateTime       *time.Time     `json:"-" gorm:"column:update_time"`                                 //更新时间
	DeleteTime       gorm.DeletedAt `json:"-" gorm:"column:delete_time"`                                 //删除时间

}

func (TwoIndicator) TableName() string {
	return TWO_INDICATOR_TABLE
}

type TwoIndicatorDao struct {
	BaseMysql
}

func NewTwoIndicatorDao(appCtx context.Context) *TwoIndicatorDao {
	return &TwoIndicatorDao{
		BaseMysql{
			db:     global.GormDB.Model(&TwoIndicator{}),
			appCtx: appCtx,
		},
	}
}

func NewTwoIndicatorDaoWithDB(db *gorm.DB, appCtx context.Context) *TwoIndicatorDao {
	return &TwoIndicatorDao{
		BaseMysql{
			db:     db.Model(&TwoIndicator{}),
			appCtx: appCtx,
		},
	}
}

func (t *TwoIndicatorDao) GetTwoIndicatorCountMap(oneIds []int) (oneIdCountMap map[int]int, err error) {
	oneIdCountMap = make(map[int]int)
	rows, err := t.db.WithContext(t.appCtx).Select("one_indicator_id,count(*)").
		Where("one_indicator_id in (?)", oneIds).Group("one_indicator_id").Rows()
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var oneId, count int
		rows.Scan(&oneId, &count)
		oneIdCountMap[oneId] = count
	}
	return
}
