// Package killigraphy parses the recognition JSON files and creates text files
// from the results.
package killigraphy

import (
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
)

var log = waterlogged.New("killigraphy")

// Transcribe creates a transcript for the specified episode number from the
// specified season number or all seasons.
func Transcribe(seasonNumber int, episodeNumber int) {
	onEpisode := func(ep *whodunit.Episode) {
		t := NewTranscript(ep)
		t.Create()
	}

	if err := whodunit.Solve(seasonNumber, episodeNumber, onEpisode); err != nil {
		log.WithError(err).Errorln("Error transcribing episode(s)")
	}
}

// Investigate logs the episode statuses.
func Investigate(status whodunit.AssetStatus) {
	table := whodunit.NewStatusTable(whodunit.AssetTypeTranscript, status)
	table.Log()
}
