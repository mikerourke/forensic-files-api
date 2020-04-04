// Package hearnoevil sends the audio files to the speech-to-text service for
// recognition.
package hearnoevil

import (
	"bufio"
	"errors"
	"os/exec"
	"strings"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

// Perpetrator contains properties and methods used to start recognition
// jobs.
type Perpetrator struct {
	S2T         *s2tInstance
	CallbackURL string
}

var log = waterlogged.New("hearnoevil")

// NewPerpetrator returns a new instance of Perpetrator.
func NewPerpetrator(callbackURL string) *Perpetrator {
	env := crimeseen.NewEnv()
	sts := newS2TInstance(env)
	p := &Perpetrator{S2T: sts}

	if callbackURL != "" {
		p.RegisterCallbackURL(callbackURL)
	} else {
		callbackURL = env.CallbackURL()
	}
	p.CallbackURL = callbackURL

	return p
}

// RegisterCallbackURL registers a callback URL with the speech to text service
// that will receive recognition job responses.
func (p *Perpetrator) RegisterCallbackURL(callbackURL string) {
	result, _, err := p.S2T.RegisterCallback(
		&speechtotextv1.RegisterCallbackOptions{
			CallbackURL: core.StringPtr(callbackURL),
		},
	)

	if result != nil {
		log.WithFields(logrus.Fields{
			"url":    *result.URL,
			"status": *result.Status,
		}).Infoln("Callback registration complete")
	}

	if err != nil {
		log.WithError(err).Fatalln("Error registering callback URL")
	}

	log.WithField("url", callbackURL).Infoln("Callback URL registered")
}

// Recognize makes a call to the speech-to-text service to create a recognition
// job for a single episode in the specified season or all episodes if the season
// was not specified.
func (p *Perpetrator) Recognize(seasonNumber int, episodeNumber int) {
	p.checkNgrok()

	if seasonNumber == 0 {
		log.Fatalln("You must specify a season number")
	}

	s := whodunit.NewSeason(seasonNumber)
	if err := s.PopulateEpisodes(); err != nil {
		log.WithError(err).Fatalln("Could not get season episodes")
	}

	if episodeNumber == 0 {
		p.recognizeSeason(s)
	} else {
		ep := s.Episode(episodeNumber)
		p.recognizeEpisode(ep)
	}
}

// LogStatusTable logs the episode statuses.
func (p *Perpetrator) LogStatusTable(status whodunit.AssetStatus) {
	totalCount := 0
	table := whodunit.NewStatusTable(whodunit.AssetTypeRecognition, status)
	jep := p.jobEpisodeMap()
	for season := 1; season <= whodunit.SeasonCount; season++ {
		s := whodunit.NewSeason(season)
		if err := s.PopulateEpisodes(); err != nil {
			panic("Could not get season episodes")
		}

		for _, ep := range s.AllEpisodes() {
			je := jep[ep.Name()]
			if je != nil {
				ep.SetAssetStatus(je.AssetStatus(whodunit.AssetTypeRecognition))
			}

			if table.AddRow(ep) {
				totalCount++
			}
		}
	}

	table.RenderTable(totalCount)
}

// StartCallbackServer starts the callback server to receive responses from the
// speech-to-text service.
func (p *Perpetrator) StartCallbackServer() {
	cs := newCallbackServer()
	cs.Start()
}

func (p *Perpetrator) recognizeSeason(s *whodunit.Season) {
	for _, ep := range s.AllEpisodes() {
		p.recognizeEpisode(ep)
	}
}

func (p *Perpetrator) recognizeEpisode(ep *whodunit.Episode) {
	r := NewRecognition(ep)
	r.StartJob(p.S2T, p.CallbackURL)
}

func (p *Perpetrator) jobEpisodeMap() map[string]*whodunit.Episode {
	result, _, err := p.S2T.CheckJobs(&speechtotextv1.CheckJobsOptions{})
	if err != nil {
		log.WithError(err).Fatalln("Error getting recognition jobs")
	}

	epMap := make(map[string]*whodunit.Episode, 0)
	for _, job := range result.Recognitions {
		name := *job.UserToken
		ep, err := whodunit.NewEpisodeFromName(name)
		if err != nil {
			log.WithError(err).Fatalln("Error parsing episode name")
		}

		if strings.Contains(*job.Status, "compl") {
			ep.SetAssetStatus(whodunit.AssetStatusComplete)
		} else {
			ep.SetAssetStatus(whodunit.AssetStatusInProcess)
		}

		epMap[name] = ep
	}

	return epMap
}

func (p *Perpetrator) checkNgrok() {
	err := p.findNgrokProcess()
	if err != nil {
		log.WithError(err).Fatalln("ngrok is not running, run `ngrok http 9000`")
	}
}

func (p *Perpetrator) findNgrokProcess() error {
	cmd := exec.Command("ps", "aux")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Start the command after having set up the pipe.
	if err := cmd.Start(); err != nil {
		return err
	}

	// Read command's stdout line by line.
	in := bufio.NewScanner(stdout)

	for in.Scan() {
		text := in.Text()
		if strings.Contains(text, "ngrok http 9000") {
			return nil
		}
	}

	return errors.New("ngrok is not running")
}
