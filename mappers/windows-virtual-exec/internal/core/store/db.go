package store

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(sqliteFile string) {
	db, err := gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
}
