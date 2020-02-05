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
	for _, episode := range allEpisodes {
		if episode.VideoHash != "" {
			downloadEpisode(episode)
		}
	}
}

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
	jsonBytes := readJSONFile()
	var jsonContents jsonEpisodesBySeason

	err := json.Unmarshal(jsonBytes, &jsonContents)
	if err != nil {
		log.WithField("err", err).Fatal("Error unmarshalling JSON")
	}

	var allEpisodes []*episode

	for season, jsonEpisodes := range jsonContents {
		for i, jsonEpisode := range jsonEpisodes {
			seasonNumber, _ := strconv.Atoi(season)
			nameItems := strings.Split(jsonEpisode.Name, " | ")

			episode := &episode{
				Title:         nameItems[2],
				SeasonNumber:  seasonNumber,
				EpisodeNumber: i + 1,
				VideoHash:     extractHash(jsonEpisode),
			}

			allEpisodes = append(allEpisodes, episode)
		}
	}

	return allEpisodes
}

func readJSONFile() []byte {
	log.Info("reading JSON file with YouTube URLs")
	youtubeLinks := filepath.Join(crimeseen.AssetsPath, "youtube-links.json")

	jsonFile, err := os.Open(youtubeLinks)

	if err != nil {
		log.WithField("error", err).Fatal("Error opening JSON file")
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.WithField("error", err).Fatal("Error reading JSON file")
	}

	return byteValue
}

func extractHash(ep jsonEpisode) string {
	if ep.URL == "" {
		return ""
	}

	parsedURL, _ := url.Parse(ep.URL)
	q := parsedURL.Query()
	return q.Get("v")
}

func downloadEpisode(ep *episode) {
	outPath := outputFilePath(ep)

	if crimeseen.FileExists(outPath) {
		return
	}

	log.WithFields(logrus.Fields{
		"season":  ep.SeasonNumber,
		"episode": ep.EpisodeNumber,
		"title":   ep.Title,
		"path":    outPath,
	}).Info("Downloading video from YouTube")

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
		}).Error("Error downloading video")
	}

	// We're hedging our bets here to make sure we don't exceed some kind of
	// rate limit:
	log.Println("Download successful, waiting 1 minute")
	time.Sleep(time.Minute * 1)
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

func seasonDirPath(season int) string {
	err := crimeseen.Mkdirp(crimeseen.VideosPath)
	if err != nil {
		log.WithField("error", err).Fatal("Error creating output directory")
	}

	seasonName := strconv.Itoa(season)
	fullDirPath := filepath.Join("assets", "videos", "season-"+seasonName)

	if err := crimeseen.Mkdirp(fullDirPath); err != nil {
		log.WithFields(logrus.Fields{
			"season": season,
			"error":  err,
		}).Fatal("Error creating season directory")
	}

	return fullDirPath
}
