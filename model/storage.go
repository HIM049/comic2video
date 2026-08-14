package model

import "gorm.io/gorm"

// type Comic struct {
// 	gorm.Model
// 	Source    string `json:"source" gorm:"uniqueIndex:idx_comic_unique,priority:1"`
// 	SeriesId  string `json:"series_id" gorm:"uniqueIndex:idx_comic_unique,priority:2"`
// 	ChapterId string `json:"chapter_id" gorm:"uniqueIndex:idx_comic_unique,priority:3"`
// }

type ComicSource string

const (
	SourceCopymanga ComicSource = "copymanga"
	SourceI8comic   ComicSource = "18comic"
)

type ComicCache struct {
	gorm.Model
	Uid       string      `json:"uid" gorm:"not null;uniqueIndex"`
	Source    ComicSource `json:"source" gorm:"uniqueIndex:idx_comic_unique,priority:1"`
	ComicId   string      `json:"comic_id" gorm:"uniqueIndex:idx_comic_unique,priority:2"`
	ComicUrl  string      `json:"comic_url" gorm:"not null"`
	VideoPath string      `json:"video_path" gorm:"not null"`
}
