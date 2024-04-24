package global

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math"
	"strconv"
	"teachers-awards/common/errorz"
	"teachers-awards/common/jwt"
	"teachers-awards/common/util"
	"time"
)

const (
	ctxKey      = "context"
	userInfoKey = "UserInfo"
)

func GetContext(c *gin.Context) context.Context {
	if v, ok := c.Get(ctxKey); ok {
		if ctx, ok := v.(context.Context); ok {
			return ctx
		}
	}
	return context.Background()
}

func SetContext(c *gin.Context, ctx context.Context) {
	c.Set(ctxKey, ctx)
}

func GetUserInfo(ctx context.Context) *UserInfo {
	if userInfo, ok := ctx.Value(userInfoKey).(*UserInfo); ok {
		return userInfo
	}
	return &UserInfo{}
}

func SetUserInfo(ctx context.Context, userinfo *UserInfo) context.Context {
	return context.WithValue(ctx, userInfoKey, userinfo)
}

func GetUserInfoFromClaims(claims *jwt.CustomClaims) *UserInfo {
	var userInfo UserInfo
	json.Unmarshal(claims.DataJson, &userInfo)
	return &userInfo
}

// 记录非空错误
func RecordNotNilError(err error) {
	if err != nil {
		PrintRecordError(err, errorz.GetErrorCallerList(err))
	}
}

func PrintSendRequestError(url, method string, header, body, response interface{}, err error) {
	fields := []zap.Field{
		zap.String("url", url),
		zap.String("method", method),
		zap.Any("header", header),
		zap.Any("body", body),
		zap.Any("response", response),
		zap.String("error", err.Error()),
	}
	Logger.Error("SendRequest", fields...)
}

func PrintRecordError(err error, any interface{}) {
	fields := []zap.Field{
		zap.String("error", err.Error()),
		zap.Any("detail", any),
	}
	Logger.Error("Record Error", fields...)
}

// 角色列表转成数据库需要保存的角色值
func RolesToRoleVal(roles []int) (roleVal int) {
	for _, role := range roles {
		roleVal += int(math.Pow(float64(2), float64(role-1)))
	}
	return roleVal
}

// 数据库需要保存的角色值转成角色列表
func RoleValToRoles(roleVal int) (roles []int) {
	w := 1
	for roleVal > 0 {
		role := roleVal % 2
		if role == 1 {
			roles = append(roles, w)
		}
		roleVal = roleVal / 2
		w++
	}
	return
}

// 中台性转当前性别
func ZtSexToCurSex(ztSexCode string) int {
	if ztSexCode == "gender|1002" {
		return 2
	}
	return 1
}

// 中台证件号转当前身份证号
func ZtCardNumberToIdentityCard(ztCardType, cardNumber string) string {
	if ztCardType == "card_type|1001" {
		return cardNumber
	}
	return ""
}

// 中台生日转当前生日
func ZtBirthdayToCurBirthday(birthday string) string {
	if len(birthday) > 0 {
		num, err := strconv.ParseInt(birthday, 10, 64)
		if err == nil {
			unix := time.Unix(num, 0)
			return unix.Format(util.YMD)
		}
	}
	return birthday
}

func GetZtConfig(from string) (ztConfig ZTConfig, err error) {
	result, err := RedisClient.Get(context.Background(), ZtConfigKey+from).Result()
	if err != nil {
		err = errorz.CodeError(errorz.RESP_ERR, err)
		return
	}
	err = json.Unmarshal([]byte(result), &ztConfig)
	if err != nil {
		err = errorz.CodeError(errorz.ERR_UNMARSHAL, err)
	}
	//if from == "600000" {
	//	ztConfig.ZtDomain = "https://middleground.readboy.com"
	//}
	return
}
