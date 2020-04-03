package hearnoevil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

type callbackServer struct{}

func newCallbackServer() *callbackServer {
	return &callbackServer{}
}

// Start starts an HTTP server that listens for responses from the
// speech-to-text service. The server runs on port 9000 and is used to validate
// registered callback URLs or write recognition results to JSON files.
func (cs *callbackServer) Start() {
	log.Infoln("Starting callback URL server on port 9000")

	handler := func(w http.ResponseWriter, r *http.Request) {
		challengeString := r.URL.Query().Get("challenge_string")
		if challengeString != "" {
			cs.onRegister(w, challengeString)
			return
		}
		cs.onResponse(r)
	}

	http.HandleFunc("/", handler)
	log.Fatalln(http.ListenAndServe(":9000", nil))
}

// onRegister responds to the request to register a new callback
// URL and adheres to the requirements specified in the IBM cloud documentation
// at https://cloud.ibm.com/docs/services/speech-to-text?topic=speech-to-text-async#register.
func (cs *callbackServer) onRegister(w http.ResponseWriter, challengeString string) {
	log.Infoln("Received callback registration, sending response")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(challengeString))
	if err != nil {
		log.WithError(err).Errorln("Error writing response")
	}
}

// onResponse writes the results of a recognition job to a JSON file in
// `/assets/recognitions`.
func (cs *callbackServer) onResponse(r *http.Request) {
	log.Infoln("Recognition response received")

	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.WithError(err).Errorln("Error reading body of request")
	}

	var jobContents speechtotextv1.RecognitionJob
	err = json.Unmarshal(bytes, &jobContents)
	if err != nil {
		log.WithError(err).Errorln("Error unmarshalling JSON")
	}

	// If the UserToken isn't present on the job (it should always be, but just
	// in case), use a random ID:
	userToken := *jobContents.UserToken
	if userToken == "" {
		id, _ := uuid.NewUUID()
		userToken = id.String()
	}

	log.WithField("file", userToken).Infoln("Writing results to file")

	err = crimeseen.WriteJSONToAssets(
		"recognitions",
		userToken+".json",
		jobContents.Results,
	)

	if err != nil {
		log.WithError(err).Errorln("Error writing JSON to assets")
	}

	log.WithField("file", userToken).Infoln(
		"Successfully wrote recognition to JSON")
}
