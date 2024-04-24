package services

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/common/util"
	"teachers-awards/dao"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/model/resp"
)

type IndicatorService struct {
	appCtx context.Context
}

func NewIndicatorService(appCtx context.Context) *IndicatorService {
	return &IndicatorService{appCtx: appCtx}
}

func (s *IndicatorService) CreateOrUpdateOneIndicator(params *req.CreateOrUpdateOneIndicatorReq) (err error) {
	oneIndicatorDao := dao.NewOneIndicatorDao(s.appCtx)
	oneIndicator := dao.OneIndicator{
		OneIndicatorName: params.OneIndicatorName,
		Content:          params.Content,
	}
	nowTime := util.NowTime()
	if params.OneIndicatorId != 0 {
		oneIndicator.UpdateTime = &nowTime
		err = oneIndicatorDao.UpdateByWhere(dao.OneIndicator{OneIndicatorId: params.OneIndicatorId}, oneIndicator)
	} else {
		oneIndicator.CreateTime = &nowTime
		err = oneIndicatorDao.Create(&oneIndicator)
	}
	return
}

func (s *IndicatorService) GetOneIndicatorList(params *req.GetOneIndicatorListReq) (data *resp.GetOneIndicatorListResp, err error) {
	data = &resp.GetOneIndicatorListResp{}
	db := global.GormDB
	if params.InputName != "" {
		db = db.Where("one_indicator_name like ?", "%"+params.InputName+"%")
	}
	if params.InputContent != "" {
		db = db.Where("content like ?", "%"+params.InputContent+"%")
	}
	oneIndicatorDao := dao.NewOneIndicatorDaoWithDB(db, s.appCtx)
	var list []dao.OneIndicator
	data.Total, err = oneIndicatorDao.FindAndCountWithPageOrder(nil, &list, params.Page, params.Limit, "")
	if err != nil {
		return
	}
	oneIds := util.ListToDeduplicationList(list, func(t dao.OneIndicator) (int, int) {
		return t.OneIndicatorId, t.OneIndicatorId
	})

	if params.WithTwo {
		var twoList []dao.TwoIndicator
		twoIndicatorDao := dao.NewTwoIndicatorDaoWithDB(db.Where("one_indicator_id in (?)", oneIds), s.appCtx)
		err = twoIndicatorDao.Find(nil, &twoList)
		if err != nil {
			return
		}
		toMapList := util.ListObjToMapList(twoList, func(obj dao.TwoIndicator) (int, dao.TwoIndicator) {
			return obj.OneIndicatorId, obj
		})
		for _, v := range list {
			twos := toMapList[v.OneIndicatorId]
			data.List = append(data.List, resp.OneIndicator{OneIndicator: v, Count: len(twos), List: twos})
		}
	} else {
		twoIndicatorDao := dao.NewTwoIndicatorDao(s.appCtx)
		var twoIndicatorCountMap map[int]int
		//获取二级指标个数
		twoIndicatorCountMap, err = twoIndicatorDao.GetTwoIndicatorCountMap(oneIds)
		if err != nil {
			return
		}
		for _, v := range list {
			data.List = append(data.List, resp.OneIndicator{OneIndicator: v, Count: twoIndicatorCountMap[v.OneIndicatorId]})
		}
	}
	return
}

func (s *IndicatorService) DeleteOneIndicatorByIds(ids []int) (err error) {
	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		oneIndicatorDao := dao.NewOneIndicatorDaoWithDB(tx.Where("one_indicator_id in (?)", ids), s.appCtx)
		err := oneIndicatorDao.DeleteByWhere(nil)
		if err != nil {
			return err
		}
		twoIndicatorDao := dao.NewTwoIndicatorDaoWithDB(tx.Where("one_indicator_id in (?)", ids), s.appCtx)
		err = twoIndicatorDao.DeleteByWhere(nil)
		return err
	})
	return
}

func (s *IndicatorService) CreateOrUpdateTwoIndicator(params *req.CreateOrUpdateTwoIndicatorReq) (data *resp.CreateOrUpdateTwoIndicatorResp, err error) {
	data = &resp.CreateOrUpdateTwoIndicatorResp{}
	twoIndicatorDao := dao.NewTwoIndicatorDao(s.appCtx)
	twoIndicator := dao.TwoIndicator{
		TwoIndicatorName: params.TwoIndicatorName,
		Score:            params.Score,
		OneIndicatorId:   params.OneIndicatorId,
	}
	nowTime := util.NowTime()
	if params.TwoIndicatorId != 0 {
		twoIndicator.UpdateTime = &nowTime
		err = twoIndicatorDao.UpdateByWhere(dao.TwoIndicator{TwoIndicatorId: params.TwoIndicatorId}, twoIndicator)
	} else {
		twoIndicator.CreateTime = &nowTime
		err = twoIndicatorDao.Create(&twoIndicator)
	}
	data.TwoIndicator = twoIndicator
	return
}

func (s *IndicatorService) DeleteTwoIndicatorByIds(ids []int) (err error) {
	twoIndicatorDao := dao.NewTwoIndicatorDaoWithDB(global.GormDB.Where("two_indicator_id in (?)", ids), s.appCtx)
	err = twoIndicatorDao.DeleteByWhere(nil)
	return
}

func (s *IndicatorService) GetTwoIndicatorList(params *req.GetTwoIndicatorListReq) (data *resp.GetTwoIndicatorListResp, err error) {
	data = &resp.GetTwoIndicatorListResp{}
	db := global.GormDB
	if params.InputName != "" {
		db = db.Where("two_indicator_name like ?", "%"+params.InputName+"%")
	}
	where := dao.TwoIndicator{
		OneIndicatorId: params.OneIndicatorId,
		Score:          params.InputScore,
	}
	twoIndicatorDao := dao.NewTwoIndicatorDaoWithDB(db, s.appCtx)
	var list []dao.TwoIndicator
	data.Total, err = twoIndicatorDao.FindAndCountWithPageOrder(where, &list, params.Page, params.Limit, "")
	if err != nil {
		return
	}
	//获取一级指标名称并组装
	oneIds := util.ListToDeduplicationList(list, func(t dao.TwoIndicator) (int, int) {
		return t.OneIndicatorId, t.OneIndicatorId
	})
	oneIndicatorDao := dao.NewOneIndicatorDao(s.appCtx)
	oneIndicatorNameMap, err := oneIndicatorDao.GetOneIndicatorNameMap(oneIds)
	if err != nil {
		return
	}
	for _, v := range list {
		data.List = append(data.List, resp.TwoIndicator{TwoIndicator: v, OneIndicatorName: oneIndicatorNameMap[v.OneIndicatorId]})
	}
	return
}
