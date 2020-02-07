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

	hearCommand := app.Command(
		"hearnoevil",
		"Service that transcribes audio with speech to text.",
	)

	hearRegisterCommand := hearCommand.Command(
		"registercb",
		"Register a callback URL.",
	)

	hearRegisterCommandURLFlag := hearRegisterCommand.Arg(
		"url",
		"Callback URL to use for transcribing.",
	).Required().String()

	hearServerCommand := hearCommand.Command(
		"server",
		"Start the callback URL server (required to start transcribing).",
	)

	hearRecognize := hearCommand.Command(
		"recognize",
		"Send recognition job requests to the speech to text service.",
	)

	hearRecognizeSeasonFlag := hearRecognize.Flag(
		"season",
		"Season of the episode to recognize.",
	).Required().Int()

	hearRecognizeEpisodeFlag := hearRecognize.Flag(
		"episode",
		"Episode number to recognize.",
	).Int()

	hearRecognizeURLFlag := hearRecognize.Flag(
		"url",
		"Callback URL to use for recognition jobs.",
	).String()

	diaryCommand := app.Command(
		"videodiary",
		"Download all episodes from YouTube.",
	)

	diaryCommandMissingFlag := diaryCommand.Flag(
		"missing",
		"Log missing downloads only.",
	).Bool()

	visZeroCommand := app.Command(
		"viszero",
		"Extract audio from downloaded episodes for recognition.",
	)

	parsedCmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch parsedCmd {
	case hearRegisterCommand.FullCommand():
		hearnoevil.RegisterCallbackURL(*hearRegisterCommandURLFlag)

	case hearServerCommand.FullCommand():
		hearnoevil.StartTranscriptionServer()

	case hearRecognize.FullCommand():
		if *hearRecognizeEpisodeFlag == 0 {
			hearnoevil.CreateSeasonRecognitions(
				*hearRecognizeSeasonFlag,
				*hearRecognizeURLFlag,
			)
		} else {
			hearnoevil.CreateEpisodeRecognition(
				*hearRecognizeSeasonFlag,
				*hearRecognizeEpisodeFlag,
				*hearRecognizeURLFlag,
			)
		}

	case diaryCommand.FullCommand():
		if *diaryCommandMissingFlag {
			videodiary.LogMissingEpisodes()
		} else {
			videodiary.DownloadEpisodes()
		}

	case visZeroCommand.FullCommand():
		visibilityzero.ExtractAudio()
	}
}
