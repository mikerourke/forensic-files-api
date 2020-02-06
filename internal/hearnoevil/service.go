package hearnoevil

import (
	"os"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

var callbackURL string

func RegisterCallbackURL() {
	//service := speechToTextService()
	//secret := os.Getenv("CALLBACK_URL_SECRET")

	// TODO: Finish writing this method. We need to make sure the server is running
	//		before attempting to register a callback URL.
	//		See https://cloud.ibm.com/docs/services/speech-to-text?topic=speech-to-text-async#register
}

func speechToTextService() *speechtotextv1.SpeechToTextV1 {
	crimeseen.LoadDotEnv()

	authenticator := &core.IamAuthenticator{
		ApiKey: os.Getenv("IBM_STT_API_KEY"),
	}

	options := &speechtotextv1.SpeechToTextV1Options{
		Authenticator: authenticator,
	}

	speechToText, err := speechtotextv1.NewSpeechToTextV1(options)

	if err != nil {
		log.WithField("error", err).Fatal("Error initializing speech to text service")
	}

	err = speechToText.SetServiceURL(os.Getenv("IBM_STT_URL"))
	if err != nil {
		log.WithField("error", err).Fatal("Error setting service URL")
	}

	return speechToText
}
