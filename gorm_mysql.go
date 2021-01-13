// +build mysql

package em

import (
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
	m, err2 := Kv.List(KvDatabase)
	if err2 != nil {
		LogInfo.OutputSimplePath(err2)
		return
	}

	dsn := m[KvDatabaseUser] + ":" + m[KvDatabasePassword] + "@tcp(" + m[KvDatabaseHost] + ":" + m[KvDatabasePort] + ")/" + m[KvDatabaseDbName] + "?charset=utf8mb4&parseTime=True&loc=" + url.QueryEscape(m[KvDatabaseTimezone])

	//Connect Database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   m[KvDatabasePrefix],
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