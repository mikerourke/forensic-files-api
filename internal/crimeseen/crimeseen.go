package crimeseen

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const SeasonCount = 14

var AudioPath = filepath.Join(AssetsPath(), "audio")
var VideosPath = filepath.Join(AssetsPath(), "videos")

func AssetsPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting pwd")
	}

	return filepath.Join(pwd, "assets")
}

func Mkdirp(path string) error {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil && !IsFileExistsError(err) {
		return err
	}

	return nil
}

func IsFileExistsError(err error) bool {
	return strings.Contains(err.Error(), "file exists")
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func PaddedNumberString(value int) string {
	numString := strconv.Itoa(value)
	if value < 10 {
		return "0" + numString
	}
	return numString
}
