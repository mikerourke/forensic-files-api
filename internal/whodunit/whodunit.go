// Package whodunit acts as a wrapper around the contents of the `/assets`
// directory, mainly episodes and seasons. Since most of the internal tools
// interact with these files, whodunit makes it easier to get file information
// and metadata for the associated assets.
package whodunit

import (
	"log"
	"os"
	"path/filepath"
)

// AssetType represents which type of asset the episode is associated with.
type AssetType int

const (
	// AssetTypeAudio represents the audio file associated with the episode.
	AssetTypeAudio AssetType = iota

	// AssetTypeRecognition represents the speech recognition JSON file
	// associated with the episode.
	AssetTypeRecognition

	// AssetTypeVideo represents the video file associated with the episode.
	AssetTypeVideo
)

// AssetsDirPath is the absolute path to the `/assets` directory.
var AssetsDirPath = assetsDirPath()

// DirPath returns the absolute path to the directory associated with the
// asset type.
func (at AssetType) DirPath() string {
	switch at {
	case AssetTypeAudio:
		return filepath.Join(AssetsDirPath, "audio")
	case AssetTypeRecognition:
		return filepath.Join(AssetsDirPath, "recognitions")
	case AssetTypeVideo:
		return filepath.Join(AssetsDirPath, "videos")
	default:
		return ""
	}
}

// FileExt returns the file extension associated with the asset type.
func (at AssetType) FileExt() string {
	switch at {
	case AssetTypeAudio:
		return ".mp3"
	case AssetTypeRecognition:
		return ".json"
	case AssetTypeVideo:
		return ".mp4"
	default:
		return ""
	}
}

func assetsDirPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting pwd")
	}

	return filepath.Join(pwd, "assets")
}
