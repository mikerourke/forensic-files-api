package hearnoevil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
)

// StartCallbackURLServer starts an HTTP server that listens for responses from the
// speech-to-text service. The server runs on port 9000 and writes responses
// to a JSON file when the service sends the transcribed response to a callback
// URL.
func StartCallbackURLServer() {
	crimeseen.LoadDotEnv()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Info(r.URL.Path)
	if r.URL.Path == "/results" {
		// TODO: Write `challenge_string` from response (when registering a
		//		callback URL) to a JSON file and reference it when needed.
		//		We should probably also write the actual callback URL to this file
		//		so we can unregister it when restarting ngrok.
	}

	log.Info("Transcription received, writing to assets")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.WithField("error", err).Error("Error reading body of request")
	}

	var contents interface{}
	err = json.Unmarshal(b, &contents)
	if err != nil {
		log.WithField("error", err).Error("Error unmarshalling JSON")
	}

	err = crimeseen.WriteJSONToAssets("", "test.json", contents)
	if err != nil {
		log.WithField("error", err).Error("Error writing JSON to assets")
	}

	log.Info("Successfully wrote contents to JSON")
}
