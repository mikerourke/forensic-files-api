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

// AssetStatus is an enum that represents the status of the asset.
type AssetStatus int

const (
	// AssetStatusAny is used primarily to show all records when logging the
	// status in the terminal.
	AssetStatusAny AssetStatus = iota

	// AssetStatusPending indicates that the asset hasn't been processed yet.
	AssetStatusPending

	// AssetStatusInProcess indicates that the asset is currently being
	// processed.
	AssetStatusInProcess

	// AssetStatusComplete indicates that the asset has been processed.
	AssetStatusComplete

	// AssetStatusMissing indicates that the asset is missing.
	AssetStatusMissing
)

// AssetType represents which type of asset the episode is associated with.
type AssetType int

const (
	// AssetTypeAudio represents the audio file associated with the episode.
	AssetTypeAudio AssetType = iota

	// AssetTypeRecognition represents the speech recognition JSON file
	// associated with the episode.
	AssetTypeRecognition

	// AssetTypeTranscript represents the transcript of the episode.
	AssetTypeTranscript

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
	case AssetTypeTranscript:
		return filepath.Join(AssetsDirPath, "transcripts")
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
	case AssetTypeTranscript:
		return ".txt"
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
