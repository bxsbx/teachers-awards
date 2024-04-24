package controllers

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/services"
)

type NoticeController struct {
}

// @Summary	获取通知列表
// @Tags		Notice
// @Produce	json
// @Param		user_id				query		string	true	"用户id"
// @Param		create_time_order	query		string	false	"按时间升序或者降序"
// @Param		page				query		int		true	"页数"
// @Param		limit				query		int		true	"每页大小"
// @Response	200					{object}	resp.Response{data=resp.GetUserNoticeListResp}
// @Router		/v1/notice/user/list [GET]
func (u *NoticeController) GetUserNoticeList(c *gin.Context) {
	var params req.GetUserNoticeListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	noticeService := services.NewNoticeService(global.GetContext(c))
	data, err := noticeService.GetUserNoticeList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	用户通知全部设置为已读
// @Tags		Notice
// @Produce	json
// @Param		user_id	query		string	true	"用户id"
// @Response	200		{object}	resp.Response
// @Router		/v1/notice/user/all/read [PUT]
func (u *NoticeController) UpdateUserAllNoticeToRead(c *gin.Context) {
	var params req.UpdateUserAllNoticeToReadReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	noticeService := services.NewNoticeService(global.GetContext(c))
	err := noticeService.UpdateUserAllNoticeToRead(params.UserId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	用户通知全部删除
// @Tags		Notice
// @Produce	json
// @Param		user_id	query		string	true	"用户id"
// @Response	200		{object}	resp.Response
// @Router		/v1/notice/user/all/delete [DELETE]
func (u *NoticeController) DeleteUserAllNotice(c *gin.Context) {
	var params req.DeleteUserAllNoticeReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	noticeService := services.NewNoticeService(global.GetContext(c))
	err := noticeService.DeleteUserAllNotice(params.UserId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	删除用户通知
// @Tags		Notice
// @Produce	json
// @Param		notice_ids	query		string	true	"通知ids，英文","隔开"
// @Param		user_id		query		string	true	"用户id"
// @Response	200			{object}	resp.Response
// @Router		/v1/notice/user/delete/ids [DELETE]
func (u *NoticeController) DeleteUserNoticeByIds(c *gin.Context) {
	var params req.DeleteUserNoticeByIdsReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	intArray, err := util.StrToIntArray(params.NoticeIds, ",")
	if err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	noticeService := services.NewNoticeService(global.GetContext(c))
	err = noticeService.DeleteUserNoticeByIds(params.UserId, intArray)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}
