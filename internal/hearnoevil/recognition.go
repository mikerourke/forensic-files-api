package hearnoevil

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

type recognition struct {
	*whodunit.Episode
}

func newRecognition(ep *whodunit.Episode) *recognition {
	return &recognition{
		Episode: ep,
	}
}

func (r *recognition) StartJob(stt *s2tInstance, callbackURL string) {
	if r.AssetExists(whodunit.AssetTypeRecognition) {
		log.WithField(
			"file",
			r.AssetFileName(whodunit.AssetTypeRecognition),
		).Infoln("Skipping job, already exists")
		return
	}

	if !r.AssetExists(whodunit.AssetTypeAudio) {
		log.WithField(
			"file",
			r.AssetFileName(whodunit.AssetTypeAudio),
		).Warnln("Skipping job, audio file not found")
		return
	}

	audio, err := os.Open(r.AssetFilePath(whodunit.AssetTypeAudio))
	if err != nil {
		log.WithError(err).Errorln("Error opening audio")
		return
	}

	log.WithFields(logrus.Fields{
		"season":  r.SeasonNumber,
		"episode": r.EpisodeNumber,
	}).Infoln("Creating recognition job")
	result, _, err := stt.CreateJob(
		&speechtotextv1.CreateJobOptions{
			Audio:           audio,
			ContentType:     core.StringPtr("audio/mp3"),
			CallbackURL:     core.StringPtr(callbackURL),
			UserToken:       core.StringPtr(r.Name()),
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
