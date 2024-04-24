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

type AnnouncementService struct {
	appCtx context.Context
}

func NewAnnouncementService(appCtx context.Context) *AnnouncementService {
	return &AnnouncementService{appCtx: appCtx}
}

func (s *AnnouncementService) GetAnnouncementById(announcementId int, userId string) (data *resp.GetAnnouncementByIdResp, err error) {
	data = &resp.GetAnnouncementByIdResp{}
	announcementDao := dao.NewAnnouncementDao(s.appCtx)
	err = announcementDao.First(dao.Announcement{AnnouncementId: announcementId}, &data.Announcement)
	if err != nil {
		return
	}
	readDao := dao.NewReadDao(s.appCtx)
	count, err := readDao.Count(dao.Read{ReadId: announcementId, ReadType: 1, ReadUserId: userId})
	if err != nil {
		return
	}
	if count == 0 {
		nowTime := util.NowTime()
		err = readDao.Create(dao.Read{ReadId: announcementId, ReadType: 1, ReadUserId: userId, CreateTime: &nowTime})
	}
	return
}

func (s *AnnouncementService) SaveAnnouncement(params *req.SaveAnnouncementReq) (data *resp.SaveAnnouncementResp, err error) {
	data = &resp.SaveAnnouncementResp{}
	announcement := dao.Announcement{
		AnnouncementId: params.AnnouncementId,
		Title:          params.Title,
		Content:        params.Content,
		UserId:         params.UserId,
		UserName:       params.UserName,
		Annex:          params.Annex,
	}
	announcementDao := dao.NewAnnouncementDao(s.appCtx)
	if params.AnnouncementId != 0 {
		err = announcementDao.UpdateByWhere(dao.Announcement{AnnouncementId: params.AnnouncementId}, announcement)
	} else {
		nowTime := util.NowTime()
		announcement.CreateTime = &nowTime
		err = announcementDao.Create(&announcement)
	}
	data.AnnouncementId = announcement.AnnouncementId
	return
}

func (s *AnnouncementService) DeleteAnnouncementById(announcementId int) (err error) {
	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		announcementDao := dao.NewAnnouncementDaoWithDB(tx, s.appCtx)
		err := announcementDao.DeleteByWhere(dao.Announcement{AnnouncementId: announcementId})
		if err != nil {
			return err
		}
		err = dao.NewReadDaoWithDB(tx, s.appCtx).DeleteByWhere(dao.Read{ReadId: announcementId, ReadType: 1})
		return err
	})
	return
}

func (s *AnnouncementService) GetAnnouncementList(params *req.GetAnnouncementListReq) (data *resp.GetAnnouncementListResp, err error) {
	data = &resp.GetAnnouncementListResp{}
	db := global.GormDB
	if params.InputContent != "" {
		like := "%" + params.InputContent + "%"
		db = db.Where("title like ? or content like ?", like, like)
	}
	if params.PublishTime != nil {
		db = db.Where("create_time between ? and ?", *params.PublishTime, params.PublishTime.AddDate(0, 0, 1))
	}
	announcementDao := dao.NewAnnouncementDaoWithDB(db, s.appCtx)
	var list []dao.Announcement
	data.Total, err = announcementDao.FindAndCountWithPageOrder(nil, &list, params.Page, params.Limit, "create_time "+params.CreatTimeOrder)
	if err != nil {
		return
	}
	announcementIds := util.ListObjToListObj(list, func(obj dao.Announcement) int {
		return obj.AnnouncementId
	})
	readDao := dao.NewReadDao(s.appCtx)
	readMap, err := readDao.GetReadMap(announcementIds, 1, params.UserId)
	if err != nil {
		return
	}
	data.List = make([]resp.Announcement, len(list))
	for i, v := range list {
		data.List[i].Announcement = v
		data.List[i].IsRead = readMap[v.AnnouncementId]
	}
	return
}

func (s *AnnouncementService) AllAnnouncementRead(userId string) (err error) {
	announcementDao := dao.NewAnnouncementDaoWithDB(global.GormDB.
		Where("announcement_id not in (select read_id from `read` where read_user_id = ? and read_type = ?)", userId, 1), s.appCtx)
	var ids []int
	err = announcementDao.Pluck(nil, &ids, "announcement_id")
	if err != nil {
		return
	}
	readDao := dao.NewReadDao(s.appCtx)
	var list []dao.Read
	nowTime := util.NowTime()
	for _, id := range ids {
		list = append(list, dao.Read{ReadId: id, ReadType: 1, ReadUserId: userId, CreateTime: &nowTime})
	}
	err = readDao.BatchInsert(list)
	if err != nil {
		return
	}
	return
}
