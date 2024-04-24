package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"time"
)

const (
	ONE_INDICATOR_TABLE = "one_indicator"
)

type OneIndicator struct {
	OneIndicatorId   int            `json:"one_indicator_id" gorm:"column:one_indicator_id;primary_key"` //主键id
	OneIndicatorName string         `json:"one_indicator_name" gorm:"column:one_indicator_name"`         //一级指标名称
	Content          string         `json:"content,omitempty" gorm:"column:content"`                     //评分标准说明
	CreateTime       *time.Time     `json:"create_time,omitempty" gorm:"column:create_time"`             //创建时间
	UpdateTime       *time.Time     `json:"-" gorm:"column:update_time"`                                 //更新时间
	DeleteTime       gorm.DeletedAt `json:"-" gorm:"column:delete_time"`                                 //删除时间

}

func (OneIndicator) TableName() string {
	return ONE_INDICATOR_TABLE
}

type OneIndicatorDao struct {
	BaseMysql
}

func NewOneIndicatorDao(appCtx context.Context) *OneIndicatorDao {
	return &OneIndicatorDao{
		BaseMysql{
			db:     global.GormDB.Model(&OneIndicator{}),
			appCtx: appCtx,
		},
	}
}

func NewOneIndicatorDaoWithDB(db *gorm.DB, appCtx context.Context) *OneIndicatorDao {
	return &OneIndicatorDao{
		BaseMysql{
			db:     db.Model(&OneIndicator{}),
			appCtx: appCtx,
		},
	}
}

func (t *OneIndicatorDao) GetOneIndicatorNameMap(oneIds []int) (oneNameMap map[int]string, err error) {
	oneNameMap = make(map[int]string)
	rows, err := t.db.WithContext(t.appCtx).Select("one_indicator_id,one_indicator_name").
		Where("one_indicator_id in (?)", oneIds).Rows()
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var oneId int
		var oneName string
		rows.Scan(&oneId, &oneName)
		oneNameMap[oneId] = oneName
	}
	return
}
