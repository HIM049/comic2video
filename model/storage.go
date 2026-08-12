package model

import "gorm.io/gorm"

type Comic struct {
	gorm.Model
	Source    string `json:"source" gorm:"uniqueIndex:idx_comic_unique,priority:1"`
	SeriesId  string `json:"series_id" gorm:"uniqueIndex:idx_comic_unique,priority:2"`
	ChapterId string `json:"chapter_id" gorm:"uniqueIndex:idx_comic_unique,priority:3"`
}
