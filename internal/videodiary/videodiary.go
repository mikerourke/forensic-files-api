// Package videodiary downloads all of the Forensic Files episodes from YouTube
// and drops them in the `/assets/videos` directory.
package videodiary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
	"github.com/sirupsen/logrus"
)

type jsonEpisode struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type jsonEpisodesBySeason map[string][]jsonEpisode

type episode struct {
	Title         string
	SeasonNumber  int
	EpisodeNumber int
	VideoHash     string
}

var log = waterlogged.ServiceLogger("videodiary")

// DownloadEpisodes parses the YouTube episode URLs from the `/assets/youtube-links.json`
// file and downloads each episode to the `/assets/videos` directory.
func DownloadEpisodes() {
	checkForYouTubeDL()

	allEpisodes := parseEpisodesFromJSON()
	for _, ep := range allEpisodes {
		if ep.VideoHash != "" {
			downloadEpisode(ep, true)
		}
	}
}

// DownloadEpisode downloads the specified episode number from the specified
// season number.
func DownloadEpisode(seasonNumber int, episodeNumber int) {
	checkForYouTubeDL()

	allEpisodes := parseEpisodesFromJSON()
	for _, ep := range allEpisodes {
		if ep.SeasonNumber == seasonNumber && ep.EpisodeNumber == episodeNumber {
			downloadEpisode(ep, false)
		}
	}
}

// LogMissingEpisodes logs the episodes that haven't been downloaded to the
// command line.
func LogMissingEpisodes() {
	allEpisodes := parseEpisodesFromJSON()

	missingCount := 0
	for _, episode := range allEpisodes {
		if episode.VideoHash == "" {
			fmt.Printf(
				"Season: %v \t Episode: %v \t Title: %v\n",
				episode.SeasonNumber,
				episode.EpisodeNumber,
				episode.Title,
			)
			missingCount++
		}
	}
	fmt.Printf("Total count missing: %v\n", missingCount)
}

func checkForYouTubeDL() {
	cmd := exec.Command("youtube-dl", "--version")
	err := cmd.Run()
	if err != nil {
		panic("Could not find youtube-dl executable, it may not be installed")
	}
}

func parseEpisodesFromJSON() []*episode {
	log.Info("reading JSON file with YouTube URLs")
	jsonContents, err := readYouTubeLinksJSON()
	if err != nil {
		log.WithField("error", err).Fatalln("Error reading YouTube URLs file")
	}

	var allEpisodes []*episode

	for season, jsonEpisodes := range jsonContents {
		for i, jsonEpisode := range jsonEpisodes {
			seasonNumber, _ := strconv.Atoi(season)
			nameItems := strings.Split(jsonEpisode.Name, " | ")

			ep := &episode{
				Title:         nameItems[2],
				SeasonNumber:  seasonNumber,
				EpisodeNumber: i + 1,
				VideoHash:     extractHash(jsonEpisode),
			}

			allEpisodes = append(allEpisodes, ep)
		}
	}

	return allEpisodes
}

func readYouTubeLinksJSON() (episodes jsonEpisodesBySeason, err error) {
	jsonFile, err := os.Open(
		filepath.Join(crimeseen.AssetsDirPath, "youtube-links.json"))
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var jsonContents jsonEpisodesBySeason
	err = json.Unmarshal(bytes, &jsonContents)
	if err != nil {
		return nil, err
	}

	return jsonContents, nil

}

func extractHash(ep jsonEpisode) string {
	if ep.URL == "" {
		return ""
	}

	parsedURL, _ := url.Parse(ep.URL)
	q := parsedURL.Query()
	return q.Get("v")
}

func downloadEpisode(ep *episode, isPaused bool) {
	outPath := outputFilePath(ep)
	if crimeseen.FileExists(outPath) {
		return
	}

	log.WithFields(logrus.Fields{
		"season":  ep.SeasonNumber,
		"episode": ep.EpisodeNumber,
		"title":   ep.Title,
		"path":    outPath,
	}).Infoln("Downloading video from YouTube")

	cmd := exec.Command("youtube-dl",
		"-o", outPath,
		ep.VideoHash)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"title": ep.Title,
			"path":  outPath,
		}).Errorln("Error downloading video")
	}

	// We're hedging our bets here to make sure we don't exceed some kind of rate limit:
	if isPaused {
		log.Println("Download successful, waiting 1 minute")
		time.Sleep(time.Minute * 1)
	}
}

// outputFilePath returns the full file path with the file name in the format
// XX-YY-title-of-file.mp4 where XX is the season number and YY is the episode
// number.
func outputFilePath(ep *episode) string {
	parentDirPath := seasonDirPath(ep.SeasonNumber)
	seasonPrefix := crimeseen.PaddedNumberString(ep.SeasonNumber)
	episodePrefix := crimeseen.PaddedNumberString(ep.EpisodeNumber)
	casedTitle := strings.ToLower(ep.Title)
	casedTitle = strings.ReplaceAll(casedTitle, " ", "-")

	fileName := seasonPrefix + "-" + episodePrefix + "-" + casedTitle + ".mp4"

	return filepath.Join(parentDirPath, fileName)
}

func seasonDirPath(seasonNumber int) string {
	err := crimeseen.Mkdirp(crimeseen.VideosDirPath)
	if err != nil {
		log.WithField("error", err).Fatal("Error creating output directory")
	}

	seasonName := strconv.Itoa(seasonNumber)
	fullDirPath := filepath.Join("assets", "videos", "season-"+seasonName)

	if err := crimeseen.Mkdirp(fullDirPath); err != nil {
		log.WithFields(logrus.Fields{
			"season": seasonNumber,
			"error":  err,
		}).Fatalln("Error creating season directory")
	}

	return fullDirPath
}
