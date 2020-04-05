package hearnoevil

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/visibilityzero"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
	stv1 "github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

// Recognition represents a speech-to-text service job.
type Recognition struct {
	*whodunit.Episode
}

// NewRecognition returns a new instance of recognition.
func NewRecognition(ep *whodunit.Episode) *Recognition {
	return &Recognition{
		Episode: ep,
	}
}

// StartJob starts a new recognition job.
func (r *Recognition) StartJob(stt *s2tInstance, callbackURL string) {
	if r.Exists() {
		log.WithField("file", r.FileName()).Infoln(
			"Skipping job, already exists")
		return
	}

	a := visibilityzero.NewAudio(r.Episode)
	if !a.Exists() {
		log.WithField("file", a.FileName()).Warnln(
			"Skipping job, audio file not found")
		return
	}

	audio := a.Open()
	if audio == nil {
		return
	}

	log.WithFields(logrus.Fields{
		"season":  r.SeasonNumber,
		"episode": r.EpisodeNumber,
	}).Infoln("Creating Recognition job")
	_, _, err := stt.CreateJob(r.jobOptions(audio, callbackURL))
	if err != nil {
		log.WithError(err).Errorln("Error creating job")
		return
	}

	log.Infoln("Job successfully created")
}

func (r *Recognition) jobOptions(
	audio *os.File,
	callbackURL string,
) *stv1.CreateJobOptions {
	return &stv1.CreateJobOptions{
		Audio:           audio,
		ContentType:     core.StringPtr("audio/mp3"),
		CallbackURL:     core.StringPtr(callbackURL),
		UserToken:       core.StringPtr(r.Name()),
		Events:          core.StringPtr("recognitions.completed_with_results"),
		ProfanityFilter: core.BoolPtr(false),
		SmartFormatting: core.BoolPtr(true),
	}
}

// WriteResults writes the specified contents to a new JSON file in the
// `/recognitions` directory.
func (r *Recognition) WriteResults(contents interface{}) error {
	path := r.AssetFilePath(whodunit.AssetTypeRecognition)
	return crimeseen.WriteJSONFile(path, contents)
}

// ReadResults returns the results from the recognition JSON file.
func (r *Recognition) ReadResults() (
	[]stv1.SpeechRecognitionResult,
	error,
) {
	path := r.FilePath()
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var contents []stv1.SpeechRecognitionResults
	err = json.Unmarshal(bytes, &contents)
	if err != nil {
		return nil, err
	}

	return contents[0].Results, nil
}

// Exists return true if the audio file exists in the `/assets` directory.
func (r *Recognition) Exists() bool {
	return r.AssetExists(whodunit.AssetTypeRecognition)
}

// FilePath returns the path to the audio file in the `/assets` directory.
func (r *Recognition) FilePath() string {
	return r.AssetFilePath(whodunit.AssetTypeRecognition)
}

// FileName returns the name of the audio file in the `/assets` directory.
func (r *Recognition) FileName() string {
	return r.AssetFileName(whodunit.AssetTypeRecognition)
}
