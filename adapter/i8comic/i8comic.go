package i8comic

import (
	"comic2video/config"
	"comic2video/ffmpeg"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
)

// DownloadChapter 下载并保存章节内图片，返回下载内容的本地资源路径列表
func DownloadChapter(aid string) ([]string, error) {
	url := fmt.Sprintf("https://18comic.vip/photo/%s", aid)
	urls, err := GetChapterContent(url)
	if err != nil {
		return nil, err
	}

	sem := make(chan struct{}, config.Cfg.Service.DownloadThreads)
	errCh := make(chan error, len(urls))
	filePaths := make([]string, len(urls))
	var wg sync.WaitGroup

	// download images
	for ix, url := range urls {
		wg.Add(1)
		go func() {
			sem <- struct{}{}
			defer func() {
				wg.Done()
				<-sem
				log.Debugf("Finish download image (%d)", ix+1)
			}()

			log.Debugf("Start downloading image (%d)", ix+1)
			imagePath, err := downloadImage(url, aid)
			if err != nil {
				errCh <- err
				return
			}

			filePaths[ix] = imagePath
		}()
	}

	// wait jobs finish
	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
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
