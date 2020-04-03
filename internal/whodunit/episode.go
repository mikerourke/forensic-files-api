package whodunit

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
)

// Episode is the high-level representation of a file in the `/assets` directory.
// An Episode has an associated audio file, video file, recognition, etc.
type Episode struct {
	SeasonNumber  int    `json:"season"`
	EpisodeNumber int    `json:"episode"`
	Title         string `json:"title"`
	URL           string `json:"url"`
	season        *Season
}

// newEpisode is private because it should only be called from within Season.
func newEpisode(
	season *Season,
	episodeNumber int,
	title string,
	url string,
) *Episode {
	return &Episode{
		SeasonNumber:  season.SeasonNumber,
		EpisodeNumber: episodeNumber,
		Title:         title,
		URL:           url,
		season:        season,
	}
}

// NewEpisodeFromName returns a new instance of an Episode from parsing the
// specified name.
//
// For example, calling NewEpisodeFromName("03-02-knot-for-everyone") would
// return an Episode instance with season number 3, episode number 2, and a
// title of "knot-for-everyone".
func NewEpisodeFromName(name string) (*Episode, error) {
	// If the name is a file path, throw out the path.
	base := filepath.Base(name)

	// Split the name up by hyphens. The first and second elements of the slice
	// represent the season and episode respectively. The rest of the elements
	// need to be combined back together to form the title.
	values := strings.Split(base, "-")
	seasonNumber, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, err
	}

	episodeNumber, err := strconv.Atoi(values[1])
	if err != nil {
		return nil, err
	}

	titles := make([]string, 0)
	for i := 2; i < len(values); i++ {
		titles = append(titles, values[i])
	}
	title := strings.Join(titles, "-")

	// Ensure the file extension (if any) is stripped from the title.
	title = strings.Replace(title, filepath.Ext(title), "", -1)

	return &Episode{
		SeasonNumber:  seasonNumber,
		EpisodeNumber: episodeNumber,
		Title:         title,
		URL:           "",
		season:        NewSeason(seasonNumber),
	}, nil
}

// WriteToRecognitionFile writes the specified contents to a new JSON file in
// the `/recognitions` directory.
func (e *Episode) WriteToRecognitionFile(contents interface{}) error {
	path := e.AssetFilePath(AssetTypeRecognition)
	return crimeseen.WriteJSONFile(path, contents)
}

// DisplayTitle returns the Title property separated by spaces with title case.
func (e *Episode) DisplayTitle() string {
	return strings.Title(strings.ReplaceAll(e.Title, "-", " "))
}

// VideoHash returns the YouTube URL's video hash used to identify the video
// to download.
func (e *Episode) VideoHash() string {
	if e.URL == "" {
		return ""
	}

	parsedURL, _ := url.Parse(e.URL)
	q := parsedURL.Query()
	return q.Get("v")
}

// AssetExists returns true if the file associated with the specified asset
// type exists in the `/assets` directory.
func (e *Episode) AssetExists(assetType AssetType) bool {
	return crimeseen.FileExists(e.AssetFilePath(assetType))
}

// AssetFilePath returns the absolute path to the asset file for the episode.
func (e *Episode) AssetFilePath(assetType AssetType) string {
	return filepath.Join(assetType.DirPath(),
		e.season.DirName(), e.AssetFileName(assetType))
}

// AssetFileName returns the file name of the episode with the appropriate
// extension based on the specified asset type.
func (e *Episode) AssetFileName(assetType AssetType) string {
	return fmt.Sprintf("%s%s", e.Name(), assetType.FileExt())
}

// Name returns the name of the episode in the common format used throughout
// the `/assets` directory: xx-yy-zz, where xx is the season, yy is the
// episode number, and zz is the title.
func (e *Episode) Name() string {
	return fmt.Sprintf("%s-%s-%s",
		crimeseen.PaddedNumberString(e.SeasonNumber),
		crimeseen.PaddedNumberString(e.EpisodeNumber),
		e.Title)
}
