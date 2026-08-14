package config

import (
	"comic_converter/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open(DB_FILE_PATH), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.ComicCache{})

	DB = db
}
