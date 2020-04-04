package killigraphy

import (
	"io"
	"os"

	"github.com/mikerourke/forensic-files-api/internal/hearnoevil"
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
)

var log = waterlogged.New("killigraphy")

// Transcribe creates a transcript for the specified episode number from the
// specified season number or all seasons.
func Transcribe(seasonNumber int, episodeNumber int) {
	if seasonNumber == 0 {
		if episodeNumber != 0 {
			log.Fatalln("You must specify a season number for an episode")
		}

		transcribeAllSeasons()
		return
	}

	s := whodunit.NewSeason(seasonNumber)
	if err := s.PopulateEpisodes(); err != nil {
		log.WithError(err).Fatalln("Could not get season episodes")
	}

	if episodeNumber == 0 {
		transcribeSeason(s)
	} else {
		ep := s.Episode(episodeNumber)
		transcribeEpisode(ep)
	}
}

// LogStatusTable logs the episode statuses.
func LogStatusTable(status whodunit.AssetStatus) {
	table := whodunit.NewStatusTable(whodunit.AssetTypeTranscript, status)
	table.Log()
}

func transcribeAllSeasons() {
	for season := 1; season <= whodunit.SeasonCount; season++ {
		s := whodunit.NewSeason(season)
		if err := s.PopulateEpisodes(); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"season": season,
			}).Fatalln("Could not get season episodes")
			return
		}
		transcribeSeason(s)
	}
}

func transcribeSeason(s *whodunit.Season) {
	for _, ep := range s.AllEpisodes() {
		transcribeEpisode(ep)
	}
}

func transcribeEpisode(ep *whodunit.Episode) {
	epLogger := log.WithFields(logrus.Fields{
		"season":  ep.SeasonNumber,
		"episode": ep.EpisodeNumber,
	})

	// TODO: Uncomment this, we're overwriting for the time being.
	// if ep.AssetExists(whodunit.AssetTypeTranscript) {
	// 	epLogger.Infoln("Transcript already exists, skipping")
	// 	return
	// }

	if !ep.AssetExists(whodunit.AssetTypeRecognition) {
		return
	}

	epLogger.Infoln("Transcribing episode")
	r := hearnoevil.NewRecognition(ep)
	results, err := r.ReadResults()
	if err != nil {
		log.WithError(err).Fatalln("Error getting recognition results")
	}

	words := make([]string, 0)
	for _, result := range results {
		for _, alt := range result.Alternatives {
			words = append(words, *alt.Transcript)
		}
	}

	l := newLuminol(words)
	path := ep.AssetFilePath(whodunit.AssetTypeTranscript)
	file, err := os.Create(path)
	if err != nil {
		log.WithError(err).Fatalln("Error creating file")
	}
	defer file.Close()

	if _, err = io.WriteString(file, l.Reveal()); err != nil {
		log.WithError(err).Fatalln("Error writing contents")
	}

	if err := file.Sync(); err != nil {
		log.WithError(err).Fatalln("Error syncing file")
	}
}
