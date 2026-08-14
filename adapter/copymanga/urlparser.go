package copymanga

import (
	"fmt"
	"net/url"
	"strings"
)

func UrlParser(rawURL string) (string, string, error) {
	// https://www.mangacopy.com/comic/{seriesId}/chapter/{chapterId}
	// https://www.mangacopy.com/h5/comicContent/{seriesId}/{chapterId}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("parse chapter URL: %w", err)
	}
	if parsed.Hostname() != "www.mangacopy.com" {
		return "", "", fmt.Errorf("unsupported copymanga URL: %q", rawURL)
	}

	parts := strings.Split(strings.Trim(parsed.EscapedPath(), "/"), "/")

	var seriesID, chapterID string
	switch parts[0] {
	case "comic":
		if len(parts) < 4 {
			return "", "", fmt.Errorf("invalid copymanga URL: %q", rawURL)
		}
		seriesID, err = url.PathUnescape(parts[1])
		if err != nil {
			return "", "", fmt.Errorf("decode series ID: %w", err)
		}
		chapterID, err = url.PathUnescape(parts[3])
		if err != nil {
			return "", "", fmt.Errorf("decode chapter ID: %w", err)
		}
	case "h5":
		if len(parts) < 4 {
			return "", "", fmt.Errorf("invalid copymanga URL: %q", rawURL)
		}
		seriesID, err = url.PathUnescape(parts[2])
		if err != nil {
			return "", "", fmt.Errorf("decode series ID: %w", err)
		}
		chapterID, err = url.PathUnescape(parts[3])
		if err != nil {
			return "", "", fmt.Errorf("decode chapter ID: %w", err)
		}
	default:
		return "", "", fmt.Errorf("Unexpected url structure")
	}

	return seriesID, chapterID, nil
}
