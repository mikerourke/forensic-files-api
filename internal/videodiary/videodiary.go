// Package videodiary downloads all of the Forensic Files episodes from YouTube
// and drops them in the `/assets/videos` directory.
package videodiary

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
)

var log = waterlogged.New("videodiary")

// DownloadEpisodes parses the YouTube episode URLs from the `/assets/youtube-links.json`
// file and downloads each episode to the `/assets/videos` directory.
func DownloadEpisodes() {
	checkForYouTubeDL()

	for season := 1; season <= whodunit.SeasonCount; season++ {
		s := whodunit.NewSeason(season)
		if err := s.PopulateEpisodes(); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"season": season,
			}).Fatalln("Could not get season episodes")
			return
		}

		for _, ep := range s.AllEpisodes() {
			downloadEpisode(ep, true)
		}
	}
}

// DownloadEpisode downloads the specified episode number from the specified
// season number.
func DownloadEpisode(seasonNumber int, episodeNumber int) {
	checkForYouTubeDL()

	s := whodunit.NewSeason(seasonNumber)
	if err := s.PopulateEpisodes(); err != nil {
		log.WithError(err).Fatalln("Could not get season episodes")
	}

	if ep := s.Episode(episodeNumber); ep != nil {
		downloadEpisode(ep, false)
	} else {
		log.WithField("episode", episodeNumber).Fatalln(
			"Could not find episode")
	}
}

// LogMissingEpisodes logs the episodes that haven't been downloaded to the
// command line.
func LogMissingEpisodes() {
	missingCount := 0

	for season := 1; season <= whodunit.SeasonCount; season++ {
		s := whodunit.NewSeason(season)
		if err := s.PopulateEpisodes(); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"season": season,
			}).Fatalln("Could not get season episodes")
			return
		}

		for _, ep := range s.AllEpisodes() {
			if ep.VideoHash() == "" {
				fmt.Printf(
					"Season: %v \t Episode: %v \t Title: %v\n",
					ep.SeasonNumber,
					ep.EpisodeNumber,
					ep.DisplayTitle(),
				)
				missingCount++
			}
		}
	}

	fmt.Printf("Total count missing: %v\n", missingCount)
}

func checkForYouTubeDL() {
	cmd := exec.Command("youtube-dl", "--version")
	err := cmd.Run()
	if err != nil {
		panic("Could not find youtube-dl executable, it may not be installed")
	}
}

func downloadEpisode(ep *whodunit.Episode, isPaused bool) {
	outPath := ep.AssetFilePath(whodunit.AssetTypeVideo)
	if crimeseen.FileExists(outPath) {
		return
	}

	log.WithFields(logrus.Fields{
		"season":  ep.SeasonNumber,
		"episode": ep.EpisodeNumber,
		"title":   ep.Title,
		"path":    outPath,
	}).Infoln("Downloading video from YouTube")

	cmd := exec.Command("youtube-dl",
		"-o", outPath,
		ep.VideoHash())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"title": ep.Title,
			"path":  outPath,
		}).Errorln("Error downloading video")
	}

	// We're hedging our bets here to make sure we don't exceed some kind of rate limit:
	if isPaused {
		log.Println("Download successful, waiting 1 minute")
		time.Sleep(time.Minute * 1)
	}
}
