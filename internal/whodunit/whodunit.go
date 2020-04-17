// Package whodunit acts as a wrapper around the contents of the `/assets`
// directory, mainly episodes and seasons. Since most of the internal tools
// interact with these files, whodunit makes it easier to get file information
// and metadata for the associated assets.
package whodunit

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
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
	// AssetTypeGCPAnalysis represents the GCP entity analysis associated with the episode.
	AssetTypeGCPAnalysis AssetType = iota

	// AssetTypeIBMAnalysis represents the IBM entity analysis associated with the episode.
	AssetTypeIBMAnalysis

	// AssetTypeAudio represents the audio file associated with the episode.
	AssetTypeAudio

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

var env = crimeseen.NewEnv()

// DirPath returns the absolute path to the directory associated with the
// asset type.
func (at AssetType) DirPath() string {
	switch at {
	case AssetTypeGCPAnalysis:
		return filepath.Join(env.InvestigationsPath(), "gcp-analyses")
	case AssetTypeIBMAnalysis:
		return filepath.Join(env.InvestigationsPath(), "ibm-analyses")
	case AssetTypeAudio:
		return filepath.Join(AssetsDirPath, "audio")
	case AssetTypeRecognition:
		return filepath.Join(env.InvestigationsPath(), "recognitions")
	case AssetTypeTranscript:
		return filepath.Join(env.InvestigationsPath(), "transcripts")
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
	case AssetTypeTranscript:
		return ".txt"
	case AssetTypeVideo:
		return ".mp4"
	default:
		return ".json"
	}
}

// Solve runs the specified function per episode based on the values passed
// in for the season and episode.
func Solve(
	seasonNumber int,
	episodeNumber int,
	onEpisode func(ep *Episode),
) error {
	// Performs the specified action for all of the episodes in a season.
	onSeason := func(s *Season) error {
		if err := s.PopulateEpisodes(); err != nil {
			return err
		}

		if episodeNumber != 0 {
			onEpisode(s.Episode(episodeNumber))
		} else {
			for _, ep := range s.AllEpisodes() {
				onEpisode(ep)
			}
		}

		return nil
	}

	// If the season specified is 0 (not specified), the onEpisode action needs
	// to run on every episode in every season.
	if seasonNumber == 0 {
		// How do we know which season to process if it isn't specified?
		if episodeNumber != 0 {
			return errors.New("you must specify a season number for an episode")
		}

		for season := 1; season <= SeasonCount; season++ {
			s := NewSeason(season)
			if err := onSeason(s); err != nil {
				return err
			}
		}
		return nil
	}

	s := NewSeason(seasonNumber)
	return onSeason(s)
}

func assetsDirPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting pwd")
	}

	return filepath.Join(pwd, "assets")
}
