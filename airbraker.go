// Package airbraker is a notifier "provider" that provides
// a way to report runtime error. It uses gobrake notifier
// by default.
package airbraker

import (
	"log"
	"net/http"
)

// Notifier defines an interface for reporting error.
type Notifier interface {
	Notify(interface{}, *http.Request) error
}

// Airbrake stores a Notifier that is used in Notify function.
var Airbrake Notifier

func init() {
	airbrake := initGobraker()

	if airbrake != nil {
		Airbrake = airbrake
	} else {
		log.Println("[WARNING] No Airbrake Notifier. Logging to `log` instead.")
	}
}

// Notify proxies the reporting notice to notifier.
//
// Will returns nil if no Notifier specified.
func Notify(err interface{}, req *http.Request) error {
	if Airbrake != nil {
		return Airbrake.Notify(err, req)
	}
	log.Printf("[AIRBRAKE] %v", err)
	return nil
}
