package hearnoevil

import (
	"time"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

type sttService struct {
	*speechtotextv1.SpeechToTextV1
}

// Initialize returns an instance of the speech-to-text service that
// can be used to register callback URLs and create recognition jobs.
func newSTTService(env *crimeseen.Env) *sttService {
	authenticator := &core.IamAuthenticator{
		ApiKey: env.IBMAPIKey(),
	}

	options := &speechtotextv1.SpeechToTextV1Options{
		Authenticator: authenticator,
	}

	speechToText, err := speechtotextv1.NewSpeechToTextV1(options)

	if err != nil {
		log.WithError(err).Fatalln("Error initializing speech to text service")
	}

	// The default timeout is 30 seconds. Depending on the file, that might not
	// cut the mustard. We're increasing it to 90 seconds to make sure the files
	// go through:
	speechToText.Service.Client.Timeout = time.Second * 90

	err = speechToText.SetServiceURL(env.IBMAPIUrl())
	if err != nil {
		log.WithError(err).Fatalln("Error setting service URL")
	}

	return &sttService{speechToText}
}
