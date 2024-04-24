package client

type Role struct {
	Id       int    `json:"id"`
	RoleType int    `json:"role_type"`
	RoleName string `json:"role_name"`
	Remark   string `json:"remark"`
	State    int    `json:"state"`
}

type UserInfo struct {
	PersonId  string `json:"person_id"`
	Username  string `json:"username"`
	SexCode   string `json:"sex_code"`
	SexName   string `json:"sex_name"`
	Phone     string `json:"phone"`
	IdNumber  string `json:"id_number"`
	Avatar    string `json:"avatar"`
	Birthtime int64  `json:"birthtime"`
	RoleList  []Role `json:"role_list"`
}

type School struct {
	SchoolId   string `json:"schoolId"`
	SchoolName string `json:"schoolName"`
}

type SchoolList struct {
	Total int `json:"total"`
	Data  []struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Stage string `json:"stage"`
	} `json:"data"`
}

type Enum struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Num  int    `json:"num"`
}

type TeacherDetailInfo struct {
	PersonId     string `json:"person_id"`
	Username     string `json:"username"`
	SexCode      string `json:"sex_code"`
	SexName      string `json:"sex_name"`
	Phone        string `json:"phone"`
	CardTypeCode string `json:"card_type_code"`
	CardTypeName string `json:"card_type_name"`
	CardNumber   string `json:"card_number"`
	CountryCode  string `json:"country_code"`
	CountryName  string `json:"country_name"`
	Address      string `json:"address"`
	Birthday     string `json:"birthday"`
	Email        string `json:"email"`
	OrganId      string `json:"organ_id"`
	OrganName    string `json:"organ_name"`
	Avatar       string `json:"avatar"`
}
