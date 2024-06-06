package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func PublicRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.PublicController{}

	//获取token和用户的基本信息
	router.GET("/v1/public/get/info/token/user", api.GetTokenAndUserInfo)
	//通过code获取用户信息和token
	router.GET("/v1/public/get/user/info/by/code", api.GetTokenAndUserInfoByCode)
	// router general tag
}
