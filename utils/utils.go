// Package utils provides a BufferNotifier that implements
// airbraker.Notifier interface. It can be used for testing.
package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
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
	notices   []notice
	noticesMu sync.Mutex
}

func (b *BufferNotifier) setNotices(notices []notice) {
	b.noticesMu.Lock()
	defer b.noticesMu.Unlock()
	b.notices = notices
}

func (b *BufferNotifier) getNotices() []notice {
	b.noticesMu.Lock()
	defer b.noticesMu.Unlock()
	return b.notices
}

func (b *BufferNotifier) count() int {
	return len(b.getNotices())
}

// Notify part of airbraker.Notifier.
func (b *BufferNotifier) Notify(err interface{}, req *http.Request) error {
	b.setNotices(append(b.getNotices(), notice{Error: err, Request: req}))
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
	notifier.setNotices([]notice{})
}

// AssertNotice assert a notice that contains the given keyword
// notice will be received in the following second.
func AssertNotice(notifier *BufferNotifier, keyword string) (result bool) {
	result = true

	waitForNotice(notifier, 1, 1*time.Second)
	if got, want := notifier.count(), 1; got == want {
		if got := notifier.getNotices()[0].Error; !strings.Contains(fmt.Sprintf("%v", got), keyword) {
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
		if notifier.count() == count || time.Since(now) > timeout {
			break
		}
	}
}
