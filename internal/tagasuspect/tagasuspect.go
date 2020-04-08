package tagasuspect

import (
	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
	nluv1 "github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
)

// Detective contains properties and methods used to start entity analysis
// jobs.
type Detective struct {
	options *nluv1.NaturalLanguageUnderstandingV1Options
	service *nluv1.NaturalLanguageUnderstandingV1
}

var (
	env = crimeseen.NewEnv()
	log = waterlogged.New("tagasuspect")
)

// NewDetective returns a new instance of a detective.
func NewDetective() *Detective {
	authenticator := &core.IamAuthenticator{
		ApiKey: env.IBMLangAPIKey(),
	}
	options := &nluv1.NaturalLanguageUnderstandingV1Options{
		Version:       "2019-07-12",
		Authenticator: authenticator,
	}

	return &Detective{
		options: options,
	}
}

// OpenCase creates a new instance of an NLU service. We're doing this here
// instead of each analysis call so we can persist the service for batch jobs.
func (d *Detective) OpenCase() {
	d.interrogate()

	svc, err := nluv1.NewNaturalLanguageUnderstandingV1(d.options)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not create new IBM service")
	}

	if err := svc.SetServiceURL(env.IBMLangAPIUrl()); err != nil {
		logrus.WithError(err).Fatalln("Could not set IBM service URL")
	}

	d.service = svc
}

// Analyze submits a request to analyze the entities in a transcript associated
// with an episode.
func (d *Detective) Analyze(seasonNumber int, episodeNumber int, overwrite bool) {
	onEpisode := func(ep *whodunit.Episode) {
		a := NewAnalysis(ep)
		a.Create(d.service, overwrite)
	}

	if err := whodunit.Solve(seasonNumber, episodeNumber, onEpisode); err != nil {
		log.WithError(err).Errorln("Error analyzing episode(s)")
	}
}

// Investigate logs the episode statuses.
func Investigate(status whodunit.AssetStatus) {
	table := whodunit.NewStatusTable(whodunit.AssetTypeAnalysis, status)
	table.Log()
}

func (d *Detective) interrogate() {
	if env.IBMLangAPIKey() == "" || env.IBMLangAPIUrl() == "" {
		log.Fatalln("IBM credentials file not specified in .env file")
	}
}
