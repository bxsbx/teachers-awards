package resp

import (
	"teachers-awards/dao"
)

type GetUploadTokenResp struct {
	Token string `json:"token"`
}

type GetOperationRecordFromUaiResp struct {
	List []dao.OperationRecord `json:"list"`
}

type Subject struct {
	SubjectCode string `json:"subject_code"`
	SubjectName string `json:"subject_name"`
}

type GetSubjectEnumResp struct {
	SubjectList []Subject `json:"subject_list"`
}

type School struct {
	SchoolId   string `json:"school_id"`
	SchoolName string `json:"school_name"`
}

type GetSchoolListByPageResp struct {
	Total      int      `json:"total"`
	SchoolList []School `json:"school_list"`
}

type GetSchoolListByNameResp struct {
	SchoolList []School `json:"school_list"`
}
