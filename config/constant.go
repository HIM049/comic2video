package config

// VERSION is overridden at build time by CI via -ldflags "-X comic_converter/config.VERSION=..."
var VERSION = "testbuild"

const BASE_PATH = "./app_data"
const COMIC_PATH = "./app_data/comic"
const VIDEO_PATH = "./app_data/video"
const DB_FILE_PATH = "./app_data/storage.db"
