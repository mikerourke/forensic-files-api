package hearnoevil

import (
	"os"
	"time"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

// speechToTextService returns an instance of the speech-to-text service that
// can be used to register callback URLs and create recognition jobs.
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
		log.WithField(
			"error", err,
		).Fatal("Error initializing speech to text service")
	}

	// The default timeout is 30 seconds. Depending on the file, that might not
	// cut the mustard. We're increasing it to 90 seconds to make sure the files
	// go through:
	speechToText.Service.Client.Timeout = time.Second * 90

	err = speechToText.SetServiceURL(os.Getenv("IBM_STT_URL"))
	if err != nil {
		log.WithField("error", err).Fatal("Error setting service URL")
	}

	return speechToText
}
