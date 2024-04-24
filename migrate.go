package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"strings"
	"teachers-awards/bao/database/dm"
	"teachers-awards/global"
)

func MigrateUpdate(dbUrl, dbType string, db database.Driver) (ok bool, err error) {
	//获取文件路径
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	wd = strings.ReplaceAll(wd, "\\", "/")
	filePath := "file://" + wd + "/migration/" + dbType

	dbDriver, err := db.Open(dbUrl)
	if err != nil {
		return
	}
	defer dbDriver.Close()

	m, err := migrate.NewWithDatabaseInstance(filePath, dbType, dbDriver)
	if err != nil {
		return
	}
	defer m.Close()

	//更新到最新
	err = m.Up()
	if err != nil {
		return
	}

	return true, nil
}

// dirty = 1 跳过脏文件（无法执行的文件），继续执行下面版本的文件，慎用，仅更新，不做回滚
func main() {
	err := global.InitConfig()
	if err != nil {
		log.Fatalf("读取配置失败，err:%v", err)
	}
	// 读取配置文件
	cfg := global.ServeCfg.TeachersAwards
	if err != nil {
		log.Fatal(err)
	}
	var dbUrl string
	var db database.Driver
	switch cfg.DBType {
	case "dm":
		dbUrl = fmt.Sprintf("%s:%s@%s?schema=%s", cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBName)
		db = &dm.DM{}
	default:
		dbUrl = fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true", cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBName)
		db = &mysql.Mysql{}
	}

	ok, err := MigrateUpdate(dbUrl, cfg.DBType, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ok, "结构迁移完成")
}
