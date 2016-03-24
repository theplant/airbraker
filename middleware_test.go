package airbraker_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/theplant/airbraker"
	au "github.com/theplant/airbraker/utils"
)

var errHandlerException = errors.New("panic on handler")

func TestRecoverMiddleware(t *testing.T) {
	server := newRecoverTestServer()
	defer func() {
		server.Close()
	}()

	bufferNotifier := &au.BufferNotifier{}
	originAirbrake := airbraker.Airbrake
	airbraker.Airbrake = bufferNotifier
	defer func() {
		airbraker.Airbrake = originAirbrake
	}()

	if len(bufferNotifier.Notices) != 0 {
		t.Fatalf("Notices must be empty.")
	}

	_, err := http.Get(server.URL + "/recover")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if len(bufferNotifier.Notices) != 1 {
		t.Fatalf("Unexpected notices length, got %d.", len(bufferNotifier.Notices))
	}

	if bufferNotifier.Notices[0].Error != errHandlerException {
		t.Fatalf("Got unexpected error: %v ", bufferNotifier.Notices[0].Error)
	}
}

// newRecoverTestServer prepares a test HTTP server that has the Recover
// middleware configured at `/recover`
func newRecoverTestServer() *httptest.Server {
	engine := gin.New()

	recoverGroup := engine.Group("/recover")

	// Catch handler panics so that test can continue.
	recoverGroup.Use(func(ctx *gin.Context) {
		defer func() {
			recover()
		}()
		ctx.Next()
	})

	// Recover middleware is executed first.
	recoverGroup.Use(airbraker.Recover())

	recoverGroup.Any("/", func(context *gin.Context) {
		panic(errHandlerException)
	})

	mux := http.NewServeMux()
	mux.Handle("/", engine)
	server := httptest.NewServer(mux)

	return server
}
