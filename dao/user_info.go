package dao

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
	"time"
)

const (
	USER_INFO_TABLE = "user_info"
)

type UserInfo struct {
	UserId       string     `json:"user_id" gorm:"column:user_id;primary_key"` //用户id
	UserName     string     `json:"user_name" gorm:"column:user_name"`         //用户名称
	UserSex      int        `json:"user_sex" gorm:"column:user_sex"`           //1：男，2：女
	Birthday     string     `json:"birthday" gorm:"column:birthday"`           //出生日期
	IdentityCard string     `json:"identity_card" gorm:"column:identity_card"` //身份证号
	Phone        string     `json:"phone" gorm:"column:phone"`                 //手机号
	Avatar       string     `json:"avatar" gorm:"column:avatar"`               //头像
	SubjectCode  string     `json:"subject_code" gorm:"column:subject_code"`   //科目code
	SchoolId     string     `json:"school_id" gorm:"column:school_id"`         //学校id
	SchoolName   string     `json:"school_name" gorm:"column:school_name"`     //学校名称
	DeclareType  int        `json:"declare_type" gorm:"column:declare_type"`   //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
	ExportAuth   int        `json:"-" gorm:"column:export_auth"`               //专家是否授权 1：未授权 2：已授权
	AuthDay      *time.Time `json:"-" gorm:"column:auth_day"`                  //授权日期，格式2006-01-02
	Role         int        `json:"-" gorm:"column:role"`                      //角色，1：学校，2：专家，4：教育局，8：教师，16：超级管理员，多个则相加
	Roles        []int      `json:"roles" gorm:"-"`                            //角色，1：学校，2：专家，3：教育局，4：教师，5：超级管理员
	Year         int        `json:"year" gorm:"column:year"`                   //年份
	CreateTime   *time.Time `json:"create_time" gorm:"column:create_time"`     //创建时间
	UpdateTime   *time.Time `json:"-" gorm:"column:update_time"`               //更新时间

	ExpertAuthIndicatorList []ExpertAuthIndicator `json:"-" gorm:"foreignKey:UserId;references:UserId;"` //专家指标权限

}

func (UserInfo) TableName() string {
	return USER_INFO_TABLE
}

type UserInfoDao struct {
	BaseMysql
}

func NewUserInfoDao(appCtx context.Context) *UserInfoDao {
	return &UserInfoDao{
		BaseMysql{
			db:     global.GormDB.Model(&UserInfo{}),
			appCtx: appCtx,
		},
	}
}

func NewUserInfoDaoWithDB(db *gorm.DB, appCtx context.Context) *UserInfoDao {
	return &UserInfoDao{
		BaseMysql{
			db:     db.Model(&UserInfo{}),
			appCtx: appCtx,
		},
	}
}

func (t *UserInfoDao) GetEveryYearTeacherNumByGroup(schoolId string) (yearTeacherNumMap map[int]int, err error) {
	yearTeacherNumMap = make(map[int]int)
	rows, err := t.db.WithContext(t.appCtx).
		Select("year,COUNT(user_id) AS num").
		Where("role & ? > 0", 8).
		Where("school_id = ?", schoolId).
		Group("year").Rows()
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	defer rows.Close()
	for rows.Next() {
		var year int
		var teacherNum int
		rows.Scan(&year, &teacherNum)
		yearTeacherNumMap[year] = teacherNum
	}
	return
}

func (t *UserInfoDao) GetSchoolNumByYear(year int) (count int64, err error) {
	err = t.db.WithContext(t.appCtx).
		Where("year <= ?", year).
		Where("role & ? > 0", 8).
		Group("school_id").Count(&count).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserInfoDao) GetTeacherNumByYear(year int, schoolId string) (count int64, err error) {
	db := t.db.WithContext(t.appCtx)
	if year != 0 {
		db = db.Where("year <= ?", year)
	}
	err = db.Where(UserInfo{SchoolId: schoolId}).
		Where("role & ? > 0", 8).
		Count(&count).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}
