// Package utils provides a BufferNotifier that implements
// airbraker.Notifier interface. It can be used for testing.
package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/theplant/airbraker"
)

// notice structurizes arguments of Notify function.
type notice struct {
	Error   interface{}
	Request *http.Request
}

// BufferNotifier implements airbraker.Notifier interface.
// It stores all notified notices.
type BufferNotifier struct {
	Notices []notice
}

// Notify part of airbraker.Notifier.
func (b *BufferNotifier) Notify(err interface{}, req *http.Request) error {
	b.Notices = append(b.Notices, notice{Error: err, Request: req})

	return nil
}

// NewBufferNotifier returns a new BufferNotifier and
// the default airbraker.Notifier.
func NewBufferNotifier() (bufferNotifier *BufferNotifier, originNotifier airbraker.Notifier) {
	bufferNotifier = &BufferNotifier{}

	originNotifier = airbraker.Airbrake

	SetNotifier(bufferNotifier)

	return
}

// SetNotifier sets the airbraker.Airbrake to the given
// airbraker.Notifier.
func SetNotifier(notifier airbraker.Notifier) {
	airbraker.Airbrake = notifier
}

// ClearBuffer clears the given notifier's buffer.
func ClearBuffer(notifier *BufferNotifier) {
	notifier.Notices = []notice{}
}

// AssertNotice assert a notice that contains the given keyword
// notice will be received in the following second.
func AssertNotice(notifier *BufferNotifier, keyword string) (result bool) {
	result = true

	waitForNotice(notifier, 1, 1*time.Second)
	if got, want := len(notifier.Notices), 1; got == want {
		if got := notifier.Notices[0].Error; !strings.Contains(fmt.Sprintf("%v", got), keyword) {
			result = false
			log.Printf(`got unexpected notices: "%v",  want: "%v"`, got, keyword)
		}
	}

	return
}

// waitForNotice is waiting util the receiving notice reach the
// given count.
//
// Will break immediately if timeout.
func waitForNotice(notifier *BufferNotifier, count int, timeout time.Duration) {
	now := time.Now()
	for {
		if len(notifier.Notices) == count || time.Since(now) > timeout {
			break
		}
	}
}
