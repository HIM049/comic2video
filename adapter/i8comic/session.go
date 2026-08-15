package i8comic

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	cookieMu       sync.RWMutex
	sessionCookies string
)

// InitSession 通过站点主页以初始化会话 cookie。
func InitSession() error {
	req, err := http.NewRequest("GET", "https://18comic.vip/", nil)
	if err != nil {
		return err
	}
	setBrowserHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("init session: unexpected status code: %d", resp.StatusCode)
	}

	parts := make([]string, 0, 4)
	for _, c := range resp.Cookies() {
		if c.Name == "theme" {
			continue
		}
		parts = append(parts, c.Name+"="+c.Value)
	}
	cookieMu.Lock()
	sessionCookies = strings.Join(parts, "; ")
	cookieMu.Unlock()
	return nil
}

func sessionCookie() string {
	cookieMu.RLock()
	defer cookieMu.RUnlock()
	return sessionCookies
}

// setBrowserHeaders 设置模拟浏览器导航的基础请求头
func setBrowserHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:153.0) Gecko/20100101 Firefox/153.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,zh-HK;q=0.7,en-US;q=0.6,en;q=0.5")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
}
