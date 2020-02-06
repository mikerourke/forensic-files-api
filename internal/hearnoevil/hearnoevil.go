// Package hearnoevil sends the audio files to the speech-to-text service for
// transcribing.
package hearnoevil

import (
	"os/exec"

	"github.com/mikerourke/forensic-files-api/internal/waterlogged"
)

var log = waterlogged.ServiceLogger("hearnoevil")

func TranscribeEpisodes() {
	checkForNgrok()
}

func checkForNgrok() {
	cmd := exec.Command("ngrok", "--version")
	err := cmd.Run()
	if err != nil {
		panic("Could not find ngrok executable, it may not be installed")
	}
}
