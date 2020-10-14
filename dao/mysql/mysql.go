package mysql

import (
	"fmt"
	"netdisk/model"

	"gorm.io/driver/mysql"

	"netdisk/setting"

	"gorm.io/gorm"
)

var Db *gorm.DB

//var SqlDb ,_=Db.DB()

// Init 初始化MySQL连接
func Init(cfg *setting.MySQLConfig) (err error) {
	// "user:password@tcp(host:port)/dbname"
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return
	}
	// 迁移 schema
	Db.AutoMigrate(&model.User{}, &model.FileInfo{})

	//// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	//db.SetMaxIdleConns(10)
	//
	//// SetMaxOpenConns 设置打开数据库连接的最大数量。
	//sqlDB.SetMaxOpenConns(100)
	//
	//// SetConnMaxLifetime 设置了连接可复用的最大时间。
	//sqlDB.SetConnMaxLifetime(time.Hour)

	//SqlDb.Close()
	return
}

//// Close 关闭MySQL连接
//func Close() {
//	_ = DB.DB().Close()
//}
