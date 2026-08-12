package copymanga

import (
	"encoding/json"
	"io"
	"net/http"
)

type url struct {
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

	var urls []url
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
