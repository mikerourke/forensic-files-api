// Package hearnoevil sends the audio files to the speech-to-text service for
// recognition.
package hearnoevil

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
	stv1 "github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

// Eyewitness contains properties and methods used to start recognition
// jobs.
type Eyewitness struct {
	s2t         *s2tInstance
	callbackURL string
}

var log = waterlogged.New("hearnoevil")

// NewEyewitness returns a new instance of Eyewitness.
func NewEyewitness(callbackURL string) *Eyewitness {
	env := crimeseen.NewEnv()
	sts := newS2TInstance(env)
	ew := &Eyewitness{s2t: sts}

	if callbackURL != "" {
		ew.RegisterCallbackURL(callbackURL)
	} else {
		callbackURL = env.CallbackURL()
	}
	ew.callbackURL = callbackURL

	return ew
}

// RegisterCallbackURL registers a callback URL with the speech to text service
// that will receive recognition job responses.
func (ew *Eyewitness) RegisterCallbackURL(callbackURL string) {
	result, _, err := ew.s2t.RegisterCallback(
		&stv1.RegisterCallbackOptions{
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
func (ew *Eyewitness) Recognize(seasonNumber int, episodeNumber int) {
	ew.interrogate()

	onEpisode := func(ep *whodunit.Episode) {
		r := NewRecognition(ep)
		r.StartJob(ew.s2t, ew.callbackURL)
	}

	if err := whodunit.Solve(seasonNumber, episodeNumber, onEpisode); err != nil {
		log.WithError(err).Errorln("Error recognizing episode(s)")
	}
}

// Investigate logs the episode statuses.
func (ew *Eyewitness) Investigate(status whodunit.AssetStatus) {
	totalCount := 0
	table := whodunit.NewStatusTable(whodunit.AssetTypeRecognition, status)
	jep := ew.jobEpisodeMap()
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
func (ew *Eyewitness) StartCallbackServer() {
	cs := newCallbackServer()
	cs.Start()
}

func (ew *Eyewitness) jobEpisodeMap() map[string]*whodunit.Episode {
	result, _, err := ew.s2t.CheckJobs(&stv1.CheckJobsOptions{})
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

func (ew *Eyewitness) interrogate() {
	const ErrorMessage = "ngrok is not running, run `ngrok http 9000`"

	cmd := exec.Command("ps", "aux")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.WithError(err).Fatalln(ErrorMessage)
	}

	// Start the command after having set up the pipe.
	if err := cmd.Start(); err != nil {
		log.WithError(err).Fatalln(ErrorMessage)
	}

	// Read command's stdout line by line.
	in := bufio.NewScanner(stdout)

	for in.Scan() {
		text := in.Text()
		if strings.Contains(text, "ngrok http 9000") {
			return
		}
	}

	log.WithError(err).Fatalln(ErrorMessage)
}
