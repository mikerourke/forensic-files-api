package tagasuspect

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/killigraphy"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
	nluv1 "github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

type Analysis struct {
	*whodunit.Episode
	detective *Detective
	assetType whodunit.AssetType
}

type AnalysisEntity struct {
	Name     string  `json:"name"`
	Type     string  `json:"entityType"`
	Salience float32 `json:"salience"`
}

// newAnalysis returns a new instance of an analysis.
func newAnalysis(ep *whodunit.Episode, d *Detective) *Analysis {
	return &Analysis{
		Episode:   ep,
		detective: d,
		assetType: assetTypeForCloudService(d.cloudService),
	}
}

// WriteCSV converts the entities to CSV records and writes the results to a file.
func (a *Analysis) WriteCSV(outputDir string) {
	if !a.Exists() {
		log.WithField("file", a.FileName()).Warnln(
			"Analysis does not exist, skipping")
		return
	}

	entities, err := a.ReadResults()
	if err != nil {
		log.WithFields(logrus.Fields{
			"file":  a.FileName(),
			"error": err,
		}).Fatalln("Unable to read the JSON file")
	}

	records := make([][]string, 0)
	records = append(records, []string{"name", "type", "salience"})

	for _, entity := range entities {
		salience := strconv.FormatFloat(float64(entity.Salience), 'f', 11, 32)
		records = append(records, []string{entity.Name, entity.Type, salience})
	}

	f, err := os.Create(a.csvFilePath(outputDir))
	if err != nil {
		log.WithError(err).Fatalln("Error creating CSV file")
	}

	defer f.Close()
	w := csv.NewWriter(f)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.WithError(err).Fatalln("Error writing record to CSV")
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalln(err)
	}

	log.WithField("file", a.FileName()).Infoln(
		"Successfully created CSV file")
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

	var result interface{}
	var err error
	if a.detective.cloudService == CloudServiceGCP {
		result, err = a.gcpAPIResult(t.Read())
	} else {
		result, err = a.ibmAPIResult(t.Read())
	}
	if err != nil {
		log.WithError(err).Errorln("Error submitting analysis request")
	}

	if err := crimeseen.WriteJSONFile(a.FilePath(), result); err != nil {
		log.WithError(err).Errorln("Error writing analysis file")
		return
	}

	log.Infoln("Analysis successfully written")
}

func (a *Analysis) gcpAPIResult(contents string) (interface{}, error) {
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

	analysisEntities := make([]AnalysisEntity, 0)
	for _, entity := range resp.Entities {
		analysisEntity := AnalysisEntity{
			Name:     entity.Name,
			Type:     entity.Type.String(),
			Salience: entity.Salience,
		}
		analysisEntities = append(analysisEntities, analysisEntity)
	}

	return analysisEntities, nil
}

func (a *Analysis) ibmAPIResult(contents string) (interface{}, error) {
	result, _, err := a.detective.service.Analyze(
		&nluv1.AnalyzeOptions{
			Text: &contents,
			Features: &nluv1.Features{
				Relations: &nluv1.RelationsOptions{},
				Entities: &nluv1.EntitiesOptions{
					Limit: core.Int64Ptr(10000),
				},
				Categories: &nluv1.CategoriesOptions{
					Limit: core.Int64Ptr(100),
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ReadResults returns the results of the analysis as an array of entity
// records.
func (a *Analysis) ReadResults() ([]AnalysisEntity, error) {
	path := a.FilePath()
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var contents []AnalysisEntity
	err = json.Unmarshal(bytes, &contents)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (a *Analysis) csvFilePath(outputDir string) string {
	fileName := strings.Replace(a.FileName(), ".json", ".csv", -1)
	return filepath.Join(outputDir, fileName)
}

// Exists return true if the analysis file exists in the `/assets` directory.
func (a *Analysis) Exists() bool {
	return a.AssetExists(a.assetType)
}

// FilePath returns the path to the analysis file in the `/assets` directory.
func (a *Analysis) FilePath() string {
	return a.AssetFilePath(a.assetType)
}

// FileName returns the name of the analysis file in the `/assets` directory.
func (a *Analysis) FileName() string {
	return a.AssetFileName(a.assetType)
}
