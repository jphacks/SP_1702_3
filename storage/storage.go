package storage

import (
	"fmt"

	"github.com/labstack/gommon/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgresdriver
	"github.com/youtangai/Otomo_backend/config"
	"github.com/youtangai/Otomo_backend/model"
)

var (
	//Storage is storage context
	Storage *gorm.DB
)

func init() {
	var err error
	Storage, err = gorm.Open("postgres", ConnectionString())
	if err != nil {
		log.Print(err)
	}
	if !Storage.HasTable(&model.IndoorInfo{}) {
		Storage.AutoMigrate(&model.IndoorInfo{})
	}
	if !Storage.HasTable(&model.Soul{}) {
		Storage.AutoMigrate(&model.Soul{})
		Storage.Create(&model.Soul{UserID: 1, OnDevice: true})
	}
}

//GetDBContext is return DBcontext
func GetDBContext() *gorm.DB {
	return Storage
}

//ConnectionString is get connection string
func ConnectionString() string {
	user := config.DBUser()
	host := config.DBHost()
	dbname := config.DBName()
	passwd := config.DBPasswd()
	connectString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, passwd)
	return connectString
}
