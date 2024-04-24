package services

import (
	"context"
	"gorm.io/gorm"
	"sort"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/dao"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/model/resp"
)

type UserActivityService struct {
	appCtx context.Context
}

func NewUserActivityService(appCtx context.Context) *UserActivityService {
	return &UserActivityService{appCtx: appCtx}
}

func (s *UserActivityService) GetActivityListToUser(params *req.GetActivityListToUserReq) (data *resp.GetActivityListToUserResp, err error) {
	data = &resp.GetActivityListToUserResp{}
	activityDao := dao.NewActivityDao(s.appCtx)
	data.Total, err = activityDao.FindAndCountWithPageOrder(nil, &data.List, params.Page, params.Limit, "create_time desc", "activity_id", "activity_name", "year", "start_time", "end_time")
	activityIds := util.ListObjToListObj(data.List, func(obj resp.ActivityListToUser) int {
		return obj.ActivityId
	})
	userActivityDao := dao.NewUserActivityDaoWithDB(global.GormDB.Where("activity_id in (?) and user_id = ?", activityIds, params.UserId), s.appCtx)
	var userActivityList []dao.UserActivity
	err = userActivityDao.Find(nil, &userActivityList, "activity_id")
	if err != nil {
		return
	}
	objToMap := util.ListObjToMap(userActivityList, func(obj dao.UserActivity) (int, interface{}) {
		return obj.ActivityId, struct{}{}
	})
	nowTime := util.NowTime()
	for i, v := range data.List {
		if nowTime.Before(*v.StartTime) {
			data.List[i].Status = global.ActivityWaitStart
		} else if nowTime.After(*v.EndTime) {
			data.List[i].Status = global.ActivityEnd
		} else {
			data.List[i].Status = global.ActivityOngoing
		}
		if _, ok := objToMap[v.ActivityId]; ok {
			data.List[i].IsDeclare = true
		}
	}
	return
}

func (s *UserActivityService) GetActivityInfoToUser(activityId int) (data *resp.GetActivityInfoToUserResp, err error) {
	data = &resp.GetActivityInfoToUserResp{}
	activityDao := dao.NewActivityDao(s.appCtx)
	err = activityDao.First(dao.Activity{ActivityId: activityId}, &data.Activity)
	return
}

func (s *UserActivityService) CreateUserActivityDeclare(params *req.CreateUserActivityDeclareReq) (err error) {
	activityService := NewActivityService(s.appCtx)
	err = activityService.activityIsOngoing(params.ActivityId)
	if err != nil {
		return err
	}
	nowTime := util.NowTime()
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	userActivity := userActivityDao.GetOnlyUserActivity(params.ActivityId, params.UserId)
	util.ObjToObjByReflect(&params.UserActivity, &userActivity)
	if userActivity.UserActivityId != 0 {
		userActivity.UpdateTime = &nowTime
		err = userActivityDao.UpdateByWhere(dao.UserActivity{UserActivityId: userActivity.UserActivityId}, userActivity)
	} else {
		userActivity.CreateTime = &nowTime
		RankPrize := 0
		userActivity.RankPrize = &RankPrize
		err = userActivityDao.Create(&userActivity)
	}
	if err != nil {
		return
	}

	var list []dao.UserActivityIndicator
	for _, v := range params.List {
		util.TimeToLocation(v.AwardDate)
		util.TimeToLocation(v.CertificateStartDate)
		util.TimeToLocation(v.CertificateEndDate)
		finishReviewNum := 0
		list = append(list, dao.UserActivityIndicator{
			UserActivityId:       userActivity.UserActivityId,
			TwoIndicatorId:       v.TwoIndicatorId,
			AwardDate:            v.AwardDate,
			CertificateType:      v.CertificateType,
			CertificateUrl:       v.CertificateUrl,
			CertificateStartDate: v.CertificateStartDate,
			CertificateEndDate:   v.CertificateEndDate,
			Status:               global.ReviewStatusCommit,
			ReviewProcess:        global.ProcessSchool,
			FinishReviewNum:      &finishReviewNum,
			CreateTime:           &nowTime,
		})
	}
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	err = userActivityIndicatorDao.BatchInsert(list)
	//操作记录
	if err == nil {
		go func() {
			otherService := NewOtherService(s.appCtx)
			for _, v := range list {
				global.RecordNotNilError(otherService.OperationUaiRecord(v.UserActivityIndicatorId, 1, 4))
			}
		}()
	}
	return
}

func (s *UserActivityService) GetUserActivityDeclareDetail(activityId int, userId string) (data *resp.GetUserActivityDeclareDetailResp, err error) {
	data = &resp.GetUserActivityDeclareDetailResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	userActivity := userActivityDao.GetOnlyUserActivity(activityId, userId)
	if userActivity.UserActivityId == 0 {
		return
	}
	util.ObjToObjByReflect(&userActivity, &data.UserActivity)
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	var list []dao.UserActivityIndicator
	err = userActivityIndicatorDao.Find(dao.UserActivityIndicator{UserActivityId: userActivity.UserActivityId}, &list)
	if err != nil {
		return
	}
	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	twoToOneMap, err := activityTwoIndicatorDao.GetTwoToOneMap(activityId)
	if err != nil {
		return
	}
	data.List = make([]resp.UserActivityIndicator, len(list))
	for i, v := range list {
		data.List[i].UserActivityIndicator = v
		data.List[i].OneIndicatorId = twoToOneMap[v.TwoIndicatorId]
	}
	return
}

func (s *UserActivityService) GetUserActivityDeclareStatusList(activityId int, userId string) (data *resp.GetUserActivityDeclareStatusListResp, err error) {
	data = &resp.GetUserActivityDeclareStatusListResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	userActivity := userActivityDao.GetOnlyUserActivity(activityId, userId)
	if userActivity.UserActivityId == 0 {
		return
	}
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	var userActivityIndicatorList []dao.UserActivityIndicator
	err = userActivityIndicatorDao.Find(dao.UserActivityIndicator{UserActivityId: userActivity.UserActivityId}, &userActivityIndicatorList, "user_activity_indicator_id", "two_indicator_id", "status", "create_time")
	if err != nil {
		return
	}

	//二级指标排列顺序：优先显示<通过>、<未通过>、<已提交>
	//
	//其次按申报时间顺序排列
	sort.Slice(userActivityIndicatorList, func(i, j int) bool {
		if userActivityIndicatorList[i].Status != global.ReviewStatusPass && userActivityIndicatorList[j].Status == global.ReviewStatusPass {
			return false
		} else if userActivityIndicatorList[i].Status == global.ReviewStatusCommit && userActivityIndicatorList[j].Status != global.ReviewStatusCommit {
			return false
		} else if userActivityIndicatorList[i].Status != global.ReviewStatusCommit && userActivityIndicatorList[j].Status == global.ReviewStatusCommit {
			return false
		} else if userActivityIndicatorList[i].CreateTime.After(*userActivityIndicatorList[j].CreateTime) {
			return false
		}
		return true
	})

	activityOneIndicatorDao := dao.NewActivityOneIndicatorDao(s.appCtx)
	activityOneIndicatorMap, err := activityOneIndicatorDao.GetActivityOneIndicatorMap(activityId)
	if err != nil {
		return
	}
	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	activityTwoIndicatorMap, err := activityTwoIndicatorDao.GetActivityTwoIndicatorMap(activityId, util.ListObjToListObj(userActivityIndicatorList, func(obj dao.UserActivityIndicator) int {
		return obj.TwoIndicatorId
	}))
	if err != nil {
		return
	}
	oneIndexMap := make(map[int]int)
	for _, v := range userActivityIndicatorList {
		data.DeclareNum++
		switch v.Status {
		case global.ReviewStatusCommit:
			data.WaitReview++
		case global.ReviewStatusPass:
			data.PassNum++
		default:
			data.NoPassNum++
		}
		two, ok1 := activityTwoIndicatorMap[v.TwoIndicatorId]
		one, ok2 := activityOneIndicatorMap[two.OneIndicatorId]
		if !ok1 || !ok2 {
			return nil, errorz.CodeMsg(errorz.RESP_ERR, "活动指标数据不存在")
		}
		index, ok := oneIndexMap[one.OneIndicatorId]
		if !ok {
			data.List = append(data.List, resp.OneIndicatorDeclareStatusList{
				OneIndicatorId:   one.OneIndicatorId,
				OneIndicatorName: one.OneIndicatorName,
				Content:          one.Content,
			})
			index = len(data.List) - 1
			oneIndexMap[one.OneIndicatorId] = index
		}
		data.List[index].DeclareStatusList = append(data.List[index].DeclareStatusList, resp.DeclareStatus{
			UserActivityIndicatorId: v.UserActivityIndicatorId,
			TwoIndicatorId:          two.TwoIndicatorId,
			TwoIndicatorName:        two.TwoIndicatorName,
			Score:                   two.Score,
			Status:                  v.Status,
			DeclareTime:             v.CreateTime,
		})
	}
	//一级指标按指标库顺序排列
	sort.Slice(data.List, func(i, j int) bool {
		return data.List[i].OneIndicatorId < data.List[j].OneIndicatorId
	})
	return
}

func (s *UserActivityService) GetUserActivityIndicator(userActivityIndicatorId int64) (data *resp.GetUserActivityIndicatorResp, err error) {
	data = &resp.GetUserActivityIndicatorResp{}
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	var userActivityIndicator dao.UserActivityIndicator
	err = userActivityIndicatorDao.First(dao.UserActivityIndicator{UserActivityIndicatorId: userActivityIndicatorId}, &userActivityIndicator)
	if err != nil {
		return
	}
	data.UserActivityIndicator = userActivityIndicator

	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	var userActivity dao.UserActivity
	err = userActivityDao.First(dao.UserActivityIndicator{UserActivityId: userActivityIndicator.UserActivityId}, &userActivity)
	if err != nil {
		return
	}
	data.UserId = userActivity.UserId
	data.UserName = userActivity.UserName
	activityService := NewActivityService(s.appCtx)
	data.ActivityIndicatorInfo, err = activityService.GetActivityIndicatorByTwoId(userActivity.ActivityId, userActivityIndicator.TwoIndicatorId)
	if err != nil {
		return
	}
	judgesVerifyDao := dao.NewJudgesVerifyDao(s.appCtx)
	err = judgesVerifyDao.Find(dao.JudgesVerify{UserActivityIndicatorId: userActivityIndicatorId}, &data.List)
	return
}

func (s *UserActivityService) userActivityIndicatorIsPass(userActivityIndicatorId int64) (err error) {
	var curItem dao.UserActivityIndicator
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	err = userActivityIndicatorDao.First(dao.UserActivityIndicator{UserActivityIndicatorId: userActivityIndicatorId}, &curItem)
	if err != nil {
		return err
	}
	if curItem.Status == global.ReviewStatusPass {
		return errorz.CodeMsg(errorz.RESP_ERR, "当前项目已经通过，无法进行该操作")
	}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	var userActivity dao.UserActivity
	err = userActivityDao.First(dao.UserActivity{UserActivityId: curItem.UserActivityId}, &userActivity)
	if err != nil {
		return
	}
	activityService := NewActivityService(s.appCtx)
	err = activityService.activityIsOngoing(userActivity.ActivityId)
	return
}

func (s *UserActivityService) UpdateUserActivityIndicator(params *req.UpdateUserActivityIndicatorReq) (err error) {
	err = s.userActivityIndicatorIsPass(params.UserActivityIndicatorId)
	if err != nil {
		return
	}
	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		nowTime := util.NowTime()
		userActivityIndicatorDao := dao.NewUserActivityIndicatorDaoWithDB(tx, s.appCtx)
		var newItem dao.UserActivityIndicator
		util.ObjToObjByReflect(&params.UserActivityIndicator, &newItem)
		newItem.Status = global.ReviewStatusCommit
		newItem.ReviewProcess = global.ProcessSchool
		newItem.UpdateTime = &nowTime
		finishReviewNum := 0
		newItem.FinishReviewNum = &finishReviewNum
		err = userActivityIndicatorDao.UpdateByWhere(dao.UserActivityIndicator{UserActivityIndicatorId: params.UserActivityIndicatorId}, newItem)
		if err != nil {
			return err
		}
		judgesVerifyDao := dao.NewJudgesVerifyDaoWithDB(tx, s.appCtx)
		err = judgesVerifyDao.DeleteByWhere(dao.JudgesVerify{UserActivityIndicatorId: params.UserActivityIndicatorId})
		return err
	})
	//操作记录
	if err == nil {
		go global.RecordNotNilError(NewOtherService(s.appCtx).OperationUaiRecord(params.UserActivityIndicatorId, 2, 4))
	}
	return
}

func (s *UserActivityService) DeleteUserActivityIndicator(userActivityIndicatorId int64) (err error) {
	err = s.userActivityIndicatorIsPass(userActivityIndicatorId)
	if err != nil {
		return
	}
	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		userActivityIndicatorDao := dao.NewUserActivityIndicatorDaoWithDB(tx, s.appCtx)
		err := userActivityIndicatorDao.DeleteByWhere(dao.UserActivityIndicator{UserActivityIndicatorId: userActivityIndicatorId})
		if err != nil {
			return err
		}
		judgesVerifyDao := dao.NewJudgesVerifyDaoWithDB(tx, s.appCtx)
		err = judgesVerifyDao.DeleteByWhere(dao.JudgesVerify{UserActivityIndicatorId: userActivityIndicatorId})
		return err
	})
	//操作记录
	if err == nil {
		go global.RecordNotNilError(NewOtherService(s.appCtx).OperationUaiRecord(userActivityIndicatorId, 3, 4))
	}
	return
}

func (s *UserActivityService) GetUserActivityDeclareResult(activityId int, userId string) (data *resp.GetUserActivityDeclareResultResp, err error) {
	data = &resp.GetUserActivityDeclareResultResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	var userActivity dao.UserActivity
	err = userActivityDao.First(dao.UserActivity{ActivityId: activityId, UserId: userId}, &userActivity)
	data.Prize = userActivity.Prize
	data.RankPrize = *userActivity.RankPrize
	data.FinalScore = userActivity.FinalScore

	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	var userActivityIndicatorList []dao.UserActivityIndicator
	err = userActivityIndicatorDao.Find(dao.UserActivityIndicator{UserActivityId: userActivity.UserActivityId, Status: 2}, &userActivityIndicatorList, "two_indicator_id")
	if err != nil {
		return
	}
	var twoIds []int
	for _, v := range userActivityIndicatorList {
		twoIds = append(twoIds, v.TwoIndicatorId)
		data.DeclareNum++
	}
	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	twoIndicatorMap, err := activityTwoIndicatorDao.GetActivityTwoIndicatorMap(activityId, twoIds)
	if err != nil {
		return
	}
	activityOneIndicatorDao := dao.NewActivityOneIndicatorDao(s.appCtx)
	activityOneIndicatorMap, err := activityOneIndicatorDao.GetActivityOneIndicatorMap(activityId)
	if err != nil {
		return
	}
	for _, v := range twoIndicatorMap {
		one := activityOneIndicatorMap[v.OneIndicatorId]
		data.List = append(data.List, resp.ActivityIndicatorInfo{
			OneIndicatorId:   one.OneIndicatorId,
			OneIndicatorName: one.OneIndicatorName,
			TwoIndicatorId:   v.TwoIndicatorId,
			TwoIndicatorName: v.TwoIndicatorName,
			Score:            v.Score,
		})
	}
	sort.Slice(data.List, func(i, j int) bool {
		return data.List[i].OneIndicatorId < data.List[j].OneIndicatorId || (data.List[i].OneIndicatorId == data.List[j].OneIndicatorId && data.List[i].TwoIndicatorId < data.List[j].TwoIndicatorId)
	})
	return
}

func (s *UserActivityService) GetUserDeclareRecordListByYear(params *req.GetUserDeclareRecordListByYearReq) (data *resp.GetUserDeclareRecordListByYearResp, err error) {
	data = &resp.GetUserDeclareRecordListByYearResp{}
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	total, list, err := userActivityIndicatorDao.GetUserDeclareRecordList(params.DeclareUserId, params.Year, params.Page, params.Limit)
	if err != nil {
		return
	}
	data.Total = total
	activityIds := util.ListObjToListObj(list, func(obj dao.UserDeclareRecord) int {
		return obj.ActivityId
	})
	twoIds := util.ListObjToListObj(list, func(obj dao.UserDeclareRecord) int {
		return obj.TwoIndicatorId
	})
	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	twoIndicatorsMap, err := activityTwoIndicatorDao.GetTwoIndicatorInfoMap(activityIds, twoIds)
	if err != nil {
		return
	}
	data.List = make([]resp.UserDeclareRecord, len(list))
	for i, v := range list {
		data.List[i].TwoIndicatorInfo = twoIndicatorsMap[v.ActivityId][v.TwoIndicatorId]
		data.List[i].UserActivityIndicatorId = v.UserActivityIndicatorId
		data.List[i].UserActivityId = v.UserActivityId
		data.List[i].CreateTime = v.CreateTime
		if v.Status == global.ReviewStatusExpertNoPass {
			v.Status = global.ReviewStatusCommit
		}
		data.List[i].Status = v.Status
	}
	//教育局
	if params.Role == global.RoleEdb {
		ids := util.ListObjToListObj(list, func(obj dao.UserDeclareRecord) int64 {
			return obj.UserActivityIndicatorId
		})
		verifyList, err := dao.NewJudgesVerifyDao(s.appCtx).GetJudgesVerifyListByIds(ids)
		if err != nil {
			return nil, err
		}
		judgesVerifyListMap := util.ListObjToMapList(verifyList, func(obj dao.JudgesVerify) (int64, resp.JudgesVerifyPass) {
			return obj.UserActivityIndicatorId, resp.JudgesVerifyPass{
				JudgesId:   obj.JudgesId,
				JudgesName: obj.JudgesName,
				JudgesType: obj.JudgesType,
				IsPass:     obj.IsPass,
			}
		})
		for i, v := range data.List {
			data.List[i].JudgesVerifyList = judgesVerifyListMap[v.UserActivityIndicatorId]
		}
	}
	return
}

func (s *UserActivityService) GetUserHistoryDeclareResultList(params *req.GetUserHistoryDeclareResultListReq) (data *resp.GetUserHistoryDeclareResultListResp, err error) {
	data = &resp.GetUserHistoryDeclareResultListResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	data.Total, data.List, err = userActivityDao.GetUserHistoryDeclareResultList(params.DeclareUserId, params.Page, params.Limit)
	return
}

func (s *UserActivityService) GetUserDeclareResultApp(params *req.GetUserDeclareResultAppReq) (data *resp.GetUserDeclareResultAppResp, err error) {
	data = &resp.GetUserDeclareResultAppResp{}
	activityDao := dao.NewActivityDao(s.appCtx)
	var activity dao.Activity
	err = activityDao.First(dao.Activity{ActivityId: params.ActivityId}, &activity)
	if err != nil {
		return
	}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	var userActivity dao.UserActivity
	err = userActivityDao.First(dao.UserActivity{ActivityId: params.ActivityId, UserId: params.UserId}, &userActivity)
	if err != nil {
		return
	}
	userActivityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	var twoIds []int
	err = userActivityIndicatorDao.Pluck(dao.UserActivityIndicator{UserActivityId: userActivity.UserActivityId}, &twoIds, "two_indicator_id")
	if err != nil {
		return
	}
	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	twoIndicatorMap, err := activityTwoIndicatorDao.GetActivityTwoIndicatorMap(params.ActivityId, twoIds)
	if err != nil {
		return
	}
	for _, v := range twoIndicatorMap {
		data.TeacherScore += v.Score
	}
	judgesVerifyDao := dao.NewJudgesVerifyDao(s.appCtx)
	var list []dao.JudgesVerify
	err = judgesVerifyDao.Find(dao.JudgesVerify{UserActivityId: userActivity.UserActivityId}, &list, "judges_type,is_pass,score")
	if err != nil {
		return
	}
	for _, v := range list {
		switch v.JudgesType {
		case global.RoleSchool:
			if v.IsPass == global.PassYes {
				data.SchoolScore += v.Score
			}
		case global.RoleExpert:
			if v.IsPass == global.PassYes {
				data.ExpertPassScore += v.Score
			} else {
				data.ExpertNoPassScore += v.Score
			}
		case global.RoleEdb:
			if v.IsPass == global.PassYes {
				data.EdbScore += v.Score
			}
		}
	}
	data.ActivityId = activity.ActivityId
	data.ActivityName = activity.ActivityName
	data.FinalScore = userActivity.FinalScore
	data.Rank = userActivity.Rank
	data.RankPrize = *userActivity.RankPrize
	data.Prize = userActivity.Prize
	return
}

func (s *UserActivityService) GetUserDeclaresToEdb(userActivityId int64) (data *resp.GetUserDeclaresToEdbResp, err error) {
	data = &resp.GetUserDeclaresToEdbResp{}
	var userActivity dao.UserActivity
	err = dao.NewUserActivityDao(s.appCtx).First(dao.UserActivity{UserActivityId: userActivityId}, &userActivity)
	if err != nil {
		return
	}
	data.FinalScore = userActivity.FinalScore
	var list []dao.UserActivityIndicator
	err = dao.NewUserActivityIndicatorDaoWithDB(global.GormDB.Where("review_process >= ?", 3), s.appCtx).
		Find(dao.UserActivityIndicator{UserActivityId: userActivityId}, &list, "user_activity_indicator_id,two_indicator_id,certificate_url,review_process")
	if err != nil {
		return
	}
	ids := util.ListObjToListObj(list, func(obj dao.UserActivityIndicator) int64 {
		return obj.UserActivityIndicatorId
	})
	verifyList, err := dao.NewJudgesVerifyDao(s.appCtx).GetJudgesVerifyListByIds(ids)
	if err != nil {
		return
	}
	judgesVerifyListMap := util.ListObjToMapList(verifyList, func(obj dao.JudgesVerify) (int64, resp.JudgesVerifyPass) {
		return obj.UserActivityIndicatorId, resp.JudgesVerifyPass{
			JudgesId:   obj.JudgesId,
			JudgesName: obj.JudgesName,
			JudgesType: obj.JudgesType,
			IsPass:     obj.IsPass,
		}
	})
	twoIds := util.ListObjToListObj(list, func(obj dao.UserActivityIndicator) int {
		return obj.TwoIndicatorId
	})

	twoIndicatorMap, err := dao.NewActivityTwoIndicatorDao(s.appCtx).GetActivityTwoIndicatorMap(userActivity.ActivityId, twoIds)
	if err != nil {
		return
	}
	oneIndicatorMap, err := dao.NewActivityOneIndicatorDao(s.appCtx).GetActivityOneIndicatorMap(userActivity.ActivityId)
	if err != nil {
		return
	}
	for _, v := range list {
		two := twoIndicatorMap[v.TwoIndicatorId]
		one := oneIndicatorMap[two.OneIndicatorId]
		info := dao.TwoIndicatorInfo{
			TwoIndicatorId:   two.TwoIndicatorId,
			TwoIndicatorName: two.TwoIndicatorName,
			Score:            two.Score,
			OneIndicatorId:   one.OneIndicatorId,
			OneIndicatorName: one.OneIndicatorName,
		}
		data.List = append(data.List, resp.UserDeclaresToEdb{
			TwoIndicatorInfo:        info,
			UserActivityIndicatorId: v.UserActivityIndicatorId,
			CertificateUrl:          v.CertificateUrl,
			JudgesVerifyList:        judgesVerifyListMap[v.UserActivityIndicatorId],
			ReviewProcess:           v.ReviewProcess,
		})
	}
	sort.Slice(data.List, func(i, j int) bool {
		return data.List[i].ReviewProcess < data.List[j].ReviewProcess ||
			(data.List[i].ReviewProcess == data.List[j].ReviewProcess && data.List[i].OneIndicatorId < data.List[j].OneIndicatorId)
	})
	oneList := util.MapToList(oneIndicatorMap, func(k int, v dao.ActivityOneIndicator) dao.ActivityOneIndicator {
		return v
	})
	data.OneIndicatorList = oneList
	return
}
