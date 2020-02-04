package common

import (
	"log"
	"os"
	"path/filepath"
)

func AudioPath() string {
	return filepath.Join(AssetsPath(), "audio")
}

func VideosPath() string {
	return filepath.Join(AssetsPath(), "videos")
}

func AssetsPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting pwd")
	}

	return filepath.Join(pwd, "assets")
}
