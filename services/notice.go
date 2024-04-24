package services

import (
	"context"
	"teachers-awards/dao"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/model/resp"
)

type NoticeService struct {
	appCtx context.Context
}

func NewNoticeService(appCtx context.Context) *NoticeService {
	return &NoticeService{appCtx: appCtx}
}

func (s *NoticeService) GetUserNoticeList(params *req.GetUserNoticeListReq) (data *resp.GetUserNoticeListResp, err error) {
	data = &resp.GetUserNoticeListResp{}
	noticeDao := dao.NewNoticeDao(s.appCtx)
	data.Total, err = noticeDao.FindAndCountWithPageOrder(dao.Notice{UserId: params.UserId}, &data.List, params.Page, params.Limit, "create_time "+params.CreatTimeOrder)
	return
}

func (s *NoticeService) UpdateUserAllNoticeToRead(userId string) (err error) {
	noticeDao := dao.NewNoticeDao(s.appCtx)
	err = noticeDao.UpdateByWhere(dao.Notice{UserId: userId}, dao.Notice{IsRead: 1})
	return
}

func (s *NoticeService) DeleteUserAllNotice(userId string) (err error) {
	noticeDao := dao.NewNoticeDao(s.appCtx)
	err = noticeDao.DeleteByWhere(dao.Notice{UserId: userId})
	return
}

func (s *NoticeService) DeleteUserNoticeByIds(userId string, noticeIds []int) (err error) {
	noticeDao := dao.NewNoticeDaoWithDB(global.GormDB.Where("notice_id in (?)", noticeIds), s.appCtx)
	err = noticeDao.DeleteByWhere(dao.Notice{UserId: userId})
	return
}
