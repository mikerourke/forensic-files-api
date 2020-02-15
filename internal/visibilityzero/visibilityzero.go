// Package visibilityzero is used to loop through the downloaded episodes and
// extract the audio that will be sent to the speech-to-text service.
package visibilityzero

import (
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

var (
	log            = waterlogged.ServiceLogger("visibilityzero")
	processedCount int
)

// ExtractAudio loops through all of the `/video` season directories, extracts the
// audio from the .mp4 file using ffmpeg, and drops it into the `/assets/audio`
// directory for the corresponding season.
func ExtractAudio() {
	checkForFFmpeg()
	extractAudioFromAllSeasons()
}

func checkForFFmpeg() {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		panic("Could not find ffmpeg executable, it may not be installed")
	}
}

func extractAudioFromAllSeasons() {
	err := crimeseen.Mkdirp(filepath.Join(crimeseen.AudioDirPath))
	if err != nil {
		log.WithField("error", err).Fatal("Error creating audio directory")
	}

	processedCount = 0

	for season := 1; season <= crimeseen.SeasonCount; season++ {
		seasonDir := "season-" + strconv.Itoa(season)
		createSeasonAudioDir(seasonDir)

		err = filepath.Walk(
			filepath.Join(crimeseen.VideosDirPath, seasonDir),
			extractAudioFromSeason,
		)

		if err != nil {
			log.WithFields(logrus.Fields{
				"season": season,
				"error":  err,
			}).Fatal("Error walking season video directory")
		}
	}
}

func extractAudioFromSeason(path string, info os.FileInfo, err error) error {
	if strings.HasSuffix(path, ".mp4") {
		// Every 10 videos, take a 5 minute breather. ffmpeg makes the
		// fans go bananas on my laptop:
		if processedCount != 0 && processedCount%10 == 0 {
			log.Infoln("Taking a breather or else I'm going to take off")
			time.Sleep(time.Minute * 5)
		}

		audioPath := audioFilePath(path)

		if !crimeseen.FileExists(audioPath) {
			extractAudioFromEpisode(path, audioPath)
			processedCount++
		}
	}

	if err != nil {
		log.WithFields(logrus.Fields{
			"name":  info.Name(),
			"error": err,
		}).Error("Error in walk function")
		return err
	}

	return nil
}

func extractAudioFromEpisode(videoPath string, audioPath string) {
	log.WithFields(logrus.Fields{
		"video": filepath.Base(videoPath),
	}).Info("Extracting audio from video file")

	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		audioPath)

	err := cmd.Run()

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"video": filepath.Base(videoPath),
		}).Error("Error extracting audio")
	}
}

func createSeasonAudioDir(seasonDir string) {
	err := crimeseen.Mkdirp(filepath.Join(crimeseen.AudioDirPath, seasonDir))
	if err != nil {
		log.WithFields(logrus.Fields{
			"name":  seasonDir,
			"error": err,
		}).Fatal("Error creating season audio directory")
	}
}

func audioFilePath(videoPath string) string {
	dir, file := filepath.Split(videoPath)
	seasonDir := filepath.Base(dir)
	mp3FileName := strings.Replace(file, ".mp4", ".mp3", -1)
	return filepath.Join(crimeseen.AudioDirPath, seasonDir, mp3FileName)
}
