// +build postgresql

package em

import (
	library "github.com/Etpmls/Etpmls-Micro/library"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

const (
	FUZZY_SEARCH = "ILIKE"
)

func (this *Register) RunDatabase() {
	dsn := "host=" + library.Config.Database.Host + " user=" + library.Config.Database.User + " password=" + library.Config.Database.Password + " dbname=" + library.Config.Database.Name + " port=" + library.Config.Database.Port + " sslmode=disable TimeZone=" + library.Config.Database.TimeZone

	//Connect Database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: library.Config.Database.Prefix,
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