package router

import (
	"comic_converter/config"
	"comic_converter/copymanga"
	"comic_converter/model"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

const COPYMANGA = "copymanga"

func InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", defaultHandler)
	r.Get("/copymanga", copymangaHandler)

	return r
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Comic Converter - %s", config.VERSION)
}

func copymangaHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	series, chapter, err := copymanga.UrlParser(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var comic model.Comic
	err = config.DB.
		Where("source = ?", COPYMANGA).
		Where("series_id = ?", series).
		Where("chapter_id = ?", chapter).
		First(&comic).
		Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// other errors
			http.Error(w, fmt.Sprintf("Failed to load resource %v", err), http.StatusInternalServerError)
			return
		} else {
			// notfound
			log.Info("Resourse not cached, downloading...")
			w.WriteHeader(http.StatusAccepted)
			err = copymanga.DownloadAndConvertCopymanga(series, chapter)
			if err != nil {
				log.Errorf("Failed to convert resource %v", err)
				http.Error(w, fmt.Sprintf("Failed to convert resource %v", err), http.StatusInternalServerError)
				return
			}

			newComic := model.Comic{
				Source:    COPYMANGA,
				SeriesId:  series,
				ChapterId: chapter,
			}

			err = config.DB.Create(&newComic).Error
			if err != nil {
				log.Errorf("Failed to create record %v", err)
				http.Error(w, fmt.Sprintf("Failed to create record %v", err), http.StatusInternalServerError)
				return
			}
			comic = newComic
		}
	}

	// find and open video file
	videoPath := filepath.Join(config.CP_VIDEO_PATH, series)
	path := filepath.Join(videoPath, fmt.Sprintf("%s.mp4", chapter))
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
