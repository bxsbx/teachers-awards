package global

var DeclareNameMap = map[int]string{
	1: "幼教",
	2: "小学",
	3: "初中",
	4: "高中",
	5: "职校",
	6: "教研",
}

var SexNameMap = map[int]string{
	1: "男",
	2: "女",
}

var RankPrizeNameMap = map[int]string{
	1: "一等奖",
	2: "二等奖",
	3: "三等奖",
}

var StatusNameMap = map[int]string{
	1: "待审核",
	2: "已通过",
	3: "学校未通过",
	4: "专家未通过",
	5: "未通过",
	6: "结果已评定，未审核",
}
var ReviewStatusNameMap = map[int]string{
	1: "待审核",
	2: "未通过",
	3: "通过",
}

var PassNameMap = map[int]string{
	0: "不通过",
	1: "通过",
}

var RoleNameMap = map[int]string{
	1: "学校",
	2: "专家",
	3: "教育局",
	4: "教师",
	5: "超级管理员",
}

var RouterAuthMap = map[string][]int{
	"/v1/announcement/save":                {3},
	"/v1/announcement/delete":              {3},
	"/v1/indicator/one/save":               {3},
	"/v1/indicator/one/delete/ids":         {3},
	"/v1/indicator/two/save":               {3},
	"/v1/indicator/two/delete/ids":         {3},
	"/v1/activity/save":                    {3},
	"/v1/activity/one/indicator/delete":    {3},
	"/v1/activity/two/indicator/delete":    {3},
	"/v1/activity/delete":                  {3},
	"/v1/user/activity/declare/create":     {4},
	"/v1/user/activity/indicator/update":   {4},
	"/v1/user/activity/indicator/delete":   {4},
	"/v1/activity/review/commit":           {1, 2, 3},
	"/v1/activity/review/set/awards":       {3},
	"/v1/activity/review/commit/result":    {3},
	"/v1/activity/review/update/two/id":    {3},
	"/v1/activity/review/edb/declare/user": {3},
	"/v1/activity/review/batch/pass":       {1, 2, 3},
	"/v1/user/set/role":                    {5},
	"/v1/user/set/expert/auth":             {3},
	"/v1/user/cancel/expert/auth":          {3},
}
