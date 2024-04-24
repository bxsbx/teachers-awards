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
	USER_ACTIVITY_INDICATOR_TABLE = "user_activity_indicator"
)

type UserActivityIndicator struct {
	UserActivityIndicatorId int64                 `json:"user_activity_indicator_id" gorm:"column:user_activity_indicator_id;primary_key"` //主键id
	UserActivityId          int64                 `json:"user_activity_id" gorm:"column:user_activity_id"`                                 //用户申报活动的id
	TwoIndicatorId          int                   `json:"two_indicator_id" gorm:"column:two_indicator_id"`                                 //二级指标id
	AwardDate               *time.Time            `json:"award_date" gorm:"column:award_date"`                                             //获奖日期
	CertificateType         int                   `json:"certificate_type" gorm:"column:certificate_type"`                                 //1.证书；需要填写证书有效期，2.证明；不需要有效期
	CertificateUrl          string                `json:"certificate_url" gorm:"column:certificate_url"`                                   //证书url
	CertificateStartDate    *time.Time            `json:"certificate_start_date" gorm:"column:certificate_start_date"`                     //证书有效期——开始时间
	CertificateEndDate      *time.Time            `json:"certificate_end_date" gorm:"column:certificate_end_date"`                         //证书有效期——结束时间
	Status                  int                   `json:"status" gorm:"column:status"`                                                     //审核状态:1.已提交 ,2.已通过,3.学校未通过,4.专家未通过,5.教育局未通过,6.活动已结束，未审批,
	FinishReviewNum         *int                  `json:"finish_review_num" gorm:"column:finish_review_num"`                               //当前已审核人数
	ReviewProcess           int                   `json:"review_process" gorm:"column:review_process"`                                     //当前审核进程状态，1：学校，2：专家，3：教育局，4：结束
	CreateTime              *time.Time            `json:"create_time" gorm:"column:create_time"`                                           //创建时间
	UpdateTime              *time.Time            `json:"-" gorm:"column:update_time"`                                                     //更新时间
	DeleteTime              *time.Time            `json:"-" gorm:"column:delete_time"`                                                     //删除时间
	DeleteAt                soft_delete.DeletedAt `json:"-" gorm:"column:delete_at;softDelete:milli,DeletedAtField:DeleteTime"`            //删除时间戳（毫秒）

}

func (UserActivityIndicator) TableName() string {
	return USER_ACTIVITY_INDICATOR_TABLE
}

type UserActivityIndicatorDao struct {
	BaseMysql
}

func NewUserActivityIndicatorDao(appCtx context.Context) *UserActivityIndicatorDao {
	return &UserActivityIndicatorDao{
		BaseMysql{
			db:     global.GormDB.Model(&UserActivityIndicator{}),
			appCtx: appCtx,
		},
	}
}

func NewUserActivityIndicatorDaoWithDB(db *gorm.DB, appCtx context.Context) *UserActivityIndicatorDao {
	return &UserActivityIndicatorDao{
		BaseMysql{
			db:     db.Model(&UserActivityIndicator{}),
			appCtx: appCtx,
		},
	}
}

func (t *UserActivityIndicatorDao) GetWaitReviewList(userId string, reviewProcess, page, limit int) (total int64, list []WaitReview, err error) {
	db := t.db.WithContext(t.appCtx).Table("user_activity_indicator uai").
		Joins("JOIN user_activity ua ON ua.user_activity_id = uai.user_activity_id").
		Where(UserActivityIndicator{ReviewProcess: reviewProcess, Status: global.ReviewStatusCommit}).
		Where("ua.rank = ?", 0)
	//专家审核进程，无权限的过滤掉
	if reviewProcess == global.ProcessExpert {
		db = db.Where("exists (select 1 from expert_auth_indicator where user_id = ? and two_indicator_id = uai.two_indicator_id)", userId)
	}
	//Where("ua.activity_id in (?)", activityIds)
	err = db.Count(&total).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	err = db.Select("activity_id,user_id,user_name,user_activity_indicator_id,two_indicator_id,uai.create_time").
		Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

// reviewStatus --1:待审核，2：未通过，3：已通过，4：已审核"
func (t *UserActivityIndicatorDao) getReviewListSql(params *req.GetReviewListReq, reviewStatus int) (total int64, db *gorm.DB, err error) {
	where := UserActivity{
		UserSex:     params.UserSex,
		SubjectCode: params.SubjectCode,
		SchoolId:    params.SchoolId,
	}
	//nowTime := util.NowTime()
	db = t.db.WithContext(t.appCtx).Table("user_activity ua").
		Joins("JOIN user_activity_indicator uai ON uai.user_activity_id = ua.user_activity_id").
		Joins("JOIN activity a ON ua.activity_id = a.activity_id").
		Where("ua.rank = ?", 0).
		//Where("a.start_time <= ? and a.end_time > ?", &nowTime, &nowTime).
		Where("a.delete_time IS NULL AND ua.delete_at = ? AND uai.delete_at = ?", 0, 0).
		Where(where)

	if params.UserName != "" {
		db = db.Where("ua.user_name like ?", "%"+params.UserName+"%")
	}
	if params.SchoolName != "" {
		db = db.Where("ua.school_name like ?", "%"+params.SchoolName+"%")
	}

	if params.ReviewProcess != 0 {
		db = db.Where("uai.review_process >= ?", params.ReviewProcess)
	}
	if params.Year != 0 {
		db = db.Where("a.year  = ?", params.Year)
	}

	var wi []string
	var warg []interface{}
	if params.OneIndicatorName != "" {
		wi = append(wi, "aoi.one_indicator_name like ?")
		warg = append(warg, "%"+params.OneIndicatorName+"%")
	}
	if params.TwoIndicatorName != "" {
		wi = append(wi, "ati.two_indicator_name like ?")
		warg = append(warg, "%"+params.TwoIndicatorName+"%")
	}

	if len(wi) > 0 {
		db = db.Where(`uai.two_indicator_id IN (
				select two_indicator_id from activity_two_indicator ati 
					JOIN activity_one_indicator aoi ON aoi.activity_id = ati.activity_id and aoi.one_indicator_id = ati.one_indicator_id where `+
			strings.Join(wi, " and ")+`)`, warg...)
	}

	exist := "select 1 from judges_verify jv where jv.user_activity_indicator_id = uai.user_activity_indicator_id and jv.judges_type = ? and jv.delete_at = ?"
	var jvarg []interface{}
	jvarg = append(jvarg, params.ReviewProcess, 0)

	if params.ReviewProcess == global.ProcessExpert {
		exist += " and jv.judges_id = ?"
		jvarg = append(jvarg, params.JudgesId)
	}

	if reviewStatus == global.ReviewNoPass || reviewStatus == global.ReviewPass {
		exist += " and jv.is_pass = ?"
		jvarg = append(jvarg, reviewStatus-2)
	}

	if reviewStatus == global.ReviewWait {
		if params.ReviewProcess == global.ProcessExpert {
			exist = "NOT EXISTS (" + exist + ") and uai.finish_review_num < a.review_num-1"
			exist += " and EXISTS (select 1 from expert_auth_indicator where user_id = ? and two_indicator_id = uai.two_indicator_id)"
			jvarg = append(jvarg, params.JudgesId)
		} else {
			exist = "NOT EXISTS (" + exist + ")"
		}
	} else if reviewStatus > global.ReviewWait {
		exist = "EXISTS (" + exist + ")"
	}

	db = db.Where(exist, jvarg...)

	err = db.Count(&total).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}

	db = db.Select(`a.activity_id,a.year,ua.user_id,ua.user_name,ua.user_sex,ua.subject_code,ua.school_id,ua.school_name,ua.declare_type,
			uai.user_activity_indicator_id,uai.two_indicator_id,uai.status,uai.create_time,
			(` + strconv.Itoa(reviewStatus) + `) as review_status`)
	return
}

func (t *UserActivityIndicatorDao) GetReviewList(params *req.GetReviewListReq) (total int64, list []ReviewItem, err error) {
	var sql1, sql2, sql3 *gorm.DB
	var total1, total2, total3 int64
	wg := util.NewWaitGroup(0)
	if params.ReviewStatus == global.ReviewWait || params.ReviewStatus == global.ReviewNoPass || params.ReviewStatus == global.ReviewPass {
		wg.Go(func() (err error) {
			total1, sql1, err = t.getReviewListSql(params, params.ReviewStatus)
			return
		})
	} else if params.ReviewStatus == global.ReviewFinish {
		wg.Go(func() (err error) {
			total2, sql2, err = t.getReviewListSql(params, global.ReviewNoPass)
			return
		})
		wg.Go(func() (err error) {
			total3, sql3, err = t.getReviewListSql(params, global.ReviewPass)
			return
		})
	} else {
		wg.Go(func() (err error) {
			total1, sql1, err = t.getReviewListSql(params, global.ReviewWait)
			return
		})
		wg.Go(func() (err error) {
			total2, sql2, err = t.getReviewListSql(params, global.ReviewNoPass)
			return
		})
		wg.Go(func() (err error) {
			total3, sql3, err = t.getReviewListSql(params, global.ReviewPass)
			return
		})
	}
	if err = wg.Wait(); err != nil {
		return
	}
	total = total1 + total2 + total3
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
	if sql3 != nil {
		sql = append(sql, sql3)
		repeat = append(repeat, "(?)")
	}
	err = db.Raw("SELECT * FROM ("+strings.Join(repeat, " UNION ALL ")+") AS temp"+" ORDER BY (CASE review_status WHEN 1 THEN 1 ELSE 2 END) ASC,create_time DESC "+pageSql, sql...).
		Scan(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityIndicatorDao) GetIndicatorInfo(uaiId int64) (twoIndicator TwoIndicatorInfo, err error) {
	var userActivityIndicator UserActivityIndicator
	err = t.db.WithContext(t.appCtx).Model(&UserActivityIndicator{}).
		Where(UserActivityIndicator{UserActivityIndicatorId: uaiId}).First(&userActivityIndicator).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	var userActivity UserActivity
	err = t.db.WithContext(t.appCtx).Model(&UserActivity{}).
		Where(UserActivity{UserActivityId: userActivityIndicator.UserActivityId}).First(&userActivity).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	var activityTwoIndicator ActivityTwoIndicator
	err = t.db.WithContext(t.appCtx).Model(&ActivityTwoIndicator{}).
		Where(ActivityTwoIndicator{ActivityId: userActivity.ActivityId, TwoIndicatorId: userActivityIndicator.TwoIndicatorId}).First(&activityTwoIndicator).Error
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
	util.ObjToObjByReflect(&activityTwoIndicator, &twoIndicator)
	twoIndicator.OneIndicatorName = activityOneIndicator.OneIndicatorName
	return
}

func (t *UserActivityIndicatorDao) GetUserDeclareRecordList(userId string, year int, page, limit int) (total int64, list []UserDeclareRecord, err error) {
	db := t.db.WithContext(t.appCtx).Table("user_activity_indicator uai").
		Joins("JOIN user_activity ua ON ua.user_activity_id =  uai.user_activity_id").
		Where("activity_id in (select activity_id from activity where year = ? and delete_time IS NULL)", year).
		Where("ua.user_id = ? and ua.delete_at = ?", userId, 0)
	err = db.Count(&total).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
		return
	}
	if page > 0 && limit > 0 {
		db = db.Offset((page - 1) * limit).Limit(limit)
	}
	err = db.Select("uai.user_activity_indicator_id,uai.status,uai.two_indicator_id,uai.review_process,uai.create_time,ua.activity_id,ua.user_activity_id").
		Order("uai.create_time desc").Find(&list).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}

func (t *UserActivityIndicatorDao) GetDeclareNum(year int, schoolId string) (count int64, err error) {
	db := t.db.WithContext(t.appCtx).Table("user_activity_indicator uai").
		Joins("inner join user_activity ua on ua.user_activity_id = uai.user_activity_id")
	if year > 0 {
		db = db.Where("ua.activity_id in (select activity_id from activity where year = ? and delete_time is null)", year)
	}
	if schoolId != "" {
		db = db.Where("ua.school_id = ?", schoolId)
	}
	err = db.Count(&count).Error
	if err != nil {
		err = errorz.CodeError(errorz.DB_SELECT_ERR, err)
	}
	return
}
