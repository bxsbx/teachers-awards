package services

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/dao"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/model/resp"
)

type ActivityService struct {
	appCtx context.Context
}

func NewActivityService(appCtx context.Context) *ActivityService {
	return &ActivityService{appCtx: appCtx}
}

func (s *ActivityService) activityIsStart(activityId int) error {
	nowTime := util.NowTime()
	activityDao := dao.NewActivityDao(s.appCtx)
	activity := dao.Activity{}
	err := activityDao.First(dao.Activity{ActivityId: activityId}, &activity)
	if err != nil {
		return err
	}
	if nowTime.After(*activity.EndTime) {
		return errorz.CodeMsg(errorz.RESP_ERR, "活动已结束，不允许该操作")
	} else if !nowTime.Before(*activity.StartTime) {
		return errorz.CodeMsg(errorz.RESP_ERR, "活动在进行中，不允许该操作")
	}
	return nil
}

// 活动是否在进行中
func (s *ActivityService) activityIsOngoing(activityId int) error {
	nowTime := util.NowTime()
	activityDao := dao.NewActivityDao(s.appCtx)
	activity := dao.Activity{}
	err := activityDao.First(dao.Activity{ActivityId: activityId}, &activity)
	if err != nil {
		return err
	}
	if nowTime.After(*activity.EndTime) {
		return errorz.CodeMsg(errorz.RESP_ERR, "活动已结束，不允许该操作")
	} else if nowTime.Before(*activity.StartTime) {
		return errorz.CodeMsg(errorz.RESP_ERR, "活动未开始，不允许该操作")
	}
	return nil
}

func (s *ActivityService) activityIsEnd(activityId int) error {
	nowTime := util.NowTime()
	activityDao := dao.NewActivityDao(s.appCtx)
	activity := dao.Activity{}
	err := activityDao.First(dao.Activity{ActivityId: activityId}, &activity)
	if err != nil {
		return err
	}
	if !nowTime.After(*activity.EndTime) {
		return errorz.CodeMsg(errorz.RESP_ERR, "活动还未结束，不允许该操作")
	}
	return nil
}

func (s *ActivityService) CreateOrUpdateActivity(params *req.CreateOrUpdateActivityReq, twoIds []int) (err error) {
	nowTime := util.NowTime()
	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		activityDao := dao.NewActivityDaoWithDB(tx, s.appCtx)
		newActivity := dao.Activity{
			ActivityName: params.ActivityName,
			Year:         params.Year,
			Description:  params.Description,
			Url:          params.Url,
			ReviewNum:    4,
			StartTime:    params.StartTime,
			EndTime:      params.EndTime,
		}
		twoIndicatorDao := dao.NewTwoIndicatorDaoWithDB(tx.Where("two_indicator_id in (?)", twoIds), s.appCtx)
		var twoIndicatorList []dao.TwoIndicator
		err := twoIndicatorDao.Find(nil, &twoIndicatorList)
		if err != nil {
			return err
		}
		oneIds := util.ListToDeduplicationList(twoIndicatorList, func(t dao.TwoIndicator) (int, int) {
			return t.OneIndicatorId, t.OneIndicatorId
		})
		oneIndicatorDao := dao.NewOneIndicatorDaoWithDB(tx.Where("one_indicator_id in (?)", oneIds), s.appCtx)
		var oneIndicatorList []dao.OneIndicator
		err = oneIndicatorDao.Find(nil, &oneIndicatorList)
		if err != nil {
			return err
		}
		if params.ActivityId != 0 {
			err := s.activityIsStart(params.ActivityId)
			if err != nil {
				return err
			}
			newActivity.UpdateTime = &nowTime
			err = activityDao.UpdateByWhere(dao.Activity{ActivityId: params.ActivityId}, newActivity)
			if err != nil {
				return err
			}

			activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDaoWithDB(tx, s.appCtx)
			var activityTwoIndicatorList []dao.ActivityTwoIndicator
			err = activityTwoIndicatorDao.Find(dao.ActivityTwoIndicator{ActivityId: params.ActivityId}, &activityTwoIndicatorList, "two_indicator_id,one_indicator_id")
			if err != nil {
				return err
			}
			existsOneIds := util.ListToDeduplicationList(activityTwoIndicatorList, func(t dao.ActivityTwoIndicator) (int, int) {
				return t.OneIndicatorId, t.OneIndicatorId
			})

			existsTwoIds := util.ListObjToListObj(activityTwoIndicatorList, func(t dao.ActivityTwoIndicator) int {
				return t.TwoIndicatorId
			})
			setTwoIds1, _, setTwoIds3 := util.GetThreeSetFromList(twoIds, existsTwoIds, func(t int) (int, int) {
				return t, t
			})
			twoIds = setTwoIds1
			setOneIds1, _, setOneIds3 := util.GetThreeSetFromList(oneIds, existsOneIds, func(t int) (int, int) {
				return t, t
			})
			oneIds = setOneIds1
			if len(setTwoIds3) > 0 {
				err = dao.NewActivityTwoIndicatorDaoWithDB(tx.Where("two_indicator_id in (?)", setTwoIds3), s.appCtx).
					DeleteByWhere(dao.ActivityTwoIndicator{ActivityId: params.ActivityId})
				if err != nil {
					return err
				}
			}
			if len(setOneIds3) > 0 {
				err = dao.NewActivityOneIndicatorDaoWithDB(tx.Where("one_indicator_id in (?)", setOneIds3), s.appCtx).
					DeleteByWhere(dao.ActivityOneIndicator{ActivityId: params.ActivityId})
				if err != nil {
					return err
				}
			}
			newActivity.ActivityId = params.ActivityId
		} else {
			years, err := activityDao.GetYearsByGroup()
			if err != nil {
				return err
			}
			for _, v := range years {
				if v == newActivity.Year {
					return errorz.CodeMsg(errorz.RESP_ERR, "该年度活动已开启过，暂时不支持开启多个活动")
				}
			}
			newActivity.CreateTime = &nowTime
			err = activityDao.Create(&newActivity)
			if err != nil {
				return err
			}
		}
		twoMap := util.ListObjToMap(twoIndicatorList, func(obj dao.TwoIndicator) (int, dao.TwoIndicator) {
			return obj.TwoIndicatorId, obj
		})
		oneMap := util.ListObjToMap(oneIndicatorList, func(obj dao.OneIndicator) (int, dao.OneIndicator) {
			return obj.OneIndicatorId, obj
		})

		var activityTwoIndicatorList []dao.ActivityTwoIndicator
		for _, id := range twoIds {
			if v, ok := twoMap[id]; ok {
				activityTwoIndicatorList = append(activityTwoIndicatorList, dao.ActivityTwoIndicator{
					ActivityId:       newActivity.ActivityId,
					TwoIndicatorId:   v.TwoIndicatorId,
					TwoIndicatorName: v.TwoIndicatorName,
					Score:            v.Score,
					OneIndicatorId:   v.OneIndicatorId,
					CreateTime:       &nowTime,
				})
			}
		}

		var activityOneIndicatorList []dao.ActivityOneIndicator
		for _, id := range oneIds {
			if v, ok := oneMap[id]; ok {
				activityOneIndicatorList = append(activityOneIndicatorList, dao.ActivityOneIndicator{
					ActivityId:       newActivity.ActivityId,
					OneIndicatorId:   v.OneIndicatorId,
					OneIndicatorName: v.OneIndicatorName,
					Content:          v.Content,
					CreateTime:       &nowTime,
				})
			}
		}
		err = dao.NewActivityOneIndicatorDaoWithDB(tx, s.appCtx).BatchInsert(activityOneIndicatorList)
		if err != nil {
			return err
		}
		err = dao.NewActivityTwoIndicatorDaoWithDB(tx, s.appCtx).BatchInsert(activityTwoIndicatorList)
		return err
	})
	return
}

func (s *ActivityService) GetActivityDetail(activityId int) (data *resp.GetActivityDetailResp, err error) {
	data = &resp.GetActivityDetailResp{}
	activityDao := dao.NewActivityDao(s.appCtx)
	activityOneIndicatorDao := dao.NewActivityOneIndicatorDao(s.appCtx)
	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	var activity dao.Activity
	err = activityDao.First(dao.Activity{ActivityId: activityId}, &activity)
	if err != nil {
		return
	}
	var activityOneIndicatorList []dao.ActivityOneIndicator
	err = activityOneIndicatorDao.Find(dao.ActivityOneIndicator{ActivityId: activityId}, &activityOneIndicatorList)
	if err != nil {
		return
	}
	var activityTwoIndicatorList []dao.ActivityTwoIndicator
	err = activityTwoIndicatorDao.Find(dao.ActivityTwoIndicator{ActivityId: activityId}, &activityTwoIndicatorList)
	if err != nil {
		return
	}
	var oneList []resp.ActivityOneIndicator
	for _, one := range activityOneIndicatorList {
		var twoList []dao.ActivityTwoIndicator
		for _, two := range activityTwoIndicatorList {
			if two.OneIndicatorId == one.OneIndicatorId {
				twoList = append(twoList, two)
			}
		}
		oneList = append(oneList, resp.ActivityOneIndicator{
			ActivityOneIndicator:  one,
			ActivityTwoIndicators: twoList,
		})
	}
	data.Activity = activity
	data.OneIndicators = oneList
	return
}

func (s *ActivityService) DeleteActivityOneIndicator(activityId, oneIndicatorId int) (err error) {
	err = s.activityIsStart(activityId)
	if err != nil {
		return err
	}
	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		activityOneIndicatorDao := dao.NewActivityOneIndicatorDaoWithDB(tx, s.appCtx)
		err := activityOneIndicatorDao.DeleteByWhere(dao.ActivityOneIndicator{ActivityId: activityId, OneIndicatorId: oneIndicatorId})
		if err != nil {
			return err
		}
		activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDaoWithDB(tx, s.appCtx)
		err = activityTwoIndicatorDao.DeleteByWhere(dao.ActivityTwoIndicator{ActivityId: activityId, OneIndicatorId: oneIndicatorId})
		return err
	})
	return
}

func (s *ActivityService) DeleteActivityTwoIndicator(activityId, twoIndicatorId int) (err error) {
	err = s.activityIsStart(activityId)
	if err != nil {
		return err
	}
	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		where := dao.ActivityTwoIndicator{ActivityId: activityId, TwoIndicatorId: twoIndicatorId}
		activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDaoWithDB(tx, s.appCtx)
		var activityTwoIndicator dao.ActivityTwoIndicator
		err := activityTwoIndicatorDao.First(where, &activityTwoIndicator)
		if err != nil {
			return err
		}
		err = activityTwoIndicatorDao.DeleteByWhere(where)
		if err != nil {
			return err
		}
		count, err := activityTwoIndicatorDao.Count(dao.ActivityTwoIndicator{ActivityId: activityId, OneIndicatorId: activityTwoIndicator.OneIndicatorId})
		if err != nil {
			return err
		}
		if count == 0 {
			activityOneIndicatorDao := dao.NewActivityOneIndicatorDaoWithDB(tx, s.appCtx)
			err := activityOneIndicatorDao.DeleteByWhere(dao.ActivityOneIndicator{ActivityId: activityId, OneIndicatorId: activityTwoIndicator.OneIndicatorId})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (s *ActivityService) DeleteActivity(activityId int) (err error) {
	err = s.activityIsStart(activityId)
	if err != nil {
		return err
	}
	activityDao := dao.NewActivityDao(s.appCtx)
	err = activityDao.DeleteByWhere(dao.Activity{ActivityId: activityId})
	return
}

func (s *ActivityService) GetActivityList(params *req.GetActivityListReq) (data *resp.GetActivityListResp, err error) {
	data = &resp.GetActivityListResp{}
	nowTime := util.NowTime()
	db := global.GormDB
	if params.ActivityName != "" {
		db = db.Where("activity_name like ?", "%"+params.ActivityName+"%")
	}

	db = db.Select(`activity_id,activity_name,year,start_time,end_time,create_time,
	CASE 
        WHEN start_time > ? THEN '1'
        WHEN end_time < ? THEN '3'
        ELSE '2'
    END AS status`, nowTime, nowTime)
	activityDao := dao.NewActivityDaoWithDB(db, s.appCtx)
	data.Total, err = activityDao.FindAndCountWithPageOrder(dao.Activity{Year: params.Year}, &data.List, params.Page, params.Limit, "status asc, start_time asc")
	if err != nil {
		return
	}
	activityIds := util.ListObjToListObj(data.List, func(obj resp.Activity) int {
		return obj.ActivityId
	})
	userNumMap, err := dao.NewUserActivityDao(s.appCtx).GetActivityUserNumMap(activityIds)
	for i, v := range data.List {
		data.List[i].AttendNum = userNumMap[v.ActivityId]
	}
	return
}

func (s *ActivityService) GetActivityTwoIndicatorList(activityId, oneIndicatorId int) (data *resp.GetActivityTwoIndicatorListResp, err error) {
	data = &resp.GetActivityTwoIndicatorListResp{}
	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	err = activityTwoIndicatorDao.Find(dao.ActivityTwoIndicator{ActivityId: activityId, OneIndicatorId: oneIndicatorId}, &data.List)
	return
}

func (s *ActivityService) GetActivityIndicatorByTwoId(activityId, twoId int) (indicatorInfo resp.ActivityIndicatorInfo, err error) {
	var activityTwoIndicator dao.ActivityTwoIndicator
	activityTwoIndicatorDao := dao.NewActivityTwoIndicatorDao(s.appCtx)
	err = activityTwoIndicatorDao.First(dao.ActivityTwoIndicator{ActivityId: activityId, TwoIndicatorId: twoId}, &activityTwoIndicator)
	if err != nil {
		return
	}
	indicatorInfo.TwoIndicatorId = activityTwoIndicator.TwoIndicatorId
	indicatorInfo.TwoIndicatorName = activityTwoIndicator.TwoIndicatorName
	indicatorInfo.Score = activityTwoIndicator.Score
	var activityOneIndicator dao.ActivityOneIndicator
	activityOneIndicatorDao := dao.NewActivityOneIndicatorDao(s.appCtx)
	err = activityOneIndicatorDao.First(dao.ActivityOneIndicator{ActivityId: activityId, OneIndicatorId: activityTwoIndicator.OneIndicatorId}, &activityOneIndicator)
	if err != nil {
		return
	}
	indicatorInfo.OneIndicatorId = activityOneIndicator.OneIndicatorId
	indicatorInfo.OneIndicatorName = activityOneIndicator.OneIndicatorName
	indicatorInfo.Content = activityOneIndicator.Content
	return
}

func (s *ActivityService) GetActivityYearList() (data *resp.GetActivityYearListResp, err error) {
	data = &resp.GetActivityYearListResp{}
	activityDao := dao.NewActivityDao(s.appCtx)
	data.Years, err = activityDao.GetYearsByGroup()
	return
}

func (s *ActivityService) GetLatestActivity(userId string) (data *resp.GetLatestActivityResp, err error) {
	data = &resp.GetLatestActivityResp{}
	activityDao := dao.NewActivityDaoWithDB(global.GormDB.Order("create_time desc"), s.appCtx)
	err = activityDao.First(nil, &data.Activity)
	if err != nil {
		err = errorz.CodeMsg(errorz.RESP_ERR, "活动尚未开始")
	}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	count, err := userActivityDao.Count(dao.UserActivity{UserId: userId, ActivityId: data.ActivityId})
	if count > 0 {
		data.HasDeclare = true
	}
	return
}

func (s *ActivityService) EndActivity(activityId int) (err error) {
	err = s.activityIsOngoing(activityId)
	if err != nil {
		return
	}
	activityDao := dao.NewActivityDao(s.appCtx)
	nowTime := util.NowTime()
	err = activityDao.UpdateByWhere(dao.Activity{ActivityId: activityId}, dao.Activity{EndTime: &nowTime, UpdateTime: &nowTime})
	return
}
