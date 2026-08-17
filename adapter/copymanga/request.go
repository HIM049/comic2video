package copymanga

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type imgUrl struct {
	Url string `json:"url"`
}

func GetChapterContent(chapterUrl, fullUrl string) ([]string, error) {

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", chapterUrl)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:153.0) Gecko/20100101 Firefox/153.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cct, contentKey, err := keyReader(body)
	if err != nil {
		return nil, err
	}

	list, err := decryption(cct, contentKey)

	var urls []imgUrl
	err = json.Unmarshal(list, &urls)
	if err != nil {
		return nil, err
	}

	var urlstrings []string
	for _, url := range urls {
		urlstrings = append(urlstrings, url.Url)
	}

	return urlstrings, nil
}

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
