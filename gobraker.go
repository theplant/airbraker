package airbraker

import (
	"log"
	"os"
	"strconv"

	gobrake "gopkg.in/airbrake/gobrake.v1"
)

// initGobraker loads airbrake configuration from environment
// and returns a gobrake notifier.
//
// Will returns nil if no airbrake configuration or airbrake
// configuration is invalid.
func initGobraker() (airbrake *gobrake.Notifier) {
	projectIDStr := os.Getenv("AIRBRAKE_PROJECT_ID")
	token := os.Getenv("AIRBRAKE_TOKEN")
	env := os.Getenv("AIRBRAKE_ENV")

	if env == "" {
		env = "dev"
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if token == "" {
		return
	} else if err != nil {
		log.Printf("Airbrake Configuration Error: %v\n", err)
	} else {
		airbrake = gobrake.NewNotifier(projectID, token)
		airbrake.SetContext("environment", env)
		log.Printf("Logging errors to Airbrake '%s' env on project %v\n", env, projectID)
	}

	return
}
