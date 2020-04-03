// Package visibilityzero is used to loop through the downloaded episodes and
// extract the audio that will be sent to the speech-to-text service.
package visibilityzero

import (
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
)

var (
	log            = waterlogged.New("visibilityzero")
	processedCount int
)

// ExtractAudio extracts the audio from the specified season and episode (or
// all if neither is specified) and saves it to an `.mp3` file.
func ExtractAudio(seasonNumber int, episodeNumber int) {
	checkForFFmpeg()

	if seasonNumber == 0 {
		if episodeNumber != 0 {
			log.Fatalln("You must specify a season number for an episode")
		}

		extractAudioFromAllSeasons()
		return
	}

	s := whodunit.NewSeason(seasonNumber)
	if err := s.PopulateEpisodes(); err != nil {
		log.WithError(err).Fatalln("Could not get season episodes")
	}

	if episodeNumber == 0 {
		extractAudioFromSeason(s)
	} else {
		ep := s.Episode(episodeNumber)
		extractAudioFromEpisode(ep)
	}
}

// LogStatusTable logs the episode statuses.
func LogStatusTable(status whodunit.AssetStatus) {
	table := whodunit.NewStatusTable(whodunit.AssetTypeAudio, status)
	table.Log()
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
		if strings.Contains(err.Error(), "exit status 1") {
			log.Infoln("Successfully extracted audio")
			return
		}

		log.WithFields(logrus.Fields{
			"error": err,
			"video": filepath.Base(videoPath),
		}).Errorln("Error extracting audio")
	}
}
