package hearnoevil

import (
	"os"
	"time"

	"github.com/IBM/go-sdk-core/core"
	"github.com/mikerourke/forensic-files-api/internal/crimeseen"
	"github.com/sirupsen/logrus"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

// RegisterCallbackURL registers a callback URL with the speech to text service
// that will receive responses.
func RegisterCallbackURL(callbackURL string) {
	speechToText := speechToTextService()

	result, _, err := speechToText.RegisterCallback(
		&speechtotextv1.RegisterCallbackOptions{
			CallbackURL: core.StringPtr(callbackURL),
		},
	)

	if result != nil {
		log.WithFields(logrus.Fields{
			"url":    *result.URL,
			"status": *result.Status,
		}).Info("Callback registration complete")
	}

	if err != nil {
		log.WithField("error", err).Fatal("Error registering callback URL")
	}

	log.WithField("url", callbackURL).Info("Callback URL registered")
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
