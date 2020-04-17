package tagasuspect

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
	nluv1 "github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
	"google.golang.org/api/option"
)

// CloudService represents the cloud service associated with the analysis.
type CloudService int

const (
	// CloudServiceGCP represents the GCP cloud service.
	CloudServiceGCP CloudService = iota

	// CloudServiceIBM represents the IBM cloud service.
	CloudServiceIBM
)

// Detective contains properties and methods used to start entity analysis
// jobs.
type Detective struct {
	cloudService CloudService
	client       *language.Client
	ctx          context.Context
	withFile     option.ClientOption
	service      *nluv1.NaturalLanguageUnderstandingV1
}

var (
	env = crimeseen.NewEnv()
	log = waterlogged.New("tagasuspect")
)

// NewDetective returns a new instance of a detective.
func NewDetective() *Detective {
	return &Detective{}
}

// OpenCase creates a new instance of an NLP client. We're doing this here
// instead of each analysis call so we can persist the client for batch jobs.
// See https://cloud.google.com/natural-language/docs/reference/rest
func (d *Detective) OpenCase(cloudService CloudService) {
	d.cloudService = cloudService
	if cloudService == CloudServiceGCP {
		if env.GCPCredsPath() == "" {
			log.Fatalln("GCP credentials file not specified in .env file")
		}

		ctx := context.Background()
		d.ctx = ctx
		d.withFile = option.WithCredentialsFile(env.GCPCredsPath())

		client, err := language.NewClient(d.ctx, d.withFile)
		if err != nil {
			logrus.WithError(err).Fatalln("Could not create new GCP client")
		}
		d.client = client
		return
	}

	if cloudService == CloudServiceIBM {
		if env.IBMLangAPIKey() == "" || env.IBMLangAPIUrl() == "" {
			log.Fatalln("IBM credentials file not specified in .env file")
		}

		authenticator := &core.IamAuthenticator{
			ApiKey: env.IBMLangAPIKey(),
		}
		options := &nluv1.NaturalLanguageUnderstandingV1Options{
			Version:       "2019-07-12",
			Authenticator: authenticator,
		}

		svc, err := nluv1.NewNaturalLanguageUnderstandingV1(options)
		if err != nil {
			logrus.WithError(err).Fatalln("Could not create new IBM service")
		}

		if err := svc.SetServiceURL(env.IBMLangAPIUrl()); err != nil {
			logrus.WithError(err).Fatalln("Could not set IBM service URL")
		}

		d.service = svc
	}
}

// Analyze submits a request to analyze the entities in a transcript associated
// with an episode.
func (d *Detective) Analyze(seasonNumber int, episodeNumber int, overwrite bool) {
	onEpisode := func(ep *whodunit.Episode) {
		a := newAnalysis(ep, d)
		a.Create(overwrite)
	}

	if err := whodunit.Solve(seasonNumber, episodeNumber, onEpisode); err != nil {
		log.WithError(err).Errorln("Error analyzing episode(s)")
	}
}

func (d *Detective) FileReport(seasonNumber int, episodeNumber int, outputDir string) {
	onEpisode := func(ep *whodunit.Episode) {
		a := newAnalysis(ep, d)
		a.WriteCSV(outputDir)
	}

	if err := whodunit.Solve(seasonNumber, episodeNumber, onEpisode); err != nil {
		log.WithError(err).Errorln("Error analyzing episode(s)")
	}
}

// CloseCase closes the GCP NLP client.
func (d *Detective) CloseCase() {
	if d.cloudService == CloudServiceGCP {
		if err := d.client.Close(); err != nil {
			log.WithError(err).Errorln("Error closing case")
		}
	}
	log.Println("Case closed")
}

// Investigate logs the episode statuses.
func Investigate(cloudService CloudService, status whodunit.AssetStatus) {
	assetType := assetTypeForCloudService(cloudService)
	table := whodunit.NewStatusTable(assetType, status)
	table.Log()
}

func assetTypeForCloudService(cloudService CloudService) whodunit.AssetType {
	if cloudService == CloudServiceGCP {
		return whodunit.AssetTypeGCPAnalysis
	}
	return whodunit.AssetTypeIBMAnalysis
}
