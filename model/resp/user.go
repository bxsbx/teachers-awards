package resp

import (
	"teachers-awards/dao"
	"time"
)

type GetUserInfoResp struct {
	dao.UserInfo
}

type UserInfo struct {
	UserId       string `json:"user_id"`       //用户id
	UserName     string `json:"user_name"`     //用户名称
	UserSex      int    `json:"user_sex"`      //1：男，2：女
	IdentityCard string `json:"identity_card"` //身份证号
	Phone        string `json:"phone"`         //手机号
	SchoolId     string `json:"school_id"`     //学校id
	SchoolName   string `json:"school_name"`   //学校名称
	Roles        []int  `json:"roles"`         //角色列表，1：学校，2：专家，3：教育局，4：老师，5：管理员
}

type GetUsersByNameResp struct {
	Total    int64      `json:"total"`
	UserList []UserInfo `json:"user_list"`
}

type GetUserListByWhereResp struct {
	Total int64          `json:"total"`
	List  []dao.UserInfo `json:"list"`
}

type OneIndicatorOnlyName struct {
	OneIndicatorId   int                    `json:"one_indicator_id"`
	OneIndicatorName string                 `json:"one_indicator_name"`
	TwoList          []TwoIndicatorOnlyName `json:"two_list"` //二级指标列表
}

type TwoIndicatorOnlyName struct {
	TwoIndicatorId   int    `json:"two_indicator_id"`   //二级指标id
	TwoIndicatorName string `json:"two_indicator_name"` //二级指标名称
}

type ExpertAuthList struct {
	UserId                  string                 `json:"user_id"`                    //用户id
	UserName                string                 `json:"user_name"`                  //用户名称
	UserSex                 int                    `json:"user_sex"`                   //1：男，2：女
	IdentityCard            string                 `json:"identity_card"`              //身份证号
	Phone                   string                 `json:"phone"`                      //手机号
	ExportAuth              int                    `json:"export_auth"`                //专家是否授权 0：未授权 1：已授权
	AuthDay                 *time.Time             `json:"auth_day"`                   //授权日期，格式2006-01-02
	ExpertAuthIndicatorList []OneIndicatorOnlyName `json:"expert_auth_indicator_list"` // 专家授权指标
}

type GetExpertAuthListByWhereResp struct {
	Total int64            `json:"total"`
	List  []ExpertAuthList `json:"list"`
}
