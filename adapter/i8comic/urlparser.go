package i8comic

import (
	"fmt"
	"net/url"
	"strings"
)

func UrlParser(rawURL string) (string, error) {
	// https://18comic.vip/photo/1447064

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse chapter URL: %w", err)
	}
	if parsed.Hostname() != "18comic.vip" {
		return "", fmt.Errorf("unsupported copymanga URL: %q", rawURL)
	}

	parts := strings.Split(strings.Trim(parsed.EscapedPath(), "/"), "/")

	if len(parts) < 2 {
		return "", fmt.Errorf("invalid 18comic URL: %q", rawURL)
	}

	var aid string
	switch parts[0] {
	case "photo":
		aid, err = url.PathUnescape(parts[1])
		if err != nil {
			return "", fmt.Errorf("decode aid: %w", err)
		}
	default:
		return "", fmt.Errorf("Unexpected url structure")
	}

	return aid, nil
}
