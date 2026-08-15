package i8comic

import (
	"comic2video/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// GetChapterContent 读取并分析章节页面，得到全部资源url
func GetChapterContent(url string) ([]string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	setBrowserHeaders(req)
	req.Header.Set("Referer", url)
	req.Header.Set("Cookie", sessionCookie())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	imgs, _, err := ParsePhotoPage(data)
	if err != nil {
		return nil, err
	}

	return imgs, nil
}

// downloadImage 读取url文件并处理后保存
func downloadImage(url, aid string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	setBrowserHeaders(req)
	req.Header.Set("Referer", fmt.Sprintf("https://18comic.vip/photo/%s", aid))
	req.Header.Set("Cookie", sessionCookie())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	filename := path.Base(url)
	newData, err := UnscrambleImage(data, aid, filename)

	comicPath := filepath.Join(config.I8COMIC_PATH, aid)
	os.MkdirAll(comicPath, 0755)

	imagePath := filepath.Join(comicPath, filename)
	file, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	_, err = file.Write(newData)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}
