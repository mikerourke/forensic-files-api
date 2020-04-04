package hearnoevil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

type Recognition struct {
	*whodunit.Episode
}

func NewRecognition(ep *whodunit.Episode) *Recognition {
	return &Recognition{
		Episode: ep,
	}
}

func (r *Recognition) StartJob(stt *s2tInstance, callbackURL string) {
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
	}).Infoln("Creating Recognition job")
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

// WriteToFile writes the specified contents to a new JSON file in the
// `/recognitions` directory.
func (r *Recognition) WriteToFile(contents interface{}) error {
	path := r.AssetFilePath(whodunit.AssetTypeRecognition)
	return crimeseen.WriteJSONFile(path, contents)
}

// ReadResults returns the results from the recognition JSON file.
func (r *Recognition) ReadResults() (
	[]speechtotextv1.SpeechRecognitionResult,
	error,
) {
	path := r.AssetFilePath(whodunit.AssetTypeRecognition)
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var contents []speechtotextv1.SpeechRecognitionResults
	err = json.Unmarshal(bytes, &contents)
	if err != nil {
		return nil, err
	}

	return contents[0].Results, nil
}
