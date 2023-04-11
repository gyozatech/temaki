package middlewares

import (
	"net/http"
)

func RecoverPanicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logError(err)
				logErr(w.Write([]byte(`{"code":500,"Message":"Internal Server Error"}`)))
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func logErr(n int, err error) {
	if err != nil {
		logError(err)
	}
}
