package services

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/dao"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/model/resp"
)

type ReviewDeclareService struct {
	appCtx context.Context
}

func NewReviewDeclareService(appCtx context.Context) *ReviewDeclareService {
	return &ReviewDeclareService{appCtx: appCtx}
}

func (s *ReviewDeclareService) GetWaitReviewList(params *req.GetWaitReviewListReq) (data *resp.GetWaitReviewListResp, err error) {
	data = &resp.GetWaitReviewListResp{}
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	userInfo := global.GetUserInfo(s.appCtx)
	total, waitReviews, err := userActivityIndicatorDao.GetWaitReviewList(userInfo.UserId, params.ReviewProcess, params.Page, params.Limit)
	if err != nil {
		return
	}
	data.Total = total

	activityIds := util.ListToDeduplicationList(waitReviews, func(t dao.WaitReview) (int, int) {
		return t.ActivityId, t.ActivityId
	})

	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	activityTwoIndicatorsMap, err := activityTwoIndicatorDao.GetTwoIndicatorInfoMap(activityIds, util.ListObjToListObj(waitReviews, func(obj dao.WaitReview) int {
		return obj.TwoIndicatorId
	}))
	if err != nil {
		return
	}

	data.List = make([]resp.WaitReview, len(waitReviews))
	for i, v := range waitReviews {
		activityIndicatorInfo := activityTwoIndicatorsMap[v.ActivityId][v.TwoIndicatorId]
		data.List[i].UserActivityIndicatorId = v.UserActivityIndicatorId
		data.List[i].UserId = v.UserId
		data.List[i].UserName = v.UserName
		data.List[i].OneIndicatorId = activityIndicatorInfo.OneIndicatorId
		data.List[i].OneIndicatorName = activityIndicatorInfo.OneIndicatorName
		data.List[i].TwoIndicatorId = activityIndicatorInfo.TwoIndicatorId
		data.List[i].TwoIndicatorName = activityIndicatorInfo.TwoIndicatorName
		data.List[i].DeclareTime = v.CreateTime
	}
	return
}

func (s *ReviewDeclareService) GetReviewList(params *req.GetReviewListReq) (data *resp.GetReviewListResp, err error) {
	data = &resp.GetReviewListResp{}
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	total, reviewItems, err := userActivityIndicatorDao.GetReviewList(params)
	if err != nil {
		return
	}
	data.Total = total
	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	activityIds := util.ListToDeduplicationList(reviewItems, func(obj dao.ReviewItem) (int, int) {
		return obj.ActivityId, obj.ActivityId
	})
	twoIds := util.ListToDeduplicationList(reviewItems, func(obj dao.ReviewItem) (int, int) {
		return obj.TwoIndicatorId, obj.TwoIndicatorId
	})
	twoIndicatorsMap, err := activityTwoIndicatorDao.GetTwoIndicatorInfoMap(activityIds, twoIds)
	if err != nil {
		return
	}
	data.List = make([]resp.ReviewItem, len(reviewItems))
	for i, v := range reviewItems {
		data.List[i].ReviewItem = v
		data.List[i].TwoIndicatorInfo = twoIndicatorsMap[v.ActivityId][v.TwoIndicatorId]
	}
	return
}

func (s *ReviewDeclareService) GetTwoInfo(userActivityId int64, twoId int) (twoInfo dao.ActivityTwoIndicator, err error) {
	var userActivity dao.UserActivity
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	err = userActivityDao.First(dao.UserActivity{UserActivityId: userActivityId}, &userActivity)
	if err != nil {
		return
	}

	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	err = activityTwoIndicatorDao.First(dao.ActivityTwoIndicator{ActivityId: userActivity.ActivityId, TwoIndicatorId: twoId}, &twoInfo)
	if err != nil {
		return
	}
	return
}

func (s *ReviewDeclareService) GetMaxScore(activityId, oneId, twoId int, userActivityId int64) (curMaxScore, preMaxScore float64, err error) {
	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	var twoList []dao.ActivityTwoIndicator
	err = activityTwoIndicatorDao.Find(dao.ActivityTwoIndicator{ActivityId: activityId, OneIndicatorId: oneId}, &twoList)
	if err != nil {
		return
	}

	var twoIds []int
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	err = userActivityIndicatorDao.Pluck(dao.UserActivityIndicator{UserActivityId: userActivityId, Status: global.ReviewStatusPass}, &twoIds, "TwoIndicatorId")
	if err != nil {
		return
	}
	twoMap := util.ListObjToMap(twoIds, func(twoId int) (int, struct{}) {
		return twoId, struct{}{}
	})
	twoMap[twoId] = struct{}{}
	for _, v := range twoList {
		if _, ok := twoMap[v.TwoIndicatorId]; ok {
			if v.Score > curMaxScore {
				curMaxScore = v.Score
			}
			if twoId != v.TwoIndicatorId && v.Score > preMaxScore {
				preMaxScore = v.Score
			}
		}
	}
	return
}

// 由于事务的并发性，且无法通过数据唯一索引来控制数据的唯一准确性，对于同一个UserActivityIndicatorId需要加锁来控制，
// 否则（比如两个校级角色同时提交就会出现数据问题，因此需要枷锁来控制，尽管出现的可能性很小）
// 单机服务可以直接使用本地锁，多台服务则可以使用分布式锁
func (s *ReviewDeclareService) CommitReview(params *req.CommitReviewReq) (err error) {
	lockSet := global.FuncMapLock.GetFunMap("CommitReview")
	key := fmt.Sprintf("%d", params.UserActivityIndicatorId)
	ok := lockSet.SetKey(key)
	if ok {
		defer lockSet.DelKey(key)
	} else {
		return errorz.CodeMsg(errorz.RESP_ERR, "当前有人正在进行审核提交，请稍后再试")
	}
	nowTime := util.NowTime()
	userInfo := global.GetUserInfo(s.appCtx)
	item := dao.JudgesVerify{
		UserActivityIndicatorId: params.UserActivityIndicatorId,
		JudgesId:                userInfo.UserId,
		JudgesName:              userInfo.UserName,
		IsPass:                  params.IsPass,
		Opinion:                 params.Opinion,
		CreateTime:              &nowTime,
	}

	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		userActivityIndicatorDao := dao.NewUserActivityIndicatorDaoWithDB(tx, s.appCtx)
		var userActivityIndicator dao.UserActivityIndicator
		err = userActivityIndicatorDao.First(dao.UserActivityIndicator{UserActivityIndicatorId: params.UserActivityIndicatorId}, &userActivityIndicator)
		if err != nil {
			return err
		}
		item.UserActivityId = userActivityIndicator.UserActivityId
		for _, role := range userInfo.UserRoles {
			if userActivityIndicator.ReviewProcess == role {
				item.JudgesType = role
			}
		}
		if userActivityIndicator.ReviewProcess != item.JudgesType {
			return errorz.CodeMsg(errorz.RESP_ERR, "当前进程与角色不匹配，无法进行评审")
		} else if userActivityIndicator.Status == global.ReviewStatusEnd {
			return errorz.CodeMsg(errorz.RESP_ERR, "审核结果已评定，无法进行评审")
		}

		//判断专家是否有权限
		if userActivityIndicator.ReviewProcess == global.ProcessExpert {
			var expertAuthIndicator dao.ExpertAuthIndicator
			err = dao.NewExpertAuthIndicatorDaoWithDB(tx, s.appCtx).
				First(dao.ExpertAuthIndicator{UserId: userInfo.UserId, TwoIndicatorId: userActivityIndicator.TwoIndicatorId}, &expertAuthIndicator)
			if err != nil {
				return errorz.CodeMsg(errorz.RESP_ERR, "当前专家无权限操作，请向管理员申请权限")
			}
		}

		userActivityIndicatorUpdateMap := make(map[string]interface{})
		userActivityIndicatorUpdateMap["finish_review_num"] = *userActivityIndicator.FinishReviewNum + 1
		userActivityIndicatorUpdateMap["update_time"] = &nowTime

		userActivityDao := dao.NewUserActivityDaoWithDB(tx, s.appCtx)
		activityInfo, err := userActivityDao.GetActivityInfo(userActivityIndicator.UserActivityId)
		if err != nil {
			return err
		}
		twoInfo, err := s.GetTwoInfo(userActivityIndicator.UserActivityId, userActivityIndicator.TwoIndicatorId)
		if err != nil {
			return err
		}
		item.Score = twoInfo.Score
		if item.JudgesType == global.RoleExpert {
			if *userActivityIndicator.FinishReviewNum == activityInfo.ReviewNum-2 {
				userActivityIndicatorUpdateMap["review_process"] = userActivityIndicator.ReviewProcess + 1
			}
		}

		if params.IsPass == global.PassNo {
			userActivityIndicatorUpdateMap["status"] = item.JudgesType + 2
			if item.JudgesType == global.RoleEdb {
				userActivityIndicatorUpdateMap["review_process"] = userActivityIndicator.ReviewProcess + 1
			}
		} else {
			if item.JudgesType != global.RoleExpert {
				userActivityIndicatorUpdateMap["review_process"] = userActivityIndicator.ReviewProcess + 1
			}
			if item.JudgesType == global.RoleEdb {
				userActivityIndicatorUpdateMap["status"] = global.ReviewStatusPass
				userActivityDao := dao.NewUserActivityDaoWithDB(tx, s.appCtx)
				userActivityUpdateMap := make(map[string]interface{})
				userActivityUpdateMap["update_time"] = &nowTime
				curMaxScore, preMaxScore, err := s.GetMaxScore(twoInfo.ActivityId, twoInfo.OneIndicatorId, twoInfo.TwoIndicatorId, userActivityIndicator.UserActivityId)
				if err != nil {
					return err
				}
				userActivityUpdateMap["final_score"] = gorm.Expr("final_score + ?", curMaxScore-preMaxScore)
				err = userActivityDao.UpdateByWhere(dao.UserActivity{UserActivityId: userActivityIndicator.UserActivityId}, userActivityUpdateMap)
				if err != nil {
					return err
				}
			}
		}
		judgesVerifyDao := dao.NewJudgesVerifyDaoWithDB(tx, s.appCtx)
		err = judgesVerifyDao.Create(&item)
		if err != nil {
			return err
		}
		err = userActivityIndicatorDao.UpdateByWhere(dao.UserActivityIndicator{UserActivityIndicatorId: userActivityIndicator.UserActivityIndicatorId}, userActivityIndicatorUpdateMap)
		return err
	})
	//操作记录
	if err == nil {
		go global.RecordNotNilError(s.appCtx, NewOtherService(s.appCtx).OperationUaiPassRecord(params.UserActivityIndicatorId, params.IsPass, item.JudgesType))
	}
	return
}

func (s *ReviewDeclareService) GetHistoryActivityList(params *req.GetHistoryActivityListReq) (data *resp.GetHistoryActivityListResp, err error) {
	data = &resp.GetHistoryActivityListResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	data.Total, data.List, err = userActivityDao.GetHistoryActivityList(params)
	return
}

func (s *ReviewDeclareService) GetAwardsSetList(params *req.GetAwardsSetListReq) (data *resp.GetAwardsSetListResp, err error) {
	data = &resp.GetAwardsSetListResp{}
	where := dao.UserActivity{
		ActivityId:   params.ActivityId,
		UserSex:      params.UserSex,
		SubjectCode:  params.SubjectCode,
		SchoolId:     params.SchoolId,
		RankPrize:    params.RankPrize,
		DeclareType:  params.DeclareType,
		IdentityCard: params.IdentityCard,
		FinalScore:   params.FinalScore,
	}
	db := global.GormDB.Where(where)
	if params.UserName != "" {
		db = db.Where("user_name like ?", "%"+params.UserName+"%")
	}
	if params.SchoolName != "" {
		db = db.Where("school_name like ?", "%"+params.SchoolName+"%")
	}
	wg := util.NewWaitGroup(0)
	wg.Go(func() (err error) {
		data.Total, err = dao.NewUserActivityDaoWithDB(db, s.appCtx).
			FindAndCountWithPageOrder(nil, &data.List, params.Page, params.Limit, "CASE WHEN rank_prize > 0 THEN rank_prize ELSE 100 END ASC,final_score DESC")
		return
	})
	wg.Go(func() (err error) {
		var list []dao.UserActivity
		err = dao.NewUserActivityDaoWithDB(db.Where("rank_prize > 0"), s.appCtx).
			Find(nil, &list, "rank_prize")
		for _, v := range list {
			data.AwardNum++
			switch *v.RankPrize {
			case global.PrizeFirst:
				data.FirstPrizeNum++
			case global.PrizeSecond:
				data.SecondPrizeNum++
			case global.PrizeThird:
				data.ThirdPrizeNum++
			}
		}
		return
	})
	wg.Go(func() (err error) {
		activityDao := dao.NewActivityDao(s.appCtx)
		var activity dao.Activity
		err = activityDao.First(dao.Activity{ActivityId: params.ActivityId}, &activity)
		data.ActivityId = activity.ActivityId
		data.ActivityName = activity.ActivityName
		return err
	})
	err = wg.Wait()
	return
}

func (s *ReviewDeclareService) SetAwards(params *req.SetAwardsReq) (err error) {
	nowTime := util.NowTime()
	activityService := NewActivityService(s.appCtx)
	err = activityService.activityIsEnd(params.ActivityId)
	if err != nil {
		return
	}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	err = userActivityDao.UpdateByWhere(dao.UserActivity{ActivityId: params.ActivityId, UserId: params.UserId},
		dao.UserActivity{Prize: params.Prize, RankPrize: params.RankPrize, UpdateTime: &nowTime})
	return
}

func (s *ReviewDeclareService) CommitActivityResult(params *req.CommitActivityResultReq) (err error) {
	nowTime := util.NowTime()
	activityService := NewActivityService(s.appCtx)
	err = activityService.activityIsEnd(params.ActivityId)
	if err != nil {
		return
	}
	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		var userActivityIds []int64
		err := dao.NewUserActivityDaoWithDB(tx.Order("final_score desc"), s.appCtx).
			Pluck(dao.UserActivity{ActivityId: params.ActivityId, DeclareType: params.DeclareType}, &userActivityIds, "user_activity_id")
		if err != nil {
			return err
		}
		userActivityDao := dao.NewUserActivityDaoWithDB(tx, s.appCtx)
		for i, id := range userActivityIds {
			err := userActivityDao.UpdateByWhere(dao.UserActivity{UserActivityId: id}, dao.UserActivity{Rank: i + 1, UpdateTime: &nowTime})
			if err != nil {
				return err
			}
		}
		err = dao.NewUserActivityIndicatorDaoWithDB(tx.Where("user_activity_id in (?)", userActivityIds).
			Where("review_process < ? and (status = ? or status = ?)", global.ProcessEnd, global.ReviewStatusCommit, global.ReviewStatusExpertNoPass), s.appCtx).
			UpdateByWhere(nil, dao.UserActivityIndicator{Status: global.ReviewStatusEnd, UpdateTime: &nowTime})
		return err
	})
	return
}

func (s *ReviewDeclareService) UpdateTwoIndicatorId(params *req.UpdateTwoIndicatorIdReq) (err error) {
	nowTime := util.NowTime()
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	where := dao.UserActivityIndicator{UserActivityIndicatorId: params.UserActivityIndicatorId}
	var oldInfo dao.UserActivityIndicator
	err = userActivityIndicatorDao.First(where, &oldInfo)
	if err != nil {
		return
	}
	err = userActivityIndicatorDao.UpdateByWhere(where, dao.UserActivityIndicator{
		TwoIndicatorId:       params.TwoIndicatorId,
		UpdateTime:           &nowTime,
		AwardDate:            params.AwardDate,
		CertificateType:      params.CertificateType,
		CertificateUrl:       params.CertificateUrl,
		CertificateStartDate: params.CertificateStartDate,
		CertificateEndDate:   params.CertificateEndDate,
	})
	//操作记录
	if err == nil {
		go global.RecordNotNilError(s.appCtx, NewOtherService(s.appCtx).UpdateUaiRecord(oldInfo.UserActivityIndicatorId, oldInfo.UserActivityId, oldInfo.TwoIndicatorId, params.TwoIndicatorId))
	}
	return
}

func (s *ReviewDeclareService) EdbDeclareToUser(params *req.EdbDeclareToUserReq) (err error) {
	nowTime := util.NowTime()
	var userActivityIndicatorId int64
	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		userActivityDao := dao.NewUserActivityDaoWithDB(tx, s.appCtx)
		var userActivity dao.UserActivity
		err := userActivityDao.First(dao.UserActivity{UserActivityId: params.UserActivityId}, &userActivity)
		if err != nil {
			return err
		}
		if userActivity.Rank != 0 {
			return errorz.CodeMsg(errorz.RESP_ERR, "活动结果已评定，无法进行该操作")
		}
		activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDaoWithDB(tx, s.appCtx)
		var twoInfo dao.ActivityTwoIndicator
		err = activityTwoIndicatorDao.First(dao.ActivityTwoIndicator{ActivityId: userActivity.ActivityId, TwoIndicatorId: params.TwoIndicatorId}, &twoInfo)
		if err != nil {
			return err
		}
		activityDao := dao.NewActivityDaoWithDB(tx, s.appCtx)
		var activityInfo dao.Activity
		err = activityDao.First(dao.Activity{ActivityId: userActivity.ActivityId}, &activityInfo)
		if err != nil {
			return err
		}
		item := dao.UserActivityIndicator{
			UserActivityId:       params.UserActivityId,
			TwoIndicatorId:       params.TwoIndicatorId,
			AwardDate:            params.AwardDate,
			CertificateType:      params.CertificateType,
			CertificateUrl:       params.CertificateUrl,
			CertificateStartDate: params.CertificateStartDate,
			CertificateEndDate:   params.CertificateEndDate,
			Status:               global.ReviewStatusPass,
			FinishReviewNum:      &activityInfo.ReviewNum,
			ReviewProcess:        global.ProcessEnd,
			CreateTime:           &nowTime,
		}
		userActivityIndicatorDao := dao.NewUserActivityIndicatorDaoWithDB(tx, s.appCtx)
		err = userActivityIndicatorDao.Create(&item)
		if err != nil {
			return err
		}
		userActivityIndicatorId = item.UserActivityIndicatorId
		judgesVerifyDao := dao.NewJudgesVerifyDaoWithDB(tx, s.appCtx)
		list := make([]dao.JudgesVerify, activityInfo.ReviewNum)
		for i := 0; i < activityInfo.ReviewNum; i++ {
			judgesVerify := dao.JudgesVerify{
				UserActivityIndicatorId: item.UserActivityIndicatorId,
				UserActivityId:          params.UserActivityId,
				JudgesId:                "edb" + strconv.Itoa(i+1),
				JudgesName:              "教育局",
				IsPass:                  global.PassYes,
				Score:                   twoInfo.Score,
				Opinion:                 "同意通过",
				CreateTime:              &nowTime,
			}
			if i == 0 {
				judgesVerify.JudgesType = global.RoleSchool
			} else if i == activityInfo.ReviewNum-1 {
				judgesVerify.JudgesType = global.RoleEdb
			} else {
				judgesVerify.JudgesType = global.RoleExpert
			}
			list[i] = judgesVerify
		}
		err = judgesVerifyDao.BatchInsert(list)
		if err != nil {
			return err
		}
		updateMap := make(map[string]interface{})
		updateMap["final_score"] = gorm.Expr("final_score + ?", twoInfo.Score)
		updateMap["update_time"] = &nowTime
		err = userActivityDao.UpdateByWhere(dao.UserActivity{UserActivityId: params.UserActivityId}, updateMap)
		return err
	})
	//操作记录
	if err == nil {
		go global.RecordNotNilError(s.appCtx, NewOtherService(s.appCtx).OperationUaiRecord(userActivityIndicatorId, 2, 3))
	}
	return
}

func (s *ReviewDeclareService) BatchPass(params *req.BatchPassReq) (err error) {
	wg := util.NewWaitGroup(10)
	for _, id := range params.UserActivityIndicatorIds {
		wg.Add()
		go func(curId int64) {
			defer wg.Done()
			s.CommitReview(&req.CommitReviewReq{UserActivityIndicatorId: curId, IsPass: 1, Opinion: "同意通过"})
		}(id)
	}
	wg.Wait()
	return
}

func (s *ReviewDeclareService) GetEdbReviewList(params *req.GetEdbReviewListReq) (data *resp.GetEdbReviewListResp, err error) {
	data = &resp.GetEdbReviewListResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	data.Total, data.List, err = userActivityDao.GetEdbReviewList(params)
	return
}
