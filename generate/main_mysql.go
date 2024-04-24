package main

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"teachers-awards/common/gormdb"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"text/template"
)

var DataTypeMap = map[string]string{
	"varchar":  "string",
	"text":     "string",
	"tinyint":  "int",
	"int":      "int",
	"bigint":   "int64",
	"decimal":  "float64",
	"double":   "float64",
	"float":    "float",
	"datetime": "*time.Time",
	"date":     "*time.Time",
}

func InitDB() {
	err := global.InitConfig()
	if err != nil {
		log.Fatalf("读取配置失败，err:%v", err)
	}

	cfg := global.ServeCfg.TeachersAwards
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
	global.GormDB = gormdb.NewMysqlDB(dbConfig)
}

const (
	ProjectNameMysql = "teachers-awards"
	MysqlBasePath    = "H:\\dream\\teachers-awards\\"
	TemplatePath     = MysqlBasePath + "generate\\template\\mysql.tmpl"
	DaoPath          = MysqlBasePath + "dao\\"
)

func main() {

	tableName := "expert_auth_indicator"
	pre := "t_"
	DeletePreTableName := strings.TrimPrefix(tableName, pre)
	filePath := DaoPath + DeletePreTableName + ".go"

	mysqlTemplate := MysqlTemplate{
		ProjectName:    ProjectNameMysql,
		ConstTableName: strings.ToUpper(DeletePreTableName),
		TableName:      tableName,
		UpperTableName: util.HumpNaming(DeletePreTableName),
	}
	mysqlTemplate.LowerTableName = util.FirstLower(mysqlTemplate.UpperTableName)

	InitDB()
	columns, err := GetTableInfo(global.GormDB, tableName)
	if err != nil {
		log.Fatal(err)
	}
	var whereList []string
	var nameTypeList []string
	var nameList []string
	for i, column := range columns {
		column.FieldName = util.HumpNaming(column.Name)
		column.Type = DataTypeMap[column.Type]
		name := util.FirstLower(column.FieldName)
		if column.PrimaryKey == "PRI" {
			nameTypeList = append(nameTypeList, name+" "+column.Type)
			nameList = append(nameList, name)
			whereList = append(whereList, name+" = ?")
		}
		columns[i] = column
	}
	mysqlTemplate.PrimaryWhere = "\"" + strings.Join(whereList, " and ") + "\", " + strings.Join(nameList, ",")
	mysqlTemplate.PrimaryParams = strings.Join(nameTypeList, ",")
	mysqlTemplate.Columns = columns

	exist, err := util.FileIsExist(filePath)
	if err != nil {
		log.Fatal(err)
	}
	if !exist {
		t := template.Must(template.ParseFiles(TemplatePath))
		var builder strings.Builder
		err = t.Execute(&builder, mysqlTemplate)
		if err != nil {
			log.Fatal(err)
		}
		err = util.WriteToFile(filePath, builder.String())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("创建成功")
	} else {
		fmt.Println("文件已存在")
	}
}

type MysqlTemplate struct {
	ProjectName    string
	ConstTableName string
	TableName      string
	UpperTableName string
	LowerTableName string
	PrimaryParams  string
	PrimaryWhere   string
	Columns        []TableColumn
}

type TableColumn struct {
	FieldName  string `gorm:"-"`
	Name       string `gorm:"column:name"`
	Type       string `gorm:"column:type"`
	Comment    string `gorm:"column:comment"`
	IsNull     string `gorm:"column:is_null"`
	PrimaryKey string `gorm:"column:primary_key"`
}

func GetTableInfo(db *gorm.DB, tableName string) (columns []TableColumn, err error) {
	query := fmt.Sprintf(`SELECT COLUMN_NAME AS name, DATA_TYPE AS type, COLUMN_COMMENT AS comment, IS_NULLABLE as is_null, COLUMN_KEY as primary_key 
					FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%s' 
					ORDER BY ORDINAL_POSITION`, tableName)
	err = db.Raw(query).Scan(&columns).Error
	return
}
