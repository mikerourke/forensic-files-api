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
		"Register a callback URL.").Alias("rcb")

	registerCommandURLFlag := registerCommand.Arg(
		"url",
		"Callback URL to use for transcribing.").Required().String()

	serverCommand := app.Command(
		"server",
		"Start the callback URL server (required to start transcribing).",
	).Alias("s")

	recognizeCommand := app.Command(
		"recognize",
		"Send recognition job requests to the speech to text service.",
	).Alias("rec")
	recogSeason, recogEpisode := addSeasonEpisodeFlags(recognizeCommand)

	investigateCommand := app.Command(
		"investigate",
		"Log status of asset.").Alias("log")

	investigateCommandAssetFlag := investigateCommand.Flag(
		"asset",
		"Asset to log.",
	).Short('a').Required().Enum("audio", "video", "recog", "trans")

	investigateCommandFilterFlag := investigateCommand.Flag(
		"filter",
		"Type to filter by.",
	).Short('f').Enum("pending", "complete", "in-process", "missing")

	downloadCommand := app.Command(
		"download",
		"Download episodes from YouTube.").Alias("dl")
	dlSeason, dlEpisode := addSeasonEpisodeFlags(downloadCommand)

	extractCommand := app.Command(
		"extract",
		"Extract audio from downloaded episodes for recognition.").Alias("ext")
	exSeason, exEpisode := addSeasonEpisodeFlags(extractCommand)

	transcribeCommand := app.Command(
		"transcribe",
		"Transcribes episode from recognition.").Alias("tr")
	transSeason, transEpisode := addSeasonEpisodeFlags(transcribeCommand)

	parsedCmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	p := hearnoevil.NewPerpetrator("")
	switch parsedCmd {
	case registerCommand.FullCommand():
		p.RegisterCallbackURL(*registerCommandURLFlag)

	case serverCommand.FullCommand():
		p.StartCallbackServer()

	case recognizeCommand.FullCommand():
		p.Recognize(*recogSeason, *recogEpisode)

	case investigateCommand.FullCommand():
		status := flagToAssetStatus(*investigateCommandFilterFlag)
		switch *investigateCommandAssetFlag {
		case "audio":
			visibilityzero.Investigate(status)
		case "recog":
			p.Investigate(status)
		case "trans":
			killigraphy.Investigate(status)
		case "video":
			videodiary.Investigate(status)
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
		"Season number to process.").Short('s').Int()
	episodeFlag = command.Flag(
		"episode",
		"Episode number to process.").Short('e').Int()
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
