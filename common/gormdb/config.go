package gormdb

type GormConfig struct {
	Host            string
	UserName        string
	Password        string
	DbName          string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime int
	ConnMaxIdleTime int
	DBLog           bool
}
