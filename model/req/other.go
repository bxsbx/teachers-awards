package req

type GetUploadTokenReq struct {
	UserId   string `form:"user_id"  binding:"required"`   //用户id
	FileName string `form:"file_name"  binding:"required"` //文件名
}

type GetOperationRecordFromUaiReq struct {
	UserActivityIndicatorId int64 `form:"user_activity_indicator_id"  binding:"required"` //用户活动申报单项目id
}

type GetSubjectEnumReq struct {
	Refresh int `form:"refresh"` //是否刷新学科枚举，0：否，1：是
}

type GetSchoolListByPageReq struct {
	Page  int `form:"page" binding:"required,gte=1"`          // 页数
	Limit int `form:"limit" binding:"required,gte=1,lte=100"` // 每页大小
}

type GetSchoolListByNameReq struct {
	SchoolName string `form:"school_name"  binding:"required"` //学校名称
}
