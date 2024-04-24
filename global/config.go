package global

import (
	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
)

type ServeConfig struct {
	Default             Default             `mapstructure:"default"`
	TeachersAwards      TeachersAwards      `mapstructure:"teachers-awards"`
	Trace               Trace               `mapstructure:"trace"`
	QiNiu               QiNiu               `mapstructure:"dreamEbagQiniu"`
	AccountCenterConfig AccountCenterConfig `mapstructure:"accountCenter"`
	LoginClient         LoginClient         `mapstructure:"dreamEbagLogin"`
}

type Default struct {
	AppName string `mapstructure:"appName"`
	AppPort int    `mapstructure:"appPort"`
}

type TeachersAwards struct {
	Domain          string `mapstructure:"domain"`
	DBHost          string `mapstructure:"dBHost"`
	DBUsername      string `mapstructure:"dBUsername"`
	DBPassword      string `mapstructure:"dBPassword"`
	DBName          string `mapstructure:"dBName"`
	DBMaxOpenConn   int    `mapstructure:"dbMaxOpenConn"`
	DBMaxIdleConn   int    `mapstructure:"dbMaxIdleConn"`
	DBType          string `mapstructure:"dbType"`
	TraceServerName string `mapstructure:"traceServerName"`
	TraceOpen       bool   `mapstructure:"traceOpen"`
	RedisHost       string `mapstructure:"redisHost"`
	RedisPassword   string `mapstructure:"redisPassword"`
	RedisPoolSize   int    `mapstructure:"redisPoolSize"`
}

type Trace struct {
	Domain string `mapstructure:"domain"`
}

type QiNiu struct {
	AccessKey string `mapstructure:"AccessKey"`
	SecretKey string `mapstructure:"SecretKey"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"`
}

type AccountCenterConfig struct {
	ZtDomain2       string `mapstructure:"apiDomainV2"`
	ZtClientId2     string `mapstructure:"clientIdV2"`
	ZtClientSecret2 string `mapstructure:"clientSecretV2"`
	ZtDomain        string `mapstructure:"pubApiDomain"`
	ZtClientId      string `mapstructure:"clientId"`
	ZtClientSecret  string `mapstructure:"clientSecret"`
}

type LoginClient struct {
	Domain   string `mapstructure:"domain-local"`
	ApiToken string `mapstructure:"apiToken"`
}

func InitConfig() error {
	v := viper.NewWithOptions(viper.IniLoadOptions(ini.LoadOptions{SpaceBeforeInlineComment: true}))
	v.SetConfigType("ini")
	v.SetConfigFile("conf/app.conf")
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	v.SetConfigFile("conf/dev/app.conf")
	err = v.MergeInConfig()
	if err != nil {
		return err
	}
	err = v.Unmarshal(&ServeCfg)
	return err
}
