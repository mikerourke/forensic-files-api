package hearnoevil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

// StartTranscriptionServer starts an HTTP server that listens for responses from the
// speech-to-text service. The server runs on port 9000 and writes responses
// to a JSON file when the service sends the response to a callback URL.
func StartTranscriptionServer() {
	log.Info("Loading environment variables")
	crimeseen.LoadDotEnv()
	log.Infoln("Starting callback URL server on port 9000")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// See https://cloud.ibm.com/docs/services/speech-to-text?topic=speech-to-text-async#register
	challengeString := r.URL.Query().Get("challenge_string")
	if challengeString != "" {
		log.Info("Received callback registration, sending response")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challengeString))
		return
	}

	log.Info("Transcription received")

	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.WithField("error", err).Error("Error reading body of request")
	}

	var jobContents speechtotextv1.RecognitionJob
	err = json.Unmarshal(bytes, &jobContents)
	if err != nil {
		log.WithField("error", err).Error("Error unmarshalling JSON")
	}

	// If the UserToken isn't present on the job (it should always be, but just
	// in case), use a random ID:
	userToken := *jobContents.UserToken
	if userToken == "" {
		id, _ := uuid.NewUUID()
		userToken = id.String()
	}

	log.WithField("file", userToken).Info("Writing results to file")

	err = crimeseen.WriteJSONToAssets(
		"recognitions",
		userToken+".json",
		jobContents.Results,
	)

	if err != nil {
		log.WithField("error", err).Error("Error writing JSON to assets")
	}

	log.Info("Successfully wrote contents to JSON")
}
