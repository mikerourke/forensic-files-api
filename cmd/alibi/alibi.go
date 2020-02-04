package main

import (
	"os"
	"path/filepath"

	"github.com/mikerourke/forensic-files-api/internal/hearnoevil"
	"github.com/mikerourke/forensic-files-api/internal/videodiary"
	"github.com/mikerourke/forensic-files-api/internal/visibilityzero"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "Internal tools for the Forensic Files API.")
	app.HelpFlag.Short('h')

	hearNoEvil := app.Command(
		"hearnoevil",
		"Send audio files to speech-to-text service for transcribing.",
	)

	videoDiary := app.Command(
		"videodiary",
		"Download all episodes from YouTube.",
	)
	videoDiaryMissing := videoDiary.Flag("missing", "Log missing downloads only.").Bool()

	visibilityZero := app.Command(
		"viszero",
		"Extract audio from downloaded episodes for transcription.",
	)

	parsedCmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch parsedCmd {
	case hearNoEvil.FullCommand():
		hearnoevil.TranscribeEpisodes()

	case videoDiary.FullCommand():
		if *videoDiaryMissing {
			videodiary.LogMissingEpisodes()
		} else {
			videodiary.DownloadEpisodes()
		}

	case visibilityZero.FullCommand():
		visibilityzero.ExtractAudio()
	}
}
