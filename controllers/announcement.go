package controllers

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/services"
)

type AnnouncementController struct {
}

//	@Summary	根据id获取公告详情
//	@Tags		Announcement
//	@Produce	json
//	@Param		user_id			query		string	true	"用户id"
//	@Param		announcement_id	query		int		true	"公告id"
//	@Response	200				{object}	resp.Response{data=resp.GetAnnouncementByIdResp}
//	@Router		/v1/announcement/get [GET]
func (u *AnnouncementController) GetAnnouncementById(c *gin.Context) {
	var params req.GetAnnouncementByIdReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	announcementService := services.NewAnnouncementService(global.GetContext(c))
	data, err := announcementService.GetAnnouncementById(params.AnnouncementId, params.UserId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	保存公告
//	@Tags		Announcement
//	@Produce	json
//	@Param		announcement_id	formData	int		false	"公告id"
//	@Param		user_id			formData	string	true	"用户id"
//	@Param		user_name		formData	string	true	"用户名称"
//	@Param		title			formData	string	true	"标题"
//	@Param		content			formData	string	true	"内容"
//	@Param		annex			formData	string	false	"附件"
//	@Response	200				{object}	resp.Response{data=resp.SaveAnnouncementResp}
//	@Router		/v1/announcement/save [POST]
func (u *AnnouncementController) SaveAnnouncement(c *gin.Context) {
	var params req.SaveAnnouncementReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	announcementService := services.NewAnnouncementService(global.GetContext(c))
	data, err := announcementService.SaveAnnouncement(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	删除公告
//	@Tags		Announcement
//	@Produce	json
//	@Param		announcement_id	query		int	true	"公告id"
//	@Response	200				{object}	resp.Response
//	@Router		/v1/announcement/delete [DELETE]
func (u *AnnouncementController) DeleteAnnouncementById(c *gin.Context) {
	var params req.DeleteAnnouncementByIdReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	announcementService := services.NewAnnouncementService(global.GetContext(c))
	err := announcementService.DeleteAnnouncementById(params.AnnouncementId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

//	@Summary	获取公告列表
//	@Tags		Announcement
//	@Produce	json
//	@Param		user_id				query		string	true	"用户id"
//	@Param		input_content		query		string	false	"输入的内容"
//	@Param		publish_time		query		string	false	"发布时间，格式2006-01-02"
//	@Param		create_time_order	query		string	false	"按时间升序或者降序"
//	@Param		page				query		int		true	"页数"
//	@Param		limit				query		int		true	"每页大小"
//	@Response	200					{object}	resp.Response{data=resp.GetAnnouncementListResp}
//	@Router		/v1/announcement/list [GET]
func (u *AnnouncementController) GetAnnouncementList(c *gin.Context) {
	var params req.GetAnnouncementListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	util.TimeToLocation(params.PublishTime)
	announcementService := services.NewAnnouncementService(global.GetContext(c))
	data, err := announcementService.GetAnnouncementList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

//	@Summary	公告全部设置为已读
//	@Tags		Announcement
//	@Produce	json
//	@Param		user_id	formData	string	true	"用户id"
//	@Response	200		{object}	resp.Response
//	@Router		/v1/announcement/user/all/read [PUT]
func (u *AnnouncementController) AllAnnouncementRead(c *gin.Context) {
	var params req.AllAnnouncementReadReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	announcementService := services.NewAnnouncementService(global.GetContext(c))
	err := announcementService.AllAnnouncementRead(params.UserId)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}
