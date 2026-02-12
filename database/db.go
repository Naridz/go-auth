package database

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

func ConnectDB() {
	dsn := os.Getenv("DB_NAME")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := Db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}
