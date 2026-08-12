package main

import (
	"comic_converter/copymanga"
	"comic_converter/ffmpeg"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const BASE_PATH = "app_data"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("url mismatch")
	}

	seriesId, chapterId, err := copymanga.UrlParser(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	chapterUrl := fmt.Sprintf("https://www.mangacopy.com/comic/%s", seriesId)
	fullUrl := fmt.Sprintf("%s/chapter/%s", chapterUrl, chapterId)

	finalPath := filepath.Join(BASE_PATH, "comic", seriesId, chapterId)
	os.MkdirAll(finalPath, 0755)

	urls, err := copymanga.GetChapterContent(chapterUrl, fullUrl)
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

	// var filePaths []string

	// for i := 0; i <= 10; i++ {
	// 	path := fmt.Sprintf("app_data/comic/ouxiangkuangre/7db0e387-6b71-11f0-ad8e-fa163e02432f/img-%d.jpg", i)
	// 	filePaths = append(filePaths, path)
	// }

	// err := ffmpeg.Convert(filePaths, "./output.mp4")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

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
