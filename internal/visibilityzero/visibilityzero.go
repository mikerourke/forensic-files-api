// Package visibilityzero is used to loop through the downloaded episodes and
// extract the audio that will be sent to the speech-to-text service.
package visibilityzero

import (
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
)

var (
	log            = waterlogged.New("visibilityzero")
	processedCount int
)

// ExtractAudio loops through all of the `/video` season directories, extracts the
// audio from the .mp4 file using ffmpeg, and drops it into the `/assets/audio`
// directory for the corresponding season.
func ExtractAudio() {
	checkForFFmpeg()
	extractAudioFromAllSeasons()
}

func checkForFFmpeg() {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		panic("Could not find ffmpeg executable, it may not be installed")
	}
}

func extractAudioFromAllSeasons() {
	processedCount = 0

	for season := 1; season <= whodunit.SeasonCount; season++ {
		s := whodunit.NewSeason(season)
		if err := s.PopulateEpisodes(); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"season": season,
			}).Errorln("Error getting season")
			return
		}
		extractAudioFromSeason(s)
	}
}

func extractAudioFromSeason(s *whodunit.Season) {
	for _, ep := range s.AllEpisodes() {
		// Every 10 videos, take a 5 minute breather. ffmpeg makes the
		// fans go bananas on my laptop:
		if processedCount != 0 && processedCount%10 == 0 {
			log.Infoln("Taking a breather or else I'm going to take off")
			time.Sleep(time.Minute * 5)
		}

		if !ep.AssetExists(whodunit.AssetTypeAudio) {
			extractAudioFromEpisode(ep)
			processedCount++
		}
	}
}

func extractAudioFromEpisode(ep *whodunit.Episode) {
	log.WithFields(logrus.Fields{
		"video": ep.AssetFileName(whodunit.AssetTypeVideo),
	}).Infoln("Extracting audio from video file")

	videoPath := ep.AssetFilePath(whodunit.AssetTypeVideo)
	cmd := exec.Command("ffmpeg", "-i",
		videoPath, ep.AssetFilePath(whodunit.AssetTypeAudio))
	err := cmd.Run()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"video": filepath.Base(videoPath),
		}).Errorln("Error extracting audio")
	}
}
