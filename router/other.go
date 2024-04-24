package router

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/controllers"
)

func OtherRouter(group *gin.RouterGroup) {
	router := group.Group("")
	api := &controllers.OtherController{}

	//获取上传token
	router.GET("/v1/other/upload/token", api.GetUploadToken)
	//获取uai表的的操作记录
	router.GET("/v1/other/operation/record/uai", api.GetOperationRecordFromUai)
	//获取学科枚举
	router.GET("/v1/other/subject/enum", api.GetSubjectEnum)
	//分页获取中台学校列表
	router.GET("/v1/other/school/list/page", api.GetSchoolListByPage)
	//通过学校名称模糊匹配获取学校列表
	router.GET("/v1/other/school/list/by/name", api.GetSchoolListByName)
	// router general tag
}
