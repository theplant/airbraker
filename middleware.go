package airbraker

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Recover returns Gin middleware that reports all `panic`s to Airbrake
func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		defer func() {
			r := recover()
			if r == nil {
				return
			}

			var err error
			if e, ok := r.(error); !ok {
				err = fmt.Errorf("%v", r)
			} else {
				err = e
			}

			// not using goroutine here in order to keep the whole backtrace in
			// airbrake report
			Notify(err, ctx.Request)
			panic(r)
		}()

		ctx.Next()
	}
}
