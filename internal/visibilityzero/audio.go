package visibilityzero

import (
	"os"
	"strings"
	"time"

	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/videodiary"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
)

// Audio represents the MP3 file extracted from the video file.
type Audio struct {
	*whodunit.Episode
}

// NewAudio returns a new instance of audio.
func NewAudio(ep *whodunit.Episode) *Audio {
	return &Audio{
		Episode: ep,
	}
}

// Extract extracts audio from the video file.
func (a *Audio) Extract(isPaused bool) {
	v := videodiary.NewVideo(a.Episode)
	if !v.Exists() {
		log.WithField("file", v.FileName()).Warnln(
			"Skipping job, video file not found")
		return
	}

	log.WithField("video", v.FileName()).Infoln(
		"Extracting audio from video file")

	err := crimeseen.RunCommand("ffmpeg", "-i", v.FilePath(), a.FilePath())
	if err != nil {
		if strings.Contains(err.Error(), "exit status 1") {
			log.Infoln("Successfully extracted audio")
		} else {
			log.WithFields(logrus.Fields{
				"error": err,
				"video": v.FileName(),
			}).Errorln("Error extracting audio")
		}
	}

	// Adding a 30 second delay here so my laptop doesn't melt.
	if isPaused {
		log.Println("Extraction successful, waiting 30 seconds")
		time.Sleep(time.Second * 30)
	}
}

// Open returns the audio file contents.
func (a *Audio) Open() *os.File {
	audio, err := os.Open(a.FilePath())
	if err != nil {
		log.WithError(err).Errorln("Error opening audio")
		return nil
	}

	return audio
}

// Exists return true if the audio file exists in the `/assets` directory.
func (a *Audio) Exists() bool {
	return a.AssetExists(whodunit.AssetTypeAudio)
}

// FilePath returns the path to the audio file in the `/assets` directory.
func (a *Audio) FilePath() string {
	return a.AssetFilePath(whodunit.AssetTypeAudio)
}

// FileName returns the name of the audio file in the `/assets` directory.
func (a *Audio) FileName() string {
	return a.AssetFileName(whodunit.AssetTypeAudio)
}
