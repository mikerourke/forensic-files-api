package tagasuspect

import (
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/killigraphy"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Analysis represents the entity analysis from the NLP service.
type Analysis struct {
	*whodunit.Episode
	detective *Detective
}

// newAnalysis returns a new instance of an analysis.
func newAnalysis(ep *whodunit.Episode, d *Detective) *Analysis {
	return &Analysis{
		Episode:   ep,
		detective: d,
	}
}

// Create creates a new analysis file by sending the transcript to the NLP
// service and writing the results to the `/assets` directory.
func (a *Analysis) Create(overwrite bool) {
	t := killigraphy.NewTranscript(a.Episode)
	if !t.Exists() {
		log.WithField("file", t.FileName()).Warnln(
			"Transcript not found, skipping")
		return
	}

	if a.Exists() && !overwrite {
		log.WithField("file", a.FileName()).Warnln(
			"Analysis already exists, skipping")
		return
	}

	log.WithField("file", a.FileName()).Infoln("Starting analysis")

	result, err := a.apiResult(t.Read())
	if err != nil {
		log.WithError(err).Errorln("Error submitting analysis request")
	}

	if err := crimeseen.WriteJSONFile(a.FilePath(), result); err != nil {
		log.WithError(err).Errorln("Error writing analysis file")
		return
	}

	log.Infoln("Analysis successfully written")
}

func (a *Analysis) apiResult(contents string) (interface{}, error) {
	doc := &languagepb.Document{
		Type: languagepb.Document_PLAIN_TEXT,
		Source: &languagepb.Document_Content{
			Content: contents,
		},
	}

	req := &languagepb.AnalyzeEntitiesRequest{
		Document:     doc,
		EncodingType: languagepb.EncodingType_UTF8,
	}

	resp, err := a.detective.client.AnalyzeEntities(a.detective.ctx, req)
	if err != nil {
		return nil, err
	}

	type JSONEntity struct {
		Name     string  `json:"name"`
		Type     string  `json:"entityType"`
		Salience float32 `json:"salience"`
	}

	jsonEntities := make([]JSONEntity, 0)
	for _, entity := range resp.Entities {
		jsonEntity := JSONEntity{
			Name:     entity.Name,
			Type:     entity.Type.String(),
			Salience: entity.Salience,
		}
		jsonEntities = append(jsonEntities, jsonEntity)
	}

	return jsonEntities, nil
}

// Exists return true if the analysis file exists in the `/assets` directory.
func (a *Analysis) Exists() bool {
	return crimeseen.FileExists(a.FilePath())
}

// FilePath returns the path to the analysis file in the `/assets` directory.
func (a *Analysis) FilePath() string {
	return a.AssetFilePath(whodunit.AssetTypeAnalysis)
}

// FileName returns the name of the analysis file in the `/assets` directory.
func (a *Analysis) FileName() string {
	return a.AssetFileName(whodunit.AssetTypeAnalysis)
}
