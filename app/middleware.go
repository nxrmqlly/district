package app

import (
	"log"
	"net/http"
	"slices"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for _, mw := range slices.Backward(mws) {
		h = mw(h)
	}
	return h
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *logResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (ro *Router) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &logResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // go defaults to 200 OK
		}
		next.ServeHTTP(rw, r)
		elapsed := time.Since(start)

		log.Printf("%d %s %s %s", rw.statusCode, elapsed, r.Method, r.URL.Path)
	})
}
