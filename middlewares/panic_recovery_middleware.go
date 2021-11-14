package middlewares

import (
	"github.com/gyozatech/noodlog"
	//httputil "coverage/util/http"

	"net/http"
)

var logger ErrLogger

func RecoverPanicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				InitDefaultErrorLogger()
				logger.Error(err)
				logErr(w.Write([]byte(`{"code":500,"Message":"Internal Server Error"}`)))
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func logErr(n int, err error) {
	if err != nil {
		logger.Error(err)
	}
}

type ErrLogger interface {
	Error(message ...interface{})
}

func SetLogger(l ErrLogger) {
	logger = l
}

func InitDefaultErrorLogger() {
	if logger == nil {
		logger = noodlog.NewLogger().EnableTraceCaller()
	}
}
