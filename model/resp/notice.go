package resp

import "teachers-awards/dao"

type GetUserNoticeListResp struct {
	Total int64        `json:"total"`
	List  []dao.Notice `json:"list"`
}
