// +build mysql

package em

import (
	"github.com/Etpmls/Etpmls-Micro/v2/define"
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
	var (
		host = MustGetServiceNameKvKey(define.KvServiceDatabaseHost)
		user = MustGetServiceNameKvKey(define.KvServiceDatabaseUser)
		password = MustGetServiceNameKvKey(define.KvServiceDatabasePassword)
		port = MustGetServiceNameKvKey(define.KvServiceDatabasePort)
		dbname = MustGetServiceNameKvKey(define.KvServiceDatabaseDbName)
		timezone = MustGetServiceNameKvKey(define.KvServiceDatabaseTimezone)
		prefix = MustGetServiceNameKvKey(define.KvServiceDatabasePrefix)
	)

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=" + url.QueryEscape(timezone)

	//Connect Database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
		},
	})
	if err != nil {
		LogPanic.AutoOutputDebug("Unable to connect to the database!", err)
	}

	err = DB.AutoMigrate(this.DatabaseMigrate...)
	if err != nil {
		LogInfo.OutputSimplePath("Failed to create database!", err)
	}

}