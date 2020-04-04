package videodiary

import (
	"time"

	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
)

// Video represents a video downloaded from YouTube.
type Video struct {
	*whodunit.Episode
}

// NewVideo returns a new instance of a video.
func NewVideo(ep *whodunit.Episode) *Video {
	return &Video{
		Episode: ep,
	}
}

// Download downloads the video from YouTube.
func (v *Video) Download(isPaused bool) {
	if v.Exists() {
		log.Infoln("Episode already downloaded, skipping")
		return
	}

	path := v.FilePath()
	log.WithFields(logrus.Fields{
		"season":  v.SeasonNumber,
		"episode": v.EpisodeNumber,
		"title":   v.Title,
		"path":    path,
		"url":     v.URL,
	}).Infoln("Downloading video from YouTube")

	err := crimeseen.RunCommand("youtube-dl", "-o", path, v.URL)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"title": v.Title,
			"path":  path,
		}).Errorln("Error downloading video")
	}

	// We're hedging our bets here to make sure we don't exceed some kind of
	// rate limit:
	if isPaused {
		log.Println("Download successful, waiting 1 minute")
		time.Sleep(time.Minute * 1)
	}
}

// Exists return true if the video file exists in the `/assets` directory.
func (v *Video) Exists() bool {
	return v.AssetExists(whodunit.AssetTypeVideo)
}

// FilePath returns the path to the video file in the `/assets` directory.
func (v *Video) FilePath() string {
	return v.AssetFilePath(whodunit.AssetTypeVideo)
}

// FileName returns the name of the video file in the `/assets` directory.
func (v *Video) FileName() string {
	return v.AssetFileName(whodunit.AssetTypeVideo)
}
