package whodunit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
)

// Season represents a season directory in the assets directory along with
// associated episodes.
type Season struct {
	SeasonNumber int
	EpisodeMap   map[int]*Episode
}

// SeasonCount is the count of seasons in Forensic Files.
const SeasonCount = 14

// NewSeason returns a new instance of a Season with an empty episode map.
func NewSeason(seasonNumber int) *Season {
	return &Season{
		SeasonNumber: seasonNumber,
		EpisodeMap:   make(map[int]*Episode, 0),
	}
}

// PopulateEpisodes populates the season's episode map from the contents of the
// episodes JSON file in the `/assets` directory.
func (s *Season) PopulateEpisodes() error {
	jsonFile, err := os.Open(filepath.Join(AssetsDirPath, "episodes.json"))
	if err != nil {
		return err
	}

	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var result map[string][]Episode
	if err := json.Unmarshal(byteValue, &result); err != nil {
		return err
	}

	seasonName := crimeseen.PaddedNumberString(s.SeasonNumber)
	for _, ep := range result[seasonName] {
		s.EpisodeMap[ep.EpisodeNumber] = newEpisode(s,
			ep.EpisodeNumber, ep.Title, ep.URL)
	}

	return nil
}

// EpisodeCount returns the count of episodes in the season.
func (s *Season) EpisodeCount() int {
	return len(s.EpisodeMap)
}

// AllEpisodes returns a slice of Episode instances in the season.
func (s *Season) AllEpisodes() []*Episode {
	episodes := make([]*Episode, 0)
	for _, ep := range s.EpisodeMap {
		episodes = append(episodes, ep)
	}

	// Sort by episode number.
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].EpisodeNumber < episodes[j].EpisodeNumber
	})

	return episodes
}

// Episode returns an Episode instance associated with the specified episode
// number.
func (s *Season) Episode(episodeNumber int) *Episode {
	return s.EpisodeMap[episodeNumber]
}

// EnsureDir ensures that the season directory exists in the assets sub-directory
// associated with the specified asset type.
func (s *Season) EnsureDir(assetType AssetType) error {
	return crimeseen.Mkdirp(s.AssetDirPath(assetType))
}

// AssetDirPath returns the absolute path to the season directory for the
// associated asset.
func (s *Season) AssetDirPath(assetType AssetType) string {
	return filepath.Join(assetType.DirPath(), s.DirName())
}

// DirName returns the name of the directory for the associated season number.
func (s *Season) DirName() string {
	return fmt.Sprintf("season-%d", s.SeasonNumber)
}
