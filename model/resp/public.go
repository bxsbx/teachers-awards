package resp

type GetTokenAndUserInfoResp struct {
	UserId        string `json:"user_id"`         //用户id
	UserName      string `json:"user_name"`       //用户名称
	UserRoles     []int  `json:"user_roles"`      //用户角色类型，1：学校，2：专家，3：教育局，4：老师，5：管理员
	Token         string `json:"token"`           //jwt的token（本系统独自维护)
	ExpiresTimeAt int64  `json:"expires_time_at"` //过期时间戳
}
