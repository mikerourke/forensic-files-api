// Package hearnoevil sends the audio files to the speech-to-text service for
// recognition.
package hearnoevil

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/sirupsen/logrus"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

var log = waterlogged.ServiceLogger("hearnoevil")

// RegisterCallbackURL registers a callback URL with the speech to text service
// that will receive recognition job responses.
func RegisterCallbackURL(callbackURL string) {
	speechToText := speechToTextService()

	result, _, err := speechToText.RegisterCallback(
		&speechtotextv1.RegisterCallbackOptions{
			CallbackURL: core.StringPtr(callbackURL),
		},
	)

	if result != nil {
		log.WithFields(logrus.Fields{
			"url":    *result.URL,
			"status": *result.Status,
		}).Info("Callback registration complete")
	}

	if err != nil {
		log.WithField("error", err).Fatal("Error registering callback URL")
	}

	log.WithField("url", callbackURL).Info("Callback URL registered")
}

// LogRecognitionJobs logs the last 100 jobs to the console. It displays all
// of the fields in JSON format.
func LogRecognitionJobs() {
	speechToText := speechToTextService()
	result, _, err := speechToText.CheckJobs(&speechtotextv1.CheckJobsOptions{})
	if err != nil {
		log.WithField("error", err).Fatal("Error getting recognition jobs")
	}

	recognitionJobs := result.Recognitions
	bytes, err := json.MarshalIndent(recognitionJobs, "", "  ")
	log.Print(string(bytes))
}

// CreateSeasonRecognitionJobs batches calls to the speech-to-text service to
// create recognition jobs for all of the episodes in the specified season.
// You must specify a season number, but the callback URL is read from the `.env`
// file if it isn't passed into the function.
func CreateSeasonRecognitionJobs(seasonNumber int, callbackURL string) {
	// Bails if ngrok isn't running.
	ensureNgrokIsRunning()

	speechToText := speechToTextService()

	// Bails if the callback URL is missing or invalid.
	validCallbackURL := validateCallbackURL(callbackURL)

	err := filepath.Walk(
		audioSeasonDirPath(seasonNumber),
		func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".mp3") {
				if hasExistingRecognition(path) {
					log.WithField(
						"file", filepath.Base(path),
					).Info("Skipping episode, already exists")
				} else {
					createEpisodeRecognitionJob(
						validCallbackURL,
						speechToText,
						path,
					)
				}
			}
			return nil
		},
	)

	if err != nil {
		log.WithFields(logrus.Fields{
			"error":  err,
			"season": seasonNumber,
		}).Fatal("Error creating recognition")
	}
}

// CreateEpisodeRecognitionJob makes a call to the speech-to-text service to
// create a recognition job for a single episode in the specified season. You
//must specify a season and episode number, but the callback URL is read from
// the `.env` file if it isn't passed into the function.
func CreateEpisodeRecognitionJob(
	seasonNumber int,
	episodeNumber int,
	callbackURL string,
) {
	// Bails if ngrok isn't running.
	ensureNgrokIsRunning()

	speechToText := speechToTextService()

	// Bails if the callback URL is missing or invalid.
	validCallbackURL := validateCallbackURL(callbackURL)

	audioFilePath, err := findAudioFilePath(seasonNumber, episodeNumber)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":   err,
			"season":  seasonNumber,
			"episode": episodeNumber,
		}).Fatal("Could not find audio file for episode and season")
	}

	createEpisodeRecognitionJob(validCallbackURL, speechToText, audioFilePath)
}

func validateCallbackURL(callbackURL string) string {
	if callbackURL == "" {
		callbackURL = os.Getenv("CALLBACK_URL")
	}

	if callbackURL == "" {
		log.Fatal("Could not find callback URL")
	}

	return callbackURL
}

func ensureNgrokIsRunning() {
	err := findNgrokProcess()
	if err != nil {
		log.WithField(
			"error", err,
		).Fatal("ngrok is not running, run `ngrok http 9000`")
	}
}

func findNgrokProcess() error {
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

func createEpisodeRecognitionJob(
	callbackURL string,
	speechToText *speechtotextv1.SpeechToTextV1,
	audioFilePath string,
) {
	log.WithField(
		"file", filepath.Base(audioFilePath),
	).Info("Starting recognition")

	var audioFile io.ReadCloser
	audioFile, err := os.Open(audioFilePath)

	fileName := filepath.Base(audioFilePath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"file":  fileName,
		}).Error("Error opening audio file")
	}

	jobName := strings.Replace(fileName, ".mp3", "", -1)

	result, _, err := speechToText.CreateJob(
		&speechtotextv1.CreateJobOptions{
			Audio:       audioFile,
			ContentType: core.StringPtr("audio/mp3"),
			CallbackURL: core.StringPtr(callbackURL),
			UserToken:   core.StringPtr(jobName),
			Events:      core.StringPtr("recognitions.completed_with_results"),
		},
	)

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"file":  fileName,
		}).Error("Error creating job")
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Error("Could not marshal JSON")
	}

	log.Println(string(bytes))
}

// hasExistingRecognition returns true if a recognition already exists for the
// specified audio file.
func hasExistingRecognition(audioFilePath string) bool {
	fileName := filepath.Base(audioFilePath)
	jsonFileName := strings.Replace(fileName, ".mp3", ".json", -1)
	recognitionPath := filepath.Join(crimeseen.RecognitionsDirPath, jsonFileName)
	return crimeseen.FileExists(recognitionPath)
}

func findAudioFilePath(
	seasonNumber int,
	episodeNumber int,
) (audioFilePath string, err error) {
	seasonPrefix := crimeseen.PaddedNumberString(seasonNumber)
	episodePrefix := crimeseen.PaddedNumberString(episodeNumber)
	fullPrefix := seasonPrefix + "-" + episodePrefix

	audioFilePath = ""

	err = filepath.Walk(
		audioSeasonDirPath(seasonNumber),
		func(path string, info os.FileInfo, err error) error {
			if strings.HasPrefix(filepath.Base(path), fullPrefix) {
				audioFilePath = path
			}

			return nil
		},
	)

	return audioFilePath, err
}

// audioSeasonDirPath return the full path to the specified season directory in
// the `/assets/audio` directory.
func audioSeasonDirPath(seasonNumber int) string {
	season := strconv.Itoa(seasonNumber)
	return filepath.Join(crimeseen.AudioDirPath, "season-"+season)
}
