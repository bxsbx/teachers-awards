package req

type GetSimpleSumStatsReq struct {
	Year     int    `form:"year" binding:"required"` // 年份
	SchoolId string `form:"school_id"`               // 学校id
}

type GetAwardRateReq struct {
	Year     int    `form:"year" `     // 年份
	SchoolId string `form:"school_id"` // 学校id
}

type GetDeclareRateReq struct {
	Year     int    `form:"year"`      // 年份
	SchoolId string `form:"school_id"` // 学校id
}

type GetEveryYearAwardNumReq struct {
	SchoolId string `form:"school_id"` // 学校id
}

type GetEverySchoolAwardNumReq struct {
	Year int `form:"year"` // 年份
}

type GetEveryTeacherTypeAwardNumReq struct {
	Year int `form:"year"` // 年份
}

type GetYearDeclareAwardRankReq struct {
	SchoolId string `form:"school_id"` // 学校id
}

type GetSchoolDeclareAwardRankReq struct {
	SchoolName string `form:"school_name"`                           // 学校名称
	Year       int    `form:"year" binding:"required"`               // 年份
	Page       int    `form:"page" binding:"required,gte=1"`         // 页数
	Limit      int    `form:"limit" binding:"required,gte=1,lte=50"` // 每页大小
}

type GetTeacherTypeDeclareAwardRankReq struct {
	Year int `form:"year" binding:"required"` // 年份
}

type GetResultGroupByDeclareTypeReq struct {
	ActivityId int `form:"activity_id"  binding:"required"` //活动id
}
