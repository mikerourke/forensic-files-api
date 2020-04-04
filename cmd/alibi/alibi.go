package main

import (
	"os"
	"path/filepath"

	"github.com/mikerourke/forensic-files-api/internal/hearnoevil"
	"github.com/mikerourke/forensic-files-api/internal/killigraphy"
	"github.com/mikerourke/forensic-files-api/internal/videodiary"
	"github.com/mikerourke/forensic-files-api/internal/visibilityzero"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "Internal tools for the Forensic Files API.")
	app.HelpFlag.Short('h')

	registerCommand := app.Command(
		"registercb",
		"Register a callback URL.")

	registerCommandURLFlag := registerCommand.Arg(
		"url",
		"Callback URL to use for transcribing.").Required().String()

	serverCommand := app.Command(
		"server",
		"Start the callback URL server (required to start transcribing).")

	recognizeCommand := app.Command(
		"recognize",
		"Send recognition job requests to the speech to text service.")
	recogSeason, recogEpisode := addSeasonEpisodeFlags(recognizeCommand)

	logCommand := app.Command(
		"log",
		"Log status of asset.")

	logCommandAssetFlag := logCommand.Flag(
		"asset",
		"Asset to log: audio, video, recog, trans.").Required().String()

	logCommandFilterFlag := logCommand.Flag(
		"filter",
		"Type to filter by: pending, complete, in-process, missing.").String()

	downloadCommand := app.Command(
		"download",
		"Download episodes from YouTube.")
	dlSeason, dlEpisode := addSeasonEpisodeFlags(downloadCommand)

	extractCommand := app.Command(
		"extract",
		"Extract audio from downloaded episodes for recognition.")
	exSeason, exEpisode := addSeasonEpisodeFlags(extractCommand)

	transcribeCommand := app.Command(
		"transcribe",
		"Transcribes episode from recognition.")
	transSeason, transEpisode := addSeasonEpisodeFlags(extractCommand)

	parsedCmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	p := hearnoevil.NewPerpetrator("")
	switch parsedCmd {
	case registerCommand.FullCommand():
		p.RegisterCallbackURL(*registerCommandURLFlag)

	case serverCommand.FullCommand():
		p.StartCallbackServer()

	case recognizeCommand.FullCommand():
		p.Recognize(*recogSeason, *recogEpisode)

	case logCommand.FullCommand():
		status := flagToAssetStatus(*logCommandFilterFlag)
		switch *logCommandAssetFlag {
		case "audio":
			visibilityzero.LogStatusTable(status)
		case "recog":
			p.LogStatusTable(status)
		case "trans":
			killigraphy.LogStatusTable(status)
		case "video":
			videodiary.LogStatusTable(status)
		}

	case downloadCommand.FullCommand():
		videodiary.Download(*dlSeason, *dlEpisode)

	case extractCommand.FullCommand():
		visibilityzero.ExtractAudio(*exSeason, *exEpisode)

	case transcribeCommand.FullCommand():
		killigraphy.Transcribe(*transSeason, *transEpisode)

	}
}

func addSeasonEpisodeFlags(
	command *kingpin.CmdClause,
) (seasonFlag *int, episodeFlag *int) {
	seasonFlag = command.Flag(
		"season",
		"Season number to process.").Int()

	episodeFlag = command.Flag(
		"episode",
		"Episode number to process.").Int()
	return
}

func flagToAssetStatus(value string) whodunit.AssetStatus {
	switch value {
	case "pending":
		return whodunit.AssetStatusPending
	case "in-process":
		return whodunit.AssetStatusInProcess
	case "complete":
		return whodunit.AssetStatusComplete
	case "missing":
		return whodunit.AssetStatusMissing
	}
	return whodunit.AssetStatusAny
}
