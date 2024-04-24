package resp

import "teachers-awards/dao"

type SimpleSumStats struct {
	Type       int     `json:"type"`         //1:学校，2：老师，3：获奖，4：申报
	Num        int64   `json:"num"`          //数量
	YearOnYear float64 `json:"year_on_year"` //同比
}

type GetSimpleSumStatsResp struct {
	List []SimpleSumStats `json:"list"`
}

type GetAwardRateResp struct {
	DeclareNum int64 `json:"declare_num"` //申报数量
	AwardNum   int64 `json:"award_num"`   //获奖人数
}

type GetDeclareRateResp struct {
	DeclareNum int64 `json:"declare_num"` //申报数量
	SumNum     int64 `json:"sum_num"`     //总人数
}

type GetEveryYearAwardNumResp struct {
	List []dao.YearDeclareAwardNum `json:"list"`
}

type GetEverySchoolAwardNumResp struct {
	List []dao.SchoolAwardNum `json:"list"`
}

type GetEveryTeacherTypeAwardNumResp struct {
	List []dao.TeacherTypeAwardNum `json:"list"`
}

type DeclareAwardRank struct {
	dao.DeclareAwardRankNum
	DeclareOnYear float64 `json:"declare_on_year"` //申报同比
	AwardOnYear   float64 `json:"award_on_year"`   //获奖同比
	AwardRate     float64 `json:"award_rate"`      //获奖占比
}

type YearDeclareAwardRank struct {
	Num  int `json:"num"`  //数量
	Year int `json:"year"` //年份
	DeclareAwardRank
}

type GetYearDeclareAwardRankResp struct {
	List []YearDeclareAwardRank `json:"list"`
}

type SchoolDeclareAwardRank struct {
	SchoolId   string `json:"school_id"`   //学校id
	SchoolName string `json:"school_name"` //学校名称
	DeclareAwardRank
}

type GetSchoolDeclareAwardRankResp struct {
	Total int64                    `json:"total"`
	List  []SchoolDeclareAwardRank `json:"list"`
}

type TeacherTypeDeclareAwardRank struct {
	DeclareType int `json:"declare_type"` //1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研
	DeclareAwardRank
}

type GetTeacherTypeDeclareAwardRankResp struct {
	List []TeacherTypeDeclareAwardRank `json:"list"`
}

type GetResultGroupByDeclareTypeResp struct {
	SchoolNum  int64                          `json:"school_num"`
	DeclareNum int64                          `json:"declare_num"`
	List       []dao.ResultGroupByDeclareType `json:"list"`
}
