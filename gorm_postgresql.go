// +build postgresql

package em

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

const (
	FUZZY_SEARCH = "ILIKE"
)

func (this *Register) RunDatabase() {
	m, err2 := Kv.List(KvDatabase)
	if err2 != nil {
		LogInfo.OutputSimplePath(err2)
		return
	}

	dsn := "host=" + m[KvDatabaseHost] + " user=" + m[KvDatabaseUser] + " password=" + m[KvDatabasePassword] + " dbname=" + m[KvDatabaseDbName] + " port=" + m[KvDatabasePort] + " sslmode=disable TimeZone=" + m[KvDatabaseTimezone]

	//Connect Database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: m[KvDatabasePrefix],
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