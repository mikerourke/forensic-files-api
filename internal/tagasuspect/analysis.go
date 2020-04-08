package tagasuspect

import (
	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/killigraphy"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	nluv1 "github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
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
func (a *Analysis) Create(svc *nluv1.NaturalLanguageUnderstandingV1, overwrite bool) {
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

	contents := t.Read()
	result, _, err := svc.Analyze(
		&nluv1.AnalyzeOptions{
			Text: &contents,
			Features: &nluv1.Features{
				Categories: &nluv1.CategoriesOptions{},
				Entities: &nluv1.EntitiesOptions{
					Limit: core.Int64Ptr(10000),
				},
				Keywords: &nluv1.KeywordsOptions{
					Limit: core.Int64Ptr(10000),
				},
			},
		},
	)
	if err != nil {
		log.WithError(err).Errorln("Error submitting analysis request")
	}

	if err := crimeseen.WriteJSONFile(a.FilePath(), result); err != nil {
		log.WithError(err).Errorln("Error writing analysis file")
		return
	}

	log.Infoln("Analysis successfully written")
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
