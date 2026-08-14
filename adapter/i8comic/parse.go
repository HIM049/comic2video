package i8comic

import (
	"fmt"
	"regexp"
)

var (
	cdnHostRe    = regexp.MustCompile(`cdnImgHost:\s*"([^"]*)"`)
	containerRe  = regexp.MustCompile(`<div class="row thumb-overlay-albums">`)
	imgTagRe     = regexp.MustCompile(`<img[^>]*\bid="album_photo_[^"]*"[^>]*>`)
	attrRe       = regexp.MustCompile(`data-original="([^"]*)"`)
	nextSeriesRe = regexp.MustCompile(`id="next_series"[\s\S]{0,500}?<a[^>]*href="([^"]*)"`)
	photoAIDRe   = regexp.MustCompile(`/photo/(\d+)`)
)

// ParsePhotoPage 解析 photo 页 HTML
func ParsePhotoPage(body []byte) (images []string, nextChapterURL string, err error) {
	// 匹配图片列表
	start := 0
	if loc := containerRe.FindIndex(body); loc != nil {
		start = loc[1]
	}
	imgTags := imgTagRe.FindAll(body[start:], -1)
	if len(imgTags) == 0 {
		return nil, "", fmt.Errorf("no comic images found in page")
	}
	for _, tag := range imgTags {
		if m := attrRe.FindSubmatch(tag); m != nil {
			images = append(images, string(m[1]))
		}
	}

	// 匹配下一话
	if m := nextSeriesRe.FindSubmatch(body); m != nil {
		if aid := photoAIDRe.FindSubmatch(m[1]); aid != nil {
			nextChapterURL = "/photo/" + string(aid[1])
		}
	}
	return images, nextChapterURL, nil
}
