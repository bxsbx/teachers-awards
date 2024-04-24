package resp

import (
	"teachers-awards/dao"
)

type GetAnnouncementByIdResp struct {
	dao.Announcement
}

type SaveAnnouncementResp struct {
	AnnouncementId int `json:"announcement_id"` // 公告id
}

type Announcement struct {
	dao.Announcement
	IsRead bool `json:"is_read"`
}

type GetAnnouncementListResp struct {
	Total int64          `json:"total"`
	List  []Announcement `json:"list"`
}
