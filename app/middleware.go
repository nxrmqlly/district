package app

import (
	"context"
	"log"
	"net/http"
	"slices"
	"time"
)

type contextKey struct{}

var (
	sessionKey contextKey
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

func (ro *Router) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("district_session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		se, err := ro.auth.GetSession(r.Context(), cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), sessionKey, se)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func (ro *Router) RequireAuth(next http.Handler) http.Handler {
	return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := SessionFromContext(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return 
		}
		
		next.ServeHTTP(w, r)
	})
}