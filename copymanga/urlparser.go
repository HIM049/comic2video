package copymanga

import (
	"fmt"
	neturl "net/url"
	"strings"
)

func UrlParser(rawURL string) (seriesID string, chapterID string, err error) {
	// https://www.mangacopy.com/comic/{seriesId}/chapter/{chapterId}
	parsed, err := neturl.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("parse chapter URL: %w", err)
	}
	if parsed.Hostname() != "www.mangacopy.com" {
		return "", "", fmt.Errorf("unsupported CopyManga URL: %q", rawURL)
	}

	parts := strings.Split(strings.Trim(parsed.EscapedPath(), "/"), "/")
	if parts[0] != "comic" || parts[2] != "chapter" || parts[1] == "" || parts[3] == "" {
		return "", "", fmt.Errorf("expected https://www.mangacopy.com/comic/{seriesId}/chapter/{chapterId}")
	}

	seriesID, err = neturl.PathUnescape(parts[1])
	if err != nil {
		return "", "", fmt.Errorf("decode series ID: %w", err)
	}
	chapterID, err = neturl.PathUnescape(parts[3])
	if err != nil {
		return "", "", fmt.Errorf("decode chapter ID: %w", err)
	}
	return seriesID, chapterID, nil
}
