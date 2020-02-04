package ytdownload

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
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
	SortIndex     int
	VideoID       string
}

// DownloadVideos parses the YouTube episode URLs from the `/assets/youtube-links.json` file
// and downloads each episode to the `/videos` directory.
func DownloadVideos() {
	jsonBytes := readJsonFile()
	var parsedContents jsonEpisodesBySeason

	err := json.Unmarshal(jsonBytes, &parsedContents)
	if err != nil {
		log.Fatal(err)
	}

	allEpisodes := parseEpisodesFromJson(parsedContents)
	episodeCount := len(allEpisodes)
	for i, episode := range allEpisodes {
		fmt.Printf("Downloading %v of %v", i + 1, episodeCount)
		downloadEpisode(episode)
	}
}

func readJsonFile() []byte {
	pwd, _ := os.Getwd()
	youtubeLinks := filepath.Join(pwd, "assets", "youtube-links.json")

	jsonFile, err := os.Open(youtubeLinks)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue
}

func parseEpisodesFromJson(jsonContents jsonEpisodesBySeason) []*episode {
	var allEpisodes []*episode

	sortIndex := 0
	for season, jsonEpisodes := range jsonContents {
		for i, jsonEpisode := range jsonEpisodes {
			seasonNumber, _ := strconv.Atoi(season)
			nameItems := strings.Split(jsonEpisode.Name, " | ")
			episodeId := videoID(jsonEpisode)

			episode := &episode{
				Title:         nameItems[2],
				SeasonNumber:  seasonNumber,
				EpisodeNumber: i + 1,
				SortIndex:     sortIndex,
				VideoID:       episodeId,
			}

			allEpisodes = append(allEpisodes, episode)
			sortIndex += 1
		}
	}

	sort.Slice(allEpisodes, func(i, j int) bool {
		return allEpisodes[i].SortIndex < allEpisodes[j].SortIndex
	})

	return allEpisodes
}

func downloadEpisode(ep *episode) {
	outPath := outputFilePath(ep)

	fmt.Printf(
		"Downloading season %v, episode %v: %v to %v",
		ep.SeasonNumber,
		ep.EpisodeNumber,
		ep.Title,
		outPath)

	cmd := exec.Command("youtube-dl",
		"-o", outPath,
		ep.VideoID)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.Printf("error downloading %s: %s", ep.Title, err)
	}

	fmt.Printf("Download successful, waiting 5 minutes")
	time.Sleep(time.Minute * 5)
}

func outputFilePath(ep *episode) string {
	parentDirPath := seasonDirPath(ep.SeasonNumber)
	episodePrefix := episodeNumberForFile(ep)
	fileName := episodePrefix + "-" + ep.Title + ".%(ext)s"

	return filepath.Join(parentDirPath, fileName)
}

func episodeNumberForFile(ep *episode) string {
	epString := strconv.Itoa(ep.EpisodeNumber)

	if ep.EpisodeNumber < 10 {
		return "0" + epString
	}

	return epString
}

func seasonDirPath(season int) string {
	err := ensureOutputDirExists()
	if err != nil {
		log.Fatal(err)
	}

	seasonName := strconv.Itoa(season)
	fullDirPath := filepath.Join("videos", "season-"+seasonName)

	err = os.Mkdir(fullDirPath, os.ModePerm)

	if err != nil && !isExistsError(err) {
		log.Fatal(err)
	}

	return fullDirPath
}

func ensureOutputDirExists() error {
	pwd, _ := os.Getwd()
	outputDir := filepath.Join(pwd, "videos")

	err := os.Mkdir(outputDir, os.ModePerm)

	if err != nil && !isExistsError(err) {
		return err
	}

	return nil
}

func isExistsError(err error) bool {
	return strings.Contains(err.Error(), "file exists")
}

func videoID(ep jsonEpisode) string {
	parsedURL, _ := url.Parse(ep.URL)
	q := parsedURL.Query()
	return q.Get("v")
}
