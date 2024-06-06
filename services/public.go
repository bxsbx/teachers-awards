package services

import (
	"context"
	"encoding/json"
	"teachers-awards/client"
	"teachers-awards/common/errorz"
	"teachers-awards/common/jwt"
	"teachers-awards/common/util"
	"teachers-awards/dao"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/model/resp"
	"time"
)

type PublicService struct {
	appCtx context.Context
}

func NewPublicService(appCtx context.Context) *PublicService {
	return &PublicService{appCtx: appCtx}
}

func (s *PublicService) GetInfoByUserId(from string, userId string) (data *resp.GetTokenAndUserInfoResp, err error) {
	platformClient := client.NewMiddlePlatformClientWithFrom(s.appCtx, from)
	teacherDetailInfo, err := platformClient.GetTeacherDetailInfoById(userId)
	if err != nil {
		return nil, err
	}

	data = &resp.GetTokenAndUserInfoResp{}
	userInfo := global.UserInfo{From: from}
	userInfoDao := dao.NewUserInfoDao(s.appCtx)
	var user dao.UserInfo
	userInfoDao.First(dao.UserInfo{UserId: teacherDetailInfo.PersonId}, &user)
	if user.UserId != "" {
		userInfo.UserRoles = global.RoleValToRoles(user.Role)
		updateMap := make(map[string]interface{})
		card := global.ZtCardNumberToIdentityCard(teacherDetailInfo.CardTypeCode, teacherDetailInfo.CardNumber)
		if card != user.IdentityCard {
			updateMap["identity_card"] = card
		}
		if teacherDetailInfo.Phone != user.Phone {
			updateMap["phone"] = teacherDetailInfo.Phone
		}
		if teacherDetailInfo.OrganId != user.SchoolId || teacherDetailInfo.OrganName != user.SchoolName {
			updateMap["school_id"] = teacherDetailInfo.OrganId
			updateMap["school_name"] = teacherDetailInfo.OrganName
		}
		if len(updateMap) > 0 {
			err := userInfoDao.UpdateByWhere(dao.UserInfo{UserId: teacherDetailInfo.PersonId}, updateMap)
			if err != nil {
				return nil, err
			}
		}
	} else {
		user.UserId = teacherDetailInfo.PersonId
		user.UserName = teacherDetailInfo.Username
		user.UserSex = global.ZtSexToCurSex(teacherDetailInfo.SexCode)
		user.Phone = teacherDetailInfo.Phone
		user.Birthday = global.ZtBirthdayToCurBirthday(teacherDetailInfo.Birthday)
		user.SchoolId = teacherDetailInfo.OrganId
		user.SchoolName = teacherDetailInfo.OrganName
		user.Avatar = teacherDetailInfo.Avatar
		user.IdentityCard = global.ZtCardNumberToIdentityCard(teacherDetailInfo.CardTypeCode, teacherDetailInfo.CardNumber)
		user.Role = global.Role2Teacher
		nowTime := util.NowTime()
		user.CreateTime = &nowTime
		user.Year = nowTime.Year()
		err = userInfoDao.Create(&user)
		if err != nil {
			return nil, err
		}
		userInfo.UserRoles = []int{global.RoleTeacher}
	}
	userInfo.UserId = user.UserId
	userInfo.UserName = user.UserName

	claims := jwt.CustomClaims{}
	claims.ExpiresAt = time.Now().Unix() + global.Jwt.ExpiresTime // 设置过期时间
	claims.NotBefore = time.Now().Unix() - 1000                   // 签名生效时间
	claims.Issuer = "teachers-awards"                             //签名的发行者
	claims.DataJson, _ = json.Marshal(userInfo)

	data.UserId = userInfo.UserId
	data.UserName = userInfo.UserName
	data.UserRoles = userInfo.UserRoles
	data.Token, err = global.Jwt.CreateToken(claims)
	if err != nil {
		return
	}
	data.ExpiresTimeAt = claims.ExpiresAt
	//err = global.RedisClient.Set(s.appCtx, global.JwtKey+userInfo.UserId, data.Token, global.Jwt.ExpiresTime).Err()
	return
}

func (s *PublicService) GetTokenAndUserInfo(params *req.GetTokenAndUserInfoReq) (data *resp.GetTokenAndUserInfoResp, err error) {
	data = &resp.GetTokenAndUserInfoResp{}
	//switch params.UserId {
	//case "111":
	//	userInfo.UserRoles = []int{1}
	//case "221", "222":
	//	userInfo.UserRoles = []int{2}
	//case "331":
	//	userInfo.UserRoles = []int{3}
	//default:
	//	userInfo.UserRoles = []int{4}
	//}
	//userInfo.From = "600000"
	var from string
	if len(params.JXYToken) > 6 {
		from = params.JXYToken[:6]
	}

	loginClient := client.NewLoginClient(s.appCtx)
	check, err := loginClient.CheckAccessToken(params.UserId, 2, params.Platform, params.JXYToken)
	if err != nil {
		return nil, err
	}
	if check.F_responseNo != 10000 {
		return nil, errorz.Code(check.F_responseNo)
	}
	return s.GetInfoByUserId(from, params.UserId)
}

func (s *PublicService) GetTokenAndUserInfoByCode(params *req.GetTokenAndUserInfoByCodeReq) (data *resp.GetTokenAndUserInfoResp, err error) {
	data = &resp.GetTokenAndUserInfoResp{}
	platformClient := client.NewMiddlePlatformClientWithFrom(s.appCtx, params.From)
	ztTokenInfo, err := platformClient.GetTokenByCode(params.Code, params.RedirectUrl)
	if err != nil {
		return nil, err
	}
	tokenInfo, err := platformClient.VerifyToken(ztTokenInfo.AccessToken)
	if err != nil {
		return nil, err
	}
	return s.GetInfoByUserId(params.From, tokenInfo.PersonId)
}
