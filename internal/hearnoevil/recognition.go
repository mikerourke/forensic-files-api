package hearnoevil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/sirupsen/logrus"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

type recognition struct {
	Season    int
	Episode   int
	audioPath string
}

func newRecognition(season int, episode int) *recognition {
	r := &recognition{Season: season, Episode: episode}
	r.audioPath = r.findAudioFile()
	return r
}

func (r *recognition) StartJob(stt *sttService, callbackURL string) {
	if r.wasCompleted() {
		log.WithField("file", r.audioFileName()).Infoln(
			"Skipping job, already exists")
		return
	}

	audio, err := os.Open(r.audioPath)
	if err != nil {
		log.WithError(err).Errorln("Error opening audio")
		return
	}

	log.WithFields(logrus.Fields{
		"episode": r.Episode,
		"season":  r.Season,
	}).Infoln("Creating recognition job")
	result, _, err := stt.CreateJob(
		&speechtotextv1.CreateJobOptions{
			Audio:           audio,
			ContentType:     core.StringPtr("audio/mp3"),
			CallbackURL:     core.StringPtr(callbackURL),
			UserToken:       core.StringPtr(r.jobName()),
			Events:          core.StringPtr("recognitions.completed_with_results"),
			ProfanityFilter: core.BoolPtr(false),
		},
	)
	if err != nil {
		log.WithError(err).Errorln("Error creating job")
		return
	}

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.WithError(err).Errorln("Error marshalling JSON")
	}

	fmt.Println(string(b))
}

func (r *recognition) SetAudioFilePath(path string) {
	r.audioPath = path
}

func (r *recognition) SeasonDir() string {
	s := fmt.Sprintf("season-%d", r.Season)
	return filepath.Join(crimeseen.AudioDirPath, s)
}

func (r *recognition) wasCompleted() bool {
	return crimeseen.FileExists(r.resultsPath())
}

func (r *recognition) resultsPath() string {
	jf := strings.Replace(r.audioFileName(), ".mp3", ".json", -1)
	return filepath.Join(crimeseen.RecognitionsDirPath, jf)
}

func (r *recognition) jobName() string {
	return strings.Replace(r.audioFileName(), ".mp3", "", -1)
}

func (r *recognition) audioFileName() string {
	return filepath.Base(r.audioPath)
}

func (r *recognition) findAudioFile() string {
	sPre := crimeseen.PaddedNumberString(r.Season)
	ePre := crimeseen.PaddedNumberString(r.Episode)
	fullPre := sPre + "-" + ePre

	ep := ""
	filepath.Walk(
		r.SeasonDir(),
		func(path string, info os.FileInfo, err error) error {
			if strings.HasPrefix(filepath.Base(path), fullPre) {
				ep = path
			}

			return nil
		},
	)

	return ep
}
