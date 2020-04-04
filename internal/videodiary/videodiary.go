// Package videodiary downloads all of the Forensic Files episodes from YouTube
// and drops them in the `/assets/videos` directory.
package videodiary

import (
	"os"
	"os/exec"
	"time"

	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
)

var log = waterlogged.New("videodiary")

// Download downloads the specified episode number from the specified season
// number or all seasons.
func Download(seasonNumber int, episodeNumber int) {
	checkForYouTubeDL()

	if seasonNumber == 0 {
		if episodeNumber != 0 {
			log.Fatalln("You must specify a season number for an episode")
		}

		downloadAllSeasons()
		return
	}

	s := whodunit.NewSeason(seasonNumber)
	if err := s.PopulateEpisodes(); err != nil {
		log.WithError(err).Fatalln("Could not get season episodes")
	}

	if episodeNumber == 0 {
		downloadSeason(s, false)
	} else {
		ep := s.Episode(episodeNumber)
		downloadEpisode(ep, false)
	}
}

// LogStatusTable logs the episode statuses.
func LogStatusTable(status whodunit.AssetStatus) {
	table := whodunit.NewStatusTable(whodunit.AssetTypeVideo, status)
	table.Log()
}

func checkForYouTubeDL() {
	cmd := exec.Command("youtube-dl", "--version")
	err := cmd.Run()
	if err != nil {
		panic("Could not find youtube-dl executable, it may not be installed")
	}
}

func downloadAllSeasons() {
	for season := 1; season <= whodunit.SeasonCount; season++ {
		s := whodunit.NewSeason(season)
		if err := s.PopulateEpisodes(); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"season": season,
			}).Fatalln("Could not get season episodes")
			return
		}
		downloadSeason(s, true)
	}
}

func downloadSeason(s *whodunit.Season, isPaused bool) {
	for _, ep := range s.AllEpisodes() {
		downloadEpisode(ep, isPaused)
	}
}

func downloadEpisode(ep *whodunit.Episode, isPaused bool) {
	if ep.AssetExists(whodunit.AssetTypeVideo) {
		log.Infoln("Episode already downloaded, skipping")
		return
	}

	outPath := ep.AssetFilePath(whodunit.AssetTypeVideo)
	log.WithFields(logrus.Fields{
		"season":  ep.SeasonNumber,
		"episode": ep.EpisodeNumber,
		"title":   ep.Title,
		"path":    outPath,
		"hash":    ep.VideoHash(),
	}).Infoln("Downloading video from YouTube")

	cmd := exec.Command("youtube-dl", "-o", outPath, ep.URL)
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
