// +build mysql

package em

import (
	library "github.com/Etpmls/Etpmls-Micro/library"
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
	dsn := library.Config.Database.User + ":" + library.Config.Database.Password + "@tcp(" + library.Config.Database.Host + ":" + library.Config.Database.Port + ")/" + library.Config.Database.Name + "?charset=utf8mb4&parseTime=True&loc=" + url.QueryEscape(library.Config.Database.TimeZone)

	//Connect Database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   library.Config.Database.Prefix,
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