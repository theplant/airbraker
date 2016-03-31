package airbraker_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/theplant/airbraker"
	au "github.com/theplant/airbraker/utils"
)

var (
	server              *httptest.Server
	notifier            *au.BufferNotifier
	errHandlerException = errors.New("panic on handler")
)

// TestMain create a new server instance and a buffer
// notifier for every test case.
func TestMain(m *testing.M) {
	server = newRecoverTestServer()
	notifier, _ = au.NewBufferNotifier()

	retCode := m.Run()

	server.Close()
	os.Exit(retCode)
}

func TestRecoverMiddleware(t *testing.T) {
	au.ClearBuffer(notifier)

	_, err := http.Get(server.URL + "/recover")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if !au.AssertNotice(notifier, errHandlerException.Error()) {
		t.Fatalf("Catch unexpected error: %v ", notifier)
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
