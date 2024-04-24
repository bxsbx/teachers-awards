package global

import (
	"teachers-awards/common/jwt"
)

// 无需刷新时间区（NoRefreshTime）等于ExpiresTime-RefreshTime
// NoRefreshTime+RefreshTime=ExpiresTime

var Jwt = jwt.Jwt{
	CustomSecret: "teachers-awards-jwt|#8233AK93%&*KO",
	ExpiresTime:  7 * 24 * 60 * 60,     // token有效时间区 (有效期一周)
	RefreshTime:  (7*24 - 1) * 60 * 60, // token刷新时间区
}

// redis-key
const (
	JwtKey      = "teachers-awards:jwt:"
	ZtConfigKey = "ztConfig:organization:code:"
)

const (
	AuthToken = "Auth-Token"
	ExpiresAt = "Expires-At"
)

// 角色  1：学校，2：专家，3：教育局，4：老师，5：超级管理员
const (
	RoleSchool = iota + 1
	RoleExpert
	RoleEdb
	RoleTeacher
	RoleAdmin
)

// 角色  1：学校，2：专家，4：教育局，8：老师，16：超级管理员
const (
	Role2School  = 1
	Role2Expert  = 2
	Role2Edb     = 4
	Role2Teacher = 8
	Role2Admin   = 16
)

// 是否通过，0：未通过 1：已通过
const (
	PassNo = iota
	PassYes
)

// 审核，1:待审核，2：未通过，3：已通过，4：已审核
const (
	ReviewWait = iota + 1
	ReviewNoPass
	ReviewPass
	ReviewFinish
)

// 审核进程  1：学校，2：专家，3：教育局，4：结束
const (
	ProcessSchool = iota + 1
	ProcessExpert
	ProcessEdb
	ProcessEnd
)

// 审核状态 1.已提交 ,2.已通过,3.学校未通过,4.专家未通过,5.教育局未通过,6.活动已结束，未审批
const (
	ReviewStatusCommit = iota + 1
	ReviewStatusPass
	ReviewStatusSchoolNoPass
	ReviewStatusExpertNoPass
	ReviewStatusEdbNoPass
	ReviewStatusEnd
)

// 奖项 1：一等奖 2：二等奖 3：三等奖
const (
	PrizeFirst = iota + 1
	PrizeSecond
	PrizeThird
)

// 活动状态  1:待开始，2：进行中，3：已结束
const (
	ActivityWaitStart = iota + 1
	ActivityOngoing
	ActivityEnd
)

// 专家是否授权 1：未授权 2：已授权
const (
	ExportAuthNo = iota + 1
	ExportAuthYes
)
