package copymanga

import (
	"comic2video/config"
	"comic2video/ffmpeg"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
)

func downloadImage(chapterUrl, path, imgUrl string) error {
	req, err := http.NewRequest("GET", imgUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Referer", chapterUrl)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:153.0) Gecko/20100101 Firefox/153.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func DownloadAndConvert(seriesId, chapterId string) (string, error) {
	chapterUrl := fmt.Sprintf("https://www.mangacopy.com/comic/%s", seriesId)
	fullUrl := fmt.Sprintf("%s/chapter/%s", chapterUrl, chapterId)

	comicPath := filepath.Join(config.CPCOMIC_PATH, seriesId, chapterId)
	os.MkdirAll(comicPath, 0755)

	urls, err := GetChapterContent(chapterUrl, fullUrl)
	if err != nil {
		return "", err
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
			path := filepath.Join(comicPath, fmt.Sprintf("img-%d.jpg", ix))

			if err := downloadImage(chapterUrl, path, url); err != nil {
				errCh <- err
				return
			}

			filePaths[ix] = path
		}()
	}

	// wait jobs finish
	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
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
