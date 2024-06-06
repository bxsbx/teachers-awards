package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"time"
)

const (
	JUDGES_VERIFY_TABLE = "judges_verify"
)

type JudgesVerify struct {
	UserIndicatorPassId     int64                 `json:"user_indicator_pass_id" gorm:"column:user_indicator_pass_id;primary_key"` //主键id
	UserActivityIndicatorId int64                 `json:"user_activity_indicator_id" gorm:"column:user_activity_indicator_id"`     //用户活动申报二级指标id 	//活动id
	UserActivityId          int64                 `json:"user_activity_id" gorm:"column:user_activity_id"`                         //0：未通过，1：通过
	JudgesId                string                `json:"judges_id" gorm:"column:judges_id"`                                       //评委id
	JudgesName              string                `json:"judges_name" gorm:"column:judges_name"`                                   //评委姓名
	JudgesType              int                   `json:"judges_type" gorm:"column:judges_type"`                                   //评委类型，1：学校，2：专家，3：教育局
	IsPass                  int                   `json:"is_pass" gorm:"column:is_pass"`                                           //0：未通过，1：通过
	Score                   float64               `json:"score" gorm:"column:score"`                                               //得分
	Opinion                 string                `json:"opinion" gorm:"column:opinion"`                                           //审核意见
	CreateTime              *time.Time            `json:"create_time" gorm:"column:create_time"`                                   //创建时间
	DeleteTime              *time.Time            `json:"-" gorm:"column:delete_time"`                                             //删除时间
	DeleteAt                soft_delete.DeletedAt `json:"-" gorm:"column:delete_at;softDelete:milli,DeletedAtField:DeleteTime"`    //删除时间戳（毫秒）

}

func (JudgesVerify) TableName() string {
	return JUDGES_VERIFY_TABLE
}

type JudgesVerifyDao struct {
	BaseMysql
}

func NewJudgesVerifyDao(appCtx context.Context) *JudgesVerifyDao {
	return &JudgesVerifyDao{
		BaseMysql{
			db:     global.GormDB.Model(&JudgesVerify{}),
			appCtx: appCtx,
		},
	}
}

func NewJudgesVerifyDaoWithDB(db *gorm.DB, appCtx context.Context) *JudgesVerifyDao {
	return &JudgesVerifyDao{
		BaseMysql{
			db:     db.Model(&JudgesVerify{}),
			appCtx: appCtx,
		},
	}
}

func (t *JudgesVerifyDao) GetJudgesVerifyListByIds(uaiIds []int64) (verifyList []JudgesVerify, err error) {
	err = t.db.WithContext(t.appCtx).Where("user_activity_indicator_id in (?)", uaiIds).
		Select("user_activity_indicator_id,judges_id,judges_name,judges_type,is_pass").
		Find(&verifyList).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	return
}
