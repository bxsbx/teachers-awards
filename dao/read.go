package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"time"
)

const (
	READ_TABLE = "read"
)

type Read struct {
	ReadId     int        `json:"read_id" gorm:"column:read_id;primary_key"`           //已读id
	ReadType   int        `json:"read_type" gorm:"column:read_type;primary_key"`       //1：公告
	ReadUserId string     `json:"read_user_id" gorm:"column:read_user_id;primary_key"` //已读用户id
	CreateTime *time.Time `json:"create_time" gorm:"column:create_time"`               //创建时间

}

func (Read) TableName() string {
	return READ_TABLE
}

type ReadDao struct {
	BaseMysql
}

func NewReadDao(appCtx context.Context) *ReadDao {
	return &ReadDao{
		BaseMysql{
			db:     global.GormDB.Model(&Read{}),
			appCtx: appCtx,
		},
	}
}

func NewReadDaoWithDB(db *gorm.DB, appCtx context.Context) *ReadDao {
	return &ReadDao{
		BaseMysql{
			db:     db.Model(&Read{}),
			appCtx: appCtx,
		},
	}
}

func (t *ReadDao) GetReadMap(readIds []int, readType int, userId string) (readMap map[int]bool, err error) {
	readMap = make(map[int]bool)
	rows, err := t.db.WithContext(t.appCtx).Select("read_id").
		Where("read_id in (?) and read_type = ? and read_user_id = ?", readIds, readType, userId).Rows()
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		readMap[id] = true
	}
	return
}
