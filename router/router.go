package router

import (
	"comic_converter/config"
	"comic_converter/copymanga"
	"comic_converter/i8comic"
	"comic_converter/model"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", defaultHandler)
	r.Get("/m2v", genericHandler)

	return r
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Comic Converter - %s", config.VERSION)
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
		vp, err := copymangaHandler(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		videoPath = vp
	case "18comic.vip":
		vp, err := i8comicHandler(u)
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

func copymangaHandler(url string) (string, error) {
	series, chapter, err := copymanga.UrlParser(url)
	if err != nil {
		return "", err
	}

	var comic model.ComicCache
	err = config.DB.
		Where("source = ?", model.SourceCopymanga).
		Where("comic_id = ?", chapter).
		First(&comic).
		Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// other errors
			return "", err
		} else {
			// notfound
			log.Info("Resourse not cached, downloading...")
			videoPath, err := copymanga.DownloadAndConvertCopymanga(series, chapter)
			if err != nil {
				log.Errorf("Failed to convert resource %v", err)
				return "", err
			}

			newComic := model.ComicCache{
				Uid:       uuid.NewString(),
				Source:    model.SourceCopymanga,
				ComicId:   chapter,
				ComicUrl:  url,
				VideoPath: videoPath,
			}

			err = config.DB.Create(&newComic).Error
			if err != nil {
				log.Errorf("Failed to create record %v", err)
				return "", err
			}
			comic = newComic
		}
	}

	return comic.VideoPath, nil
}

func i8comicHandler(url string) (string, error) {
	aid, err := i8comic.UrlParser(url)
	if err != nil {
		return "", err
	}

	var comic model.ComicCache
	err = config.DB.
		Where("source = ?", model.SourceI8comic).
		Where("comic_id = ?", aid).
		First(&comic).
		Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// other errors
			return "", err
		} else {
			// notfound
			log.Info("Resourse not cached, downloading...")
			videoPath, err := i8comic.DownloadAndConvert(aid)
			if err != nil {
				log.Errorf("Failed to convert resource %v", err)
				return "", err
			}

			newComic := model.ComicCache{
				Uid:       uuid.NewString(),
				Source:    model.SourceI8comic,
				ComicId:   aid,
				ComicUrl:  url,
				VideoPath: videoPath,
			}

			err = config.DB.Create(&newComic).Error
			if err != nil {
				log.Errorf("Failed to create record %v", err)
				return "", err
			}
			comic = newComic
		}
	}

	return comic.VideoPath, nil
}
