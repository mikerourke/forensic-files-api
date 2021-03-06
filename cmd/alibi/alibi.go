package main

import (
	"os"
	"path/filepath"

	"github.com/mikerourke/forensic-files-api/internal/hearnoevil"
	"github.com/mikerourke/forensic-files-api/internal/killigraphy"
	"github.com/mikerourke/forensic-files-api/internal/tagasuspect"
	"github.com/mikerourke/forensic-files-api/internal/videodiary"
	"github.com/mikerourke/forensic-files-api/internal/visibilityzero"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "Internal tools for the Forensic Files API.")
	app.HelpFlag.Short('h')

	overwriteFlag := app.Flag(
		"overwrite",
		"Overwrite existing file").Short('x').Bool()

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

	investigateAssetFlag := investigateCommand.Flag(
		"asset",
		"Asset to log.",
	).Short('a').Required().Enum("analysis", "audio", "video", "recog", "trans")

	investigateServiceFlag := investigateCommand.Flag(
		"service",
		"Service to use for the analysis.",
	).Default("gcp").Short('u').Enum("gcp", "ibm")

	investigateFilterFlag := investigateCommand.Flag(
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

	analyzeCommand := app.Command("analyze",
		"Create a new entity analysis.").Alias("an")
	analyzeSeason, analyzeEpisode := addSeasonEpisodeFlags(analyzeCommand)

	analyzeServiceFlag := analyzeCommand.Flag(
		"service",
		"Service to use for the analysis.",
	).Short('u').Required().Enum("gcp", "ibm")

	analyzeCSVFlag := analyzeCommand.Flag(
		"csv",
		"Output a CSV file to the specified directory.").Short('c').ExistingDir()

	parsedCmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	ew := hearnoevil.NewEyewitness("")
	d := tagasuspect.NewDetective()
	switch parsedCmd {
	case registerCommand.FullCommand():
		ew.RegisterCallbackURL(*registerCommandURLFlag)

	case serverCommand.FullCommand():
		ew.StartCallbackServer()

	case recognizeCommand.FullCommand():
		ew.Recognize(*recogSeason, *recogEpisode)

	case investigateCommand.FullCommand():
		status := whodunit.AssetStatusAny
		if investigateFilterFlag != nil {
			status = flagToAssetStatus(*investigateFilterFlag)
		}
		switch *investigateAssetFlag {
		case "analysis":
			cloudService := flagToCloudService(*investigateServiceFlag)
			tagasuspect.Investigate(cloudService, status)
		case "audio":
			visibilityzero.Investigate(status)
		case "recog":
			ew.Investigate(status)
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

	case analyzeCommand.FullCommand():
		season := *analyzeSeason
		episode := *analyzeEpisode
		if *analyzeCSVFlag != "" {
			d.FileReport(season, episode, *analyzeCSVFlag)
		} else {
			cloudService := flagToCloudService(*analyzeServiceFlag)
			d.OpenCase(cloudService)
			defer d.CloseCase()
			d.Analyze(season, episode, *overwriteFlag)
		}
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

func flagToCloudService(value string) tagasuspect.CloudService {
	if value == "gcp" {
		return tagasuspect.CloudServiceGCP
	}
	return tagasuspect.CloudServiceIBM
}
