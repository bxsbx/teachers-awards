package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func UserRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.UserController{}

	//获取用户信息
	router.GET("/v1/user/get/info", api.GetUserInfo)
	//保存当前用户信息
	router.POST("/v1/user/save/info", api.SaveUserInfo)
	//通过用户名称模糊匹配获取用户列表
	router.GET("/v1/user/list/by/name", api.GetUsersByName)
	//设置用户角色
	router.POST("/v1/user/set/role", api.SetRoleToUser)
	//获取已设置角色的用户列表
	router.GET("/v1/user/get/role/list", api.GetUserListByWhere)
	//获取专家角色列表
	router.GET("/v1/user/get/expert/auth/list", api.GetExpertAuthListByWhere)
	//给专家授权指标项
	router.POST("/v1/user/set/expert/auth", api.SetExpertAuth)
	//取消专家授权
	router.DELETE("/v1/user/cancel/expert/auth", api.CancelExpertAuth)
	// router general tag
}
