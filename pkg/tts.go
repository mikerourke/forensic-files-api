package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/go-sdk-core/core"
	"github.com/alexsasharegan/dotenv"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
)

func LoadTTS() {
	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	authenticator := &core.IamAuthenticator{
		ApiKey: os.Getenv("IBM_TTS_API_KEY"),
	}

	options := &speechtotextv1.SpeechToTextV1Options{
		Authenticator: authenticator,
	}

	speechToText, err := speechtotextv1.NewSpeechToTextV1(options)

	if err != nil {
		panic(err)
	}

	speechToText.SetServiceURL(os.Getenv("IBM_TTS_URL"))

	result, _, err := speechToText.ListModels(
		&speechtotextv1.ListModelsOptions{},
	)
	if err != nil {
		panic(err)
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
}
