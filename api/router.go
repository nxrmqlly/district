package api

import "net/http"

type Router struct {
	mux *http.ServeMux
}

func New() *Router {
	ro := Router{
		mux: http.NewServeMux(),
	}
	ro.routes()

	return &ro
}

func (ro *Router) routes() {
	ro.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello, world"))
	})
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ro.mux.ServeHTTP(w, r)
}

func (ro *Router) Handle(pattern string, f http.HandlerFunc, mws ...Middleware) {
	ro.mux.Handle(pattern, Chain(f, mws...))
}
