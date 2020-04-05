package tagasuspect

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

// Detective contains properties and methods used to start entity analysis
// jobs.
type Detective struct {
	client   *language.Client
	ctx      context.Context
	withFile option.ClientOption
}

var (
	env = crimeseen.NewEnv()
	log = waterlogged.New("tagasuspect")
)

// NewDetective returns a new instance of a detective.
func NewDetective() *Detective {
	ctx := context.Background()
	return &Detective{
		ctx:      ctx,
		withFile: option.WithCredentialsFile(env.GCPCredsPath()),
	}
}

// OpenCase creates a new instance of an NLP client. We're doing this here
// instead of each analysis call so we can persist the client for batch jobs.
func (d *Detective) OpenCase() {
	d.interrogate()

	client, err := language.NewClient(d.ctx, d.withFile)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not create new GCP client")
	}
	d.client = client
}

// Analyze submits a request to analyze the entities in a transcript associated
// with an episode.
func (d *Detective) Analyze(seasonNumber int, episodeNumber int) {
	onEpisode := func(ep *whodunit.Episode) {
		a := NewAnalysis(ep)
		a.Create(d.ctx, d.client)
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

// CloseCase closes the GCP NLP client.
func (d *Detective) CloseCase() {
	if err := d.client.Close(); err != nil {
		log.WithError(err).Errorln("Error closing case")
	}
	log.Println("Case closed")
}

func (d *Detective) interrogate() {
	if env.GCPCredsPath() == "" {
		log.Fatalln("GCP credentials file not specified in .env file")
	}
}
