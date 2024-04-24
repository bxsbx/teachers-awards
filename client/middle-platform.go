package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"teachers-awards/global"
	http2 "teachers-awards/http"
	"teachers-awards/model/client"
)

const (
	Authorization = "Authorization"
)

type MiddlePlatformClient struct {
	From   string
	AppCtx context.Context
	Client *http2.Client
}

func NewMiddlePlatformClient(appCtx context.Context) *MiddlePlatformClient {
	userInfo := global.GetUserInfo(appCtx)
	return &MiddlePlatformClient{From: userInfo.From, AppCtx: appCtx, Client: http2.GetClient(http2.MiddlePlatform)}
}

func NewMiddlePlatformClientWithFrom(appCtx context.Context, from string) *MiddlePlatformClient {
	return &MiddlePlatformClient{From: from, AppCtx: appCtx, Client: http2.GetClient(http2.MiddlePlatform)}
}

func (l *MiddlePlatformClient) GetUserInfoByZtToken(ztToken string) (userInfo client.UserInfo, err error) {
	//cfg := global.AccountCenterMap[l.From]
	cfg, err := global.GetZtConfig(l.From)
	if err != nil {
		return
	}
	path := "/service/v1/person/info/person/info"
	method := http.MethodPost

	headerMap := make(map[string]string)
	headerMap[http2.CONTENT_TYPE] = http2.BODY_FORM
	headerMap[Authorization] = global.GetToken(cfg.ZtClientId, cfg.ZtClientSecret, path, method, headerMap, global.VALID_TIME)

	form := url.Values{}
	form.Add("token", ztToken)

	err = l.Client.Request(cfg.ZtDomain+path, method, headerMap, nil, strings.NewReader(form.Encode()), &userInfo, l.AppCtx)
	return
}

func (l *MiddlePlatformClient) GetSchoolListByName(schoolName string) (schoolList []client.School, err error) {
	//cfg := global.AccountCenterMap[l.From]
	cfg, err := global.GetZtConfig(l.From)
	if err != nil {
		return
	}
	path := "/service/v1/school/list/school/fuzzy_search"
	method := http.MethodPost

	headerMap := make(map[string]string)
	headerMap[http2.CONTENT_TYPE] = http2.BODY_JSON
	headerMap[Authorization] = global.GetToken(cfg.ZtClientId, cfg.ZtClientSecret, path, method, headerMap, global.VALID_TIME)

	form := make(map[string]string)
	form["school_name"] = schoolName
	marshal, _ := json.Marshal(form)
	err = l.Client.Request(cfg.ZtDomain+path, method, headerMap, nil, strings.NewReader(string(marshal)), &schoolList, l.AppCtx)
	return
}

func (l *MiddlePlatformClient) GetEnumInfo(enumType string) (enumList []client.Enum, err error) {
	//cfg := global.AccountCenterMap[l.From]
	cfg, err := global.GetZtConfig(l.From)
	if err != nil {
		return
	}
	path := "/service/v1/system/info/enum"
	method := http.MethodPost

	headerMap := make(map[string]string)
	headerMap[http2.CONTENT_TYPE] = http2.BODY_FORM
	headerMap[Authorization] = global.GetToken(cfg.ZtClientId, cfg.ZtClientSecret, path, method, headerMap, global.VALID_TIME)

	form := url.Values{}
	form.Add("enum_type", enumType)

	err = l.Client.Request(cfg.ZtDomain+path, method, headerMap, nil, strings.NewReader(form.Encode()), &enumList, l.AppCtx)
	return
}

func (l *MiddlePlatformClient) GetSchoolListByPage(page, pageSize int) (schoolList client.SchoolList, err error) {
	//cfg := global.AccountCenterMap[l.From]
	cfg, err := global.GetZtConfig(l.From)
	if err != nil {
		return
	}
	path := "/service/v1/school/search"
	method := http.MethodPost

	headerMap := make(map[string]string)
	headerMap[http2.CONTENT_TYPE] = http2.BODY_JSON
	headerMap[Authorization] = global.GetToken(cfg.ZtClientId, cfg.ZtClientSecret, path, method, headerMap, global.VALID_TIME)

	body := make(map[string]string)
	body["page"] = strconv.Itoa(page)
	body["page_size"] = strconv.Itoa(pageSize)
	marshal, _ := json.Marshal(body)
	err = l.Client.Request(cfg.ZtDomain+path, method, headerMap, nil, strings.NewReader(string(marshal)), &schoolList, l.AppCtx)
	return
}

func (l *MiddlePlatformClient) GetTeacherDetailInfoById(teacherId string) (teacherDetailInfo client.TeacherDetailInfo, err error) {
	//cfg := global.AccountCenterMap[l.From]
	cfg, err := global.GetZtConfig(l.From)
	if err != nil {
		return
	}
	path := "/service/v1/person/info/teacher/info"
	method := http.MethodPost

	headerMap := make(map[string]string)
	headerMap[http2.CONTENT_TYPE] = http2.BODY_FORM
	headerMap[Authorization] = global.GetToken(cfg.ZtClientId, cfg.ZtClientSecret, path, method, headerMap, global.VALID_TIME)

	form := url.Values{}
	form.Add("teacher_id", teacherId)

	err = http2.DefaultClient().Request(cfg.ZtDomain+path, method, headerMap, nil, strings.NewReader(form.Encode()), &teacherDetailInfo, l.AppCtx)
	return
}

func (l *MiddlePlatformClient) GetTeachersByName(teacherName string, page, limit int) (teacherDetailInfos []client.TeacherDetailInfo, total int64, err error) {
	//cfg := global.AccountCenterMap[l.From]
	cfg, err := global.GetZtConfig(l.From)
	if err != nil {
		return
	}
	path := "/service/v1/person/info/teacher/search_name"
	method := http.MethodPost

	headerMap := make(map[string]string)
	headerMap[http2.CONTENT_TYPE] = http2.BODY_FORM
	headerMap[Authorization] = global.GetToken(cfg.ZtClientId, cfg.ZtClientSecret, path, method, headerMap, global.VALID_TIME)

	form := url.Values{}
	form.Add("name", teacherName)
	form.Add("page", strconv.Itoa(page))
	form.Add("page_size", strconv.Itoa(limit))

	type temp struct {
		Total int64                      `json:"total"`
		Code  int                        `json:"code"`
		Msg   string                     `json:"msg"`
		Data  []client.TeacherDetailInfo `json:"data"`
	}
	resData := temp{}
	err = http2.DefaultClient().Request(cfg.ZtDomain+path, method, headerMap, nil, strings.NewReader(form.Encode()), &resData, l.AppCtx)
	return resData.Data, resData.Total, err
}
