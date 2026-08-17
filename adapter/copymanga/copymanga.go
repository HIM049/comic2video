package copymanga

import (
	"comic2video/adapter"
	"comic2video/config"
	"comic2video/ffmpeg"
	"fmt"
	"os"
	"path/filepath"
)

func DownloadChapter(seriesId, chapterId string) ([]string, error) {

	chapterUrl := fmt.Sprintf("https://www.mangacopy.com/comic/%s", seriesId)
	fullUrl := fmt.Sprintf("%s/chapter/%s", chapterUrl, chapterId)

	comicPath := filepath.Join(config.CPCOMIC_PATH, seriesId, chapterId)
	os.MkdirAll(comicPath, 0755)

	urls, err := GetChapterContent(chapterUrl, fullUrl)
	if err != nil {
		return nil, err
	}

	filePaths, err := adapter.DownloadImages(urls, func(i int, url string) (string, error) {
		path := filepath.Join(comicPath, fmt.Sprintf("img-%d.jpg", i))
		if err := downloadImage(chapterUrl, path, url); err != nil {
			return "", err
		}
		return path, nil
	})
	if err != nil {
		return nil, err
	}

	return filePaths, nil
}

func DownloadAndConvert(seriesId, chapterId string) (string, error) {
	filePaths, err := DownloadChapter(seriesId, chapterId)
	if err != nil {
		return "", err
	}

	videoPath := filepath.Join(config.CP_VIDEO_PATH, seriesId)
	os.MkdirAll(videoPath, 0755)
	path := filepath.Join(videoPath, fmt.Sprintf("%s.mp4", chapterId))

	err = ffmpeg.GenerateVideo(filePaths, path)
	if err != nil {
		return "", err
	}

	return path, nil
}
