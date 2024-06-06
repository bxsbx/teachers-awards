package client

import (
	"context"
	"strconv"
	"teachers-awards/global"
	http2 "teachers-awards/http"
	"teachers-awards/model/client"
)

type LoginClient struct {
	AppCtx context.Context
	Client *http2.Client
}

func NewLoginClient(appCtx context.Context) *LoginClient {
	return &LoginClient{AppCtx: appCtx, Client: http2.GetClient(http2.Login)}
}

func (l *LoginClient) CheckAccessToken(userId string, userType, platform int, token string) (check client.CheckAccessToken, err error) {
	cfg := global.ServeCfg.LoginClient
	path := "/v1/auth/access"

	form := make(map[string]string)
	form["userId"] = userId
	form["userType"] = strconv.Itoa(userType)
	form["platform"] = strconv.Itoa(platform)
	form["ACCESSTOKEN"] = token

	err = l.Client.PostForm(cfg.Domain+path, form, &check, l.AppCtx)
	return
}
