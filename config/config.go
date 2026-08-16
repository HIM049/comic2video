package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Cfg Config

type Service struct {
	FfmpegPath      string `json:"ffmpeg_path"`
	FfmpegThreads   string `json:"ffmpeg_threads"`
	DownloadThreads int    `json:"download_threads"`
	BatchSize       int    `json:"batch_size"`
}

type Config struct {
	Service *Service `json:"service"`
}

func Init() {
	godotenv.Load()

	cfg := Config{
		Service: initService(),
	}

	Cfg = cfg
}

func initService() *Service {
	return &Service{
		FfmpegPath:      getEnvOrDefault("FFMPEG_PATH", "ffmpeg"),
		FfmpegThreads:   getEnvOrDefault("FFMPEG_THREADS", "8"),
		DownloadThreads: getIntEnvOrDefault("DOWNLOAD_THREADS", 1),
		BatchSize:       getIntEnvOrDefault("BATCH_SIZE", 0),
	}
}

// getEnv reads an environment variable, panics if not set.
func getEnv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	panic("environment variable not found: " + key)
}

// getEnv reads an environment variable, panics if not set.
func getEnvOrDefault(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return def
}

// getEnv reads an environment variable, panics if not set.
func getIntEnvOrDefault(key string, def int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	value, err := strconv.Atoi(val)
	if err != nil {
		value = def
	}

	return value
}
