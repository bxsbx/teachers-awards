package req

type ExportReviewListReq struct {
	JudgesId         string `form:"judges_id"`                                     //评委id
	UserName         string `form:"user_name"`                                     //用户姓名
	UserSex          int    `form:"user_sex"`                                      //1：男，2：女
	SubjectCode      string `form:"subject_code"`                                  //科目code
	SchoolId         string `form:"school_id"`                                     //学校id
	SchoolName       string `form:"school_name"`                                   //学校名称
	Year             int    `form:"year"`                                          //年份
	OneIndicatorName string `form:"one_indicator_name"`                            //一级指标名称
	TwoIndicatorName string `form:"two_indicator_name"`                            //二级指标名称
	ReviewStatus     int    `form:"review_status" binding:"oneof=0 1 2 3 4"`       //审核状态 0:全部 1:待审核，2：未通过，3：已通过，4：已审核"
	ReviewProcess    int    `form:"review_process" binding:"required,oneof=1 2 3"` // 当前审核进程，1：学校，2：专家
}

type ExportHistoryActivityListReq struct {
	UserName     string `form:"user_name"`     //用户姓名
	UserSex      int    `form:"user_sex"`      //1：男，2：女
	SubjectCode  string `form:"subject_code"`  //科目code
	SchoolId     string `form:"school_id"`     //学校id
	SchoolName   string `form:"school_name"`   //学校名称
	Year         int    `form:"year"`          //年份
	DeclareType  int    `form:"declare_type"`  //申报类型
	IdentityCard string `form:"identity_card"` //身份证号
	RankPrize    *int   `form:"rank_prize"`    //0：无，1：一等奖，2：二等奖，3：三等奖
	Rank         int    `form:"rank"`          //排名
	FinalScore   int    `form:"final_score"`   //最终得分（各项通过的审核）
}

type ExportUserDeclareRecordListByYearReq struct {
	Year          int    `form:"year" binding:"required"`            //年份
	DeclareUserId string `form:"declare_user_id" binding:"required"` //申报的用户id
	Role          int    `form:"role" binding:"required"`            //角色类型，1：学校，2：专家，3：教育局，4：老师
}

type ExportEdbReviewListReq struct {
	UserName    string `form:"user_name"`                    //用户姓名
	UserSex     int    `form:"user_sex"`                     //1：男，2：女
	SubjectCode string `form:"subject_code"`                 //科目code
	SchoolId    string `form:"school_id"`                    //学校id
	SchoolName  string `form:"school_name"`                  //学校名称
	Status      int    `form:"status" binding:"oneof=0 1 2"` //审核状态 0:全部 1:待审核，2：已审核
}

type ExportAwardsSetListReq struct {
	ActivityId   int    `form:"activity_id"  binding:"required"` //活动id
	UserName     string `form:"user_name"`                       //用户姓名
	UserSex      int    `form:"user_sex"`                        //1：男，2：女
	SubjectCode  string `form:"subject_code"`                    //科目code
	SchoolId     string `form:"school_id"`                       //学校id
	SchoolName   string `form:"school_name"`                     //学校名称
	DeclareType  int    `form:"declare_type"`                    //申报类型
	IdentityCard string `form:"identity_card"`                   //身份证号
	RankPrize    *int   `form:"rank_prize"`                      //名次，0：无，1：一等奖，2：二等奖，3：三等奖
	FinalScore   int    `form:"final_score"`                     //最终得分（各项通过的审核）
}
