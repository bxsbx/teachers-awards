package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func NoticeRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.NoticeController{}

	//获取通知列表
	router.GET("/v1/notice/user/list", api.GetUserNoticeList)
	//用户通知全部设置为已读
	router.PUT("/v1/notice/user/all/read", api.UpdateUserAllNoticeToRead)
	//用户通知全部删除
	router.DELETE("/v1/notice/user/all/delete", api.DeleteUserAllNotice)
	//删除用户通知
	router.DELETE("/v1/notice/user/delete/ids", api.DeleteUserNoticeByIds)
	// router general tag
}
