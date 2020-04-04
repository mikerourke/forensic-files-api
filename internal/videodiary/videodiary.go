// Package videodiary downloads all of the Forensic Files episodes from YouTube
// and drops them in the `/assets/videos` directory.
package videodiary

import (
	"os/exec"

	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
)

var log = waterlogged.New("videodiary")

// Download downloads the specified episode number from the specified season
// number or all seasons.
func Download(seasonNumber int, episodeNumber int) {
	interrogate()

	onEpisode := func(ep *whodunit.Episode) {
		v := NewVideo(ep)
		v.Download(seasonNumber == 0)
	}

	if err := whodunit.Solve(seasonNumber, episodeNumber, onEpisode); err != nil {
		log.WithError(err).Errorln("Error downloading episode(s)")
	}
}

// Investigate logs the episode statuses.
func Investigate(status whodunit.AssetStatus) {
	table := whodunit.NewStatusTable(whodunit.AssetTypeVideo, status)
	table.Log()
}

func interrogate() {
	cmd := exec.Command("youtube-dl", "--version")
	err := cmd.Run()
	if err != nil {
		log.Fatalln("Could not find youtube-dl executable, it may not be installed")
	}
}
