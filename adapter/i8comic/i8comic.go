package i8comic

import (
	"comic2video/adapter"
	"comic2video/config"
	"comic2video/ffmpeg"
	"fmt"
	"os"
	"path/filepath"
)

// DownloadChapter 下载并保存章节内图片，返回下载内容的本地资源路径列表
func DownloadChapter(aid string) ([]string, error) {
	url := fmt.Sprintf("https://18comic.vip/photo/%s", aid)
	urls, err := GetChapterContent(url)
	if err != nil {
		return nil, err
	}

	filePaths, err := adapter.DownloadImages(urls, func(i int, url string) (string, error) {
		return downloadImage(url, aid)
	})
	if err != nil {
		return nil, err
	}
	return filePaths, nil
}

func DownloadAndConvert(aid string) (string, error) {
	imgPaths, err := DownloadChapter(aid)
	if err != nil {
		return "", err
	}

	os.MkdirAll(config.I8_VIDEO_PATH, 0755)
	path := filepath.Join(config.I8_VIDEO_PATH, fmt.Sprintf("%s.mp4", aid))

	err = ffmpeg.GenerateVideo(imgPaths, path)
	if err != nil {
		return "", err
	}

	return path, nil
}
