package global

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"teachers-awards/common/lock"
)

type UserInfo struct {
	From      string `json:"from"` //用户数据来源，600000：中台1.0，6000001：中台2.0
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserRoles []int  `json:"user_roles"`
}

type AccountCenter struct {
	ZtDomain       string
	ZtClientId     string
	ZtClientSecret string
}

type ZTConfig struct {
	ZtDomain          string `json:"apiDomain"`
	ZtClientId        string `json:"clientId"`
	ZtClientSecret    string `json:"clientSecret"`
	LoginClientId     string `json:"loginClientId"`
	LoginClientSecret string `json:"loginClientSecret"`
	LoginDomain       string `json:"loginDomain"`
	PlatfromCode      string `json:"platfromCode"`
}

type SubjectEnum struct {
	sync.Mutex
	SubjectMap map[string]string
}

var (
	ServeCfg         ServeConfig
	GormDB           *gorm.DB
	Logger           *zap.Logger
	RedisClient      *redis.Client
	SubjectEnumMap   *SubjectEnum
	AccountCenterMap map[string]AccountCenter
	FuncMapLock      *lock.FuncLockMap
)
