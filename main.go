package main

import (
	"comic_converter/config"
	"comic_converter/copymanga"
	"comic_converter/i8comic"
	"comic_converter/router"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{ReportTimestamp: true})
	log.SetDefault(logger)
	log.SetLevel(log.DebugLevel)

	switch os.Args[1] {
	case "serve":
		os.MkdirAll(config.BASE_PATH, 0755)
		config.InitDB()

		r := router.InitRouter()

		log.Info("Server running on :8080")
		http.ListenAndServe(":8080", r)

	case "copymanga":
		s, c, err := copymanga.UrlParser(os.Args[2])
		if err != nil {
			log.Errorf("Failed to parse url: %v", err)
			return
		}
		err = copymanga.DownloadAndConvertCopymanga(s, c)
		if err != nil {
			log.Errorf("Failed to convert: %v", err)
			return
		}

	case "18comic":
		aid, err := i8comic.UrlParser(os.Args[2])
		if err != nil {
			log.Errorf("Failed to parse url: %v", err)
			return
		}
		err = i8comic.DownloadAndConvert(aid)
		if err != nil {
			log.Errorf("Failed to test: %v", err)
			return
		}
	default:
		panic("unknown command")

	}
}
