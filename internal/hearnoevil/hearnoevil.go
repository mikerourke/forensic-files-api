// Package hearnoevil sends the audio files to the speech-to-text service for
// recognition.
package hearnoevil

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/sirupsen/logrus"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

// Perpetrator contains properties and methods used to start recognition
// jobs.
type Perpetrator struct {
	STT         *sttService
	CallbackURL string
}

var log = waterlogged.ServiceLogger("hearnoevil")

// NewPerpetrator returns a new instance of Perpetrator.
func NewPerpetrator(callbackURL string) *Perpetrator {
	env := crimeseen.NewEnv()
	sts := newSTTService(env)
	p := &Perpetrator{STT: sts}

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
	result, _, err := p.STT.RegisterCallback(
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

// LogRecognitionJobs logs the last 100 jobs to the console. It displays all
// of the fields in JSON format.
func (p *Perpetrator) LogRecognitionJobs() {
	result, _, err := p.STT.CheckJobs(&speechtotextv1.CheckJobsOptions{})
	if err != nil {
		log.WithError(err).Fatalln("Error getting recognition jobs")
	}

	recognitionJobs := result.Recognitions
	b, err := json.MarshalIndent(recognitionJobs, "", "  ")
	log.Println(string(b))
}

// CreateSeasonRecognitionJobs batches calls to the speech-to-text service to
// create recognition jobs for all of the episodes in the specified season.
// You must specify a season number, but the callback URL is read from the `.env`
// file if it isn't passed into the function.
func (p *Perpetrator) CreateSeasonRecognitionJobs(season int) {
	p.checkNgrok()

	r := newRecognition(season, 0)
	err := filepath.Walk(
		r.SeasonDir(),
		func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".mp3") {
				r.SetAudioFilePath(path)
				r.StartJob(p.STT, p.CallbackURL)
			}
			return nil
		},
	)

	if err != nil {
		log.WithFields(logrus.Fields{
			"error":  err,
			"season": season,
		}).Fatalln("Error creating recognition")
	}
}

// CreateEpisodeRecognitionJob makes a call to the speech-to-text service to
// create a recognition job for a single episode in the specified season. You
// must specify a season and episode number, but the callback URL is read from
// the `.env` file if it isn't passed into the function.
func (p *Perpetrator) CreateEpisodeRecognitionJob(season int, episode int) {
	p.checkNgrok()

	r := newRecognition(season, episode)
	r.StartJob(p.STT, p.CallbackURL)
}

// StartCallbackServer starts the callback server to receive responses from the
// speech-to-text service.
func (p *Perpetrator) StartCallbackServer() {
	cs := newCallbackServer()
	cs.Start()
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
