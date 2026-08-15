package router

import (
	"comic2video/adapter/copymanga"
	"comic2video/adapter/i8comic"
	"comic2video/config"
	"comic2video/model"
	"comic2video/service"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", defaultHandler)
	r.Get("/c2v", genericHandler)

	return r
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Comic2Video - %s\nhttps://github.com/HIM049/comic2video (gonna public)", config.VERSION)
}

func genericHandler(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query().Get("url")
	parsed, err := url.Parse(u)
	if err != nil {
		return
	}

	var videoPath string
	switch parsed.Hostname() {
	case "www.mangacopy.com":
		vp, err := copymangaHandler(r.Context(), u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		videoPath = vp
	case "18comic.vip":
		vp, err := i8comicHandler(r.Context(), u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		videoPath = vp
	default:
		http.Error(w, "Unsupported hostname: "+parsed.Hostname(), http.StatusBadRequest)
		return
	}

	// find and open video file
	file, err := os.Open(videoPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// prepare modtime
	var t time.Time
	info, err := file.Stat()
	if err != nil {
		t = time.Now()
	}
	t = info.ModTime()

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "comic", t, file)
}

func copymangaHandler(ctx context.Context, url string) (string, error) {
	series, chapter, err := copymanga.UrlParser(url)
	if err != nil {
		return "", err
	}

	videoPath, err := service.GetOrCreateVideo(ctx, model.SourceCopymanga, chapter, url,
		func() (string, error) { return copymanga.DownloadAndConvert(series, chapter) })
	if err != nil {
		log.Errorf("Failed to convert resource %v", err)
		return "", err
	}
	return videoPath, nil
}

func i8comicHandler(ctx context.Context, url string) (string, error) {
	aid, err := i8comic.UrlParser(url)
	if err != nil {
		return "", err
	}

	videoPath, err := service.GetOrCreateVideo(ctx, model.SourceI8comic, aid, url,
		func() (string, error) { return i8comic.DownloadAndConvert(aid) })
	if err != nil {
		log.Errorf("Failed to convert resource %v", err)
		return "", err
	}
	return videoPath, nil
}
