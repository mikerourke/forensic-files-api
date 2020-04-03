// Package crimeseen contains common paths and utility functions that come in
// handy across the codebase.
package crimeseen

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// SeasonCount is the count of seasons in Forensic Files.
const SeasonCount = 14

// AssetsDirPath is the absolute path to the `/assets` directory.
var AssetsDirPath = assetsDirPath()

// AudioDirPath is the absolute path to the `/assets/audio` directory.
var AudioDirPath = filepath.Join(AssetsDirPath, "audio")

// VideosDirPath is the absolute path to the `/assets/videos` directory.
var VideosDirPath = filepath.Join(AssetsDirPath, "videos")

// RecognitionsDirPath is the absolute path to the `/assets/recognitions` directory.
var RecognitionsDirPath = filepath.Join(AssetsDirPath, "recognitions")

func assetsDirPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting pwd")
	}

	return filepath.Join(pwd, "assets")
}

// Mkdirp creates the specified directory path if it doesn't already exist.
func Mkdirp(path string) error {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil && !IsFileExistsError(err) {
		return err
	}

	return nil
}

// IsFileExistsError returns true if the specified error is due to a file already
// existing on the filesystem.
func IsFileExistsError(err error) bool {
	return strings.Contains(err.Error(), "file exists")
}

// FileExists returns true if the specified file path exists.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// PaddedNumberString takes an integer value and returns a string value with a
// "0" prefix if the value is less than 10.
func PaddedNumberString(value int) string {
	numString := strconv.Itoa(value)
	if value < 10 {
		return "0" + numString
	}
	return numString
}

// WriteJSONToAssets writes the specified contents to a directory in the `/assets`
// directory. To write to the top level `/assets` directory and not a subdirectory,
// pass an empty string as the first argument.
func WriteJSONToAssets(
	dirName string,
	fileName string,
	contents interface{},
) error {
	outputPath := filepath.Join(AssetsDirPath, dirName, fileName)

	bytes, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, string(bytes))
	if err != nil {
		return err
	}

	return file.Sync()
}

// ReadJSONFromAssets reads the specified path from the `/assets` directory.
// The path can either be a filename (for the root `/assets` directory) or
// include the subdirectory.
func ReadJSONFromAssets(path string) (interface{}, error) {
	jsonFile, err := os.Open(filepath.Join(AssetsDirPath, path))
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var contents interface{}
	err = json.Unmarshal(bytes, &contents)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
