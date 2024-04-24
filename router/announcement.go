package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func AnnouncementRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.AnnouncementController{}

	//根据id获取公告详情
	router.GET("/v1/announcement/get", api.GetAnnouncementById)
	//保存公告
	router.POST("/v1/announcement/save", api.SaveAnnouncement)
	//删除公告
	router.DELETE("/v1/announcement/delete", api.DeleteAnnouncementById)
	//获取公告列表
	router.GET("/v1/announcement/list", api.GetAnnouncementList)
	//公告全部设置为已读
	router.PUT("/v1/announcement/user/all/read", api.AllAnnouncementRead)
	// router general tag
}
