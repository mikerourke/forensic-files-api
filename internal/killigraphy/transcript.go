package killigraphy

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mikerourke/forensic-files-api/internal/hearnoevil"
	"github.com/mikerourke/forensic-files-api/internal/whodunit"
)

// Transcript represents the text file extrapolated from the recognition.
type Transcript struct {
	*whodunit.Episode
}

// NewTranscript returns a new instance of a transcript
func NewTranscript(ep *whodunit.Episode) *Transcript {
	return &Transcript{
		Episode: ep,
	}
}

// Read returns the contents of the transcript file.
func (t *Transcript) Read() string {
	contents, err := ioutil.ReadFile(t.FilePath())
	if err != nil {
		log.WithError(err).Errorln("Error reading transcript file")
		return ""
	}

	return string(contents)
}

// Create creates a transcript file from a recognition.
func (t *Transcript) Create() {
	if t.Exists() {
		log.WithField("file", t.FileName()).Warnln(
			"Transcript already exists, skipping")
		return
	}

	contents := t.recognitionContents()
	if contents == "" {
		return
	}

	file, err := os.Create(t.FilePath())
	if err != nil {
		log.WithError(err).Errorln("Error creating transcript file")
		return
	}
	defer file.Close()

	if _, err = io.WriteString(file, contents); err != nil {
		log.WithError(err).Errorln("Error writing transcript file")
		return
	}

	if err := file.Sync(); err != nil {
		log.WithError(err).Errorln("Error syncing transcript file")
	}

	log.WithField("file", t.FileName()).Infoln("Transcript successfully written")
}

func (t *Transcript) recognitionContents() string {
	r := hearnoevil.NewRecognition(t.Episode)
	if !r.Exists() {
		log.WithField("file", r.FileName()).Warnln(
			"Recognition not found, skipping")
		return ""
	}

	if t.Exists() {
		log.WithField("file", t.FileName()).Warnln(
			"Transcript already exists, skipping")
	}

	results, err := r.ReadResults()
	if err != nil {
		log.WithError(err).Fatalln("Error getting recognition results")
	}

	lines := make([]string, 0)
	for _, result := range results {
		for _, alt := range result.Alternatives {
			content := *alt.Transcript
			words := strings.Fields(content)
			var confidence float32
			if alt.Confidence != nil {
				confidence = float32(*alt.Confidence)
			}

			if confidence >= 0.7 && len(words) > 2 {
				validWords := make([]string, 0)
				for i, word := range words {
					if !strings.Contains(word, "%HESITATION") {
						if i == 0 {
							validWords = append(validWords, strings.Title(word))
						} else {
							validWords = append(validWords, word)
						}
					}
				}

				validWords = append(validWords, ".")
				validLine := strings.Join(validWords, " ")
				validLine = strings.ReplaceAll(validLine, " .", ".")
				validLine = strings.ReplaceAll(validLine, " ,", ",")
				lines = append(lines, validLine)
			}
		}
	}

	return strings.Join(lines, "\n")
}

// Exists return true if the transcript file exists in the `/assets` directory.
func (t *Transcript) Exists() bool {
	return t.AssetExists(whodunit.AssetTypeTranscript)
}

// FilePath returns the path to the transcript file in the `/assets` directory.
func (t *Transcript) FilePath() string {
	return t.AssetFilePath(whodunit.AssetTypeTranscript)
}

// FileName returns the name of the transcript file in the `/assets` directory.
func (t *Transcript) FileName() string {
	return t.AssetFileName(whodunit.AssetTypeTranscript)
}
