package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"tiktok_v2/pkg/viper"
)

var (
	_db    *gorm.DB            //后续的代码中并没有使用到 _db 变量本身，所以使用下划线的命名规范
	config = viper.Init("_db") //返回的是config结构体

)

func getDsn(key string) string {
	host := config.Viper.GetString(fmt.Sprintf("%s.host", key))
	port := config.Viper.GetInt(fmt.Sprintf("%s.port", key))
	database := config.Viper.GetString(fmt.Sprintf("%s.database", key))
	username := config.Viper.GetString(fmt.Sprintf("%s.username", key))
	password := config.Viper.GetString(fmt.Sprintf("%s.password", key))
	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	return dsn
}
func init() {
	var err error

	sourceDsn := getDsn("mysql.source")
	_db, err = gorm.Open(mysql.Open(sourceDsn), &gorm.Config{
		//当开启自动预编译时，ORM框架（如GORM）在执行数据库查询之前，会将查询语句进行预先编译，
		//并将编译后的语句缓存起来。这样，在后续的查询中，如果查询语句相同，ORM框架可以直接使用已经编译好的语句，而无需再次解析和编译查询语句。
		//预编译还有助于防止SQL注入攻击。由于查询参数是在编译时绑定到预编译语句中的，而不是在查询时动态拼接SQL字符串
		Logger:                 logger.Default.LogMode(logger.Info),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	//配置db resolver 实现读写分离
	replica1Dsn := getDsn("mysql.replica1")
	replica2Dsn := getDsn("mysql.replica2")
	err = _db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(sourceDsn)},
		Replicas: []gorm.Dialector{mysql.Open(replica1Dsn), mysql.Open(replica2Dsn)},
		Policy:   dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}))
	if err != nil {
		log.Println(err)
	}
	//创建表
	err = _db.AutoMigrate(&Video{}, &User{})
	if err != nil {
		log.Fatal(err)
	}

	db, err := _db.DB()
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)

}

func GetDB() *gorm.DB {
	return _db
}
