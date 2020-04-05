package hearnoevil

import (
	"time"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	stv1 "github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

type s2tInstance struct {
	*stv1.SpeechToTextV1
}

// newS2TInstance returns an instance of the speech-to-text service that
// can be used to register callback URLs and create recognition jobs.
func newS2TInstance(env *crimeseen.Env) *s2tInstance {
	authenticator := &core.IamAuthenticator{
		ApiKey: env.IBMSpeechToTextAPIKey(),
	}

	options := &stv1.SpeechToTextV1Options{
		Authenticator: authenticator,
	}

	speechToText, err := stv1.NewSpeechToTextV1(options)

	if err != nil {
		log.WithError(err).Fatalln("Error initializing speech to text service")
	}

	// The default timeout is 30 seconds. Depending on the file, that might not
	// cut the mustard. We're increasing it to 90 seconds to make sure the files
	// go through:
	speechToText.Service.Client.Timeout = time.Second * 90

	err = speechToText.SetServiceURL(env.IBMSpeechToTextAPIUrl())
	if err != nil {
		log.WithError(err).Fatalln("Error setting service URL")
	}

	return &s2tInstance{speechToText}
}
