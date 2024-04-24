package gormdb

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"teachers-awards/common/gormdb/dm"
	"time"
)

func NewMysqlDB(cfg GormConfig) *gorm.DB {
	openUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=Asia%sShanghai&parseTime=true", cfg.UserName, cfg.Password, cfg.Host, cfg.DbName, "%2F")
	db, err := gorm.Open(mysql.Open(openUrl), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库初始化失败, err:%v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)
	return db
}

func NewDmDB(cfg GormConfig) *gorm.DB {
	openUrl := fmt.Sprintf("dm://%s:%s@%s?schema=%s&compatibleMode=Mysql", cfg.UserName, cfg.Password, cfg.Host, cfg.DbName)
	db, err := gorm.Open(dm.Open(openUrl), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库初始化失败, err:%v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)

	return db
}
