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

func CreateSeasonRecognitions(seasonNumber int, callbackURL string) {
	speechToText := speechToTextService()
	validCallbackURL := validateCallbackURL(callbackURL)

	err := filepath.Walk(seasonsPath(seasonNumber),
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
		})

	if err != nil {
		log.WithFields(logrus.Fields{
			"error":  err,
			"season": seasonNumber,
		}).Fatal("Error creating recognition")
	}
}

func CreateEpisodeRecognition(
	seasonNumber int,
	episodeNumber int,
	callbackURL string,
) {
	err := ensureNgrokIsRunning()
	if err != nil {
		log.WithField("error", err).Fatal(
			"ngrok is not running, run `ngrok http 9000`",
		)
	}

	speechToText := speechToTextService()
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

func ensureNgrokIsRunning() error {
	cmd := exec.Command("ps", "aux")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// start the command after having set up the pipe
	if err := cmd.Start(); err != nil {
		return err
	}

	// read command's stdout line by line
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

func findAudioFilePath(
	seasonNumber int,
	episodeNumber int,
) (audioFilePath string, findErr error) {
	seasonPrefix := crimeseen.PaddedNumberString(seasonNumber)
	episodePrefix := crimeseen.PaddedNumberString(episodeNumber)
	fullPrefix := seasonPrefix + "-" + episodePrefix

	audioFilePath = ""

	err := filepath.Walk(seasonsPath(seasonNumber),
		func(path string, info os.FileInfo, err error) error {
			if strings.HasPrefix(filepath.Base(path), fullPrefix) {
				audioFilePath = path
			}

			return nil
		})

	if err != nil {
		findErr = err
	}

	return audioFilePath, findErr
}

func hasExistingRecognition(audioFilePath string) bool {
	fileName := filepath.Base(audioFilePath)
	jsonFileName := strings.Replace(fileName, ".mp3", ".json", -1)
	recognitionPath := filepath.Join(crimeseen.RecognitionsPath, jsonFileName)
	return crimeseen.FileExists(recognitionPath)
}

func seasonsPath(seasonNumber int) string {
	season := strconv.Itoa(seasonNumber)
	return filepath.Join(crimeseen.AudioPath, "season-"+season)
}
