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
	OPERATION_RECORD_TABLE = "operation_record"
)

type OperationRecord struct {
	OperationId    int64      `json:"operation_id" gorm:"column:operation_id;primary_key"` //操作id
	RelationalId   int64      `json:"relational_id" gorm:"column:relational_id"`           //关联的id
	RelationalType int        `json:"relational_type" gorm:"column:relational_type"`       //关联表，1：user_activity_indicator
	UserId         string     `json:"user_id" gorm:"column:user_id"`                       //用户id
	UserName       string     `json:"user_name" gorm:"column:user_name"`                   //用户姓名
	OperationRole  int        `json:"operation_role" gorm:"column:operation_role"`         //角色，1：学校，2：专家，3：教育局，4：教师
	OperationType  int        `json:"operation_type" gorm:"column:operation_type"`         //操作类型，1：添加，2：修改，3：删除
	Description    string     `json:"description" gorm:"column:description"`               //操作说明
	CreateTime     *time.Time `json:"create_time" gorm:"column:create_time"`               //创建时间
}

func (OperationRecord) TableName() string {
	return OPERATION_RECORD_TABLE
}

type OperationRecordDao struct {
	BaseMysql
}

func NewOperationRecordDao(appCtx context.Context) *OperationRecordDao {
	return &OperationRecordDao{
		BaseMysql{
			db:     global.GormDB.Model(&OperationRecord{}),
			appCtx: appCtx,
		},
	}
}

func NewOperationRecordDaoWithDB(db *gorm.DB, appCtx context.Context) *OperationRecordDao {
	return &OperationRecordDao{
		BaseMysql{
			db:     db.Model(&OperationRecord{}),
			appCtx: appCtx,
		},
	}
}

func (t *OperationRecordDao) InsertOperationRecordToUAI(relationalId int64, relationalType, operationType, operationRole int, description string) (err error) {
	nowTime := util.NowTime()
	userInfo := global.GetUserInfo(t.appCtx)
	err = t.db.WithContext(t.appCtx).Create(&OperationRecord{
		RelationalId:   relationalId,
		RelationalType: relationalType,
		UserId:         userInfo.UserId,
		UserName:       userInfo.UserName,
		OperationRole:  operationRole,
		OperationType:  operationType,
		Description:    description,
		CreateTime:     &nowTime,
	}).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_CREATE_ERR, err)
		return
	}
	return
}
