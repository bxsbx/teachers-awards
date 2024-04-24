package global

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"io"
	"os"
	"teachers-awards/common/gormdb"
	"teachers-awards/common/lock"
	"teachers-awards/common/tracer"
	"teachers-awards/common/zaplog"
	"time"
)

// 初始化全局变量
func InitGlobalVal() {
	//日志初始化
	Logger = zaplog.InitZap(zap.InfoLevel, os.Stdout)

	//初始化mysql数据库
	InitGorm()
	//初始化redis
	InitRedis()
	//InitAccountCenterConfig()
	//初始化变量
	SubjectEnumMap = &SubjectEnum{SubjectMap: nil}
	FuncMapLock = lock.NewFuncLockMap()

}

func InitAccountCenterConfig() {
	cfg := ServeCfg.AccountCenterConfig
	AccountCenterMap = make(map[string]AccountCenter)
	AccountCenterMap["600000"] = AccountCenter{
		ZtDomain:       cfg.ZtDomain,
		ZtClientId:     cfg.ZtClientId,
		ZtClientSecret: cfg.ZtClientSecret,
	}
	AccountCenterMap["600001"] = AccountCenter{
		ZtDomain:       cfg.ZtDomain2,
		ZtClientId:     cfg.ZtClientId2,
		ZtClientSecret: cfg.ZtClientSecret2,
	}
}

func InitGorm() {
	cfg := ServeCfg.TeachersAwards
	dbConfig := gormdb.GormConfig{
		Host:            cfg.DBHost,
		UserName:        cfg.DBUsername,
		Password:        cfg.DBPassword,
		DbName:          cfg.DBName,
		MaxOpenConn:     cfg.DBMaxOpenConn,
		MaxIdleConn:     cfg.DBMaxIdleConn,
		ConnMaxLifetime: 300,
		ConnMaxIdleTime: 300,
		DBLog:           true,
	}
	switch cfg.DBType {
	case "dm":
		GormDB = gormdb.NewDmDB(dbConfig)
	default:
		GormDB = gormdb.NewMysqlDB(dbConfig)
	}
}

func InitRedis() {
	cfg := ServeCfg.TeachersAwards
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisHost,
		DB:           5,
		Password:     cfg.RedisPassword,
		PoolSize:     cfg.RedisPoolSize,
		MinIdleConns: 2,
	})
	RedisClient = client
}

func InitTract() io.Closer {
	cfg := ServeCfg
	tracerConfig := tracer.TracerConfig{
		IsOpenTracing:       cfg.TeachersAwards.TraceOpen,
		ServiceName:         cfg.TeachersAwards.TraceServerName,
		CollectorEndpoint:   cfg.Trace.Domain,
		SamplerType:         "const",
		SamplerParam:        1,
		LogSpans:            true,
		BufferFlushInterval: 1 * time.Second,
	}
	return tracer.NewTracer(tracerConfig)
}
