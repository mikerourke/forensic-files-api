// Package visibilityzero is used to loop through the downloaded episodes and
// extract the audio that will be sent to the speech-to-text service.
package visibilityzero

import (
	"os/exec"

	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
)

var log = waterlogged.New("visibilityzero")

// ExtractAudio extracts the audio from the specified season and episode (or
// all if neither is specified) and saves it to an `.mp3` file.
func ExtractAudio(seasonNumber int, episodeNumber int) {
	interrogate()

	onEpisode := func(ep *whodunit.Episode) {
		a := NewAudio(ep)
		a.Extract(episodeNumber == 0)
	}

	if err := whodunit.Solve(seasonNumber, episodeNumber, onEpisode); err != nil {
		log.WithError(err).Errorln("Error extracting audio from episode(s)")
	}
}

// Investigate logs the episode statuses.
func Investigate(status whodunit.AssetStatus) {
	table := whodunit.NewStatusTable(whodunit.AssetTypeAudio, status)
	table.Log()
}

func interrogate() {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		log.Fatalln("Could not find ffmpeg executable, it may not be installed")
	}
}
