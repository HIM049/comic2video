package copymanga

import (
	"comic_converter/config"
	"comic_converter/ffmpeg"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloader(chapterUrl, path, imgUrl string) error {
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

func DownloadCopymanga(url string) {
	seriesId, chapterId, err := UrlParser(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	chapterUrl := fmt.Sprintf("https://www.mangacopy.com/comic/%s", seriesId)
	fullUrl := fmt.Sprintf("%s/chapter/%s", chapterUrl, chapterId)

	finalPath := filepath.Join(config.BASE_PATH, "comic", seriesId, chapterId)
	os.MkdirAll(finalPath, 0755)

	urls, err := GetChapterContent(chapterUrl, fullUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	var filePaths []string
	for ix, url := range urls {
		path := filepath.Join(finalPath, fmt.Sprintf("img-%d.jpg", ix))
		filePaths = append(filePaths, path)

		err := downloader(url, path, url)
		if err != nil {
			fmt.Println(err)
		}
	}

	path := filepath.Join(finalPath, "output.mp4")
	err = ffmpeg.Convert(filePaths, path)
	if err != nil {
		fmt.Println(err)
	}
}
