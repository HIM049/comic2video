package main

import (
	"comic2video/adapter/copymanga"
	"comic2video/adapter/i8comic"
	"comic2video/config"
	"comic2video/router"
	"comic2video/service"
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{ReportTimestamp: true})
	log.SetDefault(logger)
	log.SetLevel(log.DebugLevel)

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "version":
		log.Infof("Version: %s", config.VERSION)
		return
	case "serve":
		os.MkdirAll(config.BASE_PATH, 0755)
		config.Init()
		config.InitDB()
		service.Init()

		if err := i8comic.InitSession(); err != nil {
			log.Errorf("Failed to init 18comic session: %v", err)
			return
		}

		r := router.InitRouter()

		log.Info("Server running on :8080")
		http.ListenAndServe(":8080", r)
		return

	case "copymanga":
		s, c, err := copymanga.UrlParser(os.Args[2])
		if err != nil {
			log.Errorf("Failed to parse url: %v", err)
			return
		}
		vp, err := copymanga.DownloadAndConvert(s, c)
		if err != nil {
			log.Errorf("Failed to convert: %v", err)
			return
		}
		log.Infof("Download Successfully. Video Path: %s", vp)
		return

	case "18comic":
		aid, err := i8comic.UrlParser(os.Args[2])
		if err != nil {
			log.Errorf("Failed to parse url: %v", err)
			return
		}
		if err := i8comic.InitSession(); err != nil {
			log.Errorf("Failed to init 18comic session: %v", err)
			return
		}
		vp, err := i8comic.DownloadAndConvert(aid)
		if err != nil {
			log.Errorf("Failed to test: %v", err)
			return
		}
		log.Infof("Download Successfully. Video Path: %s", vp)
		return

	default:
		printHelp()
		return
	}
}

func printHelp() {

	help := `Comic2Video Menu
Usage:
  comic_parser <command> [args]
Commands:
  version  Print the app version
  serve    Start the server
  copymanga <url>  Convert a copymanga comic
  18comic <url>  Convert an 18comic comic
  `
	fmt.Println(help)
}
