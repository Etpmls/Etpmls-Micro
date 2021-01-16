// +build postgresql

package em

import (
	"github.com/Etpmls/Etpmls-Micro/define"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

const (
	FUZZY_SEARCH = "ILIKE"
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
	dsn := "host=" + this.panicIfMapValueEmpty(host, m) + " user=" + this.panicIfMapValueEmpty(user, m) + " password=" + this.panicIfMapValueEmpty(password, m) + " dbname=" + this.panicIfMapValueEmpty(dbname, m) + " port=" + this.panicIfMapValueEmpty(port, m) + " sslmode=disable TimeZone=" + this.panicIfMapValueEmpty(timezone, m)

	//Connect Database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: this.panicIfMapValueEmpty(prefix, m),
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