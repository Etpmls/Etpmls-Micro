// +build mysql

package em

import (
	"github.com/Etpmls/Etpmls-Micro/define"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/url"
)

var DB *gorm.DB

const (
	FUZZY_SEARCH = "LIKE"
)

func (this *Register) RunDatabase() {
	m, err2 := Kv.List(define.MakeServiceConfField(em_library.Config.Service.RpcName, define.KvServiceDatabase))
	if err2 != nil {
		LogInfo.OutputSimplePath(err2)
		return
	}

	var (
		host = define.MakeServiceConfField(em_library.Config.Service.RpcName, define.KvServiceDatabaseHost)
		user = define.MakeServiceConfField(em_library.Config.Service.RpcName, define.KvServiceDatabaseUser)
		password = define.MakeServiceConfField(em_library.Config.Service.RpcName, define.KvServiceDatabasePassword)
		port = define.MakeServiceConfField(em_library.Config.Service.RpcName, define.KvServiceDatabasePort)
		dbname = define.MakeServiceConfField(em_library.Config.Service.RpcName, define.KvServiceDatabaseDbName)
		timezone = define.MakeServiceConfField(em_library.Config.Service.RpcName, define.KvServiceDatabaseTimezone)
		prefix = define.MakeServiceConfField(em_library.Config.Service.RpcName, define.KvServiceDatabasePrefix)
	)
	dsn := this.panicIfMapValueEmpty(user, m) + ":" + this.panicIfMapValueEmpty(password, m) + "@tcp(" + this.panicIfMapValueEmpty(host, m) + ":" + this.panicIfMapValueEmpty(port, m) + ")/" + this.panicIfMapValueEmpty(dbname, m) + "?charset=utf8mb4&parseTime=True&loc=" + url.QueryEscape(this.panicIfMapValueEmpty(timezone, m))

	//Connect Database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   this.panicIfMapValueEmpty(prefix, m),
		},
	})
	if err != nil {
		LogPanic.AutoOutputDebug("Unable to connect to the database!", err)
	}

	err = DB.AutoMigrate(this.DatabaseMigrate...)
	if err != nil {
		LogInfo.AutoOutputDebug("Failed to create database!", err)
	}

}