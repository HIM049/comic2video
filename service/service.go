package service

import (
	"comic2video/config"
	"comic2video/model"
	"context"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var queue *Queue

func Init() {
	queue = NewQueue()
}

// GetOrCreateVideo loads a video from cache or queues a download task.
func GetOrCreateVideo(ctx context.Context, source model.ComicSource, comicId, comicUrl string, download func() (string, error)) (string, error) {
	var comic model.ComicCache
	err := config.DB.
		Where("source = ?", source).
		Where("comic_id = ?", comicId).
		First(&comic).
		Error
	if err == nil {
		return comic.VideoPath, nil
	}
	if err != gorm.ErrRecordNotFound {
		return "", err
	}

	// Submit task if not cached
	log.Info("Resource not cached, queuing...")
	err = queue.Submit(ctx, TaskKey{Source: string(source), ComicId: comicId}, func() error {
		videoPath, err := download()
		if err != nil {
			return err
		}

		newComic := model.ComicCache{
			Uid:       uuid.NewString(),
			Source:    source,
			ComicId:   comicId,
			ComicUrl:  comicUrl,
			VideoPath: videoPath,
		}
		return config.DB.Create(&newComic).Error
	})
	if err != nil {
		return "", err
	}

	// Load from DB after task completes
	err = config.DB.
		Where("source = ?", source).
		Where("comic_id = ?", comicId).
		First(&comic).
		Error
	if err != nil {
		return "", err
	}
	return comic.VideoPath, nil
}
