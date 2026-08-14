package config

// VERSION is overridden at build time by CI via -ldflags "-X comic_converter/config.VERSION=..."
var VERSION = "testbuild"

const BASE_PATH = "./app_data"
const CPCOMIC_PATH = "./app_data/copycomic"
const I8COMIC_PATH = "./app_data/i8comic"
const CP_VIDEO_PATH = "./app_data/cpvideo"
const I8_VIDEO_PATH = "./app_data/i8video"
const DB_FILE_PATH = "./app_data/storage.db"
