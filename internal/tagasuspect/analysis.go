package tagasuspect

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/killigraphy"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Analysis represents the entity analysis from the NLP service.
type Analysis struct {
	*whodunit.Episode
}

// NewAnalysis returns a new instance of an analysis.
func NewAnalysis(ep *whodunit.Episode) *Analysis {
	return &Analysis{
		Episode: ep,
	}
}

// Create creates a new analysis file by sending the transcript to the NLP
// service and writing the results to the `/assets` directory.
func (a *Analysis) Create(ctx context.Context, client *language.Client) {
	t := killigraphy.NewTranscript(a.Episode)
	if !t.Exists() {
		log.WithField("file", t.FileName()).Warnln(
			"Transcript not found, skipping")
		return
	}

	if a.Exists() {
		log.WithField("file", a.FileName()).Warnln(
			"Analysis already exists, skipping")
	}

	doc := &languagepb.Document{
		Type: languagepb.Document_PLAIN_TEXT,
		Source: &languagepb.Document_Content{
			Content: t.Read(),
		},
	}

	req := &languagepb.AnalyzeEntitiesRequest{
		Document:     doc,
		EncodingType: languagepb.EncodingType_UTF8,
	}

	resp, err := client.AnalyzeEntities(ctx, req)
	if err != nil {
		log.WithError(err).Errorln("Error analyzing entities")
		return
	}

	if err := crimeseen.WriteJSONFile(a.FilePath(), resp.Entities); err != nil {
		log.WithError(err).Errorln("Error writing analysis file")
	}
}

// Exists return true if the analysis file exists in the `/assets` directory.
func (a *Analysis) Exists() bool {
	return a.AssetExists(whodunit.AssetTypeAnalysis)
}

// FilePath returns the path to the analysis file in the `/assets` directory.
func (a *Analysis) FilePath() string {
	return a.AssetFilePath(whodunit.AssetTypeAnalysis)
}

// FileName returns the name of the analysis file in the `/assets` directory.
func (a *Analysis) FileName() string {
	return a.AssetFileName(whodunit.AssetTypeAnalysis)
}
