package crimeseen

import (
	"os"

	"github.com/alexsasharegan/dotenv"
)

// Env contains methods used to get data from the environment.
type Env struct{}

// NewEnv returns an instance of Env, which contains methods to get specified
// environment variables. It panics if it fails because the environment variables
// are usually a hard requirement when running the functions that utilize them.
func NewEnv() *Env {
	err := dotenv.Load()
	if err != nil {
		panic("Failed to load .env: " + err.Error())
	}

	return &Env{}
}

// CallbackURL returns the callback URL used for getting responses from various
// services.
func (e *Env) CallbackURL() string {
	return os.Getenv("CALLBACK_URL")
}

// IBMAPIKey returns the API key for the IBM speech-to-text service.
func (e *Env) IBMAPIKey() string {
	return os.Getenv("IBM_STT_API_KEY")
}

// IBMAPIUrl returns the URL for the IBM speech-to-text service.
func (e *Env) IBMAPIUrl() string {
	return os.Getenv("IBM_STT_URL")
}
