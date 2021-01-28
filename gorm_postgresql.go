// +build postgresql

package em

import (
	"github.com/Etpmls/Etpmls-Micro/define"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

const (
	FUZZY_SEARCH = "ILIKE"
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
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=" + timezone

	//Connect Database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: prefix,
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