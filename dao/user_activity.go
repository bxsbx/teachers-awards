package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"strconv"
	"strings"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"time"
)

const (
	USER_ACTIVITY_TABLE = "user_activity"
)

type UserActivity struct {
	UserActivityId int64                 `json:"user_activity_id" gorm:"column:user_activity_id;primary_key"`          //主键id
	ActivityId     int                   `json:"activity_id" gorm:"column:activity_id"`                                //活动id
	UserId         string                `json:"user_id" gorm:"column:user_id"`                                        //用户id
	UserName       string                `json:"user_name" gorm:"column:user_name"`                                    //用户姓名
	UserSex        int                   `json:"user_sex" gorm:"column:user_sex"`                                      //1：男，2：女
	Birthday       string                `json:"birthday" gorm:"column:birthday"`                                      //出生日期
	IdentityCard   string                `json:"identity_card" gorm:"column:identity_card"`                            //身份证号
	Phone          string                `json:"phone" gorm:"column:phone"`                                            //手机号
	SubjectCode    string                `json:"subject_code" gorm:"column:subject_code"`                              //科目code
	SchoolId       string                `json:"school_id" gorm:"column:school_id"`                                    //学校id
	SchoolName     string                `json:"school_name" gorm:"column:school_name"`                                //学校名称
	DeclareType    int                   `json:"declare_type" gorm:"column:declare_type"`                              //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
	FinalScore     float64               `json:"final_score" gorm:"column:final_score"`                                //最终得分（各项通过的审核）
	Rank           int                   `json:"rank" gorm:"column:rank"`                                              //排名
	RankPrize      *int                  `json:"rank_prize" gorm:"column:rank_prize"`                                  //0：无，1：一等奖，2：二等奖，3：三等奖
	Prize          int                   `json:"prize" gorm:"column:prize"`                                            //奖金
	CreateTime     *time.Time            `json:"create_time" gorm:"column:create_time"`                                //创建时间
	UpdateTime     *time.Time            `json:"-" gorm:"column:update_time"`                                          //更新时间
	DeleteTime     *time.Time            `json:"-" gorm:"column:delete_time"`                                          //删除时间
	DeleteAt       soft_delete.DeletedAt `json:"-" gorm:"column:delete_at;softDelete:milli,DeletedAtField:DeleteTime"` //删除时间戳（毫秒）

}

func (UserActivity) TableName() string {
	return USER_ACTIVITY_TABLE
}

type UserActivityDao struct {
	BaseMysql
}

func NewUserActivityDao(appCtx context.Context) *UserActivityDao {
	return &UserActivityDao{
		BaseMysql{
			db:     global.GormDB.Model(&UserActivity{}),
			appCtx: appCtx,
		},
	}
}

func NewUserActivityDaoWithDB(db *gorm.DB, appCtx context.Context) *UserActivityDao {
	return &UserActivityDao{
		BaseMysql{
			db:     db.Model(&UserActivity{}),
			appCtx: appCtx,
		},
	}
}

func (t *UserActivityDao) GetActivityUserNumMap(activityIds []int) (activityUserNumMap map[int]int, err error) {
	activityUserNumMap = make(map[int]int)
	rows, err := t.db.WithContext(t.appCtx).Select("activity_id,count(1) as count").
		Where("activity_id in (?)", activityIds).Group("activity_id").Rows()
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var activityId, count int
		rows.Scan(&activityId, &count)
		activityUserNumMap[activityId] = count
	}
	return
}

func (t *UserActivityDao) GetOnlyUserActivity(activityId int, userId string) (userActivity UserActivity) {
	err := t.db.WithContext(t.appCtx).Where(UserActivity{ActivityId: activityId, UserId: userId}).First(&userActivity).Error
	if err == gorm.ErrRecordNotFound {
		t.db.Error = nil
	}
	return
}

func (t *UserActivityDao) whereYear(db *gorm.DB, year int) *gorm.DB {
	if year > 0 {
		db = db.Where("activity_id in (select activity_id from activity where year = ? and delete_time is null)", year)
	}
	return db
}

func (t *UserActivityDao) GetAwardNum(year int, userActivity UserActivity) (count int64, err error) {
	db := t.db.WithContext(t.appCtx)
	err = t.whereYear(db, year).Where("rank_prize > ?", 0).
		Where(userActivity).Count(&count).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) GetDeclareNum(year int, userActivity UserActivity) (count int64, err error) {
	db := t.db.WithContext(t.appCtx)
	err = t.whereYear(db, year).Where(userActivity).Count(&count).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

//func (t *UserActivityDao) GetDeclarePersonNum(year int, userActivity UserActivity) (count int64, err error) {
//	db := t.db.WithContext(t.appCtx)
//	err = t.whereYear(db, year).Where(userActivity).Group("user_id").Count(&count).Error
//	if err != nil {
//		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
//	}
//	return
//}

func (t *UserActivityDao) GetEveryYearDeclareAwardNum(schoolId string) (list []YearDeclareAwardNum, err error) {
	err = t.db.WithContext(t.appCtx).Table("user_activity ua").
		Joins("JOIN activity a on a.activity_id = ua.activity_id").
		Where(UserActivity{SchoolId: schoolId}).
		Where("a.delete_time IS NULL").
		Select(`a.year, 
			COUNT(CASE WHEN ua.rank_prize > ? THEN 1 ELSE NULL END) AS award_num,
			COUNT(1) as declare_num`, 0).
		Group("a.year").Order("a.year asc").Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) GetEverySchoolAwardNum(year int) (list []SchoolAwardNum, err error) {
	db := t.db.WithContext(t.appCtx)
	err = t.whereYear(db, year).Select("school_id,school_name,count(1) as award_num").
		Group("school_id,school_name").Order("award_num desc").Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) GetEveryTeacherTypeAwardNum(year int) (list []TeacherTypeAwardNum, err error) {
	db := t.db.WithContext(t.appCtx)
	err = t.whereYear(db, year).Select("declare_type,count(1) as award_num").
		Group("declare_type").Order("award_num desc").Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) selectDeclareAwardNum() string {
	return `COUNT(CASE WHEN ua.rank_prize > 0 THEN 1 ELSE NULL END) AS award_num,
						COUNT(CASE WHEN ua.rank_prize = 1 THEN 1 ELSE NULL END) AS first_prize_num,
						COUNT(CASE WHEN ua.rank_prize = 2 THEN 1 ELSE NULL END) AS second_prize_num,
						COUNT(CASE WHEN ua.rank_prize = 3 THEN 1 ELSE NULL END) AS third_prize_num,
						COUNT(1) AS declare_num`
}

func (t *UserActivityDao) GetDeclareAwardNumGroupByYear(schoolId string) (list []YearDeclareAwardRankNum, err error) {
	err = t.db.WithContext(t.appCtx).Table("user_activity ua").
		Where(UserActivity{SchoolId: schoolId}).
		Joins("JOIN activity a on a.activity_id = ua.activity_id").
		Where("a.delete_time IS NULL").
		Select("a.year," + t.selectDeclareAwardNum()).
		Group("a.year").
		Order("a.year desc").Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) GetDeclareAwardRankGroupBySchool(year int, schoolName string, page, limit int) (total int64, list []SchoolDeclareAwardRankNum, err error) {
	db := t.db.WithContext(t.appCtx).Table("user_activity ua")
	if schoolName != "" {
		db = db.Where("ua.school_name like ?", "%"+schoolName+"%")
	}
	db = db.Joins("JOIN activity a on a.activity_id = ua.activity_id").
		Where("a.year = ? AND a.delete_time IS NULL", year).
		Group("ua.school_id,ua.school_name")
	err = db.Count(&total).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	err = db.Select("ua.school_id,ua.school_name," + t.selectDeclareAwardNum()).
		Order("award_num desc").
		Offset((page - 1) * limit).Limit(limit).
		Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) GetDeclareAwardRankGroupByDeclareType(year int) (list []TypeDeclareAwardRankNum, err error) {
	err = t.db.WithContext(t.appCtx).Table("user_activity ua").
		Joins("JOIN activity a on a.activity_id = ua.activity_id").
		Where("a.year = ? AND a.delete_time IS NULL", year).
		Group("ua.declare_type").
		Select("ua.declare_type," + t.selectDeclareAwardNum()).
		Order("award_num desc").
		Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) GetActivityInfo(userActivityId int64) (activityInfo Activity, err error) {
	var item UserActivity
	err = t.db.WithContext(t.appCtx).
		Where(UserActivity{UserActivityId: userActivityId}).
		First(&item).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	err = t.db.WithContext(t.appCtx).Model(&Activity{}).
		Where(Activity{ActivityId: item.ActivityId}).First(&activityInfo).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) GetIndicatorInfo(userActivityId int64, twoId int) (twoIndicatorInfo TwoIndicatorInfo, err error) {
	var userActivity UserActivity
	err = t.db.WithContext(t.appCtx).Model(&UserActivity{}).
		Where(UserActivity{UserActivityId: userActivityId}).First(&userActivity).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	var activityTwoIndicator ActivityTwoIndicator
	err = t.db.WithContext(t.appCtx).Model(&ActivityTwoIndicator{}).
		Where(ActivityTwoIndicator{ActivityId: userActivity.ActivityId, TwoIndicatorId: twoId}).First(&activityTwoIndicator).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	var activityOneIndicator ActivityOneIndicator
	err = t.db.WithContext(t.appCtx).Model(&ActivityOneIndicator{}).
		Where(ActivityOneIndicator{ActivityId: userActivity.ActivityId, OneIndicatorId: activityTwoIndicator.OneIndicatorId}).First(&activityOneIndicator).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	util.ObjToObjByReflect(&activityTwoIndicator, &twoIndicatorInfo)
	twoIndicatorInfo.OneIndicatorName = activityOneIndicator.OneIndicatorName
	return
}

func (t *UserActivityDao) GetHistoryActivityList(params *req.GetHistoryActivityListReq) (total int64, list []ActivityAndUserInfo, err error) {
	where := UserActivity{
		UserSex:      params.UserSex,
		SubjectCode:  params.SubjectCode,
		SchoolId:     params.SchoolId,
		IdentityCard: params.IdentityCard,
		RankPrize:    params.RankPrize,
		Rank:         params.Rank,
		FinalScore:   params.FinalScore,
		DeclareType:  params.DeclareType,
	}

	db := t.db.WithContext(t.appCtx).Table("user_activity ua").
		Joins("JOIN activity a on a.activity_id = ua.activity_id").
		Where("ua.rank != ?", 0).
		Where("a.delete_time IS NULL").
		Where(where)

	if params.UserName != "" {
		db = db.Where("ua.user_name like ?", "%"+params.UserName+"%")
	}
	if params.SchoolName != "" {
		db = db.Where("ua.school_name like ?", "%"+params.SchoolName+"%")
	}
	if params.Year != 0 {
		db = db.Where("a.year = ?", params.Year)
	}
	err = db.Count(&total).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	if params.Page > 0 && params.Limit > 0 {
		db = db.Offset((params.Page - 1) * params.Limit).Limit(params.Limit)
	}
	err = db.Select("ua.*,a.activity_name,a.year,a.create_time").
		Order("a.create_time desc,ua.rank").Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) GetUserHistoryDeclareResultList(userId string, page, limit int) (total int64, list []UserHistoryDeclareResult, err error) {
	db := t.db.WithContext(t.appCtx).Table("user_activity ua").
		Joins("JOIN activity a ON ua.activity_id =  a.activity_id").
		Where("ua.rank != ?", 0).
		Where("ua.user_id = ?", userId).
		Where("a.delete_time IS NULL")
	err = db.Count(&total).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	err = db.Offset((page - 1) * limit).Limit(limit).Select("a.activity_id,a.activity_name,ua.final_score,ua.rank,ua.rank_prize,ua.prize").
		Order("a.create_time desc").Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) GetResultGroupByDeclareType(activityId int) (list []ResultGroupByDeclareType, err error) {
	err = t.db.WithContext(t.appCtx).Table("user_activity ua").
		Where("ua.activity_id = ?", activityId).
		Group("ua.declare_type").
		Select("ua.declare_type,COUNT(DISTINCT(school_id)) AS school_num," + t.selectDeclareAwardNum() +
			`,COUNT(CASE WHEN ua.rank > 0 THEN 1 ELSE NULL END) AS rank_num`).
		Order("ua.declare_type asc").
		Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

// status --1:待审核，2：已审核
func (t *UserActivityDao) getEdbReviewListSql(params *req.GetEdbReviewListReq, status int) (total int64, db *gorm.DB, err error) {
	where := UserActivity{
		UserSex:     params.UserSex,
		SubjectCode: params.SubjectCode,
		SchoolId:    params.SchoolId,
	}
	db = t.db.WithContext(t.appCtx).Table("user_activity ua").
		Where("ua.rank = ?", 0).
		Where(where)

	if params.UserName != "" {
		db = db.Where("ua.user_name like ?", "%"+params.UserName+"%")
	}
	if params.SchoolName != "" {
		db = db.Where("ua.school_name like ?", "%"+params.SchoolName+"%")
	}

	exist := "select 1 from user_activity_indicator uai where uai.user_activity_id = ua.user_activity_id and uai.review_process = ? and uai.delete_at = ?"
	var arg []interface{}
	arg = append(arg, 3, 0)

	if status == 1 {
		exist = "exists (" + exist + ")"
	} else if status == 2 {
		exist = "not exists (" + exist + ")"
	}

	db = db.Where(exist, arg...)
	err = db.Count(&total).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}

	db = db.Select(`ua.user_activity_id,ua.user_id,ua.user_name,ua.identity_card,ua.user_sex,ua.subject_code,
				ua.school_id,ua.school_name,ua.declare_type,ua.final_score,(` + strconv.Itoa(status) + `) as status,ua.create_time`)
	return
}

func (t *UserActivityDao) GetEdbReviewList(params *req.GetEdbReviewListReq) (total int64, list []EdbReviewItem, err error) {
	var sql1, sql2 *gorm.DB
	var total1, total2 int64
	wg := util.NewWaitGroup(0)
	if params.Status == 1 || params.Status == 2 {
		wg.Go(func() (err error) {
			total1, sql1, err = t.getEdbReviewListSql(params, params.Status)
			return
		})
	} else {
		wg.Go(func() (err error) {
			total1, sql1, err = t.getEdbReviewListSql(params, 1)
			return
		})
		wg.Go(func() (err error) {
			total2, sql2, err = t.getEdbReviewListSql(params, 2)
			return
		})
	}
	if err = wg.Wait(); err != nil {
		return
	}
	total = total1 + total2
	if params.OnlyCount {
		return
	}
	db := t.db.WithContext(t.appCtx)
	pageSql := ""
	if params.Page > 0 && params.Limit > 0 {
		pageSql = fmt.Sprintf("limit %v,%v", (params.Page-1)*params.Limit, params.Limit)
	}
	var sql []interface{}
	var repeat []string
	if sql1 != nil {
		sql = append(sql, sql1)
		repeat = append(repeat, "(?)")
	}
	if sql2 != nil {
		sql = append(sql, sql2)
		repeat = append(repeat, "(?)")
	}
	err = db.Raw("SELECT * FROM ("+strings.Join(repeat, " UNION ALL ")+") AS temp"+" ORDER BY (CASE status WHEN 1 THEN 1 ELSE 2 END) ASC,create_time DESC "+pageSql, sql...).
		Scan(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityDao) GetEveryYearSchoolNumByGroup(years []int) (yearSchoolNumMap map[int]int, err error) {
	yearSchoolNumMap = make(map[int]int)
	rows, err := t.db.WithContext(t.appCtx).Table("user_activity ua").
		Joins("JOIN activity a ON ua.activity_id =  a.activity_id").
		Select("a.year,COUNT(DISTINCT(school_id)) AS num").
		Where("year IN (?)", years).
		Group("year").Rows()
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	defer rows.Close()
	for rows.Next() {
		var year int
		var schoolNum int
		rows.Scan(&year, &schoolNum)
		yearSchoolNumMap[year] = schoolNum
	}
	return
}
